// Copyright 2023 @6543. All rights reserved.
// SPDX-License-Identifier: MIT

package logfile

import (
	"fmt"
	"io"
	"os"
	"sync"
)

type wrapper struct {
	fileName string
	origFile *os.File
	lock     sync.RWMutex
	err      error
}

func (w *wrapper) Close() error {
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

func OpenFile(name string, flag int, perm os.FileMode) (io.ReadWriteCloser, error) {
	file, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return nil, err
	}
	return &wrapper{
		fileName: name,
		origFile: file,
	}, nil
}
