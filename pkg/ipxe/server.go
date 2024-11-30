package ipxe

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"toni.systems/goisoboot/pkg/config"
)

var (
	ErrMissingIP   = errors.New("missing IP")
	ErrMissingPort = errors.New("missing port")

	dhcpTemplate = `
next-server {{ .serverIP }};
filename "ipxe.efi";

class "pxeclients" {
	if exists user-class and option user-class = "iPXE" {
	    filename "http://{{ .serverIP }}:{{ .serverPort }}/ipxe";
	} else {
	    filename "ipxe.efi";
	}
}`
)

type Server interface {
	Run() error
	IPXE(w http.ResponseWriter, r *http.Request)
	Kernel(w http.ResponseWriter, r *http.Request)
	Initrd(w http.ResponseWriter, r *http.Request)
	Squashfs(w http.ResponseWriter, r *http.Request)
}

type server struct {
	ip      string
	port    int
	router  *mux.Router
	rootDir string
	ipxe    []config.IPXE
}

func New(options ...Option) (Server, error) {
	s := &server{
		rootDir: "./",
		ipxe:    []config.IPXE{},
	}

	s.router = router(s)

	for _, o := range options {
		o(s)
	}

	if s.ip == "" {
		return nil, ErrMissingIP
	}
	if s.port == 0 {
		return nil, ErrMissingPort
	}

	return s, nil
}

func (s *server) Run() error {
	addr := fmt.Sprintf(":%d", s.port)

	srv := &http.Server{
		Addr:         addr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      s.router,
	}

	log.Infof("Starting HTTP server on %s", addr)
	s.printInfo()

	return srv.ListenAndServe()
}

func (s *server) printInfo() {
	log.Infof("Setup ISC DHCP Server with the following configuration:")
	tmpl, _ := template.New("dhcp").Parse(dhcpTemplate)
	var buff bytes.Buffer
	_ = tmpl.Execute(&buff, map[string]string{
		"serverIP":   s.ip,
		"serverPort": strconv.Itoa(s.port),
	})

	log.Infof(buff.String())
}
