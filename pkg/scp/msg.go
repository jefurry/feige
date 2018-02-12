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

package scp

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
)

const (
	copyMessage  = 'C'
	errorMessage = 0x1
	warnMessage  = 0x2
)

//Message is scp control message
type message struct {
	mtype    byte
	merror   error
	mode     string
	size     int64
	fileName string
}

func (m *message) readByte(reader io.Reader) (byte, error) {
	buff := make([]byte, 1)
	if _, err := reader.Read(buff); err != nil {
		return 0, err
	}

	return buff[0], nil
}

func (m *message) readOpCode(reader io.Reader) error {
	var err error
	m.mtype, err = m.readByte(reader)

	return err
}

func (m *message) readError(reader io.Reader) error {
	msg, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	m.merror = errors.New(strings.TrimSpace(string(msg)))
	return nil
}

func (m *message) readLine(reader io.Reader) (string, error) {
	line := ""
	b, err := m.readByte(reader)
	if err != nil {
		return "", err
	}

	for b != 10 {
		line += string(b)
		b, err = m.readByte(reader)
		if err != nil {
			return "", err
		}
	}

	return line, nil
}

func (m *message) readCopy(reader io.Reader) error {
	line, err := m.readLine(reader)
	if err != nil {
		return err
	}

	parts := strings.Split(line, " ")
	if len(parts) < 2 {
		return errors.New("Invalid copy line: " + line)
	}

	m.mode = parts[0]
	m.size, err = strconv.ParseInt(parts[1], 10, 0)
	if err != nil {
		return err
	}

	m.fileName = parts[2]

	return nil
}

func (m *message) readFrom(reader io.Reader) error {
	if err := m.readOpCode(reader); err != nil {
		return err
	}

	switch m.mtype {
	case copyMessage:
		if err := m.readCopy(reader); err != nil {
			return err
		}
	case errorMessage, warnMessage:
		if err := m.readError(reader); err != nil {
			return err
		}
	default:
		return errors.New(fmt.Sprintf("Unsupported opcode: %v", m.mtype))
	}

	return nil
}

func newMessageFromReader(reader io.Reader) (*message, error) {
	m := new(message)
	if err := m.readFrom(reader); err != nil {
		return nil, err
	}

	return m, nil
}
