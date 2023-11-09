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
	"bytes"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func replaceMDInPlace(fileToRead string) error {
	var err error
	var fName string
	var absDir string
	var fileBytes []byte
	var res string

	var rFile, wFile string

	rFile, err = filepath.Abs(fileToRead)
	if err == nil {
		absDir, fName = filepath.Split(rFile)

		if outputDir != "." { // Actual file to replace is in alternate directory.
			wFile, err = filepath.Abs(filepath.Join(outputDir, fName))
			if err == nil {
				if _, statErr := os.Stat(wFile); statErr == nil {
					// Only read alt file if it exists otherwise read named file.
					rFile = wFile
				}
			}
		} else {
			wFile = rFile
		}
	}

	if verbose {
		if rFile == wFile {
			log.Printf("in place replacing of %s", rFile)
		} else {
			log.Printf(
				"in place replacing (alt dir) of %s to %s", rFile, wFile,
			)
		}
	}

	if err == nil {
		err = os.Chdir(absDir)
	}

	if err == nil {
		fileBytes, err = os.ReadFile(rFile) //nolint:gosec // Ok.
	}

	if err == nil {
		fileData := string(bytes.TrimRight(fileBytes, "\n"))
		res, err = cleanMarkDownDocument(fileData)
	}

	if err == nil {
		res = strings.TrimRight(res, "\n")
		res, err = updateMarkDownDocument(res)
	}

	if err == nil {
		err = writeFile(wFile, res)
	}

	return err
}
