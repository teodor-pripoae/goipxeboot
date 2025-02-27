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
	"toni.systems/goipxeboot/pkg/config"
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
	Health(w http.ResponseWriter, r *http.Request)
}

type server struct {
	ip       string
	port     int
	router   *mux.Router
	rootDir  string
	ipxe     []config.IPXE
	matchers map[string]Matcher
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

	matchers, err := loadMatchers(s.ipxe)
	if err != nil {
		return nil, err
	}

	s.matchers = matchers

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

	log.Info(buff.String())
}

func loadMatchers(ipxe []config.IPXE) (map[string]Matcher, error) {
	patternsMap := make(map[string]Matcher)

	for _, i := range ipxe {
		for _, p := range i.IPs {
			if _, ok := patternsMap[p]; ok {
				continue
			}

			matcher, err := DetectMatcher(p, DefaultMatchers)
			if err != nil {
				return nil, fmt.Errorf("failed to detect matcher for %s: %w", p, err)
			}

			patternsMap[p] = matcher
		}
	}

	return patternsMap, nil
}
