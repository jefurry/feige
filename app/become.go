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

package app

const (
	DEFAULT_BECOME_METHOD     = "sudo"
	DEFAULT_BECOME_EXECUTABLE = "/bin/sh"
	DEFAULT_BECOME_EXE        = "/usr/bin/sudo"
	DEFAULT_BECOME_USER       = "root"
	DEFAULT_BECOME_FLAGS      = ""

	DEFAULT_BECOME_SCP = "/usr/bin/scp"
)

var (
	BECOME_METHODS = []string{
		"sudo", "su",
	}
)

var (
	BECOME_ERROR_STRINGS = map[string][]string{
		"sudo": []string{
			"Sorry, try again.",
		},
		"su": []string{
			"su: incorrect password",
			"Authentication failure",
		},
	}

	BECOME_UNKNOWN_USER = map[string][]string{
		"sudo": []string{
			"sudo: unknown user: %s",
			"%s is not in the sudoers file.  This incident will be reported.",
		},
		"su": []string{
			"su: user %s does not exist",
		},
	}

	BECOME_UNKNOWN_EXE = map[string][]string{
		"sudo": []string{
			"-bash: %s: No such file or directory",
		},
		"su": []string{
			"-bash: %s: No such file or directory",
		},
	}

	BECOME_UNKNOWN_EXECUTABLE = map[string][]string{
		"sudo": []string{
			"sudo: %s: command not found",
		},
		"su": []string{
			"su: %s: command not found",
		},
	}
)
