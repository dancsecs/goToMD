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

type expandGlobals struct {
	forceOverwrite bool
	verbose        bool
}

func setupExpandGlobals(
	chk *szTest.Chk, override expandGlobals,
) {
	chk.T().Helper()
	setupTest(chk, true, false, override.forceOverwrite, override.verbose)
}

func setupExpandDirs(makeTarget bool) (string, func(), error) {
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

func getExpandFiles() (string, []string, []string, error) {
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

func Test_ProcessExpand_NoTargetNoForceNoVerbose(t *testing.T) {
	chk := szTest.CaptureLogAndStdout(t)
	defer chk.Release()

	setupExpandGlobals(
		chk, expandGlobals{forceOverwrite: false, verbose: false},
	)
	_, cFunc, err := setupExpandDirs(false)
	chk.NoErr(err)
	defer cFunc()

	err = expandMD(sampleGoProjectPath + "README.md.gtm")
	chk.NoErr(err)

	_, got, wnt, err := getExpandFiles()
	chk.NoErr(err)
	chk.StrSlice(got, wnt)

	chk.Log()

	chk.Stdout()
}

func Test_ProcessExpand_NoTargetForceNoVerbose(t *testing.T) {
	chk := szTest.CaptureLogAndStdout(t)
	defer chk.Release()

	setupExpandGlobals(
		chk, expandGlobals{forceOverwrite: true, verbose: false},
	)
	_, cFunc, err := setupExpandDirs(false)
	chk.NoErr(err)
	defer cFunc()

	err = expandMD(sampleGoProjectPath + "README.md.gtm")
	chk.NoErr(err)

	_, got, wnt, err := getExpandFiles()
	chk.NoErr(err)
	chk.StrSlice(got, wnt)

	chk.Log()

	chk.Stdout()
}

func Test_ProcessExpand_NoTargetNoForceVerbose(t *testing.T) {
	chk := szTest.CaptureLogAndStdout(t)
	defer chk.Release()

	setupExpandGlobals(
		chk, expandGlobals{forceOverwrite: false, verbose: true},
	)
	_, cFunc, err := setupExpandDirs(false)
	chk.NoErr(err)
	defer cFunc()

	err = expandMD(sampleGoProjectPath + "README.md.gtm")
	chk.NoErr(err)

	tFile, got, wnt, err := getExpandFiles()
	chk.NoErr(err)
	chk.StrSlice(got, wnt)

	chk.Log(
		"Expanding "+fullSampleGoProjectPath+string(os.PathSeparator)+
			"README.md.gtm to: "+tFile,
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

func Test_ProcessExpand_NoTargetForceVerbose(t *testing.T) {
	chk := szTest.CaptureLogAndStdout(t)
	defer chk.Release()

	setupExpandGlobals(
		chk, expandGlobals{forceOverwrite: true, verbose: true},
	)
	_, cFunc, err := setupExpandDirs(false)
	chk.NoErr(err)
	defer cFunc()

	err = expandMD(sampleGoProjectPath + "README.md.gtm")
	chk.NoErr(err)

	tFile, got, wnt, err := getExpandFiles()
	chk.NoErr(err)
	chk.StrSlice(got, wnt)

	chk.Log(
		"Expanding "+fullSampleGoProjectPath+string(os.PathSeparator)+
			"README.md.gtm to: "+tFile,
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

func Test_ProcessExpand_CancelOverwriteNoForceNoVerbose(t *testing.T) {
	chk := szTest.CaptureLogAndStdout(t)
	defer chk.Release()

	setupExpandGlobals(
		chk, expandGlobals{forceOverwrite: false, verbose: false},
	)
	_, cFunc, err := setupExpandDirs(true)
	chk.NoErr(err)
	defer cFunc()

	chk.SetStdinData("N\n")

	err = expandMD(sampleGoProjectPath + "README.md.gtm")
	chk.NoErr(err)

	tFile, got, wnt, err := getExpandFiles()
	chk.NoErr(err)
	chk.StrSlice(got, wnt)

	chk.Log()

	chk.Stdout(
		"Confirm overwrite of " + tFile + " (Y to overwrite)? overwrite cancelled",
	)
}

func Test_ProcessExpand_CancelOverwriteTargetForceNoVerbose(t *testing.T) {
	chk := szTest.CaptureLogAndStdout(t)
	defer chk.Release()

	setupExpandGlobals(
		chk, expandGlobals{forceOverwrite: true, verbose: false},
	)
	_, cFunc, err := setupExpandDirs(true)
	chk.NoErr(err)
	defer cFunc()

	chk.SetStdinData("N\n")

	err = expandMD(sampleGoProjectPath + "README.md.gtm")
	chk.NoErr(err)

	_, got, wnt, err := getExpandFiles()
	chk.NoErr(err)
	chk.StrSlice(got, wnt)

	chk.Log()

	chk.Stdout()
}

func Test_ProcessExpand_CancelOverwriteNoForceVerbose(t *testing.T) {
	chk := szTest.CaptureLogAndStdout(t)
	defer chk.Release()

	setupExpandGlobals(
		chk, expandGlobals{forceOverwrite: false, verbose: true},
	)
	_, cFunc, err := setupExpandDirs(true)
	chk.NoErr(err)
	defer cFunc()

	chk.SetStdinData("N\n")

	err = expandMD(sampleGoProjectPath + "README.md.gtm")
	chk.NoErr(err)

	tFile, got, wnt, err := getExpandFiles()
	chk.NoErr(err)
	chk.StrSlice(got, wnt)

	chk.Log(
		"Expanding "+fullSampleGoProjectPath+string(os.PathSeparator)+
			"README.md.gtm to: "+tFile,
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

func Test_ProcessExpand_CancelOverwriteForceVerbose(t *testing.T) {
	chk := szTest.CaptureLogAndStdout(t)
	defer chk.Release()

	setupExpandGlobals(
		chk, expandGlobals{forceOverwrite: true, verbose: true},
	)
	_, cFunc, err := setupExpandDirs(true)
	chk.NoErr(err)
	defer cFunc()

	chk.SetStdinData("N\n")

	err = expandMD(sampleGoProjectPath + "README.md.gtm")
	chk.NoErr(err)

	tFile, got, wnt, err := getExpandFiles()
	chk.NoErr(err)
	chk.StrSlice(got, wnt)

	chk.Log(
		"Expanding "+fullSampleGoProjectPath+string(os.PathSeparator)+
			"README.md.gtm to: "+tFile,
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

func Test_ProcessExpand_OverwriteNoForceNoVerbose(t *testing.T) {
	chk := szTest.CaptureLogAndStdout(t)
	defer chk.Release()

	setupExpandGlobals(
		chk, expandGlobals{forceOverwrite: false, verbose: false},
	)
	_, cFunc, err := setupExpandDirs(true)
	chk.NoErr(err)
	defer cFunc()

	chk.SetStdinData("Y\n")

	err = expandMD(sampleGoProjectPath + "README.md.gtm")
	chk.NoErr(err)

	tFile, got, wnt, err := getExpandFiles()
	chk.NoErr(err)
	chk.StrSlice(got, wnt)

	chk.Log()

	chk.Stdout(
		"Confirm overwrite of " + tFile + " (Y to overwrite)?\\s",
	)
}

func Test_ProcessExpand_OverwriteTargetForceNoVerbose(t *testing.T) {
	chk := szTest.CaptureLogAndStdout(t)
	defer chk.Release()

	setupExpandGlobals(
		chk, expandGlobals{forceOverwrite: true, verbose: false},
	)
	_, cFunc, err := setupExpandDirs(true)
	chk.NoErr(err)
	defer cFunc()

	chk.SetStdinData("Y\n")

	err = expandMD(sampleGoProjectPath + "README.md.gtm")
	chk.NoErr(err)

	_, got, wnt, err := getExpandFiles()
	chk.NoErr(err)
	chk.StrSlice(got, wnt)

	chk.Log()

	chk.Stdout()
}

func Test_ProcessExpand_OverwriteNoForceVerbose(t *testing.T) {
	chk := szTest.CaptureLogAndStdout(t)
	defer chk.Release()

	setupExpandGlobals(
		chk, expandGlobals{forceOverwrite: false, verbose: true},
	)
	_, cFunc, err := setupExpandDirs(true)
	chk.NoErr(err)
	defer cFunc()

	chk.SetStdinData("Y\n")

	err = expandMD(sampleGoProjectPath + "README.md.gtm")
	chk.NoErr(err)

	tFile, got, wnt, err := getExpandFiles()
	chk.NoErr(err)
	chk.StrSlice(got, wnt)

	chk.Log(
		"Expanding "+fullSampleGoProjectPath+string(os.PathSeparator)+
			"README.md.gtm to: "+tFile,
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

func Test_ProcessExpand_OverwriteForceVerbose(t *testing.T) {
	chk := szTest.CaptureLogAndStdout(t)
	defer chk.Release()

	setupExpandGlobals(
		chk, expandGlobals{forceOverwrite: true, verbose: true},
	)
	_, cFunc, err := setupExpandDirs(true)
	chk.NoErr(err)
	defer cFunc()

	chk.SetStdinData("Y\n")

	err = expandMD(sampleGoProjectPath + "README.md.gtm")
	chk.NoErr(err)

	tFile, got, wnt, err := getExpandFiles()
	chk.NoErr(err)
	chk.StrSlice(got, wnt)

	chk.Log(
		"Expanding "+fullSampleGoProjectPath+string(os.PathSeparator)+
			"README.md.gtm to: "+tFile,
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