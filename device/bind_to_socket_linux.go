// +build linux,!android

package device

import "errors"

// taken from android wireguard integration - protecting socket with android tooling
func BindToSocketFd(bind Bind) (int32, error) {
	native, ok := bind.(*NativeBind)
	if !ok {
		return -1, errors.New("cannot cast to NativeBind")
	}

	return native.sock4
}
