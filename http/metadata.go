package http

import (
	"encoding/json"
	"github.com/gorilla/handlers"
	"net/http"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/davepgreene/propsd/sources"
)

type metadataHandler struct{}

func newMetadataHandler() http.Handler {
	return handlers.MethodHandler{
		"GET": &metadataHandler{},
	}
}

func (h *metadataHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	s, err := session.NewSession()
	if err == nil {
		m := sources.NewMetadataSource(*s)
		m.Get()

		w.WriteHeader(http.StatusNotImplemented)
		b, _ := json.Marshal(m.Properties())

		w.Write(b)
	}
}
