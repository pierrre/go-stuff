//go:build windows

package lowpriority

import (
	"golang.org/x/sys/windows"
)

var procSetThreadPriority = windows.NewLazySystemDLL("kernel32.dll").NewProc("SetThreadPriority")

func setThreadLowPriority() error {
	const threadPriorityIdle = -15
	priority := int32(threadPriorityIdle)
	r1, _, err := procSetThreadPriority.Call(uintptr(windows.CurrentThread()), uintptr(priority))
	if r1 == 0 {
		return err //nolint:wrapcheck // Not needed.
	}
	return nil
}
