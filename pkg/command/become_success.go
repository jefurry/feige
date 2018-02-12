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

package command

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
)

func (c *Command) FindBecomeSuccess(buffer []byte) int {
	return bytes.Index(buffer, []byte(c.SuccessKey()))
}

// handler the stdout outputs and parse
func (c *Command) HandleBecomeSuccess(stdin io.Writer, stdout io.Reader) ([]byte, error) {
	if c.SuccessKey() == "" {
		return nil, ErrCmdKeyIsEmpty
	}

	succKey := []byte(c.SuccessKey())
	klen := len(succKey)

	b := make([]byte, 4096)
	buffer := bytes.NewBuffer(nil)

	passPrompt := false

	var out []byte
	var n int
	var err error
	i, j := -1, -1
	for {
		n, err = stdout.Read(b)
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
		}
		if n <= 0 {
			break
		}

		_, err = buffer.Write(b[:n])
		if err != nil && err != io.EOF {
			return nil, err
		}

		if err := c.CheckErrors(buffer.Bytes()); err != nil {
			return nil, err
		}

		if passPrompt == false {
			ok, err := c.HandlePasswordPrompt(buffer.Bytes(), stdin)
			if err != nil {
				return nil, err
			}
			if ok {
				passPrompt = true
			}
		}

		if i < 0 {
			// find the first succKey
			i = c.FindBecomeSuccess(buffer.Bytes())
			if i < 0 {
				continue
			}

			if _, err := stdin.Write([]byte("exit\n")); err != nil {
				return nil, err
			}

			// skip (i + klen) count bytes
			buffer.Next(i + klen)
		}

		if i >= 0 && j < 0 {
			// find the second succKey
			j = c.FindBecomeSuccess(buffer.Bytes())
			if j < 0 {
				continue
			}

			buffer.Truncate(j)
		}

		if i > -1 && j > -1 {
			// find the programe exit code
			k := bytes.LastIndex(buffer.Bytes(), []byte("("))
			if k < 0 {
				continue
			}

			bs := buffer.Bytes()
			out = bytes.TrimSpace(bs[:k])

			n := bytes.LastIndex(buffer.Bytes(), []byte(")"))
			if n < 0 {
				continue
			}

			statusStr := string(bytes.TrimSpace(bs[k+1 : n]))
			status, err := strconv.Atoi(statusStr)
			if err != nil {
				return nil, ErrCommandExitCodeInvalid
			}
			if status != 0 {
				return nil, errors.New(string(bytes.TrimSpace(out)))
			}

			buffer.Truncate(0)

			return out, nil
		}
	}

	if buffer.Len() > 0 {
		return nil, errors.New(fmt.Sprintf("command execute failed: '%s'", c.cmd))
	}

	return nil, nil
}
