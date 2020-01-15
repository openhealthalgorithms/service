package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/openhealthalgorithms/service/pkg/config"
	"github.com/pkg/errors"

	"github.com/openhealthalgorithms/service/database"
	"github.com/openhealthalgorithms/service/pkg"
	"github.com/openhealthalgorithms/service/pkg/algorithms/hearts"
	"github.com/openhealthalgorithms/service/pkg/tools"
	"github.com/openhealthalgorithms/service/pkg/types"

	// pq for postgresql driver
	_ "github.com/lib/pq"
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
    currentSettings := config.CurrentSettings()

    dbFile = currentSettings.LogFile

    return NewServiceWithPort(currentSettings.Port)
}

// NewServiceWithPort method
func NewServiceWithPort(port string) Service {
    return NewServiceWithPortAddress(port, "0.0.0.0")
}

// NewServiceWithPortAddress method
func NewServiceWithPortAddress(port, addr string) Service {
    router := http.NewServeMux()
    router.HandleFunc("/api/algorithm", algorithmRequestHandler)
    router.HandleFunc("/api/algorithm/", algorithmRequestHandler)
    router.HandleFunc("/api/version", versionRequestHandler)
    router.HandleFunc("/", defaultHandler)

    return Service{Port: port, Addr: addr, Router: router}
}

// StartHTTPServer method
func (s *Service) StartHTTPServer() {
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
    result := &versionResponse{Version: pkg.GetVersion()}

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

    currentSettings := config.CurrentSettings()
    v := types.NewValuesCtx()
    v.Params.Set("params", paramObj)

    if currentSettings.CloudEnable {
        projectName := strings.Replace(r.URL.Path, "/api/algorithm", "", 1)
        if len(projectName) > 1 {
            projectName = projectName[1:]

            // check for authorization
            authorizationToken := r.Header.Get("Authorization")
            if len(authorizationToken) == 0 {
                respondError(w, errors.New("authorization token missing"), http.StatusNotAcceptable)
                return
            }
            if !strings.HasPrefix(authorizationToken, "Bearer ") || len(authorizationToken) != 71 {
                respondError(w, errors.New("invalid token format. should be in the format of 'Bearer YOUR_TOKEN'"), http.StatusUnauthorized)
                return
            }

            bearerToken := strings.TrimPrefix(authorizationToken, "Bearer ")
            if len(bearerToken) != 64 {
                respondError(w, errors.New("invalid token for the api"), http.StatusUnauthorized)
                return
            }

            // check for api token in the database and get the project name
            projectForToken, err := checkAPIToken(bearerToken,
                currentSettings.CloudDBHost,
                currentSettings.CloudDBName,
                currentSettings.CloudDBUser,
                currentSettings.CloudDBPassword,
            )
            if err != nil {
                respondError(w, err, http.StatusUnauthorized)
                return
            }

            if projectForToken != projectName {
                respondError(w, errors.New("invalid token for the project"), http.StatusUnauthorized)
                return
            }
        } else {
            projectName = ""
        }
        v.Params.Set("cloud", "yes")
        v.Params.Set("project", projectName)
        v.Params.Set("bucket", currentSettings.CloudBucket)
        v.Params.Set("configfile", currentSettings.CloudConfigFile)
    } else {
        v.Params.Set("cloud", "no")
        v.Params.Set("guide", currentSettings.GuidelineFile)
        v.Params.Set("guidecontent", currentSettings.GuidelineContentFile)
        v.Params.Set("goal", currentSettings.GoalFile)
        v.Params.Set("goalcontent", currentSettings.GoalContentFile)
    }

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

func checkAPIToken(token, host, dbname, user, password string) (string, error) {
    // psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
    psqlInfo := fmt.Sprintf("postgres://%s:%s@%s:5432/%s?sslmode=disable",
        user,
        password,
        host,
        dbname)
    db, err := sql.Open("postgres", psqlInfo)
    if err != nil {
        return "", err
    }
    defer db.Close()

    sqlStatement := `SELECT projects.project_id AS projectname FROM integrations
LEFT JOIN projects ON (integrations.project_id = projects.id) WHERE integrations.api_key = $1
AND integrations.deleted_at IS null`
    projectName := ""
    err = db.QueryRow(sqlStatement, token).Scan(&projectName)
    if err != nil {
        return "", errors.New("no project found")
    }

    if len(projectName) > 0 {
        return projectName, nil
    }

    return "", errors.New("no project found")
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
