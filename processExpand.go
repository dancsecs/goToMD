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
)

func expandMD(cleanPath string) error {
	var err error
	var absCleanPath string
	var absDir string
	var srcFile, dstFile string
	var fileBytes []byte
	var res string

	absCleanPath, err = filepath.Abs(cleanPath)
	if err == nil {
		absDir = filepath.Dir(absCleanPath)
		srcFile = filepath.Base(absCleanPath)
		dstFile = filepath.Join(outputDir, srcFile[:len(srcFile)-len(".gtm")])
	}

	if verbose {
		log.Printf("Expanding %s to: %s", absCleanPath, dstFile)
	}

	if err == nil {
		err = os.Chdir(absDir)
	}

	if err == nil {
		fileBytes, err = os.ReadFile(srcFile) //nolint:gosec // Ok.
	}

	if err == nil {
		fileData := string(bytes.TrimRight(fileBytes, "\n"))
		res, err = updateMarkDownDocument(fileData)
	}

	if err == nil {
		err = writeFile(dstFile, res)
	}

	return err
}
