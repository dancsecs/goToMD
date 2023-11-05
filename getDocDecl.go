// Program runs a specific go test transforming the output
// to github compatible markdown.  This is used within this
// project to help automate keeping the README.md up to date
// when an example changes.

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
	"strings"
)

func getDocDecl(cmd string) (string, error) {
	var d *docInfo
	res := ""
	dir, action, err := parseCmd(cmd)
	if err == nil {
		items := strings.Split(action, " ")
		mi := len(items)
		for i := 0; i < mi && err == nil; i++ {
			d, err = getInfo(dir, items[i])
			if err == nil {
				if res != "" {
					res += "\n"
				}
				res += strings.Join(d.header, "\n")
			}
		}
	}
	if err == nil {
		return res, nil
	}
	return "", err
}

func expandGoDcl(cmd string) (string, error) {
	doc, err := getDocDecl(cmd)
	if err != nil {
		return "", err
	}
	return "" +
			szTestBgnPrefix + szDclPrefix + cmd + " -->\n" +
			"```go\n" + doc + "\n```\n" +
			szTestEndPrefix + szDclPrefix + cmd + " -->\n",
		nil
}

func getDocDeclSingle(cmd string) (string, error) {
	var d *docInfo
	res := ""
	dir, action, err := parseCmd(cmd)
	if err == nil {
		items := strings.Split(action, " ")
		mi := len(items)
		for i := 0; i < mi && err == nil; i++ {
			d, err = getInfo(dir, items[i])
			if err == nil {
				if res != "" {
					res += "\n"
				}
				res += d.oneLine()
			}
		}
	}
	if err == nil {
		return res, nil
	}
	return "", err
}

func expandGoDclSingle(cmd string) (string, error) {
	doc, err := getDocDeclSingle(cmd)
	if err != nil {
		return "", err
	}
	return "" +
			szTestBgnPrefix + szDclsPrefix + cmd + " -->\n" +
			"```go\n" + doc + "\n```\n" +
			szTestEndPrefix + szDclsPrefix + cmd + " -->\n",
		nil
}
