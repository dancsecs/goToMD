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
	"os"
)

func expandGoFile(cmd string) (string, error) {
	var fData []byte
	dir, fName, err := parseCmd(cmd)
	if err == nil {
		fPath := dir + string(os.PathSeparator) + fName
		fData, err = os.ReadFile(fPath) //nolint:gosec // Ok.
		if err == nil {
			return "" +
					szTestBgnPrefix + szFilePrefix + cmd + " -->\n" +
					"```bash\ncat " + fPath + "\n```\n\n" +
					"```go\n" +
					string(fData) +
					"```\n" +
					szTestEndPrefix + szFilePrefix + cmd + " -->\n",
				nil
		}
	}
	return "", err
}
