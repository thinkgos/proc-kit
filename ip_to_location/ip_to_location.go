package ip_to_location

import (
	"errors"

	"github.com/thinkgos/proc/cache"
	"golang.org/x/sync/singleflight"
)

var ErrNoLocationFound = errors.New("no location found for the IP")

type IpLocationAble interface {
	GetIp2Location(ip string) (string, error)
}

var _ IpLocationAble = (*IpLocation)(nil)

type IpLocation struct {
	LocalCache *cache.Cache
	Ils        []IpLocationAble
	group      singleflight.Group
}

func (g *IpLocation) GetIp2Location(ip string) (string, error) {
	location, ok := g.LocalCache.Get(ip)
	if ok {
		return location.(string), nil
	}
	location, err, _ := g.group.Do(ip, func() (any, error) {
		for _, v := range g.Ils {
			localtion, err := v.GetIp2Location(ip)
			if err == nil && localtion != "" {
				g.LocalCache.SetDefault(ip, localtion)
				return localtion, nil
			}
		}
		return "", ErrNoLocationFound
	})
	if err != nil {
		return "", nil
	}
	return location.(string), nil
}
