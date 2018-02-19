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
	"errors"
	"fmt"
	"github.com/jefurry/feige/pkg/command"
	"golang.org/x/crypto/ssh"
	"time"
)

var (
	ErrClientNoConnected = errors.New("ssh: client no connected")
)

type (
	Client struct {
		client *ssh.Client
		config ClientConfig
	}
)

func NewClient(config ClientConfig) *Client {
	return &Client{
		config: config,
	}
}

// connect to remote host by ssh protocol
func (sc *Client) Connect(timeout time.Duration, cipherList []string) (*Client, error) {
	cc, err := sc.config.Make(timeout, cipherList)
	if err != nil {
		return nil, err
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", sc.config.Host, sc.config.Port), cc)
	if err != nil {
		return nil, err
	}

	sc.client = client

	return sc, nil
}

func (sc *Client) NewSession() (*Session, error) {
	if !sc.Started() {
		return nil, ErrClientNoConnected
	}

	return NewSession(sc.client)
}

func (sc *Client) Run(cmd string) ([]byte, error) {
	session, err := sc.NewSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	return sc.Run(cmd)
}

func (sc *Client) Shell(cmd *command.Command, terminal Terminal) ([]byte, error) {
	session, err := sc.NewSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	return session.Shell(cmd, terminal)
}

func (sc *Client) Started() bool {
	return sc.client != nil
}

func (sc *Client) Client() *ssh.Client {
	return sc.client
}

func (sc *Client) Close() error {
	if sc.Started() {
		if err := sc.client.Close(); err != nil {
			return err
		}

		sc.client = nil
	}

	return nil
}
