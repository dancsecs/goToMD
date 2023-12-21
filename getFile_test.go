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

func Test_GetFile_GetGoFileInvalid(t *testing.T) {
	chk := szTest.CaptureNothing(t)
	defer chk.Release()

	_, err := getGoFile("TEST_DIRECTORY_DOES_NOT_EXIST/")
	chk.Err(
		err,
		"relative directory must be specified in cmd: \"TEST_DIRECTORY_DOES_NOT_EXIST/\"",
	)

	_, err = getGoTst("./sampleGoProjectOne/TEST_DOES_NOT_EXIST")
	chk.Err(err, "no tests to run")
}

func Test_GetFile_GetGoFile(t *testing.T) {
	chk := szTest.CaptureNothing(t)
	defer chk.Release()

	d, err := getGoFile("./sampleGoProjectOne/crumb.go")
	chk.NoErr(err)
	chk.Str(
		d,
		""+
			markBashCode("cat ./sampleGoProjectOne/crumb.go")+
			"\n\n"+
			markGoCode("package sampleGoProjectOne"),
	)
}

func Test_GetFile_GetGoFile2(t *testing.T) {
	chk := szTest.CaptureNothing(t)
	defer chk.Release()

	d, err := getGoFile(
		"./sampleGoProjectOne/crumb.go ./sampleGoProjectTwo/crumb.go",
	)
	chk.NoErr(err)
	chk.Str(
		d,
		""+
			markBashCode("cat ./sampleGoProjectOne/crumb.go")+
			"\n\n"+
			markGoCode("package sampleGoProjectOne")+
			"\n\n"+
			markBashCode("cat ./sampleGoProjectTwo/crumb.go")+
			"\n\n"+
			markGoCode("package sampleGoProjectTwo")+
			"",
	)
}
