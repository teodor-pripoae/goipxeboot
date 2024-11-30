package ipxe

import (
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
)

var (
	ipxeTemplate = `#!ipxe
kernel http://{{ .serverIP }}/linux/{{ .name }}/vmlinuz ip=dhcp initrd=initrd
initrd http://{{ .serverIP }}/linux/{{ .name }}/initrd
boot`
)

func (s *server) IPXE(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("ipxe").Parse(ipxeTemplate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	vars := s.getIPXEVars(r)

	err = tmpl.Execute(w, vars)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *server) Kernel(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./linux/"+mux.Vars(r)["name"]+"/vmlinuz")
}

func (s *server) Initrd(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./linux/"+mux.Vars(r)["name"]+"/initrd")
}

func (s *server) Squashfs(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./linux/"+mux.Vars(r)["name"]+"/squashfs")
}

func (s *server) getIPXEVars(r *http.Request) map[string]string {
	name := "example"

	vars := make(map[string]string)

	vars["ip"] = r.RemoteAddr
	vars["serverIP"] = s.ip
	vars["name"] = name

	return vars
}
