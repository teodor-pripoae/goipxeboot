package ipxe

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"toni.systems/goipxeboot/pkg/config"
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
)

func (s *server) IPXE(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("ipxe").Parse(ipxeTemplate)
	if err != nil {
		log.Infof("GET /ipxe 500 %v", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	vars, err := s.getIPXEVars(r)
	if err != nil {
		log.Infof("GET /ipxe 403 %v", err.Error())
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	err = tmpl.Execute(w, vars)
	if err != nil {
		log.Infof("GET /ipxe 500 %v", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Infof("GET /ipxe 200")
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
	ipxeConfig, err := s.getIPXEByIP(r)
	if err != nil {
		return nil, err
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

func (s *server) getIPXEByIP(r *http.Request) (config.IPXE, error) {
	requestAddr := r.RemoteAddr
	requestIP := strings.Split(requestAddr, ":")[0]

	for _, i := range s.ipxe {
		found := false
		for _, ip := range i.IPs {
			if ip == requestIP {
				found = true
				break
			}
		}
		if found {
			return i, nil
		}
	}

	return config.IPXE{}, fmt.Errorf("IP %s not allowed", requestIP)
}
