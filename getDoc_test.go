// Program runs a specific go test transforming the output
// to github compatible markdown.  This is used within this
// project to help automate keeping the README.md up to date
// when an example changes.

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

func Test_MarkGoCode(t *testing.T) {
	chk := szTest.CaptureNothing(t)
	defer chk.Release()

	chk.Str(
		markGoCode("ABC"),
		"```go\nABC\n```",
	)

	chk.Str(
		markGoCode("ABC\n"),
		"```go\nABC\n```",
	)
}

func Test_ExpandGoDcl_NoItems(t *testing.T) {
	chk := szTest.CaptureNothing(t)
	defer chk.Release()

	s, err := getDocDecl("./sampleGoProject/")
	chk.Err(err, "invalid action: a non-blank action is required")
	chk.Str(s, "")
}

func Test_ExpandGoDcl_Package(t *testing.T) {
	chk := szTest.CaptureNothing(t)
	defer chk.Release()

	s, err := getDocDecl("./sampleGoProject/package")
	chk.NoErr(err)
	chk.Str(
		s,
		markGoCode("package sampleGoProject\n"),
	)
}

func Test_ExpandGoDcl_InvalidItem(t *testing.T) {
	chk := szTest.CaptureNothing(t)
	defer chk.Release()

	s, err := getDocDecl("./sampleGoProject/unknownItem")
	chk.Err(err, "unknown package object: unknownItem")
	chk.Str(s, "")
}

func Test_ExpandGoDcl_OneItem(t *testing.T) {
	chk := szTest.CaptureNothing(t)
	defer chk.Release()

	s, err := getDocDecl("./sampleGoProject/TimesTwo")
	chk.AddSub(`package .*$`, "package ./sampleGoProject")
	chk.NoErr(err)
	chk.Str(
		s,
		markGoCode("func TimesTwo(i int) int\n"),
	)
}

func Test_ExpandGoDcl_TwoItems(t *testing.T) {
	chk := szTest.CaptureNothing(t)
	defer chk.Release()

	s, err := getDocDecl("./sampleGoProject/TimesTwo TimesThree")
	chk.AddSub(`package .*$`, "package ./sampleGoProject")
	chk.NoErr(err)
	chk.Str(
		s,
		markGoCode("func TimesTwo(i int) int\nfunc TimesThree(i int) int\n"),
	)
}

func Test_ExpandGoDclSingle_NoItems(t *testing.T) {
	chk := szTest.CaptureNothing(t)
	defer chk.Release()

	s, err := getDocDeclSingle("./sampleGoProject/")
	chk.Err(err, "invalid action: a non-blank action is required")
	chk.Str(s, "")
}

func Test_ExpandGoDclSingle_PackageNoItems(t *testing.T) {
	chk := szTest.CaptureNothing(t)
	defer chk.Release()

	s, err := getDocDeclSingle("./sampleGoProject/package")
	chk.NoErr(err)
	chk.Str(
		s,
		markGoCode("package sampleGoProject\n"),
	)
}

func Test_ExpandGoDclSingle_InvalidItem(t *testing.T) {
	chk := szTest.CaptureNothing(t)
	defer chk.Release()

	s, err := getDocDeclSingle("./sampleGoProject/unknownItem")
	chk.Err(err, "unknown package object: unknownItem")
	chk.Str(s, "")
}

func Test_ExpandGoDclSingle_OneItem(t *testing.T) {
	chk := szTest.CaptureNothing(t)
	defer chk.Release()

	s, err := getDocDeclSingle("./sampleGoProject/TimesTwo")
	chk.AddSub(`package .*$`, "package ./sampleGoProject")
	chk.NoErr(err)
	chk.Str(
		s,
		markGoCode("func TimesTwo(i int) int\n"),
	)
}

func Test_ExpandGoDclSingle_TwoItems(t *testing.T) {
	chk := szTest.CaptureNothing(t)
	defer chk.Release()

	s, err := getDocDeclSingle("./sampleGoProject/TimesTwo TimesThree")
	chk.AddSub(`package .*$`, "package ./sampleGoProject")
	chk.NoErr(err)
	chk.Str(
		s,
		markGoCode("func TimesTwo(i int) int\nfunc TimesThree(i int) int\n"),
	)
}

func Test_ExpandGoDclNatural_InvalidItem(t *testing.T) {
	chk := szTest.CaptureNothing(t)
	defer chk.Release()

	s, err := getDocDeclNatural("./sampleGoProject/unknownItem")
	chk.Err(err, "unknown package object: unknownItem")
	chk.Str(s, "")
}

func Test_ExpandGoDclNatural_OneItem(t *testing.T) {
	chk := szTest.CaptureNothing(t)
	defer chk.Release()

	s, err := getDocDeclNatural("./sampleGoProject/TimesTwo")
	chk.AddSub(`package .*$`, "package ./sampleGoProject")
	chk.NoErr(err)
	chk.Str(
		s,
		markGoCode(
			"// TimesTwo returns the value times two.\n"+
				"func TimesTwo(i int) int",
		),
	)
}

func Test_ExpandGoDclNatural_TwoItems(t *testing.T) {
	chk := szTest.CaptureNothing(t)
	defer chk.Release()

	s, err := getDocDeclNatural("./sampleGoProject/TimesTwo TimesThree")
	chk.AddSub(`package .*$`, "package ./sampleGoProject")
	chk.NoErr(err)
	chk.Str(
		s,
		markGoCode(
			"// TimesTwo returns the value times two.\n"+
				"func TimesTwo(i int) int\n"+
				"\n"+
				"// TimesThree returns the value times three.\n"+
				"func TimesThree(i int) int",
		),
	)
}