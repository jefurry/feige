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
	"github.com/jefurry/feige/utils"
	"golang.org/x/crypto/ssh"
	"time"
)

const (
	DEFAULT_SSH_PORT    = 22
	DEFAULT_SSH_TIMEOUT = 60 // default timeout for operations (in seconds)
)

const (
	AUTH_METHOD_WITH_PRIVATE_KEY_FILE = iota // auth by private key file
	AUTH_METHOD_WITH_PRIVATE_KEY             // auth by private key
	AUTH_METHOD_WITH_PASSWORD                // auth by password
	AUTH_METHOD_WITH_AGENT                   // auth by agent
)

type (
	ClientConfig struct {
		Host                       string
		Port                       int
		SkipHostChecking           bool
		Username, Password         string
		PrivateKey, PrivateKeyFile string
	}
)

// step by step query auth method
// if unspecified username that get current login username
// PrivateKeyFile if it is not empty
// PrivateKey if it is not empty
// Password if it is not empty
// then try agent
func (cc ClientConfig) QueryAuthMothod() (int, error) {
	if cc.Port <= 0 {
		cc.Port = DEFAULT_SSH_PORT
	}

	cc.Username = utils.FindUsername(cc.Username)

	if cc.PrivateKeyFile != "" {
		return AUTH_METHOD_WITH_PRIVATE_KEY_FILE, nil
	} else if cc.PrivateKey != "" {
		return AUTH_METHOD_WITH_PRIVATE_KEY, nil
	} else if cc.Password != "" {
		return AUTH_METHOD_WITH_PASSWORD, nil
	}

	return AUTH_METHOD_WITH_AGENT, nil
}

// get auth method
func (cc ClientConfig) AuthMethod() ([]ssh.AuthMethod, error) {
	flag, err := cc.QueryAuthMothod()
	if err != nil {
		return nil, err
	}

	authMethod := []ssh.AuthMethod{}

	if flag == AUTH_METHOD_WITH_PRIVATE_KEY_FILE {
		// with private key file
		auth, err := PublicKeyFile(cc.PrivateKeyFile)
		if err != nil {
			return nil, err
		}

		authMethod = append(authMethod, auth)
	} else if flag == AUTH_METHOD_WITH_PRIVATE_KEY {
		// with private key
		auth, err := PublicKey([]byte(cc.PrivateKey))
		if err != nil {
			return nil, err
		}

		authMethod = append(authMethod, auth)
	} else if flag == AUTH_METHOD_WITH_PASSWORD {
		// with password
		auth, err := Password(cc.Password)
		if err != nil {
			return nil, err
		}

		authMethod = append(authMethod, auth)
	} else {
		// with agent
		auth, err := Agent()
		if err != nil {
			return nil, err
		}

		authMethod = append(authMethod, auth)
	}

	return authMethod, nil
}

// make a ssh client config
func (cc ClientConfig) Make(timeout time.Duration, cipherList []string) (*ssh.ClientConfig, error) {
	if timeout <= 0 {
		timeout = DEFAULT_SSH_TIMEOUT
	}
	authMethod, err := cc.AuthMethod()
	if err != nil {
		return nil, err
	}

	var config ssh.Config
	if len(cipherList) == 0 {
		config.Ciphers = []string{
			"aes128-ctr",
			"aes192-ctr",
			"aes256-ctr",
			"aes128-gcm@openssh.com",
			"arcfour256",
			"arcfour128",
			"aes128-cbc",
			"3des-cbc",
			"aes192-cbc",
			"aes256-cbc",
		}
	} else {
		config.Ciphers = cipherList
	}

	clientConfig := &ssh.ClientConfig{
		User:    cc.Username,
		Auth:    authMethod,
		Config:  config,
		Timeout: timeout * time.Second,
	}

	if cc.SkipHostChecking {
		clientConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()
	} else {
		hostKey, err := FixedHostKey(cc.Username, cc.Host)
		if err != nil {
			return nil, err
		}

		clientConfig.HostKeyCallback = ssh.FixedHostKey(hostKey)
	}

	return clientConfig, nil
}
