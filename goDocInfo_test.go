package main

import (
	"testing"

	"github.com/dancsecs/szTest"
)

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

func Test_GetDoc_OneLine(t *testing.T) {
	chk := szTest.CaptureNothing(t)
	defer chk.Release()

	dInfo := &docInfo{}

	chk.Str(dInfo.oneLine(), "UNKNOWN DECLARATION")
}

func Test_GetInfo_Expand(t *testing.T) {
	chk := szTest.CaptureNothing(t)
	defer chk.Release()

	d, err := getInfo("./sampleGoProject", "TimesTwo")
	chk.NoErr(err)

	chk.Str(
		d.expand(szDocPrefix, "TimesTwo"),
		"<!--- goToMD::Bgn::doc::TimesTwo -->\n"+
			"```go\nfunc TimesTwo(i int) int\n```\n"+
			"\n"+
			"TimesTwo returns the value times two.\n"+
			"<!--- goToMD::End::doc::TimesTwo -->\n",
	)
}
