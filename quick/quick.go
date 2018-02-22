package main

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

	r.HandleFunc("/workflowStatuses.json", func(resp http.ResponseWriter, req *http.Request) {
		d := getWorkflowStatuses(cli)
		enc := json.NewEncoder(resp)
		enc.SetIndent("", "  ")
		err := enc.Encode(d)
		if err != nil {
			panic(err)
		}

	})

	r.HandleFunc("/submit", func(resp http.ResponseWriter, req *http.Request) {
		f, h, err := req.FormFile("workflow")
		if err != nil {
			log.Error("submit", err)
		}
		if f == nil {
			return
		}

		fi, hi, err := req.FormFile("inputs")
		if err != nil {
			log.Error("submit", err)
		}
		if fi == nil {
			return
		}

		b, _ := ioutil.ReadAll(f)
		bi, _ := ioutil.ReadAll(fi)

		type submit struct {
			App    string                 `json:"app"`
			Inputs map[string]interface{} `json:"inputs"`
		}
		s := submit{}

		enc := base64.StdEncoding.EncodeToString(b)
		hash := fmt.Sprintf("%x", md5.Sum(b))
		prefix := "data:text/plain;base64,"

		if err := json.Unmarshal(bi, &s.Inputs); err != nil {
			log.Error("submit unmarshal", err)
			return
		}
		s.App = prefix + enc

		bout, _ := json.Marshal(s)

		buf := bytes.NewBuffer(bout)
		presp, err := http.Post("http://localhost:8081/v0/engine/jobs/", "application/json", buf)
		if err != nil {
			log.Error("post err", err)
		}
		pb, _ := ioutil.ReadAll(presp.Body)
		log.Info("post resp", string(pb))

		type response struct {
			ID     string
			RootID string
		}
		bunnyResp := response{}
		json.Unmarshal(pb, &bunnyResp)

		wf := &graph.Workflow{ID: hash}
		run := &graph.Run{ID: bunnyResp.ID}

		cli.AddVertex(wf)
		cli.AddVertex(run)
		cli.AddEdge(graph.RunForWorkflow(run, wf))

		fmt.Println(string(bout), h.Filename, hash, hi.Filename, bunnyResp)
	})

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

	/*
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
	*/

	// Root web application
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("build/web"))))

	// TODO this is far too general. doesn't handle 404s.
	r.PathPrefix("/").HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		http.ServeFile(resp, req, "build/web/index.html")
	})

	srv := http.Server{
		Handler:      r,
		Addr:         ":9653",
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}

	log.Info("listening", "http://localhost:9653")
	srv.ListenAndServe()

	//log.Info("listening", "https://localhost:9653")
	//srv.ListenAndServeTLS("cert.pem", "key.pem")
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
