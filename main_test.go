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
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dancsecs/szTest"
)

const sampleGoProjectPath = "." + string(os.PathSeparator) +
	"sampleGoProject" + string(os.PathSeparator)

var fullSampleGoProjectPath string

func init() {
	fullSampleGoProjectPath, _ = filepath.Abs(sampleGoProjectPath)
}

func Test_SampleGoProjectExpandNoTarget(t *testing.T) {
	chk := szTest.CaptureNothing(t)
	defer chk.Release()

	dir := chk.CreateTmpDir()
	chk.NoErr(setup(dir, "README.md.gtm", "sample_test.go", "sample.go"))

	chk.SetupArgsAndFlags([]string{
		"programName",
		filepath.Join(dir, "README.md.gtm"),
	})

	// Now Run the main function with no -f arg requiring confirmation
	chk.NoPanic(main)

	got, wnt, err := getTestFiles(dir, "README.md")
	chk.NoErr(err)
	chk.StrSlice(got, wnt)
}

func Test_SampleGoProjectExpandTargetCancel(t *testing.T) {
	chk := szTest.CaptureStdout(t)
	defer chk.Release()

	clearPackageCache()

	dir := chk.CreateTmpDir()
	chk.NoErr(
		setup(dir, "README.md", "README.md.gtm", "sample_test.go", "sample.go"),
	)

	chk.SetupArgsAndFlags([]string{
		"programName",
		filepath.Join(dir, "README.md.gtm"),
	})

	// Setup a N to be read from os.Stdin by main() for conformation to overwrite
	// the file.
	origStdin := os.Stdin
	defer func() {
		os.Stdin = origStdin
	}()

	r, w, err := os.Pipe()
	chk.NoErr(err)
	os.Stdin = r

	fmt.Fprintln(w, "N")

	// Run command expecting the overwrite to be cancelled.
	chk.NoPanic(main)
	chk.Stdout("Confirm overwrite of README.md (Y to overwrite)? " +
		"overwrite cancelled",
	)
}

func Test_SampleGoProjectExpandTargetOverwrite(t *testing.T) {
	chk := szTest.CaptureStdout(t)
	defer chk.Release()

	clearPackageCache()

	dir := chk.CreateTmpDir()
	chk.NoErr(
		setup(dir, "README.md", "README.md.gtm", "sample_test.go", "sample.go"),
	)

	chk.SetupArgsAndFlags([]string{
		"programName",
		filepath.Join(dir, "README.md.gtm"),
	})

	// Setup a N to be read from os.Stdin by main() for conformation to overwrite
	// the file.
	origStdin := os.Stdin
	defer func() {
		os.Stdin = origStdin
	}()

	r, w, err := os.Pipe()
	chk.NoErr(err)
	os.Stdin = r

	fmt.Fprintln(w, "Y")

	// Run command expecting the overwrite to be cancelled.
	chk.NoPanic(main)

	got, wnt, err := getTestFiles(dir, "README.md")
	chk.NoErr(err)
	chk.StrSlice(got, wnt)

	chk.Stdout("Confirm overwrite of README.md (Y to overwrite)?\\s")
}

func Test_SampleGoProjectExpandTargetOverwriteDir(t *testing.T) {
	chk := szTest.CaptureStdout(t)
	defer chk.Release()

	clearPackageCache()

	dir := chk.CreateTmpDir()
	chk.NoErr(
		setup(dir, "README.md", "README.md.gtm", "sample_test.go", "sample.go"),
	)

	chk.SetupArgsAndFlags([]string{
		"programName",
		dir,
	})

	// Setup a N to be read from os.Stdin by main() for conformation to overwrite
	// the file.
	origStdin := os.Stdin
	defer func() {
		os.Stdin = origStdin
	}()

	r, w, err := os.Pipe()
	chk.NoErr(err)
	os.Stdin = r

	fmt.Fprintln(w, "Y")

	// Run command expecting the overwrite to be cancelled.
	chk.NoPanic(main)

	got, wnt, err := getTestFiles(dir, "README.md")
	chk.NoErr(err)
	chk.StrSlice(got, wnt)

	chk.Stdout("Confirm overwrite of README.md (Y to overwrite)?\\s")
}

