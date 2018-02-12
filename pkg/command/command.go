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

package command

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/jefurry/feige/app"
	"github.com/jefurry/feige/utils"
	"strings"
)

var (
	ErrCommandExitCodeInvalid = errors.New("command exit code invalid")
	ErrCmdKeyIsEmpty          = errors.New("command key is empty, maybe forget calls the MakeBecome function")
)

type (
	User struct {
		Username, Password string
	}

	Command struct {
		becomeExe     string
		becomeMethod  string
		becomeFlags   string
		cmd           string
		answer        string
		becomeUser    User
		become        bool
		sendOnly      bool
		promptHandler PromptHandler // handler stdout and write answer into stdin

		successKey string
		executable string
	}
)

var (
	DefaultCommand = NewCommand(app.DEFAULT_BECOME_EXE, app.DEFAULT_BECOME_EXECUTABLE, app.DEFAULT_BECOME_METHOD,
		app.DEFAULT_BECOME_FLAGS, "", "", User{Username: app.DEFAULT_BECOME_USER}, false, false)
)

func NewScpCommand(becomeExe, executable, becomeMethod, becomeFlags, scpExe, answer string, becomeUser User, become, sendOnly bool) *Command {
	if scpExe == "" {
		scpExe = app.DEFAULT_BECOME_SCP
	}

	return NewCommand(becomeExe, executable, becomeMethod, becomeFlags, scpExe, answer, becomeUser, become, sendOnly)
}

// becomeExe: examples /usr/bin/sudo, /bin/su etc.
// executable: examples /bin/sh, /bin/bash etc.
// becomemMethod: examples sudo, su etc.
// becomeFlags: examples -n etc.
// cmd: examples ls, pwd etc.
// answer: fill in the answer
// becomeUser: become user
// become: whether to executd as a other user
// sendonly: whether to wait for a response
func NewCommand(becomeExe, executable, becomeMethod, becomeFlags, cmd, answer string, becomeUser User, become, sendOnly bool) *Command {
	c := &Command{
		becomeExe:    becomeExe,
		becomeMethod: becomeMethod,
		becomeFlags:  becomeFlags,
		cmd:          cmd,
		answer:       answer,
		becomeUser:   becomeUser,
		become:       become,
		sendOnly:     sendOnly,
		executable:   executable,
	}

	if c.becomeExe == "" {
		c.becomeExe = app.DEFAULT_BECOME_EXE
	}
	if c.executable == "" {
		c.executable = app.DEFAULT_BECOME_EXECUTABLE
	}
	if c.becomeMethod == "" {
		c.becomeMethod = app.DEFAULT_BECOME_METHOD
	}
	if c.becomeFlags == "" {
		c.becomeFlags = app.DEFAULT_BECOME_FLAGS
	}

	return c
}

