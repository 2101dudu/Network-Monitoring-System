package nettask

import (
	"log"
	utils "nms/internal/utils"
	"sync"
)

var (
	packetsWaitingAck = make(map[byte]bool)
	pMutex            sync.Mutex
	packetID          = byte(1)
	packetMutex       sync.Mutex
)

func StartUDPAgent() {
	// include "| log.Lshortfile" in the log flags to include the file name and line of code in the log
	// log.SetFlags(log.Ltime | log.Lmicroseconds | log.Lshortfile)
	log.SetFlags(0)

	// get the IP address of the agent
	ip := utils.GetIPAddress()

	// make the agent open an UDP connection via port 9091
	agentConn := utils.ResolveUDPAddrAndListen(ip, "9091")

	// make the agent connect to the server via UDP on port 8081
	serverConn := utils.ResolveUDPAddrAndDial("localhost", "8081")

	// register the agent with the server
	registerAgent(serverConn, ip)

	// handle tasks from the server
	handleTasks(agentConn)

	// close the agent connection
	agentConn.Close()
}