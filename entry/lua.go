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

// Package entry implements LuaState Manages. It Wraps gopher-lua as runner.
package entry

import (
	"errors"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/yuin/gopher-lua"
	"os"
	"strings"
	"sync"
)

const (
	// lua state pool size default to 10
	DEFAULT_LUA_STATE_POOL_SIZE = 10
)

var (
	OldLuaPathDefault   = lua.LuaPathDefault
	ErrLuaStatePoolFull = errors.New("lua state pool fulled")
)

type (
	NewFunc func(l *lua.LState) error
)

func SetLuaPath(p string) error {
	lua.LuaPathDefault = ""

	return AddLuaPath(p)
}

// add one lua path search to LuaPathDefault
func AddLuaPath(p string) error {
	if p == "" {
		return nil
	}

	p, err := homedir.Expand(p)
	if err != nil {
		return err
	}

	var prefix = ""
	if lua.LuaPathDefault != "" {
		prefix = lua.LuaPathDefault + ";"
	}

	p = strings.TrimRight(p, "/")
	p = strings.TrimRight(p, "\\")

	ps := string(os.PathSeparator)
	p = fmt.Sprintf("%s%s%s?.lua;%s%s?%sinit.lua", prefix, p, ps, p, ps, ps)
	lua.LuaPathDefault = p

	return nil
}

type LuaStatePool struct {
	ls        []*lua.LState
	m         sync.Mutex
	size      int
	makedSize int
	whenNew   NewFunc
}

func (lp *LuaStatePool) NewLuaStatePool(whenNew NewFunc, size int) *LuaStatePool {
	if size <= 0 {
		size = DEFAULT_LUA_STATE_POOL_SIZE
	}
	lp.size = size

	return &LuaStatePool{
		ls:      make([]*lua.LState, size),
		whenNew: whenNew,
	}
}

func (lp *LuaStatePool) New() (*lua.LState, error) {
	lp.m.Lock()
	defer lp.m.Unlock()

	lp.makedSize = lp.makedSize + 1
	l := lua.NewState()

	// setting the Lua State up here.
	// load scripts, set global variables, share channels, etc...
	if lp.whenNew != nil {
		if err := lp.whenNew(l); err != nil {
			return nil, err
		}
	}

	return l, nil
}

// get one lua state from pool
func (lp *LuaStatePool) Get() (*lua.LState, error) {
	lp.m.Lock()
	defer lp.m.Unlock()

	n := len(lp.ls)
	if n == 0 {
		if lp.makedSize >= lp.size {
			return nil, ErrLuaStatePoolFull
		}

		return lp.New()
	}

	l := lp.ls[n-1]
	lp.ls = lp.ls[0 : n-1]

	return l, nil
}

// put lua state into pool
func (lp *LuaStatePool) Put(l *lua.LState) {
	lp.m.Lock()
	defer lp.m.Unlock()

	lp.ls = append(lp.ls, l)
}

func (lp *LuaStatePool) Close(l *lua.LState) {
	lp.m.Lock()
	defer lp.m.Unlock()

	l.Close()
	lp.makedSize = lp.makedSize - 1
}

// close all lua state
func (lp *LuaStatePool) Shutdown() {
	for _, l := range lp.ls {
		lp.Close(l)
	}
}
