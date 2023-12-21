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

func Test_Markdown_CleanMarkDownDocument(t *testing.T) {
	chk := szTest.CaptureNothing(t)
	defer chk.Release()

	md, err := cleanMarkDownDocument(
		szTestBgnPrefix + szDocPrefix + "action1 -->\n" +
			szTestEndPrefix + szTstPrefix + "action2 -->\n",
	)

	chk.Err(
		err,
		"out of sequence: End before begin: <!--- goToMD::End::tst::action2 -->",
	)
	chk.Str(md, "")
}

func Test_Markdown_CleanMarkDownDocumentMissingBlankAfterAuto(t *testing.T) {
	chk := szTest.CaptureNothing(t)
	defer chk.Release()

	md, err := cleanMarkDownDocument(
		szAutoPrefix + " -->\nThis is not a blnk line.",
	)

	chk.Err(
		err,
		"missing blank line in auto generated output",
	)
	chk.Str(md, "")
}

func Test_Markdown_UpdateMarkDownDocument(t *testing.T) {
	chk := szTest.CaptureNothing(t)
	defer chk.Release()

	md, err := updateMarkDownDocument("",
		szTestPrefix+szDocPrefix+"./INVALID_ROOT_DIRECTORY/action1 -->\n",
	)

	chk.Err(
		err,
		"invalid directory specified as: \"./INVALID_ROOT_DIRECTORY\"",
	)
	chk.Str(md, "")
}

func Test_Markdown_UpdateMarkDown_InvalidCommand(t *testing.T) {
	chk := szTest.CaptureNothing(t)
	defer chk.Release()

	md, err := updateMarkDownDocument("",
		szTestPrefix+"unknownCommand -->\n",
	)

	chk.Err(
		err,
		"unknown cmd: <!--- goToMD::unknownCommand -->",
	)
	chk.Str(md, "")
}

func Test_Markdown_Expand(t *testing.T) {
	chk := szTest.CaptureNothing(t)
	defer chk.Release()

	d, err := getInfo("./sampleGoProjectOne", "TimesTwo")
	chk.NoErr(err)

	chk.Str(
		expand(szDocPrefix,
			"TimesTwo",
			d.declGoLang()+"\n\n"+
				d.docMarkdown(),
		),
		"<!--- goToMD::Bgn::doc::TimesTwo -->\n"+
			"```go\nfunc TimesTwo(i int) int\n```\n"+
			"\n"+
			"TimesTwo returns the value times two.\n"+
			"<!--- goToMD::End::doc::TimesTwo -->\n",
	)
}

func Test_Markdown_Search(t *testing.T) {
	chk := szTest.CaptureNothing(t)
	defer chk.Release()

	chk.Int(action.search("a"), -1)
	chk.Int(action.search("z"), -1)
}
