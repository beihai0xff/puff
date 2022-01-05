package types

import (
	"fmt"
	"net"
)

var Peers []peer

type peer struct {
	ip   string
	port int
}

// ValidateIP 检查 ip 地址是否正确，如果不正确，该方法返回 error
func (p *peer) ValidateIP() error {
	address := net.ParseIP(p.ip)
	if address == nil {
		return fmt.Errorf("illegal ip address: %s", p.ip)
	}

	p.ip = address.String()
	return nil
}
