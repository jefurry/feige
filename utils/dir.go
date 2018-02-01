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
	"fmt"
	"os"
	"path/filepath"
)

// EnsurePath is used to make sure a path exists
func EnsurePath(path string, dir bool) error {
	if !dir {
		path = filepath.Dir(path)
	}

	err := os.MkdirAll(path, 0755)
	if err != nil {
		if !os.IsExist(err) {
			return err
		}
	}

	return nil
}

func GenHashDir(name string, ismd5 bool) []string {
	if name == "" {
		name = GenName("")
	}
	if !ismd5 {
		name = GenMd5(name)
	}

	code := GenHashCode(name)
	hs := fmt.Sprintf("%06d", code)

	return []string{hs[0:2], hs[2:4], hs[4:6]}
}
