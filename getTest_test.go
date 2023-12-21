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

	_, err := getGoTst("TEST_DIRECTORY_DOES_NOT_EXIST/")
	chk.Err(
		err,
		"relative directory must be specified in cmd: \"TEST_DIRECTORY_DOES_NOT_EXIST/\"",
	)

	_, err = getGoTst("./TEST_DOES_NOT_EXIST")
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

	s, err := getGoTst("./sampleGoProjectOne/package ./sampleGoProjectTwo/package")

	chk.NoErr(err)
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
	chk.Stdout("" +
		markBashCode("go test -v -cover ./sampleGoProjectOne") + "\n" + `
    {{latexOn}}=== RUN   Test_PASS_sampleGoProjectOne{{latexOff}}
    <br>
    {{latexOn}}--- PASS: Test_PASS_sampleGoProjectOne (0.0s){{latexOff}}
    <br>
    {{latexOn}}=== RUN   Test_FAIL_sampleGoProjectOne{{latexOff}}
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
    {{latexOn}}--- FAIL: Test_FAIL_sampleGoProjectOne (0.0s){{latexOff}}
    <br>
    {{latexOn}}FAIL{{latexOff}}
    <br>
    {{latexOn}}coverage: 100.0&#xFE6A; of statements{{latexOff}}
    <br>
    {{latexOn}}FAIL github.com/dancsecs/goToMD/sampleGoProjectOne 0.0s{{latexOff}}
    <br>
    {{latexOn}}FAIL{{latexOff}}
    <br>

    ` +
		markBashCode("go test -v -cover ./sampleGoProjectTwo") + "\n" + `
    {{latexOn}}=== RUN   Test_PASS_sampleGoProjectTwo{{latexOff}}
    <br>
    {{latexOn}}--- PASS: Test_PASS_sampleGoProjectTwo (0.0s){{latexOff}}
    <br>
    {{latexOn}}=== RUN   Test_FAIL_sampleGoProjectTwo{{latexOff}}
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
    {{latexOn}}--- FAIL: Test_FAIL_sampleGoProjectTwo (0.0s){{latexOff}}
    <br>
    {{latexOn}}FAIL{{latexOff}}
    <br>
    {{latexOn}}coverage: 100.0&#xFE6A; of statements{{latexOff}}
    <br>
    {{latexOn}}FAIL github.com/dancsecs/goToMD/sampleGoProjectTwo 0.0s{{latexOff}}
    <br>
    {{latexOn}}FAIL{{latexOff}}
    <br>
  `)
}

func Test_GetTest_BuildTestCmds(t *testing.T) {
	chk := szTest.CaptureNothing(t)
	defer chk.Release()

	dirs := []string{}
	actions := []string{}
	d, a := buildTestCmds(dirs, actions)
	chk.StrSlice(d, nil)
	chk.StrSlice(a, nil)

	actions = append(actions, "D1A1")
	d, a = buildTestCmds(dirs, actions)
	chk.StrSlice(d, nil)
	chk.StrSlice(a, nil)

	dirs = append(dirs, "D1")
	d, a = buildTestCmds(dirs, actions)
	chk.StrSlice(d, []string{"D1"})
	chk.StrSlice(a, []string{"D1A1"})

	dirs = append(dirs, "D1")
	actions = append(actions, "D1A2")
	d, a = buildTestCmds(dirs, actions)
	chk.StrSlice(d, []string{"D1"})
	chk.StrSlice(a, []string{"D1A1 D1A2"})

	dirs = append(dirs, "D2")
	actions = append(actions, "D2A1")
	d, a = buildTestCmds(dirs, actions)
	chk.StrSlice(d, []string{"D1", "D2"})
	chk.StrSlice(a, []string{"D1A1 D1A2", "D2A1"})
}
