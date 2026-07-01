package types

import (
	"net"
	"net/netip"
)

func ToPGInet(ip net.IP) *netip.Addr {

	if ip == nil {
		return nil
	}

	addr, ok := netip.AddrFromSlice(ip)

	if !ok {
		return nil
	}

	return &addr
}

func FromPGInet(addr *netip.Addr) net.IP {

	if addr == nil {
		return nil
	}

	return net.ParseIP(addr.String())
}
