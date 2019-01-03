package tun

import (
	"os"

	"github.com/mysteriumnetwork/wireguard-go/rwcancel"
)

func AndroidTunDevice(fd int) (TUNDevice, error) {
	file := os.NewFile(uintptr(fd), "/dev/tun")
	tun := &nativeTun{
		tunFile: file,
		fd:      file.Fd(),
		events:  make(chan TUNEvent, 5),
		errors:  make(chan error, 5),
		nopi:    true,
	}
	var err error
	tun.fdCancel, err = rwcancel.NewRWCancel(fd)
	if err != nil {
		return nil, err
	}

	return tun, nil
}
