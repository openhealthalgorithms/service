package service

import (
	"context"
	"encoding/json"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/openhealthalgorithms/service/pkg/algorithms/hearts"
	"github.com/openhealthalgorithms/service/pkg/tools"
	"github.com/openhealthalgorithms/service/pkg/types"
)

// errorResponse used to respond with errors.
type errorResponse struct {
	Error string `json:"error"`
}

type versionResponse struct {
	Version string `json:"version"`
}

type defaultResponse struct {
	Message string `json:"message"`
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
	router.HandleFunc("/api/version", versionRequestHandler)
	router.HandleFunc("/", defaultHandler)

	return Service{Port: port, Addr: addr, Router: router}
}

func (s *Service) StartHttpServer() {
	addr := tools.JoinStringsSep(":", s.Addr, s.Port)
	srv := &http.Server{
		Addr:           addr,
		Handler:        s.Router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Fatalln(srv.ListenAndServe())
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	result := &defaultResponse{Message: "welcome to open health algorithm service"}

	respondSuccess(w, result)
}

func versionRequestHandler(w http.ResponseWriter, r *http.Request) {
	result := &versionResponse{Version: "0.1"}

	respondSuccess(w, result)
}

func algorithmRequestHandler(w http.ResponseWriter, r *http.Request) {
	if strings.ToLower(r.Method) != "post" {
		respondError(w, errors.New("invalid method, only accepts post request"), http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()

	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondError(w, err, http.StatusInternalServerError)
		return
	}

	paramObj, err := tools.ParseParams(content)
	if err != nil {
		respondError(w, err, http.StatusUnprocessableEntity)
		return
	}

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
		respondError(w, err, http.StatusForbidden)
		return
	}

	algorithmOut["request_id"] = uuid.NewRandom().String()
	algorithmOut["hearts"] = algorithmOut["Hearts"]
	delete(algorithmOut, "Hearts")

	result := &algorithmOut

	respondSuccess(w, result)
}

// Helper functions for building responses.
func respondSuccess(w http.ResponseWriter, data interface{}) {
	respJSON, err := json.Marshal(data)
	if err != nil {
		respondError(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respJSON)
}

func respondError(w http.ResponseWriter, err error, code int) {
	resp := errorResponse{Error: err.Error()}
	respJSON, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(respJSON)
}
