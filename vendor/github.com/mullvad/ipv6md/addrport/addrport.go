package addrport

import (
	"encoding/binary"
	"errors"
	"net"
	"net/netip"

	"github.com/mullvad/ipv6md"
	"github.com/mullvad/ipv6md/utils"
)

var (
	ErrAddrPortInvalidIP    = errors.New("invalid ip")
	ErrAddrPortInvalidIPLen = errors.New("invalid ip length")
)

// Encode encodes an ipv4:port into an IPv6 address.
func Encode(addrPort string) (net.IP, error) {
	ap, err := netip.ParseAddrPort(addrPort)
	if err != nil {
		return nil, err
	}

	data := [16]byte{ipv6md.IPv6Prefix[0], ipv6md.IPv6Prefix[1]}

	binary.LittleEndian.PutUint16(data[2:4], ipv6md.AddrPort.ToUint16())

	addr := ap.Addr()
	addrBytes := addr.AsSlice()
	copy(data[4:8], addrBytes)

	port := ap.Port()
	binary.LittleEndian.PutUint16(data[8:10], port)

	return net.IP(data[:]), nil
}

// Decode assumes an IPv4 address and port has been encoded within the
// IPv6 address and returns a netip.AddrPort with the information.
func Decode(ip net.IP) (netip.AddrPort, error) {
	var addrPort netip.AddrPort

	if ip == nil {
		return addrPort, ErrAddrPortInvalidIP
	}

	data := []byte(ip.To16())
	if len(data) != 16 {
		return addrPort, ErrAddrPortInvalidIPLen
	}

	addr := utils.ToNetIPAddr(data[4:8])
	port := binary.LittleEndian.Uint16(data[8:10])
	addrPort = netip.AddrPortFrom(addr, port)

	return addrPort, nil
}
