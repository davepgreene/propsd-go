package http

import (
	"github.com/gorilla/handlers"
	"net/http"
	"encoding/json"

	"github.com/davepgreene/propsd/sources"
	"github.com/aws/aws-sdk-go/aws/session"
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
		m.Poll()

		w.WriteHeader(http.StatusNotImplemented)
		b, _ := json.Marshal(m.Properties())

		w.Write(b)
	}
}

