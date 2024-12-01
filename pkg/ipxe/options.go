package ipxe

import "toni.systems/goipxeboot/pkg/config"

type Option func(*server)

func WithPort(port int) Option {
	return func(s *server) {
		s.port = port
	}
}

func WithIP(ip string) Option {
	return func(s *server) {
		s.ip = ip
	}
}

func WithRootDir(dir string) Option {
	return func(s *server) {
		s.rootDir = dir
	}
}

func WithIPXE(ipxe []config.IPXE) Option {
	return func(s *server) {
		s.ipxe = ipxe
	}
}
