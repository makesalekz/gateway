package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func pathParamInt64(r *http.Request, key string) (int64, error) {
	vars := mux.Vars(r)
	return strconv.ParseInt(vars[key], 10, 64)
}

func queryParamInt32(r *http.Request, key string) int32 {
	v := r.URL.Query().Get(key)
	if v == "" {
		return 0
	}
	n, _ := strconv.ParseInt(v, 10, 32)
	return int32(n)
}

func queryParamInt64(r *http.Request, key string) int64 {
	v := r.URL.Query().Get(key)
	if v == "" {
		return 0
	}
	n, _ := strconv.ParseInt(v, 10, 64)
	return n
}

func queryParam(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}
