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
	"github.com/mitchellh/go-homedir"
	"path/filepath"
)

const (
	// sub director
	SUB_DIR_CONF    = "conf"
	SUB_DIR_LUALIBS = "lualibs"
	SUB_DIR_PLUGINS = "plugins"
	SUB_DIR_MODULES = "modules"
	SUB_DIR_TEMP    = "tmp"
)

var (
	programHomeDirCache string
)

func ProgramHomeDir(programName string, disableCache bool) (string, error) {
	homedir.DisableCache = disableCache
	if !disableCache {
		if programHomeDirCache != "" {
			return programHomeDirCache, nil
		}
	}

	var dir string
	var err error

	dir, err = homedir.Expand("~/." + programName)
	if err != nil {
		return "", err
	}

	if err = EnsurePath(dir, true); err != nil {
		return "", err
	}

	programHomeDirCache = dir

	return dir, nil
}

// initialize working directory structure
func InitHomeDir(programName string) error {
	home, err := ProgramHomeDir(programName, true)
	if err != nil {
		return err
	}

	subDirs := []string{
		SUB_DIR_CONF,
		SUB_DIR_LUALIBS,
		SUB_DIR_PLUGINS,
		SUB_DIR_MODULES,
		SUB_DIR_TEMP,
	}

	for _, d := range subDirs {
		p := filepath.Join(home, d)
		if err = EnsurePath(p, true); err != nil {
			return err
		}
	}

	return nil
}
