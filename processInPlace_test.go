package main

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
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

type inPlaceGlobals struct {
	forceOverwrite bool
	verbose        bool
}

func setupInPlaceGlobals(
	chk *szTest.Chk, override inPlaceGlobals,
) {
	chk.T().Helper()
	setupTest(chk, true, false, override.forceOverwrite, override.verbose)
}

func setupInPlaceDirs(makeTarget bool) (string, func(), error) {
	const fName = "README.md"
	var tFile string
	var fData []byte

	// Save state to restore when returned function is called.
	origOutputDir := outputDir
	origCWD, err := os.Getwd()

	if err != nil {
		return "", nil, err
	}

	restoreFunction := func() {
		outputDir = origOutputDir
		err := os.Chdir(origCWD)
		if err != nil {
			log.Printf("Could not retore original working directory.")
		}
	}

	if makeTarget {
		fData, err = os.ReadFile(filepath.Join(fullSampleGoProjectPath, fName))
		if err != nil {
			return "", nil, err
		}
		tFile = filepath.Join(outputDir, fName)
		err = os.WriteFile(tFile, fData, fs.FileMode(defaultPerm))
		if err != nil {
			return "", nil, err
		}
	}
	return tFile, restoreFunction, nil
}

func getInPlaceFiles() (string, []string, []string, error) {
	const fName = "README.md"
	var targetPath string
	var err error
	var gotBytes []byte
	var wntBytes []byte

	targetPath, err = filepath.Abs(filepath.Join(outputDir, fName))
	if err == nil {
		gotBytes, err = os.ReadFile(targetPath)
	}
	if err == nil {
		wntBytes, err = os.ReadFile(fName)
	}
	if err != nil {
		return "", nil, nil, err
	}
	return targetPath,
		strings.Split(string(gotBytes), "\n"),
		strings.Split(string(wntBytes), "\n"),
		nil
}

func Test_ProcessInPlace_NoTargetNoForceNoVerbose(t *testing.T) {
	chk := szTest.CaptureLogAndStdout(t)
	defer chk.Release()

	setupInPlaceGlobals(
		chk, inPlaceGlobals{forceOverwrite: false, verbose: false},
	)
	_, cFunc, err := setupInPlaceDirs(false)
	chk.NoErr(err)
	defer cFunc()

	err = replaceMDInPlace(sampleGoProjectPath + "README.md")
	chk.NoErr(err)

	_, got, wnt, err := getInPlaceFiles()
	chk.NoErr(err)
	chk.StrSlice(got, wnt)

	chk.Log()

	chk.Stdout()
}

func Test_ProcessInPlace_NoTargetForceNoVerbose(t *testing.T) {
	chk := szTest.CaptureLogAndStdout(t)
	defer chk.Release()

	setupInPlaceGlobals(
		chk, inPlaceGlobals{forceOverwrite: true, verbose: false},
	)
	_, cFunc, err := setupInPlaceDirs(false)
	chk.NoErr(err)
	defer cFunc()

	err = replaceMDInPlace(sampleGoProjectPath + "README.md")
	chk.NoErr(err)

	_, got, wnt, err := getInPlaceFiles()
	chk.NoErr(err)
	chk.StrSlice(got, wnt)

	chk.Log()

	chk.Stdout()
}

