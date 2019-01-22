// +build !linux,!windows android

package device

import "golang.org/x/sys/unix"

func setsockoptInt(fd uintptr, fwmarkIoctl, mark int) error {
	return unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, fwmarkIoctl, int(mark))
}
