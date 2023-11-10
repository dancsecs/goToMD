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
	"errors"
	"strings"
)

const szTestPrefix = "<!--- goToMD::"
const szAutoPrefix = szTestPrefix + "Auto::"
const szTestBgnPrefix = szTestPrefix + "Bgn::"
const szTestEndPrefix = szTestPrefix + "End::"
const szDocPrefix = "doc::"
const szTestDocPrefix = szTestPrefix + szDocPrefix
const szDclnPrefix = "dcln::"
const szTestDclnPrefix = szTestPrefix + szDclnPrefix
const szDclsPrefix = "dcls::"
const szTestDclsPrefix = szTestPrefix + szDclsPrefix
const szDclPrefix = "dcl::"
const szTestDclPrefix = szTestPrefix + szDclPrefix
const szTstPrefix = "tst::"
const szTestTstPrefix = szTestPrefix + szTstPrefix
const szFilePrefix = "file::"
const szTestFilePrefix = szTestPrefix + szFilePrefix

func cleanMarkDownDocument(fData string) (string, error) {
	var err error
	var skipBlank = false
	updatedFile := ""
	skipTo := ""
	lines := strings.Split(fData+"\n", "\n")
	for i, mi := 0, len(lines)-1; i < mi && err == nil; i++ {
		l := strings.TrimRight(lines[i], " ")
		switch {
		case skipBlank:
			if l != "" {
				err = errors.New("missing blank line in auto generated output")
			}
			skipBlank = false
		case skipTo != "":
			switch {
			case strings.HasPrefix(l, skipTo):
				skipTo = ""
			case strings.HasPrefix(l, szTestEndPrefix):
				err = errors.New("out of sequence: End before begin: " + l)
			}
		case strings.HasPrefix(l, szTestBgnPrefix):
			skipTo = szTestEndPrefix + l[len(szTestBgnPrefix):]
			// Add unexpanded line.
			updatedFile += szTestPrefix + l[len(szTestEndPrefix):] + "\n"
		case strings.HasPrefix(l, szAutoPrefix):
			// Do not add auto generated line or next blank line to output.
			skipBlank = true
		default:
			updatedFile += l + "\n"
		}
	}
	if err != nil {
		return "", err
	}
	return strings.TrimRight(updatedFile, "\n"), nil
}

//nolint:funlen // Ok.
func updateMarkDownDocument(fData string) (string, error) {
	var res string
	var cmd string
	var err error

	updatedFile := szAutoPrefix + " See github.com/dancsecs/goToMD "
	if !replace {
		updatedFile += "** DO NOT MODIFY ** "
	}
	updatedFile += "-->\n\n"
	lines := strings.Split(fData+"\n", "\n")
	for i, mi := 0, len(lines)-1; i < mi && err == nil; i++ {
		l := strings.TrimRight(lines[i], " ")
		switch {
		case !strings.HasPrefix(l, szTestPrefix):
			updatedFile += l + "\n"
		case strings.HasPrefix(l, szTestDocPrefix):
			cmd = l[len(szTestDocPrefix) : len(l)-len(" -->")]
			var dir, action string
			dir, action, err = parseCmd(cmd)
			if err == nil {
				var di *docInfo
				di, err = getInfo(dir, action)
				if err == nil {
					updatedFile += di.expand(szDocPrefix, cmd)
				}
			}
		case strings.HasPrefix(l, szTestDclPrefix):
			cmd = l[len(szTestDclPrefix) : len(l)-len(" -->")]
			res, err = expandGoDcl(cmd)
			if err == nil {
				updatedFile += res
			}
		case strings.HasPrefix(l, szTestDclsPrefix):
			cmd = l[len(szTestDclsPrefix) : len(l)-len(" -->")]
			res, err = expandGoDclSingle(cmd)
			if err == nil {
				updatedFile += res
			}
		case strings.HasPrefix(l, szTestDclnPrefix):
			cmd = l[len(szTestDclsPrefix) : len(l)-len(" -->")]
			res, err = expandGoDclNatural(cmd)
			if err == nil {
				updatedFile += res
			}
		case strings.HasPrefix(l, szTestTstPrefix):
			cmd := l[len(szTestTstPrefix) : len(l)-len(" -->")]
			res, err = expandGoTst(cmd)
			if err == nil {
				updatedFile += res
			}
		case strings.HasPrefix(l, szTestFilePrefix):
			cmd := l[len(szTestFilePrefix) : len(l)-len(" -->")]
			res, err = expandGoFile(cmd)
			if err == nil {
				updatedFile += res
			}
		default:
			err = errors.New("unknown cmd: " + l)
		}
	}
	if err != nil {
		return "", err
	}
	return strings.TrimRight(updatedFile, "\n"), nil
}
