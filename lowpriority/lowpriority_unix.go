//go:build unix && !darwin

package lowpriority

import (
	"golang.org/x/sys/unix"
)

func setThreadLowPriority() error {
	return unix.Setpriority(unix.PRIO_PROCESS, unix.Gettid(), 19) //nolint:wrapcheck // It's OK.
}
