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

package sftp

import (
	"io"
	//"github.com/jefurry/feige/app"
	//"github.com/jefurry/feige/pkg/command"
	//"github.com/jefurry/feige/pkg/connection/gossh"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"os"
)

const (
	PACKET_SIZE = 1 << 17
)

type (
	SftpClient struct {
		client *sftp.Client
	}
)

func NewSftpClient(client *ssh.Client, opts ...sftp.ClientOption) (*SftpClient, error) {
	opts = append(opts, sftp.MaxPacket(PACKET_SIZE))
	c, err := sftp.NewClient(client, opts...)
	if err != nil {
		return nil, err
	}

	return &SftpClient{
		client: c,
	}, nil
}

func (c *SftpClient) Close() error {
	return c.client.Close()
}

// upload file
func (c *SftpClient) Put(local, remote string) (int64, error) {
	r, err := os.Open(local)
	if err != nil {
		return 0, err
	}
	defer r.Close()

	w, err := c.client.Create(remote)
	if err != nil {
		return 0, err
	}
	defer w.Close()

	return io.Copy(w, io.LimitReader(r, 1e9))
}

// download file
func (c *SftpClient) Fetch(remote, local string) (int64, error) {
	r, err := c.client.Open(remote)
	if err != nil {
		return 0, err
	}
	defer r.Close()

	w, err := os.Create(local)
	if err != nil {
		return 0, err
	}
	defer w.Close()

	return io.Copy(w, io.LimitReader(r, 1e9))
}
