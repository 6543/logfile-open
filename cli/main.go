package main

import (
	"fmt"
	"os"
	"time"

	"github.com/6543/logfile-open"
)

func main() {
	aFile := "/tmp/a/alog"

	f, _ := os.OpenFile(aFile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0o660)
	aloc := []byte("helloooo")
	_, err := f.Write(aloc)
	fmt.Println(err)
	f.Close()

	file, err := logfile.OpenFile(aFile, 0o660)
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
