package nettask

import (
	"bufio"
	"fmt"
	"log"
	parse "nms/internal/jsonParse"
	utils "nms/internal/utils"
	"os"
	"sync"
)

var (
	packetsWaitingAck = make(map[byte]bool)
	pMutex            sync.Mutex
	packetID          = byte(1)
	packetMutex       sync.Mutex
)

var agentsIPs map[byte][4]byte
var numAgents int

func StartUDPServer(port string) {
	// include "| log.Lshortfile" in the log flags to include the file name and line of code in the log
	//log.SetFlags(log.Ltime | log.Lmicroseconds | log.Lshortfile)
	log.SetFlags(0)

	// inicializae reader
	reader := bufio.NewReader(os.Stdin)

	// Initialize the map
	agentsIPs = make(map[byte][4]byte)

	// make the server open an UDP connection via port 8081
	serverConn := utils.ResolveUDPAddrAndListen("localhost", "8081")

	// ask the user for the number of agents
	numAgents = askNumAgents(reader)

	// handle registrations from agents
	handleRegistrations(serverConn)

	// ask for json file path
	jsonPath := askJsonPath(reader)

	// parse the tasks from the json file
	jsonData := parse.GetDataFromJson(jsonPath)
	var taskList []parse.Task = parse.ParseDataFromJson(jsonData)

	// ask the user if he wants to proceed with tasks validation
	fmt.Print(utils.Blue, "Do you wish to proceed with tasks validation? (y/n): ", utils.Reset)
	input := parseString(reader)
	if input == "y" || input == "Y" {
		// validate the tasks
		parse.ValidateTaskList(taskList)
	}

	// ask the user if he wants to proceed with tasks validation
	fmt.Print(utils.Blue, "Do you wish to proceed with the tasks delegation? (y/n): ", utils.Reset)
	input = parseString(reader)
	if input == "n" || input == "N" {
		// close the server connection
		serverConn.Close()
		return
	}

	// connect and send tasks to agents
	go handleTasks(taskList)

	// receive metrics from agents
	go handleMetrics(serverConn)

	// give the user the option to consult the metrics/alerts file
	for {
		log.Println(utils.Blue, "Which file do you want to consult? (1 - Metrics file | 2- Alerts file | 3 - Exit)", utils.Reset)
		choice := parseString(reader)

		if choice == "1" {
			// consult the metrics file
			consultMetricsFile()
		}
		if choice == "2" {
			// consult the alerts file
			consultAlertsFile()
		}
		if choice == "3" {
			log.Println("Exiting...")
			break
		}
	}

	// close the server connection
	serverConn.Close()
}