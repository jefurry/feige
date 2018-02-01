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

package lua

import (
	"github.com/yuin/gopher-lua"
	"testing"
)

func TestLuaPath(t *testing.T) {
	if err := SetLuaPath(""); err != nil {
		t.Fatalf(err.Error())
	}

	if lua.LuaPathDefault != "" {
		t.Fatalf("LuaPathDefault mismatching")
	}

	if err := SetLuaPath("~/.feige/modules"); err != nil {
		t.Fatalf(err.Error())
	}
	if lua.LuaPathDefault != "/Users/jefurry/.feige/modules/?.lua;/Users/jefurry/.feige/modules/?/init.lua" {
		t.Fatalf("LuaPathDefault mismatching")
	}

	if err := AddLuaPath("~/.feige/lublibs"); err != nil {
		t.Fatalf(err.Error())
	}
	if lua.LuaPathDefault != "/Users/jefurry/.feige/modules/?.lua;/Users/jefurry/.feige/modules/?/init.lua;/Users/jefurry/.feige/lublibs/?.lua;/Users/jefurry/.feige/lublibs/?/init.lua" {
		t.Fatalf("LuaPathDefault mismatching")
	}

	lua.LuaPathDefault = OldLuaPathDefault
	if lua.LuaPathDefault != "./?.lua;/usr/local/share/lua/5.1/?.lua;/usr/local/share/lua/5.1/?/init.lua" {
		t.Fatalf("LuaPathDefault mismatching")
	}

	if err := AddLuaPath("~/.feige/modules"); err != nil {
		t.Fatalf(err.Error())
	}
	if lua.LuaPathDefault != "./?.lua;/usr/local/share/lua/5.1/?.lua;/usr/local/share/lua/5.1/?/init.lua;/Users/jefurry/.feige/modules/?.lua;/Users/jefurry/.feige/modules/?/init.lua" {
		t.Fatalf("LuaPathDefault mismatching")
	}
}
