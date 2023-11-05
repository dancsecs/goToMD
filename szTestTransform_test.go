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

func Test_MarksToMarkdownHTML(t *testing.T) {
	chk := szTest.CaptureNothing(t)
	defer chk.Release()

	s, err := marksToMarkdownHTML("")
	chk.NoErr(err)
	chk.Str(s, "")

	s, err = marksToMarkdownHTML(szTest.SettingMarkInsOn())
	chk.Err(err, "no closing mark found for \"<{INS_OFF}>\" in \"\"")
	chk.Str(s, "")

	s, err = marksToMarkdownHTML(
		szTest.SettingMarkInsOn() + "-" + szTest.SettingMarkDelOn(),
	)
	chk.Err(
		err,
		"unexpected closing mark: Got: \"<{DEL_ON}>\"  Want: \"<{INS_OFF}>\"",
	)
	chk.Str(s, "")
}
