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
	"errors"
	"golang.org/x/crypto/ssh"
)

var (
	ErrSessionNoStarted = errors.New("ssh: session no started")
)

type (
	Session struct {
		session *ssh.Session
	}
)

func NewSessionWithSession(session *ssh.Session) *Session {
	return &Session{
		session: session,
	}
}

func NewSession(client *ssh.Client) (*Session, error) {
	session, err := client.NewSession()
	if err != nil {
		return nil, err
	}

	return NewSessionWithSession(session), nil
}

// execute command on remote host
// run once
func (s *Session) Run(cmd string) ([]byte, error) {
	if !s.Started() {
		return nil, ErrSessionNoStarted
	}

	defer s.Close()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	s.session.Stdout = &stdout
	s.session.Stderr = &stderr

	if err := s.session.Run(cmd); err != nil {
		return nil, err
	}

	if stderr.Len() > 0 {
		return nil, errors.New(stderr.String())
	}

	return bytes.TrimSpace(stdout.Bytes()), nil
}

func (s *Session) Started() bool {
	return s.session != nil
}

func (s *Session) Session() *ssh.Session {
	return s.session
}

func (s *Session) Setenv(name, value string) error {
	if !s.Started() {
		return ErrSessionNoStarted
	}

	return s.session.Setenv(name, value)
}

func (s *Session) Close() error {
	if s.Started() {
		if err := s.session.Close(); err != nil {
			return err
		}

		s.session = nil
	}

	return nil
}
