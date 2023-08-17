// Copyright 2023 @6543. All rights reserved.
// SPDX-License-Identifier: MIT

//go:build !windows

package logfile

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type wrapper struct {
	ctx            context.Context
	ctxCloser      func()
	receivedSignal chan os.Signal
	filePerm       fs.FileMode
	origFile       *os.File
	lock           sync.RWMutex
	err            error
}

func (w *wrapper) Close() error {
	w.ctxCloser()
	w.err = fmt.Errorf("writer got closed")
	return nil
}

func (w *wrapper) close() {
	signal.Stop(w.receivedSignal)
	close(w.receivedSignal)
	w.err = w.origFile.Close()
	if w.err == nil {
		w.err = w.ctx.Err()
	}
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

func (w *wrapper) freeUp() {
	w.lock.Lock()
	defer w.lock.Unlock()

	w.err = w.origFile.Close()
	if w.err != nil {
		w.ctxCloser()
		return
	}

	// TODO: do we need this or is it enough to close and open it?
	time.Sleep(time.Millisecond)

	w.origFile, w.err = os.OpenFile(w.origFile.Name(), os.O_CREATE|os.O_RDWR|os.O_APPEND, w.filePerm)
	if w.err != nil {
		w.ctxCloser()
		return
	}
}

func (w *wrapper) signalListener() {
	for {
		select {
		case <-w.ctx.Done():
			w.close()
			return
		case <-w.receivedSignal:
			w.freeUp()
		}
	}
}

func OpenFile(name string, perm os.FileMode) (io.ReadWriteCloser, error) {
	return OpenFileWithContext(context.Background(), name, perm)
}

func OpenFileWithContext(ctx context.Context, name string, perm os.FileMode) (io.ReadWriteCloser, error) {
	file, err := os.OpenFile(name, os.O_CREATE|os.O_RDWR|os.O_APPEND, perm)
	if err != nil {
		return nil, err
	}

	receivedSignal := make(chan os.Signal, 1)
	signal.Notify(receivedSignal, syscall.SIGUSR1)
	newCtx, ctxCancel := context.WithCancel(ctx)

	rwc := &wrapper{
		ctx:            newCtx,
		ctxCloser:      ctxCancel,
		receivedSignal: receivedSignal,
		filePerm:       perm,
		origFile:       file,
	}

	go rwc.signalListener()

	return rwc, nil
}
