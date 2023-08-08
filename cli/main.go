package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/6543/logfile-open"
)

func main() {
	aFile := "/tmp/a/alog"

	ctx, ctxCancel := context.WithCancel(context.Background())

	fmt.Println("open file")
	file, err := logfile.OpenFileWithContext(ctx, aFile, 0o660)
	if err != nil {
		fmt.Printf("ERROR: %v", err)
		os.Exit(1)
	}
	defer file.Close()

	fmt.Println("write 1")
	_, err = file.Write([]byte("helloooo\n"))
	if err != nil {
		fmt.Printf("ERROR: %v", err)
		os.Exit(2)
	}

	time.Sleep(time.Second * 20)

	fmt.Println("write 2")
	_, err = file.Write([]byte("world\n"))
	if err != nil {
		fmt.Printf("ERROR: %v", err)
		os.Exit(2)
	}

	ctxCancel()
	// err = file.Close()
	// if err != nil {
	// 	fmt.Printf("ERROR: %v", err)
	// 	os.Exit(3)
	// }

	time.Sleep(time.Second * 20)
}
