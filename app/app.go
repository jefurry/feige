package app

import (
	"strings"
)

const (
	NAME      = "Feige"
	LONG_NAME = NAME + " IT Automation Platform"
	SITE      = "feige.io"
	VERSION   = "0.0.1"
)

var (
	NameForLower = strings.ToLower(NAME)
	NameForUpper = strings.ToUpper(NAME)
)
