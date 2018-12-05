package device

import (
	"time"

	"git.zx2c4.com/wireguard-go/tun"
)

// The purpose of this DeviceApi is to wrap actual device and export variuos method which don't exist or
// queries uses unexported fields of original device. This way we can introduce our own code with minimal modification
// of original fork and keep changes as localized as possible. Implementations are copied from uapi ipc handling functions.

type DeviceApi struct {
	device *Device
}

func UserspaceDeviceApi(tun tun.TUNDevice) *DeviceApi {

	log := NewLogger(LogLevelDebug, "[userspace-wg]")
	return &DeviceApi{
		device: NewDevice(tun, log),
	}
}

func (expDev *DeviceApi) SetListeningPort(port uint16) error {
	dev := expDev.device
	dev.net.mutex.Lock()
	dev.net.port = uint16(port)
	dev.net.mutex.Unlock()
	return dev.BindUpdate()
}

func (expDev *DeviceApi) SetPrivateKey(key NoisePrivateKey) error {
	return expDev.device.SetPrivateKey(key)
}

// Peer related stuff, exposed as "own" structure to avoid coupling
// Used for configuration (runtime populated fields are ignored)

type Statistics struct {
	Received uint64
	Sent     uint64
}

type ExternalPeer struct {
	//configurable fields
	PublicKey       NoisePublicKey
	RemoteEndpoint  Endpoint
	KeepAlivePeriod uint16 //seconds

	//readable fields
	Stats        Statistics
	LastHanshake int //seconds
}

func (expDev *DeviceApi) AddPeer(peer ExternalPeer) error {
	dev := expDev.device
	internalPeer, err := dev.NewPeer(peer.PublicKey)
	if err != nil {
		return err
	}
	internalPeer.mutex.Lock()
	defer internalPeer.mutex.Unlock()
	internalPeer.endpoint = peer.RemoteEndpoint
	internalPeer.persistentKeepaliveInterval = peer.KeepAlivePeriod
	return nil
}

func (expDev *DeviceApi) Peers() (externalPeers []ExternalPeer, err error) {
	dev := expDev.device
	dev.peers.mutex.RLock()
	defer dev.peers.mutex.Unlock()
	for _, peer := range dev.peers.keyMap {
		peer.mutex.RLock()
		externalPeers = append(externalPeers, peerToExternal(peer))
		peer.mutex.Unlock()
	}

	return externalPeers, err
}

func (expDev *DeviceApi) Wait() {
	<-expDev.device.Wait()
}

func (expDev *DeviceApi) Close() {
	expDev.device.Close()
}

func (expDev *DeviceApi) GetNetworkSocket() (int32, error) {
	return BindToSocketFd(expDev.device.net.bind)
}

func peerToExternal(peer *Peer) ExternalPeer {

	return ExternalPeer{
		PublicKey:       peer.handshake.remoteStatic,
		RemoteEndpoint:  peer.endpoint,
		KeepAlivePeriod: peer.persistentKeepaliveInterval,
		Stats: Statistics{
			Received: peer.stats.rxBytes,
			Sent:     peer.stats.txBytes,
		},
		LastHanshake: int(peer.stats.lastHandshakeNano / time.Second.Nanoseconds()),
	}
}
