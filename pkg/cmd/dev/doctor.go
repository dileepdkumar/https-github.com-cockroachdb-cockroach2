// Copyright 2020 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package main

import (
	"errors"
	"log"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
)

// makeDoctorCmd constructs the subcommand used to build the specified binaries.
func makeDoctorCmd(runE func(cmd *cobra.Command, args []string) error) *cobra.Command {
	return &cobra.Command{
		Use:     "doctor",
		Short:   "Check whether your machine is ready to build.",
		Long:    "Check whether your machine is ready to build.",
		Example: "dev doctor",
		Args:    cobra.ExactArgs(0),
		RunE:    runE,
	}
}

func (d *dev) doctor(cmd *cobra.Command, _ []string) error {
	ctx := cmd.Context()
	success := true

	// If we're running on macOS, we need to check whether XCode is installed.
	d.log.Println("doctor: running xcode check")
	if runtime.GOOS == "darwin" {
		stdout, err := d.exec.CommandContextSilent(ctx, "/usr/bin/xcodebuild", "-version")
		if err != nil {
			success = false
			log.Printf("Failed to run /usr/bin/xcodebuild -version. Got output: ")
			log.Printf("stdout:   %s", string(stdout))
			var cmderr *exec.ExitError
			if errors.As(err, &cmderr) {
				log.Printf("stderr:   %s", string(cmderr.Stderr))
			}
			log.Println(`You must have a full installation of XCode to build with Bazel.
A command-line tools instance does not suffice.
Please perform the following steps:
  1. Install XCode from the App Store.
  2. Launch Xcode.app at least once to perform one-time initialization of developer tools.
  3. Run ` + "`xcode-select -switch /Applications/Xcode.app/`.")
		}
	}

	// Check whether the dev config is active on a `bazel build`.
	if _, err := d.exec.CommandContextSilent(ctx, "bazel", "build", "//build/bazelutil:dev_toolchain_being_used"); err != nil {
		success = false
		log.Printf(`Your machine should be configured to build using the "dev" config.
Please add the following to your ~/.bazelrc:`)
		if runtime.GOOS == "darwin" && runtime.GOARCH == "amd64" {
			log.Printf("    build --config=devdarwinx86_64")
		} else {
			log.Printf("    build --config=dev")
		}
	}

	if success {
		log.Println("You are ready to build :)")
		return nil
	}
	return errors.New("please address the errors described above and try again")
}
