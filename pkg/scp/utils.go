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
	"io"
)

const (
	BUFFER_SIZE = 2 * 1024 * 1024
)

var (
	ErrScpWriteAck = errors.New("scp: Failed to write ack buffer")
	ErrScpFailed   = errors.New("scp: copy file failed")
)

func ack(writer io.Writer) error {
	var msg = []byte{0, 0}
	n, err := writer.Write(msg)
	if err != nil {
		return err
	}

	if n < len(msg) {
		return ErrScpWriteAck
	}

	return nil
}

func copyN(writer io.Writer, src io.Reader, size int64) (int64, error) {
	reader := io.LimitReader(src, size)
	var total int64

	for total < size {
		n, err := io.CopyBuffer(writer, reader, make([]byte, BUFFER_SIZE))
		if err != nil {
			return 0, err
		}

		total += n
	}

	return total, nil
}