func Test_SampleGoProjectExpandTargetOverwriteDirVerbose(t *testing.T) {
	chk := szTest.CaptureLogAndStdout(t)
	defer chk.Release()

	clearPackageCache()

	dir := chk.CreateTmpDir()
	chk.NoErr(
		setup(dir, "README.md", "README.md.gtm", "sample_test.go", "sample.go"),
	)

	chk.SetupArgsAndFlags([]string{
		"programName",
		"-v",
		dir,
	})

	// Setup a N to be read from os.Stdin by main() for conformation to overwrite
	// the file.
	origStdin := os.Stdin
	defer func() {
		os.Stdin = origStdin
	}()

	r, w, err := os.Pipe()
	chk.NoErr(err)
	os.Stdin = r

	fmt.Fprintln(w, "Y")

	// Run command expecting the overwrite to be cancelled.
	chk.NoPanic(main)

	got, wnt, err := getTestFiles(dir, "README.md")
	chk.NoErr(err)
	chk.StrSlice(got, wnt)

	pName := filepath.Join(dir, "README.md.gtm")
	chk.Stdout(
		"filesToProcess:  "+pName,
		"Confirm overwrite of README.md (Y to overwrite)?\\s",
	)

	chk.Log("Expanding "+pName+" to: README.md",
		"Loading Package info for: .",
		"getInfo(\"package\")",
		"getInfo(\"TimesTwo\")",
		"getInfo(\"TimesThree\")",
		"getInfo(\"InterfaceType\")",
		"getInfo(\"StructureType\")",
		"getInfo(\"TimesTwo\")",
		"getInfo(\"TimesTwo\")",
		"getInfo(\"TimesThree\")",
		"getInfo(\"TimesTwo\")",
		"getInfo(\"TimesTwo\")",
		"getInfo(\"TimesThree\")",
		"getInfo(\"TimesTwo\")",
		"getInfo(\"TimesTwo\")",
		"getInfo(\"TimesThree\")",
	)
}

////////////

func Test_SampleGoProjectReplaceNoTarget(t *testing.T) {
	chk := szTest.CaptureNothing(t)
	defer chk.Release()

	dir := chk.CreateTmpDir()
	chk.NoErr(setup(dir, "sample_test.go", "sample.go"))

	fName := filepath.Join(dir, "README.md")

	chk.SetupArgsAndFlags([]string{
		"programName",
		"-r",
		fName,
	})

	chk.Panic(
		main,
		"stat "+filepath.Join(dir, "README.md")+": no such file or directory",
	)

	_, _, err := getTestFiles(dir, "README.md")
	chk.Err(
		err,
		"open "+fName+": no such file or directory",
	)
}

func Test_SampleGoProjectReplaceTargetCancel(t *testing.T) {
	chk := szTest.CaptureLogAndStdout(t)
	defer chk.Release()

	clearPackageCache()

	dir := chk.CreateTmpDir()
	chk.NoErr(
		setup(dir, "README.md", "sample_test.go", "sample.go"),
	)

	fName := filepath.Join(dir, "README.md")
	chk.SetupArgsAndFlags([]string{
		"programName",
		"-v",
		"-r",
		fName,
	})

	// Setup a N to be read from os.Stdin by main() for conformation to overwrite
	// the file.
	origStdin := os.Stdin
	defer func() {
		os.Stdin = origStdin
	}()

	r, w, err := os.Pipe()
	chk.NoErr(err)
	os.Stdin = r

	fmt.Fprintln(w, "N")

	// Run command expecting the overwrite to be cancelled.
	chk.NoPanic(main)

	chk.Stdout(
		"filesToProcess:  "+fName,
		"Confirm overwrite of "+fName+" (Y to overwrite)? "+
			"overwrite cancelled")

	chk.Log(
		"",
		"in place replacing of "+fName,
		"Loading Package info for: .",
		"getInfo(\"package\")",
		"getInfo(\"TimesTwo\")",
		"getInfo(\"TimesThree\")",
		"getInfo(\"InterfaceType\")",
		"getInfo(\"StructureType\")",
		"getInfo(\"TimesTwo\")",
		"getInfo(\"TimesTwo\")",
		"getInfo(\"TimesThree\")",
		"getInfo(\"TimesTwo\")",
		"getInfo(\"TimesTwo\")",
		"getInfo(\"TimesThree\")",
		"getInfo(\"TimesTwo\")",
		"getInfo(\"TimesTwo\")",
		"getInfo(\"TimesThree\")",
	)
}

