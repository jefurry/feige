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
	"errors"
	"fmt"
)

var (
	booleans_string_true  = []string{"y", "yes", "on", "1", "true", "t"}
	booleans_string_false = []string{"n", "no", "off", "0", "false", "f"}
)

func Boolean(value interface{}, strict bool) (bool, error) {
	switch v := value.(type) {
	case bool:
		return v, nil
	case int, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
		if v == 1 {
			return true, nil
		}
		if v == 0 {
			return false, nil
		}
	case float32, float64:
		if v == 1.0 {
			return true, nil
		}
		if v == 0.0 {
			return false, nil
		}
	case string:
		if IndexOf(booleans_string_true, v) > -1 {
			return true, nil
		}
		if IndexOf(booleans_string_false, v) > -1 || !strict {
			return false, nil
		}
	}

	return false, errors.New(fmt.Sprintf("%v is not a valid boolean.", value))
}
