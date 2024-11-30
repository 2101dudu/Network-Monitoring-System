package udp

import (
	"log"
	"net"
	ack "nms/internal/packet/ack"
	"nms/internal/packet/metrics"
	"nms/internal/packet/task"
	"nms/internal/utils"
	"os/exec"
	"time"
)

func handleIperfClientTask(taskPayload []byte, agentConn *net.UDPConn, udpAddr *net.UDPAddr) {
	iperfClient, err := task.DecodeIperfClientPacket(taskPayload)
	if err != nil {
		log.Fatalln("[AGENT] [ERROR 85] Decoding ping packet")
	}

	if !task.ValidateHashIperfClientPacket(iperfClient) {
		noack := ack.NewAckBuilder().SetPacketID(iperfClient.PacketID).SetSenderID(utils.SERVERID).Build()
		hash := ack.CreateHashAckPacket(noack)
		noack.Hash = (string(hash))
		ack.EncodeAndSendAck(agentConn, udpAddr, noack)

		log.Println("[AGENT] [ERROR 100] Invalid hash in iperf client packet")
		return
	}

	newAck := ack.NewAckBuilder().SetPacketID(iperfClient.PacketID).SetSenderID(utils.SERVERID).HasAcknowledged().Build()
	hash := ack.CreateHashAckPacket(newAck)
	newAck.Hash = (string(hash))
	ack.EncodeAndSendAck(agentConn, udpAddr, newAck)

	// keep track of the start time
	startTime := time.Now()

	// execute the pingPacket's command
	cmd := exec.Command("sh", "-c", iperfClient.IperfClientCommand)

	outputData, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalln("[AGENT] [ERROR 86] Executing ping command")
	}

	preparedOutput := parseIperfOutput(iperfClient.Bandwidth, iperfClient.Jitter, iperfClient.PacketLoss, string(outputData))

	serverConn := utils.ResolveUDPAddrAndDial("localhost", "8081")

	metricsID := utils.ReadAndIncrementPacketID(&packetID, &packetMutex, true)
	newMetrics := metrics.NewMetricsBuilder().SetPacketID(metricsID).SetAgentID(agentID).SetTaskID(iperfClient.TaskID).SetTime(startTime.Format("15:04:05.000000000")).SetMetrics(preparedOutput).Build()

	hash = metrics.CreateHashMetricsPacket(newMetrics)
	newMetrics.Hash = (string(hash))

	packetData := metrics.EncodeMetrics(newMetrics)
	ack.SendPacketAndWaitForAck(metricsID, agentID, packetsWaitingAck, &pMutex, serverConn, nil, packetData, "[SERVER] [MAIN READ THREAD] Metrics packet sent", "[SERVER] [ERROR 36] Unable to send metrics packet")
}