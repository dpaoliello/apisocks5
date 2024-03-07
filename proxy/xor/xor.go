package xor

import (
	"io"
	"net"

	"github.com/mullvad/ipv6md/addrportxor"
	"github.com/mullvad/proxy/typ"
)

type xor struct {
	addrPort string
	xorBytes uint16
	xorKey   []byte
}

func New(ip net.IP) (*xor, error) {
	target, err := addrportxor.Decode(ip)
	if err != nil {
		return nil, err
	}

	return &xor{
		addrPort: target.AddrPort.String(),
		xorBytes: target.XORBytes,
		xorKey:   target.XORKey,
	}, nil
}

func (x *xor) Address() string                       { return x.addrPort }
func (x *xor) Type() typ.Type                        { return typ.XOR }
func (x *xor) FromPeer(dst io.Writer, src io.Reader) { x.forward(dst, src) }
func (x *xor) ToPeer(dst io.Writer, src io.Reader)   { x.forward(dst, src) }

func (x *xor) forward(dst io.Writer, src io.Reader) {
	buf := make([]byte, 1024*64)

	for {
		nr, err := src.Read(buf)
		if err != nil || nr <= 0 {
			break
		}

		l := int(x.xorBytes)
		if l > nr || l == 0 {
			l = nr
		}

		for i := 0; i < l; i++ {
			buf[i] ^= x.xorKey[i%len(x.xorKey)]
		}

		nw, err := dst.Write(buf[0:nr])
		if (err != nil) || (nr != nw) {
			break
		}
	}

	tcpConn, ok := dst.(*net.TCPConn)
	if !ok {
		return
	}

	tcpConn.CloseWrite()
}
