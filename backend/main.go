package main

import (
	"fmt"
	"os"

	"github.com/furqanmk/reviews-browser/cmd"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: myapp <command>")
		fmt.Println("Commands: scheduler, api")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "schedulers":
		cmd.StartSchedulers()
	case "api":
		cmd.StartAPIServer()
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}
