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
	"github.com/jefurry/feige/app"
	"strings"
)

// find become method
// return true if found else false
func CheckBecomeMethod(method string) bool {
	for _, m := range app.BECOME_METHODS {
		if m == method {
			return true
		}
	}

	return false
}

// check unknown user for become
func CheckBecomeUnknownUser(buffer []byte, method, username string) error {
	if v, ok := app.BECOME_UNKNOWN_USER[method]; ok {
		for _, s := range v {
			prompt := fmt.Sprintf(s, username)
			if strings.Index(string(buffer), prompt) > -1 {
				return errors.New(prompt)
			}
		}

		return nil
	}

	return errors.New(fmt.Sprintf("unknown become method '%s'", method))
}

func CheckBecomeUnknownExec(buffer []byte, method, becomeExe string) error {
	if v, ok := app.BECOME_UNKNOWN_EXE[method]; ok {
		for _, s := range v {
			prompt := fmt.Sprintf(s, becomeExe)
			if strings.Index(string(buffer), prompt) > -1 {
				return errors.New(prompt)
			}
		}

		return nil
	}

	return errors.New(fmt.Sprintf("unknown become method '%s'", method))
}

func CheckBecomeUnknownExecutable(buffer []byte, method, executable string) error {
	if v, ok := app.BECOME_UNKNOWN_EXECUTABLE[method]; ok {
		for _, s := range v {
			prompt := fmt.Sprintf(s, executable)
			if strings.Index(string(buffer), prompt) > -1 {
				return errors.New(prompt)
			}
		}

		return nil
	}

	return errors.New(fmt.Sprintf("unknown become method '%s'", method))
}

func CheckBecomeErrorString(buffer []byte, method, username string) error {
	if v, ok := app.BECOME_ERROR_STRINGS[method]; ok {
		for _, s := range v {
			if strings.Index(string(buffer), s) > -1 {
				return errors.New(fmt.Sprintf("incorrect password given for user '%s'", username))
			}
		}

		return nil
	}

	return errors.New(fmt.Sprintf("unknown become method '%s'", method))
}
