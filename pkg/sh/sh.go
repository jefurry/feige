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

package sh

import (
	gosh "github.com/codeskyblue/go-sh"
	"github.com/jefurry/feige/app"
	"github.com/jefurry/feige/utils"
)

type (
	SH struct {
		executable string
	}
)

func NewSH(executable string) *SH {
	sh := &SH{
		executable: executable,
	}

	if sh.executable == "" {
		sh.executable = app.DEFAULT_BECOME_EXECUTABLE
	}

	return sh
}

func (sh *SH) SetExecutable(executable string) {
	sh.executable = executable
}

func (sh *SH) Executable() string {
	return sh.executable
}

func (sh *SH) Run(cmd string) ([]byte, error) {
	parts, err := utils.ShlexSplit(cmd)
	if err != nil {
		return nil, err
	}

	c := parts[0]
	parts = parts[1:]
	iparts := make([]interface{}, len(parts))
	for i, v := range parts {
		iparts[i] = v
	}

	s := gosh.NewSession()
	out, err := s.Command(c, iparts...).Output()
	if err != nil {
		return nil, err
	}

	return out, nil
}

// support pipe io redirect
func (sh *SH) RunTrick(cmd string) ([]byte, error) {
	s := gosh.NewSession()

	return s.Command(sh.executable, "-c", cmd).Output()
}
