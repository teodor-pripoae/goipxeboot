package ipxe

import (
	"github.com/gorilla/mux"
)

func router(server Server) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/health", server.Health).Methods("GET")
	r.HandleFunc("/ipxe", server.IPXE).Methods("GET")
	r.HandleFunc("/linux/{name}/vmlinuz", server.Kernel).Methods("GET")
	r.HandleFunc("/linux/{name}/initrd", server.Initrd).Methods("GET")
	r.HandleFunc("/linux/{name}/squashfs", server.Squashfs).Methods("GET")

	return r
}
