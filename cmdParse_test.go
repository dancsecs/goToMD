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

func Test_ParseCmd(t *testing.T) {
	chk := szTest.CaptureNothing(t)
	defer chk.Release()

	dir, action, err := parseCmd("")
	chk.Err(err, "relative directory must be specified in cmd: \"\"")
	chk.Str(dir, "")
	chk.Str(action, "")

	dir, action, err = parseCmd("/action")
	chk.Err(err, "relative directory must be specified in cmd: \"/action\"")
	chk.Str(dir, "")
	chk.Str(action, "")

	dir, action, err = parseCmd("./")
	chk.Err(
		err,
		"invalid action: a non-blank action is required",
	)
	chk.Str(dir, "")
	chk.Str(action, "")

	dir, action, err = parseCmd("sampleGoProject/action")
	chk.Err(
		err,
		"relative directory must be specified in cmd: \"sampleGoProject/action\"",
	)
	chk.Str(dir, "")
	chk.Str(action, "")

	dir, action, err = parseCmd("./sampleGoProject/action")
	chk.NoErr(err)
	chk.Str(dir, "./sampleGoProject")
	chk.Str(action, "action")
}

func Test_ParseCmds1(t *testing.T) {
	chk := szTest.CaptureNothing(t)
	defer chk.Release()

	dirs, actions, err := parseCmds("")
	chk.Nil(dirs)
	chk.Nil(actions)
	chk.Err(err, "relative directory must be specified in cmd: \"\"")

	dirs, actions, err = parseCmds("/action")
	chk.Nil(dirs)
	chk.Nil(actions)
	chk.Err(err, "relative directory must be specified in cmd: \"/action\"")

	dirs, actions, err = parseCmds("./")
	chk.Nil(dirs)
	chk.Nil(actions)
	chk.Err(err, "invalid action: a non-blank action is required")

	dirs, actions, err = parseCmds("sampleGoProject/action")
	chk.Nil(dirs)
	chk.Nil(actions)
	chk.Err(
		err,
		"relative directory must be specified in cmd: \"sampleGoProject/action\"",
	)

	dirs, actions, err = parseCmds("./sampleGoProject/action")
	chk.NoErr(err)
	chk.StrSlice(dirs, []string{"./sampleGoProject"})
	chk.StrSlice(actions, []string{"action"})
}

func Test_ParseCmds2(t *testing.T) {
	chk := szTest.CaptureNothing(t)
	defer chk.Release()

	dirs, actions, err := parseCmds("./sampleGoProject/action /action2")
	chk.Nil(dirs)
	chk.Nil(actions)
	chk.Err(err, "relative directory must be specified in cmd: \"/action2\"")

	dirs, actions, err = parseCmds("./sampleGoProject/action sampleGoProject/action")
	chk.Nil(dirs)
	chk.Nil(actions)
	chk.Err(err, "relative directory must be specified in cmd: \"sampleGoProject/action\"")

	dirs, actions, err = parseCmds("./sampleGoProject/action action2")
	chk.NoErr(err)
	chk.StrSlice(dirs, []string{"./sampleGoProject", "./sampleGoProject"})
	chk.StrSlice(actions, []string{"action", "action2"})

	dirs, actions, err = parseCmds("./sampleGoProject/action ./sampleGoProject/action2")
	chk.NoErr(err)
	chk.StrSlice(dirs, []string{"./sampleGoProject", "./sampleGoProject"})
	chk.StrSlice(actions, []string{"action", "action2"})
}
