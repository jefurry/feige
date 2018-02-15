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
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"strconv"
)

const (
	DEFAULT_TERM_NAME   = "vt100"
	DEFAULT_TERM_HEIGHT = 80
	DEFAULT_TERM_WIDTH  = 40
)

const (
	DEFAULT_TERM_ECHO     = 0     // disable echoing
	DEFAULT_TERM_ISIG     = 1     // Enable signals INTR, QUIT, [D]SUSP.
	DEFAULT_TTY_OP_ISPEED = 14400 // input speed = 14.4kbaud
	DEFAULT_TTY_OP_OSPEED = 14400 // output speed = 14.4kbaud
)

func TerminalSize(h, w int) (int, int) {
	var fd int

	var th, tw int
	if h < 0 || w < 0 {
		fd = int(os.Stdin.Fd())
		if terminal.IsTerminal(fd) {
			if oldState, err := terminal.MakeRaw(fd); err == nil {
				defer terminal.Restore(fd, oldState)
				tw, th, _ = terminal.GetSize(fd)
			}
		}
	}

	if h < 0 {
		if th >= 0 {
			h = th
		} else {
			h = TerminalHeight(h)
		}
	}

	if w < 0 {
		if tw >= 0 {
			w = tw
		} else {
			w = TerminalWidth(w)
		}
	}

	return h, w
}

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

	i, err := strconv.Atoi(os.Getenv("LINES"))
	if err == nil && i >= 0 {
		return i
	}

	return DEFAULT_TERM_HEIGHT
}

func TerminalWidth(w int) int {
	if w >= 0 {
		return w
	}

	i, err := strconv.Atoi(os.Getenv("COLUMNS"))
	if err == nil && i >= 0 {
		return i
	}

	return DEFAULT_TERM_WIDTH
}

func DefaultTerminalModels() ssh.TerminalModes {
	return ssh.TerminalModes{
		ssh.ECHO:          DEFAULT_TERM_ECHO, // disable echoing
		ssh.ISIG:          DEFAULT_TERM_ISIG,
		ssh.TTY_OP_ISPEED: DEFAULT_TTY_OP_ISPEED, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: DEFAULT_TTY_OP_OSPEED, // output speed = 14.4kbaud
	}
}
