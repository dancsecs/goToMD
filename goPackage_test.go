package main

import (
	"path/filepath"
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

type docInfoTest struct {
	action  string
	header  []string
	body    []string
	doc     []string
	oneLine string
}

func Test_GetDoc_GetInfo_InvalidDirectory(t *testing.T) {
	chk := szTest.CaptureNothing(t)
	defer chk.Release()

	absDir, err := filepath.Abs("INVALID_DIRECTORY")
	chk.NoErr(err)
	_, err = getInfo("INVALID_DIRECTORY", "TimesTwo")
	chk.Err(err, "open "+absDir+": no such file or directory")
}

func Test_GetDoc_GetInfo_InvalidObject(t *testing.T) {
	chk := szTest.CaptureNothing(t)
	defer chk.Release()

	_, err := getInfo("./sampleGoProject", "DOES_NOT_EXIST")
	chk.Err(err, "unknown package object: DOES_NOT_EXIST")
}

func Test_DocInfo_RunTests(t *testing.T) {
	chk := szTest.CaptureNothing(t)
	defer chk.Release()

	var docInfoTests = []docInfoTest{
		//  ----------------------------------------------------------------------
		{
			action: "TimesTwo",
			header: []string{
				"func TimesTwo(i int) int",
			},
			body: []string{
				"func TimesTwo(i int) int {",
				"    return i + i",
				"}",
			},
			doc: []string{
				"TimesTwo returns the value times two.",
			},
			oneLine: "func TimesTwo(i int) int",
		},
		//  ----------------------------------------------------------------------
		{
			action: "TimesThree",
			header: []string{
				"func TimesThree(i int) int",
			},
			body: []string{
				"func TimesThree(i int) int {",
				"    return i + i + i",
				"}",
			},
			doc: []string{
				"TimesThree returns the value times three.",
			},
			oneLine: "func TimesThree(i int) int",
		},
		//  ----------------------------------------------------------------------
		{
			action: "ConstDeclSingleCmtSingle",
			header: nil,
			body: []string{
				"const ConstDeclSingleCmtSingle = " +
					"\"single line declaration and comment\"",
			},
			doc: []string{
				"ConstDeclSingleCmtSingle has a single line comment.",
			},
			oneLine: "" +
				"const ConstDeclSingleCmtSingle = " +
				"\"single line declaration and comment\"",
		},
		//  ----------------------------------------------------------------------
		{
			action: "ConstDeclMultiCmtSingle",
			header: nil,
			body: []string{
				"const ConstDeclMultiCmtSingle = `multi line constant",
				"definition",
				"`",
			},
			doc: []string{
				"ConstDeclMultiCmtSingle has a single line comment with a multi line decl.",
			},
			oneLine: "" +
				"const ConstDeclMultiCmtSingle = `multi line constant ...",
		},
		//  ----------------------------------------------------------------------
		{
			action: "ConstDeclConstrCmtSingle",
			header: nil,
			body: []string{
				"const ConstDeclConstrCmtSingle = `multi line constant` + \"\n\" +",
				"    ConstDeclMultiCmtSingle + \" including other constants: \n\" +",
				"    ConstDeclSingleCmtSingle + \"\n\" + `",
				"=========end of constant=============",
				"`",
			},
			doc: []string{
				"ConstDeclConstrCmtSingle has a single line comment with a multi line decl.",
			},
			oneLine: "" +
				"const ConstDeclConstrCmtSingle = `multi line constant` + \"\n\" + ...",
		},
		//  ----------------------------------------------------------------------
		{
			action: "ConstDeclConstrCmtMulti",
			header: nil,
			body: []string{
				"const ConstDeclConstrCmtMulti = `multi line constant` + \"\n\" +",
				"    ConstDeclMultiCmtSingle + \" including other constants: \n\" +",
				"    ConstDeclSingleCmtSingle + \"\n\" + `",
				"=========end of constant=============",
				"`",
			},
			doc: []string{
				"ConstDeclConstrCmtMulti has a multi line comment with",
				"a multi line decl.",
			},
			oneLine: "const ConstDeclConstrCmtMulti =" +
				" `multi line constant` + \"\n\" + ...",
		},
		//  ----------------------------------------------------------------------
		{
			action: "StructureType.GetF1",
			header: []string{
				"func (s *StructureType) GetF1(",
				"    a, b, c int,",
				") string",
			},
			body: []string{
				"func (s *StructureType) GetF1(",
				"    a, b, c int,",
				") string {",
				"    const base10 = 10",
				"    t := a + c + b",
				"    return s.F1 + strconv.FormatInt(int64(t), base10)",
				"}",
			},
			doc: []string{
				"GetF1 is a method to a structure.",
			},
			oneLine: "func (s *StructureType) GetF1(a, b, c int) string",
		},
	}

	for _, tst := range docInfoTests {
		dInfo, err := getInfo("./sampleGoProject", tst.action)
		chk.NoErr(err)
		chk.StrSlice(dInfo.header, tst.header, "HEADER For action: ", tst.action)
		chk.StrSlice(dInfo.body, tst.body, "BODY For action: ", tst.action)
		chk.StrSlice(dInfo.doc, tst.doc, "DOC For action: ", tst.action)
		chk.Str(dInfo.oneLine(), tst.oneLine, "OneLine For action: ", tst.action)
	}
}
