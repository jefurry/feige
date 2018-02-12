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

package gossh

import (
	"github.com/mitchellh/go-homedir"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"io/ioutil"
	"net"
	"os"
)

const (
	SSH_AUTH_SOCK = "SSH_AUTH_SOCK"
)

// with public key file
func PublicKeyFile(file string) (ssh.AuthMethod, error) {
	f, err := homedir.Expand(file)
	if err != nil {
		return nil, err
	}

	buffer, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}

	return PublicKey(buffer)
}

// with public key
func PublicKey(buffer []byte) (ssh.AuthMethod, error) {
	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil, err
	}

	return ssh.PublicKeys(key), nil
}

// with password
func Password(password string) (ssh.AuthMethod, error) {
	return ssh.Password(password), nil
}

// with agent
func Agent() (ssh.AuthMethod, error) {
	sshAgent, err := net.Dial("unix", os.Getenv(SSH_AUTH_SOCK))
	if err != nil {
		return nil, err
	}

	return ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers), nil
}
