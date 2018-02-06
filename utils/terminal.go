// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Feige
//
// Feige is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Feige is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with Feige.  If not, see <http://www.gnu.org/licenses/>.

package utils

import (
	"os"
	"strconv"
)

const (
	DEFAULT_TERM_NAME   = "vt100"
	DEFAULT_TERM_HEIGHT = 80
	DEFAULT_TERM_WIDTH  = 40
)

func TerminalName(name string) string {
	if name != "" {
		return name
	}

	name = os.Getenv("TERM")
	if name != "" {
		return name
	}

	return DEFAULT_TERM_NAME
}

func TerminalHeight(h int) int {
	if h >= 0 {
		return h
	}

	i, err := strconv.Atoi(os.Getenv("COLUMNS"))
	if err == nil && i >= 0 {
		return i
	}

	return DEFAULT_TERM_HEIGHT
}

func TerminalWidth(w int) int {
	if w >= 0 {
		return w
	}

	i, err := strconv.Atoi(os.Getenv("LINES"))
	if err == nil && i >= 0 {
		return i
	}

	return DEFAULT_TERM_WIDTH
}
