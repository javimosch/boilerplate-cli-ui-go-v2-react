package main

import (
	"flag"
	"fmt"
	"os"
)

const Version = "1.0.0"

func main() {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "start":
		handleStart()
	case "stop":
		handleStop()
	case "status":
		handleStatus()
	case "version":
		handleVersion()
	case "help", "--help", "-h":
		printHelp()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		printHelp()
		os.Exit(1)
	}
}

func handleStart() {
	startCmd := flag.NewFlagSet("start", flag.ExitOnError)
	port := startCmd.Int("port", 8080, "Port for HTTP server")
	daemon := startCmd.Bool("daemon", false, "Run as daemon in background")
	startCmd.Parse(os.Args[2:])

	if *daemon {
		startDaemon(*port)
	} else {
		startServer(*port)
	}
}

func handleStop() {
	stopDaemon()
}

func handleStatus() {
	checkDaemonStatus()
}

func handleVersion() {
	fmt.Printf("boilerplate-cli-ui-go v%s\n", Version)
}

func printHelp() {
	fmt.Println("boilerplate-cli-ui-go - Go CLI with embedded web UI")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  boilerplate-cli-ui-go <command> [options]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  start       Start HTTP server with web UI")
	fmt.Println("  stop        Stop daemon server")
	fmt.Println("  status      Check daemon status")
	fmt.Println("  version     Show version information")
	fmt.Println("  help        Show this help message")
	fmt.Println()
	fmt.Println("Start Options:")
	fmt.Println("  -port int   Port for HTTP server (default 8080)")
	fmt.Println("  -daemon     Run as daemon in background")
	fmt.Println()
	fmt.Println("API Endpoints:")
	fmt.Println("  GET /            Web UI")
	fmt.Println("  GET /api/status  Server status (JSON)")
	fmt.Println("  GET /api/health  Health check (JSON)")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  boilerplate-cli-ui-go start")
	fmt.Println("  boilerplate-cli-ui-go start -port 3000")
	fmt.Println("  boilerplate-cli-ui-go start -daemon")
	fmt.Println("  boilerplate-cli-ui-go stop")
	fmt.Println("  boilerplate-cli-ui-go status")
}
