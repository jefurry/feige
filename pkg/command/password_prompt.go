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
	"io"
)

func (c *Command) HandlePasswordPrompt(buffer []byte, w io.Writer) (bool, error) {
	if c.promptHandler != nil {
		ok, err := c.promptHandler(string(buffer))
		if err != nil {
			return false, err
		}

		if ok {
			if _, err := w.Write([]byte(c.becomeUser.Password + "\n")); err != nil {
				return false, err
			}

			return true, nil
		}
	}

	return false, nil
}
