package tftp

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/pin/tftp/v3"
	log "github.com/sirupsen/logrus"
)

const (
	ipxeFilename = "ipxe.efi"
	listenAddr   = ":69"
)

var (
	timeout             = 5 * time.Second
	ErrInvalidFilename  = errors.New("invalid filename, expected ipxe.efi")
	ErrFailedToOpenFile = errors.New("failed to open file")
	ErrFailedToReadFile = errors.New("failed to read file")
)

type Server interface {
	Run() error
}

type server struct {
	tftp     *tftp.Server
	filename string
	root     string
}

func New(options ...Option) Server {
	s := &server{
		filename: "./ipxe/ipxe.efi",
	}

	t := tftp.NewServer(s.ReadHandler, nil)
	t.SetTimeout(timeout)
	s.tftp = t

	for _, option := range options {
		option(s)
	}

	return s
}

func (s *server) Run() error {
	log.Infof("Starting TFTP server on %s", listenAddr)
	return s.tftp.ListenAndServe(listenAddr) // blocks until s.Shutdown() is called
}

func (s *server) ReadHandler(filename string, rf io.ReaderFrom) error {
	log.Infof("Getting file %s", filename)
	if filename != ipxeFilename {
		log.Warnf("Filename %s not allowed", filename)
		return ErrInvalidFilename
	}

	path := filepath.Join(s.root, s.filename)
	log.Infof("Opening file %s", path)

	file, err := os.Open(path)
	if err != nil {
		return errors.Join(err, ErrFailedToOpenFile)
	}

	n, err := rf.ReadFrom(file)
	if err != nil {
		return errors.Join(err, ErrFailedToReadFile)
	}

	log.Infof("File %s %d bytes sent\n", filename, n)
	return nil
}
