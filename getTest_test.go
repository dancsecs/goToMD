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
	"fmt"
	"testing"

	"github.com/dancsecs/szTest"
)

func Test_ExpandGoTst(t *testing.T) {
	chk := szTest.CaptureNothing(t)
	defer chk.Release()

	_, err := expandGoTst("TEST_DIRECTORY_DOES_NOT_EXIST/")
	chk.Err(
		err,
		"relative directory must be specified in cmd: \"TEST_DIRECTORY_DOES_NOT_EXIST/\"",
	)

	_, err = expandGoTst("./TEST_DOES_NOT_EXIST")
	chk.Err(err, "no tests to run")
}

func Test_RunTestNotDirectory(t *testing.T) {
	chk := szTest.CaptureNothing(t)
	defer chk.Release()

	f := chk.CreateTmpFile(nil)
	chk.Panic(
		func() {
			_, _, _ = runTest(f, "")
		},
		"",
	)
}

func Test_RunTest(t *testing.T) {
	chk := szTest.CaptureStdout(t)
	defer chk.Release()

	c, s, err := runTest("./sampleGoProject", "package")

	chk.NoErr(err)
	chk.Str(c, "go test -v -cover ./sampleGoProject")
	fmt.Printf("%s\n", s)

	chk.AddSub("{{insOn}}", internalTestMarkInsOn)
	chk.AddSub("{{insOff}}", internalTestMarkInsOff)
	chk.AddSub("{{delOn}}", internalTestMarkDelOn)
	chk.AddSub("{{delOff}}", internalTestMarkDelOff)
	chk.AddSub("{{chgOn}}", internalTestMarkChgOn)
	chk.AddSub("{{chgOff}}", internalTestMarkChgOff)
	chk.AddSub("{{gotOn}}", internalTestMarkGotOn)
	chk.AddSub("{{gotOff}}", internalTestMarkGotOff)
	chk.AddSub("{{wntOn}}", internalTestMarkWntOn)
	chk.AddSub("{{wntOff}}", internalTestMarkWntOff)
	chk.AddSub("{{msgOn}}", internalTestMarkMsgOn)
	chk.AddSub("{{msgOff}}", internalTestMarkMsgOff)
	chk.AddSub("{{sepOn}}", internalTestMarkSepOn)
	chk.AddSub("{{sepOff}}", internalTestMarkSepOff)
	chk.AddSub("{{latexOn}}", `$\small{\texttt{`)
	chk.AddSub("{{latexOff}}", `}}$`)
	chk.AddSub(`\t\d+\.\d+s`, "\t0.0s")
	chk.AddSub(` `, hardSpace)
	chk.AddSub(`_`, hardUnderscore)
	chk.Stdout(`
    {{latexOn}}=== RUN   Test_PASS_SampleGoProject{{latexOff}}
    <br>
    {{latexOn}}--- PASS: Test_PASS_SampleGoProject (0.0s){{latexOff}}
    <br>
    {{latexOn}}=== RUN   Test_FAIL_SampleGoProject{{latexOff}}
    <br>
    {{latexOn}}    sample_test.go:28: unexpected int:{{latexOff}}
    <br>
    {{latexOn}}        {{msgOn}}2+2=5 (is true for big values of two){{msgOff}}:{{latexOff}}
    <br>
    {{latexOn}}        {{gotOn}}GOT: {{gotOff}}{{chgOn}}4{{chgOff}}{{latexOff}}
    <br>
    {{latexOn}}        {{wntOn}}WNT: {{wntOff}}{{chgOn}}5{{chgOff}}{{latexOff}}
    <br>
    {{latexOn}}    sample_test.go:29: unexpected string:{{latexOff}}
    <br>
    {{latexOn}}        {{gotOn}}GOT: {{gotOff}}{{insOn}}New in Got{{insOff}} Similar in ({{chgOn}}1{{chgOff}}) both{{latexOff}}
    <br>
    {{latexOn}}        {{wntOn}}WNT: {{wntOff}} Similar in ({{chgOn}}2{{chgOff}}) both{{delOn}}, new in Wnt{{delOff}}{{latexOff}}
    <br>
    {{latexOn}}    sample_test.go:35: Unexpected stdout Entry: got (1 lines) - want (1 lines){{latexOff}}
    <br>
    {{latexOn}}        {{chgOn}}0{{chgOff}}:{{chgOn}}0{{chgOff}} This output line {{delOn}}is{{delOff}}{{sepOn}}/{{sepOff}}{{insOn}}will be{{insOff}} different{{latexOff}}
    <br>
    {{latexOn}}    sample_test.go:39: unexpected string:{{latexOff}}
    <br>
    {{latexOn}}        {{gotOn}}GOT: {{gotOff}}{{chgOn}}Total{{chgOff}}: 6{{latexOff}}
    <br>
    {{latexOn}}        {{wntOn}}WNT: {{wntOff}}{{chgOn}}Sum{{chgOff}}: 6{{latexOff}}
    <br>
    {{latexOn}}--- FAIL: Test_FAIL_SampleGoProject (0.0s){{latexOff}}
    <br>
    {{latexOn}}FAIL{{latexOff}}
    <br>
    {{latexOn}}coverage: 100.0&#xFE6A; of statements{{latexOff}}
    <br>
    {{latexOn}}FAIL github.com/dancsecs/goToMD/sampleGoProject 0.0s{{latexOff}}
    <br>
    {{latexOn}}FAIL{{latexOff}}
    <br>
  `)
}
