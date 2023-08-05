package main

import (
	"fmt"
	"os"
	"time"

	"github.com/6543/logfile-open"
)

func main() {
	file, err := logfile.OpenFile("/tmp/a/alog", os.O_CREATE, 0o666)
	if err != nil {
		fmt.Printf("ERROR: %v", err)
		os.Exit(1)
	}
	defer file.Close()

	// _, err = file.Write([]byte("helloooo"))
	// if err != nil {
	// 	fmt.Printf("ERROR: %v", err)
	// 	os.Exit(2)
	// }

	time.Sleep(time.Minute * 2)

	err = file.Close()
	if err != nil {
		fmt.Printf("ERROR: %v", err)
		os.Exit(3)
	}

	time.Sleep(time.Second * 20)
}
