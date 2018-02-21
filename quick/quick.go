package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/alexandrevicenzi/go-sse"
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
	graphID := "mortar"

	acli, err := aql.Connect(server, true)
	if err != nil {
		panic(err)
	}

	err = acli.AddGraph(graphID)
	if err != nil {
		panic(err)
	}

	cli := &graph.Client{Client: &acli, Graph: graphID}

	/*
		b := makeData()
		if err := cli.AddBatch(b); err != nil {
			fmt.Println("ERR", err)
		}
	*/

	r := mux.NewRouter()

	// Prometheus metrics
	r.HandleFunc("/metrics", func(resp http.ResponseWriter, req *http.Request) {
		updateMetrics(cli)
		promhttp.Handler().ServeHTTP(resp, req)
	})

	// JSON data for run/workflow/step/etc status
	r.HandleFunc("/workflowRuns.json", func(resp http.ResponseWriter, req *http.Request) {
		d := getWorkflowRuns(cli)
		enc := json.NewEncoder(resp)
		enc.SetIndent("", "  ")
		err := enc.Encode(d)
		if err != nil {
			panic(err)
		}

	})

	sses := sse.NewServer(nil)
	defer sses.Shutdown()

	go func() {
		for {
			d := getWorkflowRuns(cli)

			b, err := json.Marshal(d)
			if err != nil {
				log.Error("publishing", err)
				continue
			}

			sses.SendMessage("/sub/workflowRuns.json", sse.SimpleMessage(string(b)))
			time.Sleep(5 * time.Second)
		}
	}()

	r.Handle("/sub/workflowRuns.json", sses)

	r.HandleFunc("/workflow/{wfid}", func(resp http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		wfid, ok := vars["wfid"]
		if !ok {
			return
		}
		d := getWorkflowInfo(cli, wfid)
		enc := json.NewEncoder(resp)
		enc.SetIndent("", "  ")
		enc.Encode(d)
	})

	r.HandleFunc("/runsByStep/{wfid}", func(resp http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		wfid, ok := vars["wfid"]
		if !ok {
			return
		}
		d := getRunsByStep(cli, wfid)
		enc := json.NewEncoder(resp)
		enc.SetIndent("", "  ")
		enc.Encode(d)
	})

	// Root web application
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("build/web"))))

	// TODO this is far too general. doesn't handle 404s.
	r.PathPrefix("/").HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		http.ServeFile(resp, req, "build/web/index.html")
	})

	srv := http.Server{
		Handler:      r,
		Addr:         ":9653",
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	log.Info("listening", "https://localhost:9653")
	//srv.ListenAndServe()

	srv.ListenAndServeTLS("cert.pem", "key.pem")
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
