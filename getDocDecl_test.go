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

func Test_ExpandGoDcl_NoItems(t *testing.T) {
	chk := szTest.CaptureNothing(t)
	defer chk.Release()

	s, err := expandGoDcl("./sampleGoProject/")
	chk.Err(err, "invalid action: a non-blank action is required")
	chk.Str(s, "")
}

func Test_ExpandGoDcl_Package(t *testing.T) {
	chk := szTest.CaptureNothing(t)
	defer chk.Release()

	s, err := expandGoDcl("./sampleGoProject/package")
	chk.NoErr(err)
	chk.Str(
		s,
		"<!--- goToMD::Bgn::dcl::./sampleGoProject/package -->\n"+
			"```go\n"+
			"package sampleGoProject\n"+
			"```\n"+
			"<!--- goToMD::End::dcl::./sampleGoProject/package -->\n",
	)
}

func Test_ExpandGoDcl_InvalidItem(t *testing.T) {
	chk := szTest.CaptureNothing(t)
	defer chk.Release()

	s, err := expandGoDcl("./sampleGoProject/unknownItem")
	chk.Err(err, "unknown package object: unknownItem")
	chk.Str(s, "")
}

func Test_ExpandGoDcl_OneItem(t *testing.T) {
	chk := szTest.CaptureNothing(t)
	defer chk.Release()

	s, err := expandGoDcl("./sampleGoProject/TimesTwo")
	chk.AddSub(`package .*$`, "package ./sampleGoProject")
	chk.NoErr(err)
	chk.Str(
		s,
		"<!--- goToMD::Bgn::dcl::./sampleGoProject/TimesTwo -->\n```go\n"+
			"func TimesTwo(i int) int\n"+
			"```\n<!--- goToMD::End::dcl::./sampleGoProject/TimesTwo -->\n",
	)
}

func Test_ExpandGoDcl_TwoItems(t *testing.T) {
	chk := szTest.CaptureNothing(t)
	defer chk.Release()

	s, err := expandGoDcl("./sampleGoProject/TimesTwo TimesThree")
	chk.AddSub(`package .*$`, "package ./sampleGoProject")
	chk.NoErr(err)
	chk.Str(
		s,
		"<!--- goToMD::Bgn::dcl::./sampleGoProject/TimesTwo TimesThree -->\n```go\n"+
			"func TimesTwo(i int) int\nfunc TimesThree(i int) int\n"+
			"```\n<!--- goToMD::End::dcl::./sampleGoProject/TimesTwo TimesThree -->\n",
	)
}

func Test_ExpandGoDclSingle_NoItems(t *testing.T) {
	chk := szTest.CaptureNothing(t)
	defer chk.Release()

	s, err := expandGoDclSingle("./sampleGoProject/")
	chk.Err(err, "invalid action: a non-blank action is required")
	chk.Str(s, "")
}

func Test_ExpandGoDclSingle_PackageNoItems(t *testing.T) {
	chk := szTest.CaptureNothing(t)
	defer chk.Release()

	s, err := expandGoDclSingle("./sampleGoProject/package")
	chk.NoErr(err)
	chk.Str(
		s,
		"<!--- goToMD::Bgn::dcls::./sampleGoProject/package -->\n"+
			"```go\n"+
			"package sampleGoProject\n"+
			"```\n"+
			"<!--- goToMD::End::dcls::./sampleGoProject/package -->\n",
	)
}

func Test_ExpandGoDclSingle_InvalidItem(t *testing.T) {
	chk := szTest.CaptureNothing(t)
	defer chk.Release()

	s, err := expandGoDclSingle("./sampleGoProject/unknownItem")
	chk.Err(err, "unknown package object: unknownItem")
	chk.Str(s, "")
}

func Test_ExpandGoDclSingle_OneItem(t *testing.T) {
	chk := szTest.CaptureNothing(t)
	defer chk.Release()

	s, err := expandGoDclSingle("./sampleGoProject/TimesTwo")
	chk.AddSub(`package .*$`, "package ./sampleGoProject")
	chk.NoErr(err)
	chk.Str(
		s,
		"<!--- goToMD::Bgn::dcls::./sampleGoProject/TimesTwo -->\n```go\n"+
			"func TimesTwo(i int) int\n"+
			"```\n<!--- goToMD::End::dcls::./sampleGoProject/TimesTwo -->\n",
	)
}

func Test_ExpandGoDclSingle_TwoItems(t *testing.T) {
	chk := szTest.CaptureNothing(t)
	defer chk.Release()

	s, err := expandGoDclSingle("./sampleGoProject/TimesTwo TimesThree")
	chk.AddSub(`package .*$`, "package ./sampleGoProject")
	chk.NoErr(err)
	chk.Str(
		s,
		"<!--- goToMD::Bgn::dcls::./sampleGoProject/TimesTwo TimesThree -->\n```go\n"+
			"func TimesTwo(i int) int\nfunc TimesThree(i int) int\n"+
			"```\n<!--- goToMD::End::dcls::./sampleGoProject/TimesTwo TimesThree -->\n",
	)
}
