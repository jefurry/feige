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
	"fmt"
	"github.com/jefurry/feige/pkg/command"
	"io"
	"os"
	"path/filepath"
)

func (scp *Scp) putFile(local, remote string, writer io.Writer, reader io.Reader) (int64, error) {
	cmd := scp.cmd
	if cmd.SuccessKey() == "" {
		return 0, command.ErrCmdKeyIsEmpty
	}

	succKey := []byte(cmd.SuccessKey())
	klen := len(succKey)

	b := make([]byte, 4096)
	buffer := bytes.NewBuffer(nil)

	passPrompt := false
	startScp := false

	var n int
	var err error
	i := -1
	var total int64
	for {
		n, err = reader.Read(b)
		if err != nil {
			if err != io.EOF {
				return int64(n), err
			}
		}
		if n <= 0 {
			break
		}

		n, err = buffer.Write(b[:n])
		if err != nil && err != io.EOF {
			return total, err
		}

		if err := cmd.CheckErrors(buffer.Bytes()); err != nil {
			return total, err
		}

		if err := scp.checkScpErrors(buffer.Bytes()); err != nil {
			return 0, err
		}

		if err := scp.checkScpEnding(buffer.Bytes()); err != nil {
			if startScp {
				return total, nil
			}

			return 0, err
		}

		if passPrompt == false {
			ok, err := cmd.HandlePasswordPrompt(buffer.Bytes(), writer)
			if err != nil {
				return total, err
			}
			if ok {
				passPrompt = true
			}
		}

		if i < 0 {
			// find the first succKey
			i = cmd.FindBecomeSuccess(buffer.Bytes())
			if i < 0 {
				continue
			}

			// skip (i + klen) count bytes
			buffer.Next(i + klen)
		}

		if i > -1 && !startScp {
			startScp = true
			// copy from local to remote
			localFile, err := os.Open(local)
			if err != nil {
				return 0, err
			}
			defer localFile.Close()

			fileInfo, err := localFile.Stat()
			if err != nil {
				return 0, err
			}

			size := fileInfo.Size()
			n, err = fmt.Fprintf(writer, "C%#o %d %s\n", fileInfo.Mode().Perm(), size, filepath.Base(remote))
			if err != nil {
				return total, err
			}

			total, err = copyN(writer, localFile, size)
			if err != nil {
				return total, err
			}
			if total != size {
				return total, ErrScpFailed
			}

			err = ack(writer)
			if err != nil {
				return total, err
			}

			if _, err := writer.Write([]byte("exit\n")); err != nil {
				return total, err
			}
		}
	}

	return 0, nil
}
