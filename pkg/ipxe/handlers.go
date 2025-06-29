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
kernel http://{{ .serverIP }}:{{ .serverPort }}/linux/{{ .name }}/vmlinuz {{ .extraKernelArgs }}
initrd http://{{ .serverIP }}:{{ .serverPort }}/linux/{{ .name }}/initrd
boot`
	defaultKernelArgs = map[string]*string{
		"rd.neednet":          stringPointer("1"),
		"rd.live.debug":       stringPointer("1"),
		"rd.fstab":            stringPointer("0"),
		"rd.luks":             stringPointer("0"),
		"rd.lvm":              stringPointer("0"),
		"rd.md":               stringPointer("0"),
		"rd.net.timeout.dhcp": stringPointer("10"),
		"ip":                  stringPointer("dhcp"),
		"initrd":              stringPointer("initrd"),
		"rootfstype":          stringPointer("squashfs"),
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

func (s *server) Health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK\n")
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

	vars := make(map[string]string)

	vars["ip"] = r.RemoteAddr
	vars["serverIP"] = s.getServerIP(ipxeConfig)
	vars["serverPort"] = strconv.Itoa(s.port)
	vars["name"] = ipxeConfig.Name

	kernelArgs := map[string]*string{}
	if !ipxeConfig.SkipDefaults {
		for k, v := range defaultKernelArgs {
			kernelArgs[k] = v
		}
		kernelArgs["root"] = stringPointer(
			fmt.Sprintf("live:http://%s:%s/linux/%s/squashfs", vars["serverIP"], vars["serverPort"], vars["name"]),
		)
	}

	for k, v := range ipxeConfig.KernelArgs {
		kernelArgs[k] = v
	}

	extraKernelArgsList := make([]string, 0)
	for k, v := range kernelArgs {
		vstr := k
		if v != nil {
			vstr += "=" + *v
		}

		extraKernelArgsList = append(extraKernelArgsList, vstr)
	}

	vars["extraKernelArgs"] = strings.Join(extraKernelArgsList, " ")

	return vars, nil
}

func (s *server) getIPXEByIP(r *http.Request) (config.IPXE, error) {
	requestAddr := r.RemoteAddr
	requestIP := strings.Split(requestAddr, ":")[0]

	for _, i := range s.ipxe {
		for _, p := range i.IPs {
			matcher, ok := s.matchers[p]
			if ok && matcher(requestIP) {
				log.Infof("Found pattern %s matching ip %s, selecting configuration %s", p, requestIP, i.Name)
				return i, nil
			}
		}
	}

	return config.IPXE{}, fmt.Errorf("IP %s not allowed", requestIP)
}

func (s *server) getServerIP(ipxeConfig config.IPXE) string {
	if ipxeConfig.ServerIP != "" {
		return ipxeConfig.ServerIP
	}
	return s.ip
}

func stringPointer(s string) *string {
	ptr := s
	return &ptr
}
