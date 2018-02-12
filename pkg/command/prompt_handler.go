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
	"regexp"
	"strings"
)

var (
	SuPromptLocalizations = []string{
		"Password",
		"암호",
		"パスワード",
		"Adgangskode",
		"Contraseña",
		"Contrasenya",
		"Hasło",
		"Heslo",
		"Jelszó",
		"Lösenord",
		"Mật khẩu",
		"Mot de passe",
		"Parola",
		"Parool",
		"Pasahitza",
		"Passord",
		"Passwort",
		"Salasana",
		"Sandi",
		"Senha",
		"Wachtwoord",
		"ססמה",
		"Лозинка",
		"Парола",
		"Пароль",
		"गुप्तशब्द",
		"शब्दकूट",
		"సంకేతపదము",
		"හස්පදය",
		"密码",
		"密碼",
		"口令",
	}
)

type (
	PromptHandler func(data string) (bool, error)
)

// sudo
func SudoPromptHandler(prompt string) PromptHandler {
	return func(data string) (bool, error) {
		if i := strings.Index(data, prompt); i > -1 {
			return true, nil
		}

		return false, nil
	}
}

// su
func SuPromptHandler(prompt string) PromptHandler {
	if prompt == "" {
		prompts := make([]string, 0, len(SuPromptLocalizations))
		for _, p := range SuPromptLocalizations {
			prompts = append(prompts, `(\w+\'s )?`+p)
		}

		prompt = strings.Join(prompts, "|") + ` ?(:|：) ?`
	}

	return func(data string) (bool, error) {
		return regexp.MatchString(prompt, data)
	}
}

// pbrun
func PbrunPromptHandler(prompt string) PromptHandler {
	return func(data string) (bool, error) {
		return prompt == data, nil
	}
}

// ksu
func KsuPromptHandler(prompt string) PromptHandler {
	return func(data string) (bool, error) {
		return regexp.MatchString(prompt, data)
	}
}

// doas
func DoasPromptHandler(prompt string) PromptHandler {
	return func(data string) (bool, error) {
		return prompt == data, nil
	}
}

// dzdo
func DzdoPromptHandler(prompt string) PromptHandler {
	return func(data string) (bool, error) {
		return prompt == data, nil
	}
}
