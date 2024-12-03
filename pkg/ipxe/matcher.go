package ipxe

import (
	"errors"
	"fmt"
	"net"
)

var (
	ErrInvalidIP        = errors.New("invalid IP")
	ErrInvalidIPRange   = errors.New("invalid IP range")
	ErrIPv6NotSupported = errors.New("IPv6 not supported")

	DefaultMatchers = []MatcherBuilder{
		MatchIPRange,
		MatchStaticIP,
	}
)

type MatcherBuilder func(string) (Matcher, error)
type Matcher func(ip string) bool

func MatchStaticIP(ip string) (Matcher, error) {
	i := net.ParseIP(ip)
	if i == nil {
		return nil, ErrInvalidIP
	} else if i.To4() == nil {
		return nil, ErrIPv6NotSupported
	}

	return func(ip string) bool {
		return ip == i.String()
	}, nil
}

func MatchIPRange(ipRange string) (Matcher, error) {
	_, ipnet, err := net.ParseCIDR(ipRange)
	if err != nil {
		return nil, err
	}

	return func(ip string) bool {
		i := net.ParseIP(ip)
		if i == nil {
			return false
		} else if i.To4() == nil {
			return false
		}

		return ipnet.Contains(i)
	}, nil
}

func DetectMatcher(pattern string, builders []MatcherBuilder) (Matcher, error) {
	for _, b := range builders {
		m, err := b(pattern)
		if err == nil {
			return m, nil
		}
	}

	return nil, fmt.Errorf("no matcher found for %s", pattern)
}
