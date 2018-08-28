package service

import (
	"context"
	"encoding/json"
	"github.com/openhealthalgorithms/service/pkg/algorithms/hearts"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/openhealthalgorithms/service/pkg/tools"
	"github.com/openhealthalgorithms/service/pkg/types"
)

// errorResponse used to respond with errors.
type errorResponse struct {
	Error string
}

type Service struct {
	Port   string
	Addr   string
	Router *http.ServeMux
}

func NewService() Service {
	return NewServiceWithPort("9595")
}

func NewServiceWithPort(port string) Service {
	return NewServiceWithPortAddress(port, "0.0.0.0")
}

func NewServiceWithPortAddress(port, addr string) Service {
	router := http.NewServeMux()
	router.HandleFunc("/api/algorithm", algorithmRequestHandler)

	return Service{Port: port, Addr: addr, Router: router}
}

func (s *Service) StartHttpServer() {
	addr := tools.JoinStringsSep(":", s.Addr, s.Port)
	srv := &http.Server{Addr: addr, Handler: s.Router}

	log.Fatalln(srv.ListenAndServe())
}

func algorithmRequestHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondError(w, err, http.StatusInternalServerError)
		return
	}

	paramObj := tools.ParseParams(content)

	v := types.NewValuesCtx()
	v.Params.Set("params", paramObj)
	v.Params.Set("guide", "guideline_hearts.json")
	v.Params.Set("guidecontent", "guideline_hearts_content.json")

	ctx := context.WithValue(context.Background(), types.KeyValuesCtx, &v)
	algorithm := hearts.New()
	err = algorithm.Get(ctx)
	if err != nil {
		respondError(w, err, http.StatusNotImplemented)
		return
	}

	algorithmOut, err := algorithm.Output()
	if err != nil {
		respondError(w, err, http.StatusNotImplemented)
		return
	}

	result := &algorithmOut

	respondSuccess(w, result)
}

// Helper functions for building responses.
func respondSuccess(w http.ResponseWriter, data interface{}) {
	respJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		respondError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(respJSON)
}

func respondError(w http.ResponseWriter, err error, code int) {
	resp := errorResponse{Error: err.Error()}
	respJSON, _ := json.Marshal(resp)
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	w.Write(respJSON)
}
