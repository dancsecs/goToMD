package main

/*
   Golang To Github Markdown: goToMD.
   Copyright (C) 2023  Leslie Dancsecs

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

import (
	"testing"

	"github.com/dancsecs/szTest"
)

func Test_SampleNoArgsNotDirectory(t *testing.T) {
	chk := szTest.CaptureNothing(t)
	defer chk.Release()

	chk.SetupArgsAndFlags([]string{
		"programName",
	})

	chk.AddSub(`(?s)\n.*$`, "\\nUsage Information")
	chk.Panic(
		main,
		"at least one file or directory must be specified\\nUsage Information",
	)
}

func Test_SampleInvalidFile(t *testing.T) {
	chk := szTest.CaptureNothing(t)
	defer chk.Release()

	fPath := chk.CreateTmpFile(nil)
	chk.SetupArgsAndFlags([]string{
		"programName",
		fPath,
	})

	chk.Panic(
		main,
		"file must have extension: .md.gtm",
	)
}

func Test_InvalidDefaultPermissions(t *testing.T) {
	chk := szTest.CaptureNothing(t)
	defer chk.Release()

	fPath := chk.CreateTmpFile(nil)
	chk.SetupArgsAndFlags([]string{
		"programName",
		"-p", "0744",
		fPath,
	})

	chk.AddSub(`(?s)\n.*$`, "\\nUsage Information")
	chk.Panic(
		main,
		"invalid default permissions specified\\nUsage Information",
	)
}

func Test_InvalidCleanAndReplace(t *testing.T) {
	chk := szTest.CaptureNothing(t)
	defer chk.Release()

	fPath := chk.CreateTmpFile(nil)
	chk.SetupArgsAndFlags([]string{
		"programName",
		"-r",
		"-c",
		fPath,
	})

	chk.AddSub(`(?s)\n.*$`, "\\nUsage Information")
	chk.Panic(
		main,
		"only one of -c and -r may be specified\\nUsage Information",
	)
}

func Test_InvalidOutoutDirectory(t *testing.T) {
	chk := szTest.CaptureNothing(t)
	defer chk.Release()

	fPath := chk.CreateTmpFile(nil)
	chk.SetupArgsAndFlags([]string{
		"programName",
		"-o", "DIRECTORY_DOES_NOT_EXIST",
		fPath,
	})

	chk.AddSub(`(?s)\n.*$`, "\\nUsage Information")
	chk.Panic(
		main,
		"invalid output directory specified: "+
			"DIRECTORY_DOES_NOT_EXIST\\nUsage Information",
	)
}
