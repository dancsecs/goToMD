//nolint:funlen // Ok.
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
	"math"
	"strings"

	"github.com/dancsecs/szTest"
)

const (
	markInsOn  = "<{INS_ON}>"
	markInsOff = "<{INS_OFF}>"
	markDelOn  = "<{DEL_ON}>"
	markDelOff = "<{DEL_OFF}>"
	markChgOn  = "<{CHG_ON}>"
	markChgOff = "<{CHG_OFF}>"
	markWntOn  = "<{WNT_ON}>"
	markWntOff = "<{WNT_OFF}>"
	markGotOn  = "<{GOT_ON}>"
	markGotOff = "<{GOT_OFF}>"
	markSepOn  = "<{SEP_ON}>"
	markSepOff = "<{SEP_OFF}>"
	markMsgOn  = "<{MSG_ON}>"
	markMsgOff = "<{MSG_OFF}>"
)

const (
	internalTestMarkDelOn  = `\color{red}`
	internalTestMarkDelOff = `\color{default}`
	internalTestMarkInsOn  = `\color{green}`
	internalTestMarkInsOff = `\color{default}`
	internalTestMarkChgOn  = `\color{darkturquoise}`
	internalTestMarkChgOff = `\color{default}`
	internalTestMarkSepOn  = `\color{yellow}`
	internalTestMarkSepOff = `\color{default}`
	internalTestMarkWntOn  = `\color{cyan}`
	internalTestMarkWntOff = `\color{default}`
	internalTestMarkGotOn  = `\color{magenta}`
	internalTestMarkGotOff = `\color{default}`
	internalTestMarkMsgOn  = `\emph{`
	internalTestMarkMsgOff = `}`
)

// findNextMark searches the string for all known marks.
func findNextMark(s, expectedClose string,
) (int, string, string, string) {
	if s == "" {
		return -1, "", "", ""
	}

	markOpenIndex := math.MaxInt
	markOpen := ""
	markOpenInternal := ""
	markOpenExpectedInternal := ""

	findOnMark := func(eOpenMark, iOpenMark, iCloseMark string) {
		tmpIndex := strings.Index(s, eOpenMark)
		if tmpIndex >= 0 && tmpIndex < markOpenIndex {
			markOpenIndex = tmpIndex
			markOpen = eOpenMark
			markOpenInternal = iOpenMark
			markOpenExpectedInternal = iCloseMark
		}
	}

	findOnMark(szTest.SettingMarkInsOn(), markInsOn, markInsOff)
	findOnMark(szTest.SettingMarkDelOn(), markDelOn, markDelOff)
	findOnMark(szTest.SettingMarkChgOn(), markChgOn, markChgOff)
	findOnMark(szTest.SettingMarkWntOn(), markWntOn, markWntOff)
	findOnMark(szTest.SettingMarkGotOn(), markGotOn, markGotOff)
	findOnMark(szTest.SettingMarkSepOn(), markSepOn, markSepOff)
	findOnMark(szTest.SettingMarkMsgOn(), markMsgOn, markMsgOff)

	markCloseIndex := math.MaxInt
	markClose := ""
	markCloseInternal := ""

	findOffMark := func(mark, internalMark string) {
		tmpIndex := strings.Index(s, mark)
		if tmpIndex >= 0 &&
			tmpIndex < markOpenIndex &&
			tmpIndex <= markCloseIndex {
			if tmpIndex == markCloseIndex && markCloseInternal == expectedClose {
				return
			}
			markCloseIndex = tmpIndex
			markClose = mark
			markCloseInternal = internalMark
		}
	}

	findOffMark(szTest.SettingMarkInsOff(), markInsOff)
	findOffMark(szTest.SettingMarkDelOff(), markDelOff)
	findOffMark(szTest.SettingMarkChgOff(), markChgOff)
	findOffMark(szTest.SettingMarkWntOff(), markWntOff)
	findOffMark(szTest.SettingMarkGotOff(), markGotOff)
	findOffMark(szTest.SettingMarkSepOff(), markSepOff)
	findOffMark(szTest.SettingMarkMsgOff(), markMsgOff)

	if markOpenIndex < math.MaxInt || markCloseIndex < math.MaxInt {
		if markOpenIndex < markCloseIndex {
			return markOpenIndex,
				markOpen,
				markOpenInternal,
				markOpenExpectedInternal
		} else {
			return markCloseIndex, markClose, markCloseInternal, ""
		}
	}
	return -1, "", "", ""
}

func translateToTestSymbols(s string) string {
	s = strings.ReplaceAll(s, markDelOn, internalTestMarkDelOn)
	s = strings.ReplaceAll(s, markDelOff, internalTestMarkDelOff)
	s = strings.ReplaceAll(s, markInsOn, internalTestMarkInsOn)
	s = strings.ReplaceAll(s, markInsOff, internalTestMarkInsOff)
	s = strings.ReplaceAll(s, markChgOn, internalTestMarkChgOn)
	s = strings.ReplaceAll(s, markChgOff, internalTestMarkChgOff)
	s = strings.ReplaceAll(s, markSepOn, internalTestMarkSepOn)
	s = strings.ReplaceAll(s, markSepOff, internalTestMarkSepOff)
	s = strings.ReplaceAll(s, markWntOn, internalTestMarkWntOn)
	s = strings.ReplaceAll(s, markWntOff, internalTestMarkWntOff)
	s = strings.ReplaceAll(s, markGotOn, internalTestMarkGotOn)
	s = strings.ReplaceAll(s, markGotOff, internalTestMarkGotOff)
	s = strings.ReplaceAll(s, markMsgOn, internalTestMarkMsgOn)
	s = strings.ReplaceAll(s, markMsgOff, internalTestMarkMsgOff)
	return s
}

func marksToMarkdownHTML(source string) (string, error) {
	iCloseMarkExpected := ""
	newS := ""
	for {
		i, eNextMark, iNextMark, iNextCloseMark :=
			findNextMark(source, iCloseMarkExpected)

		// If no more marks are present then we are done.  Either return the
		// translated string with the all marks reversed or an error if we are
		// expecting a close mark.
		if i < 0 {
			if iCloseMarkExpected != "" {
				return "", fmt.Errorf(
					"no closing mark found for %q in %q",
					iCloseMarkExpected,
					source,
				)
			}
			return translateToTestSymbols(newS + source), nil
		}

		// Otherwise we found a Mark.  Move all text up to the next mark from
		// the string to the translated string.
		if i > 0 {
			newS += source[:i]
			source = source[i:]
		}

		// Add the internal representation, replacing the resolved marks.
		newS += iNextMark

		// Remove the resolved Mark from the source string
		source = source[len(eNextMark):]

		if iCloseMarkExpected != "" {
			// There is an open mark that needs to be closed.
			if iNextMark != iCloseMarkExpected {
				return "", fmt.Errorf(
					"unexpected closing mark: Got: %q  Want: %q",
					iNextMark,
					iCloseMarkExpected,
				)
			}
			iCloseMarkExpected = ""
		} else {
			iCloseMarkExpected = iNextCloseMark
		}
	}
}
