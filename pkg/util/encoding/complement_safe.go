// Copyright 2015 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

//go:build !386 && !amd64
// +build !386,!amd64

package encoding

func onesComplement(b []byte) {
	for i := range b {
		b[i] = ^b[i]
	}
}
