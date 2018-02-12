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

package utils

import (
	"github.com/flynn-archive/go-shlex"
	shellquote "github.com/kballard/go-shellquote"
	"strconv"
)

func ShlexQuote(s string) string {
	return strconv.Quote(s)
}

func ShlexSplit(s string) ([]string, error) {
	return shlex.Split(s)
}

func ShlexJoin(args ...string) string {
	return shellquote.Join(args...)
}
