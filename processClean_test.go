package main

import (
	"io/fs"
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

type cleanGlobals struct {
	forceOverwrite bool
	verbose        bool
}

func getCleanedFiles() (string, []string, []string, error) {
	const fName = "README.md.gtm"
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

func setupCleanDirs(chk *szTest.Chk, makeTarget bool) string {
	const fName = "README.md.gtm"
	var tFile string

	chk.T().Helper()
	origOutputDir := outputDir
	origCWD, err := os.Getwd()

	if chk.NoErr(err) {
		outputDir = chk.CreateTmpDir()
		chk.PushPostReleaseFunc(func() error {
			outputDir = origOutputDir
			return os.Chdir(origCWD)
		})
		if makeTarget {
			tFile = filepath.Join(outputDir, fName)
			err = os.WriteFile(tFile, nil, fs.FileMode(defaultPerm))
			chk.NoErr(err, "could not create empty target to be replaced")
		}
	}
	return tFile
}

func setupCleanGlobals(
	chk *szTest.Chk, override cleanGlobals,
) {
	chk.T().Helper()
	setupTest(chk, true, false, override.forceOverwrite, override.verbose)
}

// +-------------------------------------------------------+
// | Flag possibilities for type of test.                  |
// +------------+-----------+------------------+-----------+
// | cleanOnly  |  replace  |  forceOverwrite  |  verbose  |
// +------------+-----------+------------------+-----------+
// |  false     |   false   |     false        |   false   |
// |  false     |   true    |     false        |   false   |
// |  true      |   false   |     false        |   false   |
// +------------+-----------+------------------+-----------+
// |  false     |   false   |     false        |   true    |
// |  false     |   true    |     false        |   true    |
// |  true      |   false   |     false        |   true    |
// +------------+-----------+------------------+-----------+
// |  false     |   false   |     true         |   false   |
// |  false     |   true    |     true         |   false   |
// |  true      |   false   |     true         |   false   |
// +------------+-----------+------------------+-----------+
// |  false     |   false   |     true         |   true    |
// |  false     |   true    |     true         |   true    |
// |  true      |   false   |     true         |   true    |
// +------------+-----------+------------------+-----------+.

func Test_ProcessClean_NoTargetNoForceNoVerbose(t *testing.T) {
	chk := szTest.CaptureLog(t)
	defer chk.Release()

	setupCleanGlobals(chk, cleanGlobals{forceOverwrite: false, verbose: false})
	setupCleanDirs(chk, false)

	chk.NoErr(
		cleanMD(sampleGoProjectPath + "README.md"),
	)

	_, got, wnt, err := getCleanedFiles()
	chk.NoErr(err)
	chk.StrSlice(got, wnt)

	chk.Log()
}

func Test_ProcessClean_NoTargetForceNoVerbose(t *testing.T) {
	chk := szTest.CaptureLog(t)
	defer chk.Release()

	setupCleanGlobals(chk, cleanGlobals{forceOverwrite: true, verbose: false})
	setupCleanDirs(chk, false)

	chk.NoErr(
		cleanMD(sampleGoProjectPath + "README.md"),
	)

	_, got, wnt, err := getCleanedFiles()
	chk.NoErr(err)
	chk.StrSlice(got, wnt)

	chk.Log()
}

func Test_ProcessClean_NoTargetNoForceVerbose(t *testing.T) {
	chk := szTest.CaptureLog(t)
	defer chk.Release()

	setupCleanGlobals(chk, cleanGlobals{forceOverwrite: false, verbose: true})
	setupCleanDirs(chk, false)

	chk.NoErr(cleanMD(sampleGoProjectPath + "README.md"))

	tPath, got, wnt, err := getCleanedFiles()
	chk.NoErr(err)
	chk.StrSlice(got, wnt)

	chk.Log(
		"Cleaning README.md to: " + tPath + " in dir: " + fullSampleGoProjectPath,
	)
}

func Test_ProcessClean_NoTargetForceVerbose(t *testing.T) {
	chk := szTest.CaptureLog(t)
	defer chk.Release()

	setupCleanGlobals(chk, cleanGlobals{forceOverwrite: true, verbose: true})
	setupCleanDirs(chk, false)

	chk.NoErr(cleanMD(sampleGoProjectPath + "README.md"))

	tPath, got, wnt, err := getCleanedFiles()
	chk.NoErr(err)
	chk.StrSlice(got, wnt)

	chk.Log(
		"Cleaning README.md to: " + tPath + " in dir: " + fullSampleGoProjectPath,
	)
}

func Test_ProcessClean_CancelOverwriteNoForceNoVerbose(t *testing.T) {
	chk := szTest.CaptureLogAndStdout(t)
	defer chk.Release()

	setupCleanGlobals(chk, cleanGlobals{forceOverwrite: false, verbose: false})
	tFile := setupCleanDirs(chk, true)

	chk.SetStdinData("N\n")

	// Run command expecting the overwrite to be cancelled.
	chk.NoErr(cleanMD(sampleGoProjectPath + "README.md"))

	chk.Stdout("Confirm overwrite of " + tFile + " (Y to overwrite)? " +
		"overwrite cancelled",
	)

	chk.Log()
}

func Test_ProcessClean_CancelOverwriteForceNoVerbose(t *testing.T) {
	chk := szTest.CaptureLogAndStdout(t)
	defer chk.Release()

	setupCleanGlobals(chk, cleanGlobals{forceOverwrite: true, verbose: false})
	_ = setupCleanDirs(chk, true)

	// Run command expecting the overwrite to be cancelled.
	chk.NoErr(cleanMD(sampleGoProjectPath + "README.md"))

	chk.Stdout()

	chk.Log()
}

func Test_ProcessClean_CancelOverwriteNoForceVerbose(t *testing.T) {
	chk := szTest.CaptureLogAndStdout(t)
	defer chk.Release()

	setupCleanGlobals(chk, cleanGlobals{forceOverwrite: false, verbose: true})
	tFile := setupCleanDirs(chk, true)

	chk.SetStdinData("N\n")

	// Run command expecting the overwrite to be cancelled.
	chk.NoErr(cleanMD(sampleGoProjectPath + "README.md"))

	chk.Stdout("Confirm overwrite of " + tFile + " (Y to overwrite)? " +
		"overwrite cancelled",
	)

	chk.Log(
		"Cleaning README.md to: " +
			tFile +
			" in dir: " +
			fullSampleGoProjectPath,
	)
}

func Test_ProcessClean_CancelOverwriteForceVerbose(t *testing.T) {
	chk := szTest.CaptureLogAndStdout(t)
	defer chk.Release()

	setupCleanGlobals(chk, cleanGlobals{forceOverwrite: true, verbose: true})
	tFile := setupCleanDirs(chk, true)

	// Run command expecting the overwrite to be cancelled.
	chk.NoErr(cleanMD(sampleGoProjectPath + "README.md"))

	chk.Stdout()

	chk.Log(
		"Cleaning README.md to: " +
			tFile +
			" in dir: " +
			fullSampleGoProjectPath,
	)
}

func Test_ProcessCLean_OverwriteNoForceNoVerbose(t *testing.T) {
	chk := szTest.CaptureLogAndStdout(t)
	defer chk.Release()

	setupCleanGlobals(chk, cleanGlobals{forceOverwrite: false, verbose: false})
	tFile := setupCleanDirs(chk, true)

	chk.SetStdinData("Y\n")

	// Run command expecting the overwrite to be cancelled.
	chk.NoErr(cleanMD(sampleGoProjectPath + "README.md"))

	_, got, wnt, err := getCleanedFiles()
	chk.NoErr(err)
	chk.StrSlice(got, wnt)

	chk.Stdout("Confirm overwrite of " + tFile + " (Y to overwrite)?\\s")

	chk.Log()
}

func Test_ProcessCLean_OverwriteForceNoVerbose(t *testing.T) {
	chk := szTest.CaptureLogAndStdout(t)
	defer chk.Release()

	setupCleanGlobals(chk, cleanGlobals{forceOverwrite: true, verbose: false})
	_ = setupCleanDirs(chk, true)

	// Run command expecting the overwrite to be cancelled.
	chk.NoErr(cleanMD(sampleGoProjectPath + "README.md"))

	_, got, wnt, err := getCleanedFiles()
	chk.NoErr(err)
	chk.StrSlice(got, wnt)

	chk.Stdout()

	chk.Log()
}

func Test_ProcessCLean_OverwriteNoForceVerbose(t *testing.T) {
	chk := szTest.CaptureLogAndStdout(t)
	defer chk.Release()

	setupCleanGlobals(chk, cleanGlobals{forceOverwrite: false, verbose: true})
	tFile := setupCleanDirs(chk, true)

	chk.SetStdinData("Y\n")

	// Run command expecting the overwrite to be cancelled.
	chk.NoErr(cleanMD(sampleGoProjectPath + "README.md"))

	_, got, wnt, err := getCleanedFiles()
	chk.NoErr(err)
	chk.StrSlice(got, wnt)

	chk.Stdout("Confirm overwrite of " + tFile + " (Y to overwrite)?\\s")

	chk.Log(
		"Cleaning README.md to: " +
			tFile +
			" in dir: " +
			fullSampleGoProjectPath,
	)
}

func Test_ProcessCLean_OverwriteForceVerbose(t *testing.T) {
	chk := szTest.CaptureLogAndStdout(t)
	defer chk.Release()

	setupCleanGlobals(chk, cleanGlobals{forceOverwrite: true, verbose: true})
	tFile := setupCleanDirs(chk, true)

	// Run command expecting the overwrite to be cancelled.
	chk.NoErr(cleanMD(sampleGoProjectPath + "README.md"))

	_, got, wnt, err := getCleanedFiles()
	chk.NoErr(err)
	chk.StrSlice(got, wnt)

	chk.Stdout()

	chk.Log(
		"Cleaning README.md to: " +
			tFile +
			" in dir: " +
			fullSampleGoProjectPath,
	)
}
