package device

import (
	"net"
	"time"

	"golang.zx2c4.com/wireguard/tun"
)

// The purpose of this DeviceApi is to wrap actual device and export various methods which don't exist or
// queries uses unexported fields of original device.
//
// This way we can introduce our own code with minimal modification of original fork and keep changes as localized as
// possible.
//
// Implementations are copied from uapi ipc handling functions.

type DeviceApi struct {
	device *Device
}

func UserspaceDeviceApi(tun tun.Device) *DeviceApi {

	log := NewLogger(LogLevelDebug, "[userspace-wg]")
	return &DeviceApi{
		device: NewDevice(tun, log),
	}
}

func (expDev *DeviceApi) SetListeningPort(port uint16) error {
	dev := expDev.device
	dev.net.Lock()
	dev.net.port = uint16(port)
	dev.net.Unlock()
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
	AllowedIPs      []string
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
	internalPeer.Lock()
	defer internalPeer.Unlock()
	internalPeer.endpoint = peer.RemoteEndpoint
	internalPeer.persistentKeepaliveInterval = peer.KeepAlivePeriod
	for _, cidr := range peer.AllowedIPs {
		_, network, err := net.ParseCIDR(cidr)
		if err != nil {
			return err
		}
		ones, _ := network.Mask.Size()
		dev.allowedips.Insert(network.IP, uint(ones), internalPeer)
	}

	return nil
}

func (expDev *DeviceApi) RemovePeer(publicKey NoisePublicKey) {
	expDev.device.RemovePeer(publicKey)
}

func (expDev *DeviceApi) Peers() (externalPeers []ExternalPeer, err error) {
	dev := expDev.device
	dev.peers.RLock()
	defer dev.peers.RUnlock()
	for _, peer := range dev.peers.keyMap {
		peer.RLock()
		allowedIPs := dev.allowedips.EntriesForPeer(peer)
		externalPeers = append(externalPeers, peerToExternal(peer, allowedIPs))
		peer.RUnlock()
	}

	return externalPeers, err
}

func (expDev *DeviceApi) Wait() {
	<-expDev.device.Wait()
}

func (expDev *DeviceApi) Close() {
	expDev.device.Close()
}

func (expDev *DeviceApi) Boot() {
	expDev.device.Up()
}

func (expDev *DeviceApi) GetNetworkSocket() (int, error) {
	return BindToSocketFd(expDev.device.net.bind)
}

func peerToExternal(peer *Peer, allowedIPs []net.IPNet) ExternalPeer {

	var subnets []string
	for _, allowedIP := range allowedIPs {
		subnets = append(subnets, allowedIP.String())
	}

	return ExternalPeer{
		PublicKey:       peer.handshake.remoteStatic,
		RemoteEndpoint:  peer.endpoint,
		KeepAlivePeriod: peer.persistentKeepaliveInterval,
		AllowedIPs:      subnets,
		Stats: Statistics{
			Received: peer.stats.rxBytes,
			Sent:     peer.stats.txBytes,
		},
		LastHanshake: int(peer.stats.lastHandshakeNano / time.Second.Nanoseconds()),
	}
}