// create privilege escalation commands
// executable: examples /bin/sh, /usr/bin/lua etc.
func (c *Command) MakeBecome() (string, error) {
	randBits := utils.RandKey()
	succKey := fmt.Sprintf(`$?: %s-BECOME-SUCCESS`, randBits)

	// hack with fuck the stdout of session
	succCMD := fmt.Sprintf(`echo "%s"; %s; echo "("$?")"; echo "%s"`,
		succKey, c.cmd, succKey)
	//succCMD := fmt.Sprintf(`echo '%s'; %s; echo '%s'`, succKey, c.cmd, succKey)

	c.successKey = strings.Replace(succKey, "$?", "0", 1)

	if c.become == false {
		return succCMD, nil
	}

	command := fmt.Sprintf(`%s -c %s`, c.executable, utils.ShlexQuote(succCMD))

	becomeExe := c.becomeMethod
	if c.becomeExe != "" {
		becomeExe = c.becomeExe
	}

	becomeUsername := utils.FindUsername(c.becomeUser.Username)
	becomeCMD := ""
	if c.becomeMethod == "sudo" {
		if c.becomeUser.Password != "" {
			prompt := fmt.Sprintf(`[sudo via %s, key=%s]`, app.NameForLower, randBits)
			c.promptHandler = SudoPromptHandler(prompt)

			becomeCMD = fmt.Sprintf(`%s %s -p "%s" -u %s %s`,
				becomeExe, strings.Replace(c.becomeFlags, "-n", "", -1), prompt, becomeUsername, command)
		} else {
			becomeCMD = fmt.Sprintf(`%s %s -u %s %s`, becomeExe, c.becomeFlags, becomeUsername, command)
		}
	} else if c.becomeMethod == "su" {
		// passing code ref to examine prompt as simple string comparisson isn't good enough with su
		c.promptHandler = SuPromptHandler("")

		becomeCMD = fmt.Sprintf(`%s %s %s -c %s`, becomeExe, c.becomeFlags, becomeUsername, utils.ShlexQuote(command))
	} else if c.becomeMethod == "pbrun" {
		prompt := "Password:"
		c.promptHandler = PbrunPromptHandler(prompt)

		becomeCMD = fmt.Sprintf(`%s %s -u %s %s`, becomeExe, c.becomeFlags, becomeUsername, succCMD)
	} else if c.becomeMethod == "ksu" {
		prompt := `Kerberos password for .*@.*:`
		c.promptHandler = KsuPromptHandler(prompt)

		becomeCMD = fmt.Sprintf(`%s %s %s -e %s`, becomeExe, becomeUsername, c.becomeFlags, command)
	} else if c.becomeMethod == "pfexec" {
		// No user as it uses it's own exec_attr to figure it out
		becomeCMD = fmt.Sprintf(`%s %s "%s"`, becomeExe, c.becomeFlags, succCMD)
	} else if c.becomeMethod == "runas" {
		// windows
		if c.becomeUser.Password == "" {
			return "", errors.New(fmt.Sprintf("The 'runas' become method requires a password "+
				"(specify with the '-K' CLI arg or the '%s_become_password' variable)", app.NameForLower))
		}

		becomeCMD = c.cmd
	} else if c.becomeMethod == "doas" {
		prompt := fmt.Sprintf(`doas (%s@`, becomeUsername)
		c.promptHandler = DoasPromptHandler(prompt)

		if c.becomeUser.Password == "" {
			c.becomeFlags = c.becomeFlags + " -n "
		}
		c.becomeFlags = c.becomeFlags + fmt.Sprintf(` -u %s `, becomeUsername)

		becomeCMD = fmt.Sprintf(`%s %s echo %s && %s %s env %s=true %s`,
			becomeExe, c.becomeFlags, succKey, becomeExe, c.becomeFlags, app.NameForUpper, c.cmd)
	} else if c.becomeMethod == "dzdo" {
		if c.becomeUser.Password != "" {
			prompt := fmt.Sprintf(`[dzdo via %s, key=%s] password: `, app.NameForLower, randBits)
			c.promptHandler = DzdoPromptHandler(prompt)

			becomeCMD = fmt.Sprintf(`%s -p %s -u %s %s`, becomeExe, utils.ShlexQuote(prompt), becomeUsername, command)
		} else {
			becomeCMD = fmt.Sprintf(`%s -u %s %s`, becomeExe, becomeUsername, command)
		}
	} else {
		return "", errors.New(fmt.Sprintf("Privilege escalation method not found: %s", c.becomeMethod))
	}

	return becomeCMD, nil
}

// handler the stderr outputs
func (c *Command) HandleStderr(out []byte) error {
	if len(out) > 0 {
		return errors.New(string(bytes.TrimSpace(out)))
	}

	return nil
}

func (c *Command) SuccessKey() string {
	return c.successKey
}

func (c *Command) SendOnly() bool {
	return c.sendOnly
}

func (c *Command) CMD() string {
	return c.cmd
}

func (c *Command) SetCMD(cmd string) {
	c.cmd = cmd
}
