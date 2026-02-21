package handlers

import "net/http"

func Health(w http.ResponseWriter, r *http.Request) {
	respondOK(w, map[string]string{"status": "ok"})
}
