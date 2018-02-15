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
	"bytes"
	"github.com/jefurry/feige/pkg/command"
	"golang.org/x/crypto/ssh"
	"strings"
)

// execute shell command on remote host
func (s *Session) Shell(cmd *command.Command, terminal Terminal) ([]byte, error) {
	if !s.Started() {
		return nil, ErrSessionNoStarted
	}

	if err := terminal.RequestPty(s.session); err != nil {
		return nil, err
	}

	defer s.Close()

	stdout, err := s.session.StdoutPipe()
	if err != nil {
		return nil, err
	}

	stdin, err := s.session.StdinPipe()
	if err != nil {
		return nil, err
	}
	defer stdin.Close()

	var stderr bytes.Buffer
	s.session.Stderr = &stderr

	if err := s.session.Shell(); err != nil {
		return nil, err
	}

	becomeCMD, err := cmd.MakeBecome()
	if err != nil {
		return nil, err
	}

	becomeCMD = strings.TrimRight(becomeCMD, "\n") + "\n"
	if _, err := stdin.Write([]byte(becomeCMD)); err != nil {
		return nil, err
	}

	if cmd.SendOnly() {
		return nil, nil
	}

	buffer, err := cmd.HandleBecomeSuccess(stdin, stdout)
	if err != nil {
		return nil, err
	}

	if err := s.session.Wait(); err != nil {
		if e, ok := err.(*ssh.ExitError); !ok {
			return nil, e
		}
	}

	if err := cmd.HandleStderr(stderr.Bytes()); err != nil {
		return nil, err
	}

	return buffer, nil
}
