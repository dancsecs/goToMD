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

func getProcessedFiles(dir, fName string) (string, string, []string, []string, error) {
	var targetPath string
	var sourceDir string
	var err error
	var gotBytes []byte
	var wntBytes []byte

	targetPath, err = filepath.Abs(filepath.Join(dir, fName))
	if err == nil {
		gotBytes, err = os.ReadFile(targetPath)
	}
	if err == nil {
		sourceDir, err = filepath.Abs(".")
	}
	if err == nil {
		wntBytes, err = os.ReadFile(fName)
	}
	if err != nil {
		return "", "", nil, nil, err
	}
	return targetPath, sourceDir,
		strings.Split(string(gotBytes), "\n"),
		strings.Split(string(wntBytes), "\n"),
		nil
}

func setupDirs(chk *szTest.Chk) {
	chk.T().Helper()
	origOutputDir := outputDir
	origCWD, err := os.Getwd()

	if chk.NoErr(err) {
		outputDir = chk.CreateTmpDir()
		chk.PushPostReleaseFunc(func() error {
			outputDir = origOutputDir
			return os.Chdir(origCWD)
		})
	}
}

func setupTest(
	chk *szTest.Chk, tCleanOnly, tReplace, tForceOverwrite, tVerbose bool,
) {
	chk.T().Helper()
	origOutputDir := outputDir
	origCWD, err := os.Getwd()
	origCleanOnly := cleanOnly
	origReplace := replace
	origForceOverwrite := forceOverwrite
	origVerbose := verbose

	cleanOnly = tCleanOnly
	replace = tReplace
	forceOverwrite = tForceOverwrite
	verbose = tVerbose

	if chk.NoErr(err) {
		outputDir = chk.CreateTmpDir()
		chk.PushPostReleaseFunc(func() error {
			outputDir = origOutputDir
			cleanOnly = origCleanOnly
			replace = origReplace
			forceOverwrite = origForceOverwrite
			verbose = origVerbose
			return os.Chdir(origCWD)
		})
	}
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

func Test_Process_SampleGoProjectCleanNoForceNoVerboseNoTarget(t *testing.T) {
	chk := szTest.CaptureLog(t)
	defer chk.Release()

	setupTest(chk, true, false, false, false)

	chk.NoErr(cleanMD(sampleGoProjectPath + "README.md"))

	// tPath, sDir, got, wnt, err := getProcessedFiles(outputDir, "README.md.gtm")
	_, _, got, wnt, err := getProcessedFiles(outputDir, "README.md.gtm")
	chk.NoErr(err)
	chk.StrSlice(got, wnt)

	chk.Log()
}

func Test_Process_SampleGoProjectCleanForceNoVerboseNoTarget(t *testing.T) {
	chk := szTest.CaptureLog(t)
	defer chk.Release()

	setupTest(chk, true, false, true, false)

	chk.NoErr(cleanMD(sampleGoProjectPath + "README.md"))

	// tPath, sDir, got, wnt, err := getProcessedFiles(outputDir, "README.md.gtm")
	_, _, got, wnt, err := getProcessedFiles(outputDir, "README.md.gtm")
	chk.NoErr(err)
	chk.StrSlice(got, wnt)

	chk.Log()
}

func Test_Process_SampleGoProjectCleanNoForceVerboseNoTarget(t *testing.T) {
	chk := szTest.CaptureLog(t)
	defer chk.Release()

	setupTest(chk, true, false, false, true)

	chk.NoErr(cleanMD(sampleGoProjectPath + "README.md"))

	tPath, sDir, got, wnt, err := getProcessedFiles(outputDir, "README.md.gtm")
	chk.NoErr(err)
	chk.StrSlice(got, wnt)

	chk.Log(
		"Cleaning README.md to: " + tPath + " in dir: " + sDir,
	)
}

func Test_Process_SampleGoProjectCleanForceVerboseNoTarget(t *testing.T) {
	chk := szTest.CaptureLog(t)
	defer chk.Release()

	setupTest(chk, true, false, true, true)

	chk.NoErr(cleanMD(sampleGoProjectPath + "README.md"))

	tPath, sDir, got, wnt, err := getProcessedFiles(outputDir, "README.md.gtm")
	chk.NoErr(err)
	chk.StrSlice(got, wnt)

	chk.Log(
		"Cleaning README.md to: " + tPath + " in dir: " + sDir,
	)
}

func Test_Process_SampleGoProjectCleanTargetCancel(t *testing.T) {
	chk := szTest.CaptureLogAndStdout(t)
	defer chk.Release()

	setupDirs(chk)

	targetFile := filepath.Join(outputDir, "README.md.gtm")
	// create target to test overwrite processing
	chk.NoErr(os.WriteFile(targetFile, nil, fs.FileMode(defaultPerm)))

	chk.SetStdinData("N\n")

	// Run command expecting the overwrite to be cancelled.
	chk.NoErr(cleanMD(sampleGoProjectPath + "README.md"))

	chk.Stdout("Confirm overwrite of " + targetFile + " (Y to overwrite)? " +
		"overwrite cancelled",
	)

	chk.Log(
		"Cleaning README.md to: " +
			targetFile +
			" in dir: " +
			fullSampleGoProjectPath,
	)
}

func Test_Process_SampleGoProjectCleanTargetOverwrite(t *testing.T) {
	chk := szTest.CaptureLogAndStdout(t)
	defer chk.Release()

	setupDirs(chk)

	// create target to test overwrite processing
	chk.NoErr(
		os.WriteFile(
			filepath.Join(outputDir, "README.md.gtm"),
			nil,
			fs.FileMode(defaultPerm),
		),
	)

	chk.SetStdinData("Y\n")

	// Run command expecting the overwrite to be cancelled.
	chk.NoErr(cleanMD(sampleGoProjectPath + "README.md"))

	tPath, sDir, got, wnt, err := getProcessedFiles(outputDir, "README.md.gtm")
	chk.NoErr(err)
	chk.StrSlice(got, wnt)

	chk.Stdout("Confirm overwrite of " + tPath + " (Y to overwrite)?\\s")

	chk.Log(
		"Cleaning README.md to: " +
			tPath +
			" in dir: " +
			sDir,
	)
}
