// Copyright 2023 @6543. All rights reserved.
// SPDX-License-Identifier: MIT

package logfile

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type wrapper struct {
	ctx            context.Context
	ctxCloser      func()
	receivedSignal chan os.Signal
	fileName       string
	origFile       *os.File
	lock           sync.RWMutex
	err            error
}

func (w *wrapper) Close() error {
	w.ctxCloser()
	if w.err != nil {
		return w.err
	}
	w.lock.Lock() // lock indefinite
	w.err = fmt.Errorf("writer got closed")
	return w.origFile.Close()
}

func (w *wrapper) Write(p []byte) (n int, err error) {
	if w.err != nil {
		return 0, w.err
	}
	w.lock.Lock()
	defer w.lock.Unlock()
	return w.origFile.Write(p)
}

func (w *wrapper) Read(p []byte) (n int, err error) {
	if w.err != nil {
		return 0, w.err
	}
	w.lock.RLock()
	defer w.lock.RUnlock()
	return w.origFile.Read(p)
}

var freeUp = func(w *wrapper) {
	// do magic so logfile can be rotated
	fmt.Println("got it got it")
}

func (w *wrapper) signalListener() {
	for {
		select {
		case <-w.ctx.Done():
			signal.Stop(w.receivedSignal)
			close(w.receivedSignal)
			return
		case <-w.receivedSignal:
			freeUp(w)
		}
	}
}

func OpenFile(name string, flag int, perm os.FileMode) (io.ReadWriteCloser, error) {
	file, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return nil, err
	}

	receivedSignal := make(chan os.Signal, 1)
	signal.Notify(receivedSignal, syscall.SIGUSR1)
	ctx, ctxCancel := context.WithCancel(context.Background())

	rwc := &wrapper{
		ctx:            ctx,
		ctxCloser:      ctxCancel,
		receivedSignal: receivedSignal,
		fileName:       name,
		origFile:       file,
	}

	go rwc.signalListener()

	return rwc, nil
}
