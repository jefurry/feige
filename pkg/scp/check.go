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
	"errors"
)

func (scp *Scp) checkScpErrors(buffer []byte) error {
	errmsg := "scp: protocol error: bad mode"
	if bytes.Index(buffer, []byte(errmsg)) > -1 {
		return errors.New(errmsg)
	}

	return nil
}

func (scp *Scp) checkScpEnding(buffer []byte) error {
	errmsg := "scp: protocol error: expected control record"
	if bytes.Index(buffer, []byte(errmsg)) > -1 {
		return errors.New(errmsg)
	}

	return nil
}
