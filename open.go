// Copyright 2023 @6543. All rights reserved.
// SPDX-License-Identifier: MIT

package logfile

import "os"

func OpenFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	return os.OpenFile(name, flag, perm)
}
