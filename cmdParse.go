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
	"fmt"
	"os"
	"strings"
)

func parseCmd(cmd string) (string, string, error) {
	if !strings.HasPrefix(cmd, "./") {
		return "", "", fmt.Errorf(
			"relative directory must be specified in cmd: %q", cmd,
		)
	}

	lastSeparatorPos := strings.LastIndex(cmd, string(os.PathSeparator))
	dir := strings.TrimSpace(cmd[:lastSeparatorPos])
	action := strings.TrimSpace(cmd[lastSeparatorPos+1:])
	s, err := os.Stat(dir)
	if err != nil || !s.IsDir() {
		return "", "", fmt.Errorf("invalid directory specified as: %q", dir)
	}
	if action == "" {
		return "", "", fmt.Errorf("invalid action: a non-blank action is required")
	}
	return dir, action, nil
}
