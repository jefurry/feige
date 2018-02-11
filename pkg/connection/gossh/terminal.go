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

package gossh

import (
	"github.com/jefurry/feige/utils"
	"golang.org/x/crypto/ssh"
)

const (
	DEFAULT_TERM_ECHO     = 0     // disable echoing
	DEFAULT_TERM_ISIG     = 1     // Enable signals INTR, QUIT, [D]SUSP.
	DEFAULT_TTY_OP_ISPEED = 14400 // input speed = 14.4kbaud
	DEFAULT_TTY_OP_OSPEED = 14400 // output speed = 14.4kbaud
)

// RequestPty requests the association of a pty with the session on the remote host.
//
// termmodes := ssh.TerminalModes{
//		ssh.ECHO:          0,     // disable echoing
//		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
//		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
// }
//
func (s *Session) RequestPty(term string, h, w int, termModes ssh.TerminalModes) error {
	if !s.Started() {
		return ErrSessionNoStarted
	}

	term = utils.TerminalName(term)
	h, w = utils.TerminalSize(h, w)

	if err := s.session.RequestPty(term, h, w, termModes); err != nil {
		return err
	}

	return nil
}

func (s *Session) DefaultTerminalModels() ssh.TerminalModes {
	return ssh.TerminalModes{
		ssh.ECHO:          DEFAULT_TERM_ECHO, // disable echoing
		ssh.ISIG:          DEFAULT_TERM_ISIG,
		ssh.TTY_OP_ISPEED: DEFAULT_TTY_OP_ISPEED, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: DEFAULT_TTY_OP_OSPEED, // output speed = 14.4kbaud
	}
}
