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

func OpenFileWithContext(_ context.Context, name string, perm os.FileMode) (io.ReadWriteCloser, error) {
	return OpenFile(name, perm)
}
