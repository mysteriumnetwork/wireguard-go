// +build !linux android

package device

import "errors"

// taken from android wireguard integration - protecting socket with android tooling
func BindToSocketFd(bind Bind) (int, error) {
	native, ok := bind.(*NativeBind)
	if !ok {
		return -1, errors.New("cannot cast to NativeBind(default)")
	}

	conn, err := native.ipv4.SyscallConn()
	if err != nil {
		return -1, err
	}

	var fd int
	err = conn.Control(func(f uintptr) {
		fd = int(f)
	})
	if err != nil {
		return -1, err
	}
	return fd, nil
}
