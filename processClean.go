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

func cleanMD(completePath string) error {
	var err error
	var absCompletePath string
	var absDir string
	var srcFile, dstFile string
	var fileBytes []byte
	var res string

	absCompletePath, err = filepath.Abs(completePath)
	if err == nil {
		absDir = filepath.Dir(absCompletePath)
		srcFile = filepath.Base(absCompletePath)
		dstFile = filepath.Join(outputDir, srcFile+".gtm")
	}

	if verbose {
		log.Printf("Cleaning %s to: %s in dir: %s",
			srcFile, dstFile, absDir)
	}

	if err == nil {
		err = os.Chdir(absDir)
	}

	if err == nil {
		fileBytes, err = os.ReadFile(srcFile) //nolint:gosec // Ok.
	}

	if err == nil {
		fileData := string(bytes.TrimRight(fileBytes, "\n"))
		res, err = cleanMarkDownDocument(fileData)
	}

	okToOverwrite := forceOverwrite
	if err == nil && !forceOverwrite {
		okToOverwrite, err = confirmOverwrite(dstFile)
	}

	if err == nil && okToOverwrite {
		var f *os.File

		//nolint:gosec // Ok.
		f, err = os.OpenFile(dstFile,
			os.O_TRUNC|os.O_WRONLY|os.O_CREATE,
			os.FileMode(defaultPerm),
		)
		if err == nil {
			_, err = f.WriteString(strings.ReplaceAll(res, "\t", "    ") + "\n")
			if err == nil {
				err = f.Close()
			}
		}
	}

	return err
}
