package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/openhealthalgorithms/service/pkg/algorithms/hearts"
	"github.com/openhealthalgorithms/service/pkg/database"
	"github.com/openhealthalgorithms/service/pkg/tools"
	"github.com/openhealthalgorithms/service/pkg/types"
	"github.com/pkg/errors"
)

var (
	dbFile = filepath.Join(tools.GetCurrentDirectory(), "logs.db")

	sqlite *database.SqliteDb
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

// Service object
type Service struct {
	Port   string
	Addr   string
	Router *http.ServeMux
}

// NewService method
func NewService() Service {
	return NewServiceWithPort("9595")
}

// NewServiceWithPort method
func NewServiceWithPort(port string) Service {
	return NewServiceWithPortAddress(port, "0.0.0.0")
}

// NewServiceWithPortAddress method
func NewServiceWithPortAddress(port, addr string) Service {
	router := http.NewServeMux()
	router.HandleFunc("/api/algorithm", algorithmRequestHandler)
	router.HandleFunc("/api/version", versionRequestHandler)
	router.HandleFunc("/", defaultHandler)

	return Service{Port: port, Addr: addr, Router: router}
}

// StartHttpServer method
func (s *Service) StartHttpServer() {
	var err error

	// Check if the DB file exists
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		f, err := os.Create(dbFile)
		if err != nil {
			fmt.Println("Error:", err)
		}
		f.Close()
	}

	addr := tools.JoinStringsSep(":", s.Addr, s.Port)
	srv := &http.Server{
		Addr:           addr,
		Handler:        s.Router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	sqlite, err = database.InitDb(dbFile)
	if err != nil {
		fmt.Printf("Error in DB: %v\n", err)
		os.Exit(1)
	}

	err = sqlite.Migrate()
	if err != nil {
		fmt.Printf("Error in DB: %v\n", err)
		os.Exit(1)
	}
	defer sqlite.Closer()

	log.Fatalln(srv.ListenAndServe())
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	result := &defaultResponse{Message: "welcome to open health algorithm service"}

	respondSuccess(w, result)
}

func versionRequestHandler(w http.ResponseWriter, r *http.Request) {
	result := &versionResponse{Version: "0.4.3"}

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
	v.Params.Set("goal", "goals_hearts.json")
	v.Params.Set("goalcontent", "goals_hearts_content.json")

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

	algorithmOut["hearts"] = algorithmOut["Algorithm"]
	algorithmOut["errors"] = algorithmOut["Errors"]
	delete(algorithmOut, "Algorithm")
	delete(algorithmOut, "Errors")

	result := &algorithmOut

	tx, err := sqlite.DB.Begin()
	if err != nil {
		log.Println(err)
	}
	stmt, err := tx.Prepare("insert into logs(request, response) values(?, ?)")
	if err != nil {
		log.Println(err)
	}
	defer stmt.Close()

	requestObj, _ := json.Marshal(paramObj)
	responseObj, _ := json.Marshal(algorithmOut)
	_, err = stmt.Exec(string(requestObj), string(responseObj))
	if err != nil {
		tx.Rollback()
		log.Println(err)
	}
	tx.Commit()

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
	// enableCors(&w)
	w.WriteHeader(http.StatusOK)
	w.Write(respJSON)
}

func respondError(w http.ResponseWriter, err error, code int) {
	resp := errorResponse{Error: err.Error()}
	respJSON, _ := json.Marshal(resp)

	w.Header().Set("Content-Type", "application/json")
	// enableCors(&w)
	w.WriteHeader(code)
	w.Write(respJSON)
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
}
