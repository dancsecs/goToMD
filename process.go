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
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func confirmOverwrite(fPath string) (bool, error) {
	var ok bool
	_, err := os.Stat(fPath)
	if errors.Is(err, os.ErrNotExist) {
		return true, nil
	}
	if err == nil {
		fmt.Print("Confirm overwrite of ", fPath, " (Y to overwrite)? ")
		var response string
		if _, err = fmt.Scanln(&response); err == nil {
			ok = response == "Y"
			if !ok {
				fmt.Println("overwrite cancelled")
			}
		}
	}
	return ok, err
}

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
		dstFile = srcFile[:len(srcFile)-len(".gtm")]
	}

	if verbose {
		log.Printf("Expanding %s to: %s in dir: %s", srcFile, dstFile, absDir)
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

	okToOverwrite := forceOverwrite
	if err == nil && !forceOverwrite {
		okToOverwrite, err = confirmOverwrite(dstFile)
	}

	if err == nil && okToOverwrite {
		var f *os.File

		//nolint:gosec // Ok.
		f, err = os.OpenFile(filepath.Join(outputDir, dstFile),
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

func replaceMDInPlace(completePath string) error {
	var err error
	var absCompletePath string
	var fName string
	var absDir string
	var fileBytes []byte
	var res string

	absCompletePath, err = filepath.Abs(completePath)
	if err == nil {
		absDir = filepath.Dir(absCompletePath)
		fName = filepath.Base(absCompletePath)
	}

	if verbose {
		log.Printf("in place replacing of %s in dir: %s", fName, absDir)
	}

	if err == nil {
		err = os.Chdir(absDir)
	}

	if err == nil {
		fileBytes, err = os.ReadFile(fName) //nolint:gosec // Ok.
	}

	if err == nil {
		fileData := string(bytes.TrimRight(fileBytes, "\n"))
		res, err = cleanMarkDownDocument(fileData)
	}

	if err == nil {
		res = strings.TrimRight(res, "\n")
		res, err = updateMarkDownDocument(res)
	}

	okToOverwrite := forceOverwrite
	if err == nil && !forceOverwrite {
		okToOverwrite, err = confirmOverwrite(fName)
	}

	if err == nil && okToOverwrite {
		var f *os.File

		//nolint:gosec // Ok.
		f, err = os.OpenFile(filepath.Join(outputDir, fName),
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

func getFilesToProcess() ([]string, error) {
	var err error
	var files []os.DirEntry
	var stat os.FileInfo
	var filesToProcess []string
	var filter = ".md"

	if !cleanOnly && !replace {
		filter += ".gtm"
	}

	for i, mi := 0, flag.NArg(); i < mi && err == nil; i++ {
		stat, err = os.Stat(flag.Arg(i))
		if err == nil && stat.IsDir() {
			files, err = os.ReadDir(flag.Arg(i))
			for j, mj := 0, len(files); j < mj && err == nil; j++ {
				fName := files[j].Name()
				if strings.HasSuffix(fName, filter) {
					filesToProcess = append(filesToProcess,
						flag.Arg(i)+string(os.PathSeparator)+fName,
					)
					if verbose {
						fmt.Println("filesToProcess: ",
							flag.Arg(i)+string(os.PathSeparator)+fName,
						)
					}
				}
			}
		}

		if err == nil && !stat.IsDir() {
			if !strings.HasSuffix(stat.Name(), filter) {
				err = errors.New("file must have extension: " + filter)
			} else {
				filesToProcess = append(filesToProcess, flag.Arg(i))
				if verbose {
					fmt.Println("filesToProcess: ", flag.Arg(i))
				}
			}
		}
	}
	return filesToProcess, err
}
