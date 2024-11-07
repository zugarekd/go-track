package handlers

import (
	"encoding/json"
	"net/http"
)

type RadonProResponse struct {
	Message string `json:"message"`
}

// RadonProHandler is a simple handler that returns a JSON response
func RadonProHandler(w http.ResponseWriter, r *http.Request) {
	// "GET /log?AID=&GID=&CPM=10&ACPM=11.81&uSV=0.07&pci=0.45&model=RadonPro HTTP/1.1" 404 455 "-" "-"
	aid := r.URL.Query().Get("AID")
	println("AID: ", aid)
	gid := r.URL.Query().Get("GID")
	println("GID: ", gid)
	acpm := r.URL.Query().Get("ACPM")
	println("ACPM: ", acpm)
	usv := r.URL.Query().Get("uSV")
	println("uSV: ", usv)
	model := r.URL.Query().Get("model")
	println("Model: ", model)
	pci := r.URL.Query().Get("pci")
	println("PCI: ", pci)
	cpm := r.URL.Query().Get("CPM")
	println("CPM: ", cpm)

	response := RadonProResponse{Message: "Ok"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	w.WriteHeader(http.StatusOK)
}
