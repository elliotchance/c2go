// build:ignore

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
	go func() {
		ticker := time.NewTicker(time.Millisecond * 500)
		for range ticker.C {
			tick <- true
		}
	}()

	// run script without argguments
	go func() {
		defer func() {
			done <- true
		}()
		cmd := exec.Command(os.Args[1])
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stdout
		if err := cmd.Start(); err != nil {
			log.Fatal(err)
		}
		if err := cmd.Wait(); err != nil {
			log.Fatal(err)
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