func Test_ProcessInPlace_NoTargetNoForceVerbose(t *testing.T) {
	chk := szTest.CaptureLogAndStdout(t)
	defer chk.Release()

	setupInPlaceGlobals(
		chk, inPlaceGlobals{forceOverwrite: false, verbose: true},
	)
	_, cFunc, err := setupInPlaceDirs(false)
	chk.NoErr(err)
	defer cFunc()

	err = replaceMDInPlace(sampleGoProjectPath + "README.md")
	chk.NoErr(err)

	tFile, got, wnt, err := getInPlaceFiles()
	chk.NoErr(err)
	chk.StrSlice(got, wnt)

	chk.Log(
		"in place replacing (alt dir) of "+
			fullSampleGoProjectPath+string(os.PathSeparator)+
			"README.md to "+
			tFile,
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

	chk.Stdout()
}

func Test_ProcessInPlace_NoTargetForceVerbose(t *testing.T) {
	chk := szTest.CaptureLogAndStdout(t)
	defer chk.Release()

	setupInPlaceGlobals(
		chk, inPlaceGlobals{forceOverwrite: true, verbose: true},
	)
	_, cFunc, err := setupInPlaceDirs(false)
	chk.NoErr(err)
	defer cFunc()

	err = replaceMDInPlace(sampleGoProjectPath + "README.md")
	chk.NoErr(err)

	tFile, got, wnt, err := getInPlaceFiles()
	chk.NoErr(err)
	chk.StrSlice(got, wnt)

	chk.Log(
		"in place replacing (alt dir) of "+
			fullSampleGoProjectPath+string(os.PathSeparator)+
			"README.md to "+
			tFile,
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

	chk.Stdout()
}

func Test_ProcessInPlace_CancelOverwriteNoForceNoVerbose(t *testing.T) {
	chk := szTest.CaptureLogAndStdout(t)
	defer chk.Release()

	setupInPlaceGlobals(
		chk, inPlaceGlobals{forceOverwrite: false, verbose: false},
	)
	_, cFunc, err := setupInPlaceDirs(true)
	chk.NoErr(err)
	defer cFunc()

	chk.SetStdinData("N\n")

	err = replaceMDInPlace(sampleGoProjectPath + "README.md")
	chk.NoErr(err)

	tFile, got, wnt, err := getInPlaceFiles()
	chk.NoErr(err)
	chk.StrSlice(got, wnt)

	chk.Log()

	chk.Stdout(
		"Confirm overwrite of " + tFile + " (Y to overwrite)? overwrite cancelled",
	)
}

func Test_ProcessInPlace_CancelOverwriteForceNoVerbose(t *testing.T) {
	chk := szTest.CaptureLogAndStdout(t)
	defer chk.Release()

	setupInPlaceGlobals(
		chk, inPlaceGlobals{forceOverwrite: true, verbose: false},
	)
	_, cFunc, err := setupInPlaceDirs(true)
	chk.NoErr(err)
	defer cFunc()

	chk.SetStdinData("N\n")

	err = replaceMDInPlace(sampleGoProjectPath + "README.md")
	chk.NoErr(err)

	_, got, wnt, err := getInPlaceFiles()
	chk.NoErr(err)
	chk.StrSlice(got, wnt)

	chk.Log()

	chk.Stdout()
}

func Test_ProcessInPlace_CancelOverwriteNoForceVerbose(t *testing.T) {
	chk := szTest.CaptureLogAndStdout(t)
	defer chk.Release()

	setupInPlaceGlobals(
		chk, inPlaceGlobals{forceOverwrite: false, verbose: true},
	)
	_, cFunc, err := setupInPlaceDirs(true)
	chk.NoErr(err)
	defer cFunc()

	chk.SetStdinData("N\n")

	err = replaceMDInPlace(sampleGoProjectPath + "README.md")
	chk.NoErr(err)

	tFile, got, wnt, err := getInPlaceFiles()
	chk.NoErr(err)
	chk.StrSlice(got, wnt)

	chk.Log(
		"in place replacing of "+
			tFile,
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

	chk.Stdout(
		"Confirm overwrite of " + tFile + " (Y to overwrite)? overwrite cancelled",
	)
}

func Test_ProcessInPlace_CancelOverwriteForceVerbose(t *testing.T) {
	chk := szTest.CaptureLogAndStdout(t)
	defer chk.Release()

	setupInPlaceGlobals(
		chk, inPlaceGlobals{forceOverwrite: true, verbose: true},
	)
	_, cFunc, err := setupInPlaceDirs(true)
	chk.NoErr(err)
	defer cFunc()

	chk.SetStdinData("N\n")

	err = replaceMDInPlace(sampleGoProjectPath + "README.md")
	chk.NoErr(err)

	tFile, got, wnt, err := getInPlaceFiles()
	chk.NoErr(err)
	chk.StrSlice(got, wnt)

	chk.Log(
		"in place replacing of "+
			tFile,
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

	chk.Stdout()
}

func Test_ProcessInPlace_OverwriteNoForceNoVerbose(t *testing.T) {
	chk := szTest.CaptureLogAndStdout(t)
	defer chk.Release()

	setupInPlaceGlobals(
		chk, inPlaceGlobals{forceOverwrite: false, verbose: false},
	)
	_, cFunc, err := setupInPlaceDirs(true)
	chk.NoErr(err)
	defer cFunc()

	chk.SetStdinData("Y\n")

	err = replaceMDInPlace(sampleGoProjectPath + "README.md")
	chk.NoErr(err)

	tFile, got, wnt, err := getInPlaceFiles()
	chk.NoErr(err)
	chk.StrSlice(got, wnt)

	chk.Log()

	chk.Stdout(
		"Confirm overwrite of " + tFile + " (Y to overwrite)?\\s",
	)
}

func Test_ProcessInPlace_OverwriteForceNoVerbose(t *testing.T) {
	chk := szTest.CaptureLogAndStdout(t)
	defer chk.Release()

	setupInPlaceGlobals(
		chk, inPlaceGlobals{forceOverwrite: true, verbose: false},
	)
	_, cFunc, err := setupInPlaceDirs(true)
	chk.NoErr(err)
	defer cFunc()

	chk.SetStdinData("Y\n")

	err = replaceMDInPlace(sampleGoProjectPath + "README.md")
	chk.NoErr(err)

	_, got, wnt, err := getInPlaceFiles()
	chk.NoErr(err)
	chk.StrSlice(got, wnt)

	chk.Log()

	chk.Stdout()
}

func Test_ProcessInPlace_OverwriteNoForceVerbose(t *testing.T) {
	chk := szTest.CaptureLogAndStdout(t)
	defer chk.Release()

	setupInPlaceGlobals(
		chk, inPlaceGlobals{forceOverwrite: false, verbose: true},
	)
	_, cFunc, err := setupInPlaceDirs(true)
	chk.NoErr(err)
	defer cFunc()

	chk.SetStdinData("Y\n")

	err = replaceMDInPlace(sampleGoProjectPath + "README.md")
	chk.NoErr(err)

	tFile, got, wnt, err := getInPlaceFiles()
	chk.NoErr(err)
	chk.StrSlice(got, wnt)

	chk.Log(
		"in place replacing of "+
			tFile,
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

	chk.Stdout(
		"Confirm overwrite of " + tFile + " (Y to overwrite)?\\s",
	)
}

func Test_ProcessInPlace_OverwriteForceVerbose(t *testing.T) {
	chk := szTest.CaptureLogAndStdout(t)
	defer chk.Release()

	setupInPlaceGlobals(
		chk, inPlaceGlobals{forceOverwrite: true, verbose: true},
	)
	_, cFunc, err := setupInPlaceDirs(true)
	chk.NoErr(err)
	defer cFunc()

	chk.SetStdinData("Y\n")

	err = replaceMDInPlace(sampleGoProjectPath + "README.md")
	chk.NoErr(err)

	tFile, got, wnt, err := getInPlaceFiles()
	chk.NoErr(err)
	chk.StrSlice(got, wnt)

	chk.Log(
		"in place replacing of "+
			tFile,
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

	chk.Stdout()
}
