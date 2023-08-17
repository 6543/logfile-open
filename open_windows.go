// Copyright 2023 @6543. All rights reserved.
// SPDX-License-Identifier: MIT

//go:build windows

package logfile

import (
	"context"
	"io"
	"os"
)

// as windows do not have any concept of signals we just forward the interface without any logic

func OpenFile(name string, perm os.FileMode) (io.ReadWriteCloser, error) {
	return os.OpenFile(name, os.O_CREATE|os.O_RDWR|os.O_APPEND, perm)
}

type wrapper struct {
	ctx            context.Context
	ctxCloser      func()
	origFile       io.ReadWriteCloser
}

func (w *wrapper) Close() error {
	err := w.origFile.Close()
	w.ctxCloser()
	w.ctxCloser = func() {}
	return err
}

func (w *wrapper) Write(p []byte) (n int, err error) {
	return w.origFile.Write(p)
}

func (w *wrapper) Read(p []byte) (n int, err error) {
	return w.origFile.Read(p)
}

func (w *wrapper) contextListener() {
	for {
		select {
		case <-w.ctx.Done():
			origFile.Close()
			return
	}
}

func OpenFileWithContext(ctx context.Context, name string, perm os.FileMode) (io.ReadWriteCloser, error) {
	file, err := OpenFile(name, perm)
	if err != nil {
		return nil, err
	}

	newCtx, ctxCancel := context.WithCancel(ctx)

	rwc := &wrapper{
		ctx:            newCtx,
		ctxCloser:      ctxCancel,
		origFile:       file,
	}

	go rwc.contextListener()

	return rwc, nil
}


func OpenFileWithContext(ctx context.Context, name string, perm os.FileMode) (io.ReadWriteCloser, error) {
	file, err := OpenFile(name, perm)

	go func() {
		file := file
		for {
			select {
			case <-w.ctx.Done():
				w.close()
				return
		}
	}()

	return file, err
}
