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

package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/jefurry/feige/app"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func RandKey() string {
	b := RandBytes(32)

	return fmt.Sprintf("%s@%s-V%s-%s-%d", app.NAME, app.SITE, app.VERSION, b, time.Now().UnixNano())
}

func RandBytes(n int) []byte {
	chars := []byte{
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u',
		'v', 'w', 'x', 'y', 'z', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P',
		'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '_',
	}

	b := make([]byte, n)
	for i, l := 0, len(chars); i < n; i++ {
		rand.Seed(time.Now().UnixNano())
		b[i] = chars[rand.Intn(l)]
	}

	return b
}

// genrate md5
func GenMd5(name string) string {
	hasher := md5.New()
	hasher.Write([]byte(name))

	return hex.EncodeToString(hasher.Sum(nil))
}

func GenHashCode(str string) uint {
	var hc uint = 0
	if str == "" {
		return hc
	}

	var n uint = uint(len(str))
	var i uint
	for i = 0; i < n; i++ {
		hc = hc ^ (uint(str[i]) << (uint(i) & 0xFF))
	}

	return hc
}

// use goid as name if name is empty
func GenName(name string) string {
	if name == "" {
		rs := string(RandBytes(32))
		pid := int64(os.Getpid())
		wd := rs + "-" + strconv.FormatInt(pid, 10)

		goid, err := GoId()
		if err == nil {
			wd = wd + "-" + strconv.FormatInt(int64(goid), 10)
		}

		return wd
	}

	return name
}
