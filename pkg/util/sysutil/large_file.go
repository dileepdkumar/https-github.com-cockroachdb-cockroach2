// Copyright 2018 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package sysutil

import (
	"os"

	"github.com/cockroachdb/errors"
)

func resizeLargeFileNaive(path string, bytes int64) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	sixtyFourMB := make([]byte, 64<<20)
	for bytes > 0 {
		z := sixtyFourMB
		if bytes < int64(len(z)) {
			z = sixtyFourMB[:bytes]
		}
		if _, err := f.Write(z); err != nil {
			return errors.Wrap(err, "write")
		}
		bytes -= int64(len(z))
	}
	return errors.Wrap(f.Sync(), "fsync")
}
