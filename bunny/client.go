package bunny

import (
  "bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"
	"io/ioutil"
)

type Client struct {
  Server string
}

type CreateJobResponse struct {
  ID     string
  RootID string
  Name   string
  Status string
  App    string
  Inputs map[string]interface{}
}

// Job describes a request made to Bunny to run a workflow.
//
// TODO currently this expects Bunny to have access to the filepath given by "App".
//      this is a limitation of how Bunny loads workflow documents, and requires
//      Mortar to know about the workflow files Bunny can access.
type Job struct {
  // App is the filepath of a CWL workflow document which the job will run.
  App    string                 `json:"app"`

  /*
    Inputs is the CWL job inputs map, describing the workflow input files/values.
    e.g.
    {
      "input_file": {
          "class" : "File",
              "path": "/full/path/to/input.txt"
          }
    }
  */
  Inputs map[string]interface{} `json:"inputs"`
}

func (c *Client) CreateJob(job *Job) (*CreateJobResponse, error) {

		b, err := json.Marshal(job)
    if err != nil {
      return nil, fmt.Errorf("failed to marshal job: %s", err)
    }

		buf := bytes.NewBuffer(b)
		resp, err := http.Post(c.Server + "/v0/engine/jobs/", "application/json", buf)
		if err != nil {
      return nil, fmt.Errorf("failed to post job request: %s", err)
		}

		rb, err := ioutil.ReadAll(resp.Body)
    if err != nil {
      return nil, fmt.Errorf("failed to read bunny response body: %s", err)
    }

    if resp.StatusCode != 200 {
      return nil, fmt.Errorf("failed request to create job: %s", string(rb))
    }

		jresp := &CreateJobResponse{}

		err = json.Unmarshal(rb, jresp)
    if err != nil {
      return nil, fmt.Errorf("failed to unmarshal bunny response: %s", err)
    }
    return jresp, nil
}

// hashDoc returns the md5 hexadecimal checksum of the given string.
// used to create a content-based ID of a workflow document.
func hashDoc(s string) string {
  return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}
