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

package scp

import (
	"github.com/jefurry/feige/app"
	"github.com/jefurry/feige/pkg/command"
	"github.com/jefurry/feige/pkg/connection/gossh"
)

type (
	Scp struct {
		client   *gossh.Client
		cmd      *command.Command
		terminal gossh.Terminal
	}
)

func NewScp(client *gossh.Client, cmd *command.Command, terminal gossh.Terminal) *Scp {
	c := &Scp{
		client:   client,
		cmd:      cmd,
		terminal: terminal,
	}
	if cmd == nil {
		// with default Command
		c.cmd = command.DefaultCommand
	}

	if c.cmd.CMD() == "" {
		c.cmd.SetCMD(app.DEFAULT_BECOME_SCP)
	}

	return c
}
