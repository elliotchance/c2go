// +build ignore

package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

func main() {
	if len(os.Args[1:]) == 0 {
		fmt.Println("Please add name of script")
		fmt.Println("For example:")
		fmt.Println("$ go run counter.go ./script_name.sh")
		os.Exit(0)
	}

	// Channel for end of work
	done := make(chan bool)

	// run ticker
	tick := make(chan bool)
	go func() { // infinite time ticker
		ticker := time.NewTicker(time.Second * 2)
		for range ticker.C {
			tick <- true
		}
	}()

	// run script without argguments
	go func() {
		// send `done` is else...
		defer func() {
			done <- true
		}()
		cmd := exec.Command(os.Args[1])
		// all output of script work send to system output
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stdout
		// execution
		if err := cmd.Start(); err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		if err := cmd.Wait(); err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
	}()

	for {
		select {
		case <-tick:
			fmt.Println(time.Now())
		case <-done:
			return
		}
	}

}
