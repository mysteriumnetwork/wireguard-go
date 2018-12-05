// +build !android

package device

import "errors"

// taken from android wireguard integration - protecting socket with android tooling
func BindToSocketFd(bind Bind) (int, error) {
	native, ok := bind.(*nativeBind)
	if !ok {
		return -1, errors.New("cannot cast to NativeBind(Linux)")
	}

	return native.sock4, nil
}
