package app

const (
	DEFAULT_BECOME_METHOD     = "sudo"
	DEFAULT_BECOME_EXECUTABLE = "/bin/bash"
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
