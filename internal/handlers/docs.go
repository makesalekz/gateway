package handlers

import "net/http"

// DocsHandler serves the OpenAPI specification.
type DocsHandler struct {
	spec []byte
}

func NewDocsHandler(spec []byte) *DocsHandler {
	return &DocsHandler{spec: spec}
}

func (h *DocsHandler) ServeSpec(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/yaml")
	w.WriteHeader(http.StatusOK)
	w.Write(h.spec)
}