func Test_SampleGoProjectReplaceTargetOverwrite(t *testing.T) {
	chk := szTest.CaptureStdout(t)
	defer chk.Release()

	clearPackageCache()

	dir := chk.CreateTmpDir()
	chk.NoErr(
		setup(dir, "README.md", "sample_test.go", "sample.go"),
	)

	fName := filepath.Join(dir, "README.md")
	chk.SetupArgsAndFlags([]string{
		"programName",
		"-r",
		fName,
	})

	// Setup a N to be read from os.Stdin by main() for conformation to overwrite
	// the file.
	origStdin := os.Stdin
	defer func() {
		os.Stdin = origStdin
	}()

	r, w, err := os.Pipe()
	chk.NoErr(err)
	os.Stdin = r

	fmt.Fprintln(w, "Y")

	// Run command expecting the overwrite to be cancelled.
	chk.NoPanic(main)

	got, wnt, err := getTestFiles(dir, "README.md")
	chk.NoErr(err)
	wnt[0] = strings.ReplaceAll(wnt[0], "** DO NOT MODIFY ** ", "")
	chk.StrSlice(got, wnt)

	_, _, err = getTestFiles(dir, "README.md.gtm")
	chk.Err(
		err,
		"open "+fName+".gtm: no such file or directory",
	)

	chk.Stdout("Confirm overwrite of " + fName + " (Y to overwrite)?\\s")
}

func Test_SampleGoProjectReplaceTargetOverwriteDir(t *testing.T) {
	chk := szTest.CaptureStdout(t)
	defer chk.Release()

	clearPackageCache()

	dir := chk.CreateTmpDir()
	chk.NoErr(
		setup(dir, "README.md", "sample_test.go", "sample.go"),
	)

	chk.SetupArgsAndFlags([]string{
		"programName",
		"-r",
		dir,
	})

	// Setup a N to be read from os.Stdin by main() for conformation to overwrite
	// the file.
	origStdin := os.Stdin
	defer func() {
		os.Stdin = origStdin
	}()

	r, w, err := os.Pipe()
	chk.NoErr(err)
	os.Stdin = r

	fmt.Fprintln(w, "Y")

	// Run command expecting the overwrite to be cancelled.
	chk.NoPanic(main)

	got, wnt, err := getTestFiles(dir, "README.md")
	chk.NoErr(err)
	wnt[0] = strings.ReplaceAll(wnt[0], "** DO NOT MODIFY ** ", "")
	chk.StrSlice(got, wnt)

	_, _, err = getTestFiles(dir, "README.md.gtm")
	chk.Err(
		err,
		"open "+dir+"/README.md.gtm: no such file or directory",
	)

	chk.Stdout("Confirm overwrite of " + dir + "/README.md (Y to overwrite)?\\s")
}

func Test_SampleGoProjectReplaceTargetOverwriteDirFromClean(t *testing.T) {
	chk := szTest.CaptureStdout(t)
	defer chk.Release()

	clearPackageCache()

	dir := chk.CreateTmpDir()
	chk.NoErr(
		setup(dir, "README.md.gtm", "sample_test.go", "sample.go"),
	)

	chk.NoErr(
		os.Rename(
			filepath.Join(dir, "README.md.gtm"),
			filepath.Join(dir, "README.md"),
		),
	)

	chk.SetupArgsAndFlags([]string{
		"programName",
		"-r",
		dir,
	})

	// Setup a N to be read from os.Stdin by main() for conformation to overwrite
	// the file.
	origStdin := os.Stdin
	defer func() {
		os.Stdin = origStdin
	}()

	r, w, err := os.Pipe()
	chk.NoErr(err)
	os.Stdin = r

	fmt.Fprintln(w, "Y")

	// Run command expecting the overwrite to be cancelled.
	chk.NoPanic(main)

	got, wnt, err := getTestFiles(dir, "README.md")
	chk.NoErr(err)
	wnt[0] = strings.ReplaceAll(wnt[0], "** DO NOT MODIFY ** ", "")
	chk.StrSlice(got, wnt)

	_, _, err = getTestFiles(dir, "README.md.gtm")
	chk.Err(
		err,
		"open "+dir+"/README.md.gtm: no such file or directory",
	)

	chk.Stdout("Confirm overwrite of " + dir + "/README.md (Y to overwrite)?\\s")
}

