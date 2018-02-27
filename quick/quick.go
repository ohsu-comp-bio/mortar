package quick

import (
  "encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/bmeg/arachne/aql"
	"github.com/bmeg/arachne/protoutil"
	"github.com/gorilla/mux"
	"github.com/ohsu-comp-bio/funnel/logger"
	"github.com/ohsu-comp-bio/mortar/graph"
	"github.com/ohsu-comp-bio/mortar/bunny"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var log = logger.NewLogger("quick", logger.DefaultConfig())


func Serve() {
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
  bunnyCli := bunny.Client{Server: "http://localhost:8081"}

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
      log.Error("error", err)
      http.Error(resp, "failed to encode workflow runs", http.StatusInternalServerError)
      return
		}
	})

	r.HandleFunc("/workflowStatuses.json", func(resp http.ResponseWriter, req *http.Request) {
		d := getWorkflowStatuses(cli)
		enc := json.NewEncoder(resp)
		enc.SetIndent("", "  ")
		err := enc.Encode(d)
		if err != nil {
      log.Error("error", err)
      http.Error(resp, "failed to encode workflow status", http.StatusInternalServerError)
      return
		}
	})

  // TODO redesign REST url structure and use POST /run
  //      also, consider that workflow engines might often call this "job"
	r.HandleFunc("/submit", func(resp http.ResponseWriter, req *http.Request) {

    workflowID := req.FormValue("workflow")
    if workflowID == "" {
      http.Error(resp, "missing required workflow ID", http.StatusBadRequest)
      return
    }

    wfv, err := cli.GetVertex(workflowID)
    if err != nil {
      log.Error("submit error", err)
      http.Error(resp, "error getting workflow", http.StatusInternalServerError)
      return
    }

    data := protoutil.AsMap(wfv.Data)
    doc, ok := data["Doc"].(map[string]interface{})
    if !ok {
      log.Error("workflow load error", err)
      http.Error(resp, "error loading workflow", http.StatusInternalServerError)
      return
    }
    wf := &graph.Workflow{
      ID: wfv.Gid,
      Doc: doc,
    }

		fi, _, err := req.FormFile("inputs")
		if err != nil {
      log.Error("submit error", err)
      http.Error(resp, "error loading inputs file", http.StatusInternalServerError)
			return
		}
		if fi == nil {
      http.Error(resp, "missing required workflow inputs file", http.StatusBadRequest)
			return
		}

		bi, err := ioutil.ReadAll(fi)
		if err != nil {
      log.Error("submit error", err)
      http.Error(resp, "failed to read inputs file", http.StatusInternalServerError)
			return
		}

    docb, err := json.Marshal(doc)
    if err != nil {
      log.Error("marshaing doc to json", err)
      http.Error(resp, "failed to process workflow doc", http.StatusInternalServerError)
      return
    }

    job := &bunny.Job{
      App: encodeDoc(docb),
      Inputs: map[string]interface{}{},
    }

    err = json.Unmarshal(bi, &job.Inputs)
    if err != nil {
      log.Error("submit error", err)
      http.Error(resp, "failed while unmarshaling inputs", http.StatusBadRequest)
      return
    }

    log.Info("submitting", job)

    bresp, err := bunnyCli.CreateJob(job)
    if err != nil {
      log.Error("submit error", err)
      http.Error(resp, "failed calling CreateJob", http.StatusBadRequest)
      return
    }
		run := &graph.Run{
      ID: bresp.ID,
      Inputs: bresp.Inputs,
    }

		cli.AddVertex(run)
		cli.AddEdge(graph.RunForWorkflow(run, wf))
	})

	r.HandleFunc("/workflow/{wfid}", func(resp http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		wfid, ok := vars["wfid"]
		if !ok {
      http.Error(resp, "missing workflow ID", http.StatusBadRequest)
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

func encodeDoc(b []byte) string {
  return "data:text/plain;base64," + base64.StdEncoding.EncodeToString(b)
}
