package ip_to_location

import (
	"net/netip"

	"github.com/oschwald/geoip2-golang/v2"
)

type IpLocationViaGeoip2 struct {
	Reader *geoip2.Reader
}

func (g *IpLocationViaGeoip2) GetIp2Location(ip string) (string, error) {
	ipAddr, err := netip.ParseAddr(ip)
	if err != nil {
		return "", err
	}
	record, err := g.Reader.City(ipAddr)
	if err != nil {
		return "", err
	}
	if !record.HasData() {
		return "", ErrNoLocationFound
	}
	location := record.City.Names.SimplifiedChinese
	if len(record.Subdivisions) > 0 {
		location += record.Subdivisions[0].Names.SimplifiedChinese
	}
	return location, nil
}
