package main

import (
  "flag"
  "fmt"
  "os"
)

func main() {
  // Crate network variable
  var network Network = NetworkHelper{}

  // Subcommands initialization
  serverCommand := flag.NewFlagSet("server", flag.ExitOnError)
  forwardCommand := flag.NewFlagSet("forward", flag.ExitOnError)
  checkCommand := flag.NewFlagSet("check", flag.ExitOnError)

  // Server subcommand flag pointers
  serverPortPtr := serverCommand.Int("port", DEFUALT_PORT, "Server port to open (e.g. server -port 8080)")
  serverFilePtr := serverCommand.Bool("file", false, "Use file on server (e.g. server -file)")

  // Forward subcommand flag pointers
  forwardTargetPtr := forwardCommand.String("target", "", "Server port to listen (required) (e.g. forward -target 127.0.0.1:8080)")
  forwardPortPtr := forwardCommand.Int("port", -1, "Server port to forward to (required) (e.g. forward -port 8080)")

  // Check subcommand flag pointerss
  var checkPortList portList
  checkCommand.Var(&checkPortList, "portList", "A comma seperated list of ports to be checked (e.g. check -portList 80,8080,4000)")
  checkIpPtr := checkCommand.Bool("ip", false, "Check your internal and external IP addresses (e.g. check -ip)")

  // Verify that a subcommand has been provided
  if len(os.Args) < 2 {
    exitRoutine("'server' or 'forward' or 'check' subcommand is required\nUse 'help' subcommand to see more details")
  }

  // Switch on the subcommand then parse the flags
  switch os.Args[1] {
  case "server":
		serverCommand.Parse(os.Args[2:])
	case "forward":
		forwardCommand.Parse(os.Args[2:])
  case "check":
  	checkCommand.Parse(os.Args[2:])
  case "help":
  	serverCommand.PrintDefaults()
  	forwardCommand.PrintDefaults()
    checkCommand.PrintDefaults()
  default:
    exitRoutine("No subcommand found")
  }

  // Check which subcommand should be ran
  if serverCommand.Parsed() {
  	err := network.startHttpServer(*serverPortPtr, *serverFilePtr)
    if err != nil {
      errorRoutine(err.Error())
    }
  } else if forwardCommand.Parsed() {
  	if (*forwardTargetPtr == "" || *forwardPortPtr < 0) {
  		exitRoutine("Both target and port number options are required (see help)");
  	}
		err := network.forwarding(*forwardTargetPtr, *forwardPortPtr)
    if err != nil {
      errorRoutine(err.Error())
    }
	} else if checkCommand.Parsed() {
    if !*checkIpPtr {
    	if checkPortList == nil {
        fmt.Printf("Unavailable ports in this machine are:\n%v\n", network.getAllUnavailablePorts())
      } else {
        fmt.Printf("The only unavailable ports among %v are:\n%v\n", checkPortList, network.getAllUnavailablePortsFromList(&checkPortList))
      }
    } else {
      internalIp, err := network.getInternalIp()
      if err != nil {
        errorRoutine(err.Error())
      }
      externalIp, err := network.getExternalIp()
      if err != nil {
        errorRoutine(err.Error())
      }
      fmt.Printf("Your internal IP: %v\nYour external IP: %v\n", internalIp, externalIp)
    }
  }
}

/*
 * Exit the main program
 */
func exitRoutine(message string) {
	if message != "" {
		fmt.Println(message)
	} else {
		flag.PrintDefaults()
	}
  os.Exit(1)
}

/*
 * Print error message
 */
func errorRoutine(message string) {
  fmt.Printf("Error: %s\n", message)
}
