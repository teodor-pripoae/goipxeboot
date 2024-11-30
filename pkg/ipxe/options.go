package ipxe

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
