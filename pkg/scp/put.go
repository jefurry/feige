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
	"bytes"
	"github.com/jefurry/feige/utils"
	"golang.org/x/crypto/ssh"
	"path/filepath"
	"strings"
)

func (scp *Scp) Put(local, remote string) (int64, error) {
	sess, err := scp.client.NewSession()
	if err != nil {
		return 0, err
	}
	defer sess.Close()

	session := sess.Session()
	if err := scp.terminal.RequestPty(session); err != nil {
		return 0, nil
	}

	scp.cmd.SetCMD(utils.ShlexJoin(scp.cmd.CMD(), "-tq", filepath.Dir(remote)))

	stdout, err := session.StdoutPipe()
	if err != nil {
		return 0, err
	}

	stdin, err := session.StdinPipe()
	if err != nil {
		return 0, err
	}
	defer stdin.Close()

	var stderr bytes.Buffer
	session.Stderr = &stderr

	if err := session.Shell(); err != nil {
		return 0, err
	}

	becomeCMD, err := scp.cmd.MakeBecome()
	if err != nil {
		return 0, err
	}

	becomeCMD = strings.TrimRight(becomeCMD, "\n") + "\n"
	if _, err := stdin.Write([]byte(becomeCMD)); err != nil {
		return 0, err
	}

	// Listen the cmd stop signal
	go func() {
		stopChan := make(chan error)
		stopChan <- session.Wait()
	}()

	total, err := scp.putFile(local, remote, stdin, stdout)
	if err != nil {
		return 0, err
	}

	if err := scp.cmd.HandleStderr(stderr.Bytes()); err != nil {
		return 0, err
	}

	session.Signal(ssh.SIGINT)

	return total, nil
}
