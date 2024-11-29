package tcp

import (
	"log"
	utils "nms/internal/utils"
)

func StartTCPServer(port string, done chan struct{}) {
	listener := utils.ResolveTCPAddr("localhost", utils.SERVERTCP)

	for {
		select {
		case <-done:
			//log.Println("[TCP] Shutdown signal received. Closing TCP server.")
			return
		default:
			conn, err := listener.AcceptTCP()
			if err != nil {
				log.Println("[TCP] [ERROR] Unable to accept connection:", err)
				continue
			}
			go handleTCPConnection(conn)
		}
	}
}

/* // Função para tratar conexões TCP
func handleTCPConnection(conn net.Conn) {
	defer conn.Close()

	log.Println("[TCP] Established connection with Agent", conn.RemoteAddr())

	// decode and process registration request from agent

	regData := make([]byte, utils.BUFFERSIZE)
	_, err := conn.Read(regData)
	if err != nil {
		log.Fatalln("[TCP] [ERROR] Unable to read data:", err)
	}

	//reg, err := p.DecodeRegistration(regData[1:n])
	//if err != nil {
	//	log.Println("[TCP] [ERROR] Unable to decode registration data:", err)
	//	os.Exit(1)
	//}

	//if reg.NewID != 0 || reg.SenderIsServer {
	//	log.Println("[TCP] [ERROR] Invalid registration request parameters")
	//	// send NO_ACK
	//}

	// create, encode and send new registration request to agent

	//newReg := p.NewRegistrationBuilder().IsServer().SetNewID(1).Build()
	//newRegData := p.EncodeRegistration(newReg)

	//_, err = conn.Write(newRegData)
	//if err != nil {
	//	log.Println("[TCP] [ERROR] Unable to send new registration request", err)
	//	os.Exit(1)
	//}

	// send ACK
	log.Println("[TCP] New registration request sent")

} */