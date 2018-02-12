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
	"time"
)

type (
	Connection struct {
		client  *Client
		session *Session
	}
)

func NewConnection(clientConfig ClientConfig, timeout time.Duration, cipherList []string) (*Connection, error) {
	client := NewClient(clientConfig)
	if _, err := client.Connect(timeout, cipherList); err != nil {
		return nil, err
	}

	return &Connection{
		client: client,
	}, nil
}

func (c *Connection) OpenSession() error {
	session, err := c.client.NewSession()
	if err != nil {
		return err
	}

	c.session = session

	return nil
}

func (c *Connection) Run(cmd string) ([]byte, error) {
	if c.session == nil {
		return nil, ErrSessionNoStarted
	}

	return c.session.Run(cmd)
}

func (c *Connection) ShutdownSession() error {
	if c.session == nil {
		return ErrSessionNoStarted
	}

	if err := c.session.Close(); err != nil {
		return err
	}

	return nil
}

func (c *Connection) Close() error {
	return c.client.Close()
}
