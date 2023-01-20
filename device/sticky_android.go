//go:build android

package device

import (
	"golang.zx2c4.com/wireguard/conn"
	"golang.zx2c4.com/wireguard/rwcancel"
)

//
// As of Android 11 raw sockets are no more available:
// see: https://developer.android.com/training/articles/user-data-ids#mac-11-plus
//
func (device *Device) startRouteListener(bind conn.Bind) (*rwcancel.RWCancel, error) {
	return nil, nil
}