func Test_SampleGoProjectReplaceTargetOverwriteDirVerbose(t *testing.T) {
	chk := szTest.CaptureLogAndStdout(t)
	defer chk.Release()

	clearPackageCache()

	dir := chk.CreateTmpDir()
	chk.NoErr(
		setup(dir, "README.md", "sample_test.go", "sample.go"),
	)

	chk.SetupArgsAndFlags([]string{
		"programName",
		"-v",
		"-r",
		dir,
	})

	// Setup a N to be read from os.Stdin by main() for conformation to overwrite
	// the file.
	origStdin := os.Stdin
	defer func() {
		os.Stdin = origStdin
	}()

	r, w, err := os.Pipe()
	chk.NoErr(err)
	os.Stdin = r

	fmt.Fprintln(w, "Y")

	// Run command expecting the overwrite to be cancelled.
	chk.NoPanic(main)

	got, wnt, err := getTestFiles(dir, "README.md")
	chk.NoErr(err)
	wnt[0] = strings.ReplaceAll(wnt[0], "** DO NOT MODIFY ** ", "")
	chk.StrSlice(got, wnt)

	pName := filepath.Join(dir, "README.md")
	chk.Stdout(
		"filesToProcess:  "+pName,
		"Confirm overwrite of "+pName+" (Y to overwrite)?\\s",
	)

	chk.Log(
		"in place replacing of "+pName,
		"Loading Package info for: .",
		"getInfo(\"package\")",
		"getInfo(\"TimesTwo\")",
		"getInfo(\"TimesThree\")",
		"getInfo(\"InterfaceType\")",
		"getInfo(\"StructureType\")",
		"getInfo(\"TimesTwo\")",
		"getInfo(\"TimesTwo\")",
		"getInfo(\"TimesThree\")",
		"getInfo(\"TimesTwo\")",
		"getInfo(\"TimesTwo\")",
		"getInfo(\"TimesThree\")",
		"getInfo(\"TimesTwo\")",
		"getInfo(\"TimesTwo\")",
		"getInfo(\"TimesThree\")",
	)
}

func setup(dir string, files ...string) error {
	var err error
	var b []byte

	const ext = ".sample"

	files = append(files, "go.mod"+ext, "go.sum"+ext)
	for i, mi := 0, len(files); i < mi && err == nil; i++ {
		b, err = os.ReadFile(filepath.Join("sampleGoProject", files[i]))
		if err == nil {
			err = os.WriteFile(
				filepath.Join(dir, strings.TrimSuffix(files[i], ext)),
				b,
				os.FileMode(defaultPerm),
			)
		}
	}

	return err
}

func getTestFiles(dir, fName string) ([]string, []string, error) {
	gotBytes, err := os.ReadFile(filepath.Join(dir, fName))
	if err != nil {
		return nil, nil, err
	}
	wntBytes, err := os.ReadFile(filepath.Join("sampleGoProject", fName))
	if err != nil {
		return nil, nil, err
	}
	return strings.Split(string(gotBytes), "\n"),
		strings.Split(string(wntBytes), "\n"),
		nil
}

func Test_SampleGoProjectCleanNoTargetAlternateOut(t *testing.T) {
	chk := szTest.CaptureLogAndStdout(t)
	defer chk.Release()

	dir := chk.CreateTmpDir()
	altDir := chk.CreateTmpSubDir("altDir")

	chk.NoErr(setup(dir, "README.md", "sample_test.go", "sample.go"))

	chk.SetupArgsAndFlags([]string{
		"programName",
		"-v",
		"-l",
		"-c",
		"-o", altDir,
		filepath.Join(dir, "README.md"),
	})

	// Nor Run the main function with no -f arg requiring confirmation
	chk.NoPanic(main)

	got, wnt, err := getTestFiles(altDir, "README.md.gtm")
	chk.NoErr(err)
	chk.StrSlice(got, wnt)

	pName := filepath.Join(dir, "README.md")
	chk.Stdout(
		"filesToProcess:  "+pName,
		license,
	)
	oName := filepath.Join(altDir, "README.md.gtm")
	chk.Log(
		"Cleaning README.md to: " + oName + " in dir: " + dir,
	)
}
