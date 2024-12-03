package tftp

type Option func(*server)

func WithRootDir(dir string) Option {
	return func(s *server) {
		s.root = dir
	}
}
