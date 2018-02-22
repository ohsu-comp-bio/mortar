import React, { Component } from "react"
import ReactDOM from "react-dom"
import { BrowserRouter as Router, Route, Link } from 'react-router-dom'

var stateMap = {
  "COMPLETE": "Complete",
  "RUNNING": "Running",
  "INITIALIZING": "Initializing",
  "SYSTEM_ERROR": "Error",
  "EXECUTOR_ERROR": "Error",
  "CANCELED": "Canceled",
}

// Given a run and a step ID, return the state of that step in the run.
function stepState(run, sid) {
  var step = run.Steps[sid]
  var tasks = run.StepTasks[sid]
  var state = "Not Started"

  if (tasks && tasks.length > 0) {
    var latestID = tasks[0]
    var task = run.Tasks[latestID]
    state = stateMap[task.state]

    if (!state) {
      state = "Unknown"
    }
  }
  return state
}


class Home extends Component {
  constructor(props) {
    super(props)
    this.state = {
      "data": null,
    }
  }

  fetchData() {
    console.log("fetch")
    fetch("/workflowRuns.json")
      .then(resp => resp.json())
      .then(data => this.setState({"data": data}))
  }

  componentDidMount() {
    this.fetchData()
  }

  componentWillUnmount() {
    clearInterval(this.fetchInterval)
  }


  render() {
    if (!this.state.data) {
      return <div></div>
    }

    var data = this.state.data
    var headers = data.Columns.map(col => (<td>{col}</td>))
    var rows = []

    for (var i = 0; i < data.Rows.length; i++) {
      var row = data.Rows[i]

      var cells = data.Columns.map(col => {
        return Cell(data.Cells[row + "-" + col])
      })

      rows.push(<tr key={"row-" + row}>
        <td><Link to={"/workflow/" + row}>{row}</Link></td>
        { cells }
      </tr>)
    }

    return (<div>
      <h3><Link to="/">Mortar</Link></h3>
      <table>
        <thead>
          <tr><td></td>{headers}</tr>
        </thead>
        <tbody>{rows}</tbody>
      </table>
    </div>)
  }
}

const Time = {
  Second: 1000,
}

class Run extends Component {
  constructor(props) {
    super(props)
    this.state = {
      "data": null,
    }
  }

  fetchData() {
    fetch("/data.json")
      .then(resp => resp.json())
      .then(data => this.setState({"data": data}))
  }

  componentDidMount() {
    this.fetchData()
    this.fetchInterval = setInterval(this.fetchData, 5 * Time.Second)
  }

  componentWillUnmount() {
    clearInterval(this.fetchInterval)
  }

  render() {
    if (!this.state.data) {
      return <div>Loading</div>
    }

    var rid = this.props.match.params.rid
    if (!rid) {
      return <div>Run not found</div>
    }

    var data = this.state.data
    var run = data.Runs[rid]
    if (!run) {
      return <div>Run not found</div>
    }

    var stepIDs = Object.keys(run.Steps).sort()
    var rows = []

    for (var i = 0; i < stepIDs.length; i++) {
      var sid = stepIDs[i]
      var step = run.Steps[sid]
      var state = stepState(run, sid)

      rows.push(<tr key={"step-tr-" + sid}>
        <td>{step.gid}</td>
        <td>{state}</td>
      </tr>)
    }
    console.log(run)

    return (<div>
      <h3><Link to="/">Mortar</Link></h3>
      <table>
        <thead><tr>
          <th>Step</th>
          <th>Status</th>
        </tr></thead>
        <tbody>{rows}</tbody>
      </table>
    </div>)
  }
}

class RunsForWorkflow extends Component {
  constructor(props) {
    super(props)
    this.state = {
      "data": null,
    }
  }

  componentDidMount() {

    var wfid = this.props.match.params.wfid
    if (!wfid) {
      return
    }

    fetch("/data/wf/"+wfid)
      .then(resp => resp.json())
      .then(data => this.setState({"data": data}))
  }

  render() {
    if (!this.state.data) {
      return <div>Loading</div>
    }

    var data = this.state.data
    var header = [<th key="empty"></th>]
    var rows = []

    for (var i = 0; i < data.StepIDs.length; i++) {
      var sid = data.StepIDs[i]
      var steps = data.Steps[sid]
      header.push(<th key={"step-th-" + sid}>{sid}</th>)
    }

    for (var i = 0; i < data.RunIDs.length; i++) {
      var rid = data.RunIDs[i]
      var run = data.Runs[rid]
      var row = []

      for (var j = 0; j < data.StepIDs.length; j++) {
        var sid = data.StepIDs[j]
        var steps = data.Steps[sid]
        var state = stepState(run, sid)
        row.push(<td key={"run-" + rid + "-step-" + sid}>{state}</td>)
      }

      rows.push(<tr key={"run-" + rid}>
        <td>{rid}</td>
        {row}
      </tr>)
    }
    console.log(run)

    return (<div>
      <h3><Link to="/">Mortar</Link></h3>
      <table>
        <thead><tr>
          {header}
        </tr></thead>
        <tbody>{rows}</tbody>
      </table>
    </div>)
  }
}
const Cell = (run) => {
  var cn = ""

  if (run.Total == run.Complete) {
    cn = "complete"
  }
  return (<td key={"run-" + run.ID} className={cn}>
    <Link to={"/run/" + run.ID}>{run.State} ({run.Complete} / {run.Total})</Link>
  </td>)
}

class Submit extends Component {
  constructor(props) {
    super()
    this.state = { workflow: "", inputs: "" }
  }

  onChangeWorkflow(ev) {
    this.setState({ workflow: ev.target.files[0] })
  }

  onChangeInputs(ev) {
    this.setState({ inputs: ev.target.files[0] })
  }

  onSubmit(ev) {
    ev.preventDefault()

    if (!this.state.workflow || !this.state.inputs) {
      return
    }

    const data = new FormData()
    data.append("workflow", this.state.workflow)
    data.append("inputs", this.state.inputs)

    fetch("/submit", {
      method: "POST",
      body: data,
    }).then(resp => {
      console.log(resp)
    })
  }

  render() {
    return (<div>
      <form onSubmit={ev => this.onSubmit(ev)}>
        <p>Upload a workflow file</p>
        <input type="file" onChange={ev => this.onChangeWorkflow(ev)} />

        <p>Upload an inputs file</p>
        <input type="file" onChange={ev => this.onChangeInputs(ev)} />
        <button type="submit">Submit</button>
      </form>
    </div>)
  }
}

ReactDOM.render(
  (<Router>
    <div>
      <Route exact path="/" component={Home} />
      <Route exact path="/submitjob" component={Submit} />
      <Route path="/run/:rid" component={Run} />
      <Route path="/runs/:wfid" component={RunsForWorkflow} />
    </div>
  </Router>),
  document.getElementById('root')
)
