// +build !linux android

package device

import "golang.org/x/sys/windows"

func setsockoptInt(fd uintptr, fwmarkIoctl, mark int) error {
	return windows.SetsockoptInt(windows.Handle(fd), windows.SOL_SOCKET, fwmarkIoctl, int(mark))
}
