package ipxe

import (
	"errors"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"toni.systems/goisoboot/pkg/config"
)

var (
	ipxeTemplate = `#!ipxe
kernel http://{{ .serverIP }}:{{ .serverPort }}/linux/{{ .name }}/vmlinuz rd.neednet=1 rd.live.debug=1 ip=dhcp  root=live:http://{{ .serverIP }}:{{ .serverPort }}/linux/{{ .name }}/squashfs initrd=initrd rootfstype=squashfs {{ .extraKernelArgs }}
initrd http://{{ .serverIP }}:{{ .serverPort }}/linux/{{ .name }}/initrd
boot`
	defaultKernelArgs = map[string]string{
		"rd.fstab":            "0",
		"rd.luks":             "0",
		"rd.lvm":              "0",
		"rd.md":               "0",
		"rd.net.timeout.dhcp": "10",
	}
	ErrIPNotAllowed = errors.New("IP not allowed")
)

func (s *server) IPXE(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("ipxe").Parse(ipxeTemplate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	vars, err := s.getIPXEVars(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	err = tmpl.Execute(w, vars)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *server) Kernel(w http.ResponseWriter, r *http.Request) {
	filename := filepath.Join(s.rootDir, "linux", mux.Vars(r)["name"], "vmlinuz")
	log.Infof("Serving kernel: %s", filename)
	http.ServeFile(w, r, filename)
}

func (s *server) Initrd(w http.ResponseWriter, r *http.Request) {
	filename := filepath.Join(s.rootDir, "linux", mux.Vars(r)["name"], "initrd")
	log.Infof("Serving initrd: %s", filename)
	http.ServeFile(w, r, filename)
}

func (s *server) Squashfs(w http.ResponseWriter, r *http.Request) {
	filename := filepath.Join(s.rootDir, "linux", mux.Vars(r)["name"], "squashfs")
	log.Infof("Serving squashfs: %s", filename)
	http.ServeFile(w, r, filename)
}

func (s *server) getIPXEVars(r *http.Request) (map[string]string, error) {
	ipxeConfig := s.getIPXEByIP(r)
	if ipxeConfig == nil {
		return nil, ErrIPNotAllowed
	}

	kernelArgs := map[string]string{}
	for k, v := range defaultKernelArgs {
		kernelArgs[k] = v
	}
	for k, v := range ipxeConfig.KernelArgs {
		kernelArgs[k] = v
	}

	extraKernelArgsList := make([]string, 0)
	for k, v := range kernelArgs {
		extraKernelArgsList = append(extraKernelArgsList, k+"="+v)
	}

	vars := make(map[string]string)

	vars["ip"] = r.RemoteAddr
	vars["serverIP"] = s.ip
	vars["serverPort"] = strconv.Itoa(s.port)
	vars["name"] = ipxeConfig.Name
	vars["extraKernelArgs"] = strings.Join(extraKernelArgsList, " ")

	return vars, nil
}

func (s *server) getIPXEByIP(r *http.Request) *config.IPXE {
	requestIP := r.RemoteAddr

	for _, i := range s.ipxe {
		found := false
		for _, ip := range i.IPs {
			if ip == requestIP {
				found = true
				break
			}
		}
		if found {
			return &i
		}
	}

	return nil
}
