package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/bmeg/arachne/aql"
	"github.com/gorilla/mux"
	"github.com/ohsu-comp-bio/funnel/logger"
	"github.com/ohsu-comp-bio/mortar/graph"
	"github.com/ohsu-comp-bio/tes"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var log = logger.NewLogger("quick", logger.DefaultConfig())

func main() {

	server := "localhost:8202"
	graphID := "quick"

	acli, err := aql.Connect(server, true)
	if err != nil {
		panic(err)
	}

	err = acli.AddGraph(graphID)
	if err != nil {
		panic(err)
	}

	cli := graph.Client{Client: &acli, Graph: graphID}

	b := makeData()
	if err := cli.AddBatch(b); err != nil {
		fmt.Println("ERR", err)
	}

	r := mux.NewRouter()

	// Prometheus metrics
	r.HandleFunc("/metrics", func(resp http.ResponseWriter, req *http.Request) {
		updateMetrics(cli, graphID)
		promhttp.Handler().ServeHTTP(resp, req)
	})

	// JSON data for run/workflow/step/etc status
	r.HandleFunc("/data.json", func(resp http.ResponseWriter, req *http.Request) {
		d := getData(cli, graphID)
		enc := json.NewEncoder(resp)
		enc.SetIndent("", "  ")
		enc.Encode(d)
	})

	// JSON data for run/workflow/step/etc status
	r.HandleFunc("/data/wf/{wfid}", func(resp http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		wfid, ok := vars["wfid"]
		if !ok {
			return
		}
		d := getData2(cli, graphID, wfid)
		enc := json.NewEncoder(resp)
		enc.SetIndent("", "  ")
		enc.Encode(d)
	})

	// Root web application
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("build/web"))))

	r.PathPrefix("/").HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		http.ServeFile(resp, req, "build/web/index.html")
	})

	srv := http.Server{
		Handler:      r,
		Addr:         ":9653",
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	log.Info("listening", "http://localhost:9653")
	srv.ListenAndServe()
}

type ByCreationTime []*tes.Task

func (b ByCreationTime) Len() int {
	return len(b)
}
func (b ByCreationTime) Less(i, j int) bool {
	return b[i].CreationTime < b[j].CreationTime
}
func (b ByCreationTime) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func fmtRow(row []*aql.QueryResult) []string {
	o := []string{}
	for _, item := range row {
		switch el := item.Result.(type) {
		// TODO this type switch is not intuitive. Should be aql.Vertex/Edge
		case *aql.QueryResult_Vertex:
			o = append(o, fmt.Sprintf("V(%s, %s)", el.Vertex.Label, el.Vertex.Gid))
		case *aql.QueryResult_Edge:
			o = append(o, fmt.Sprintf("E(%s, %s)", el.Edge.Label, el.Edge.Gid))
		}
	}
	return o
}
