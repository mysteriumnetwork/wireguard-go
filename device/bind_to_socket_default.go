// +build !linux android

package device

import "errors"

// taken from android wireguard integration - protecting socket with android tooling
func BindToSocketFd(bind Bind) (int32, error) {
	native, ok := bind.(*NativeBind)
	if !ok {
		return -1, errors.New("cannot cast to NativeBind")
	}

	conn, err := native.ipv4.SyscallConn()
	if err != nil {
		return -1, err
	}

	var fd int32
	err = conn.Control(func(f uintptr) {
		fd = int32(f)
	})
	if err != nil {
		return -1, err
	}
	return fd, nil
}
