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

  componentDidMount() {
    fetch("/data.json")
      .then(resp => resp.json())
      .then(data => this.setState({"data": data}))
  }

  render() {
    if (!this.state.data) {
      return <div></div>
    }

    var data = this.state.data
    var rows = []
    var header = [<th key="empty"></th>]

    for (var i = 0; i < data.RunIDs.length; i++) {
      var rid = data.RunIDs[i]
      var run = data.Runs[rid]
      header.push(<th key={"run-th-" + rid}>{run.Name}</th>)
    }

    for (var i = 0; i < data.WorkflowIDs.length; i++) {
      var wfid = data.WorkflowIDs[i]
      var wf = data.Workflows[wfid]
      var cells = []
      cells.push(<td key={"wf-name-" + wfid}>{wf.Name}</td>)

      for (var j = 0; j < data.RunIDs.length; j++) {
        var rid = data.RunIDs[j]
        var run = data.Runs[rid]
        cells.push(Cell(rid, run))
      }
      rows.push(<tr key={"wf-" + wfid}>{cells}</tr>)
    }

    return (<div>
      <h3>Mortar</h3>
      <table>
        <thead><tr>{header}</tr></thead>
        <tbody>{rows}</tbody>
      </table>
    </div>)
  }
}


class Run extends Component {
  constructor(props) {
    super(props)
    this.state = {
      "data": null,
    }
  }

  componentDidMount() {
    fetch("/data.json")
      .then(resp => resp.json())
      .then(data => this.setState({"data": data}))
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
      <h3>Mortar Steps</h3>
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
    fetch("/data2.json")
      .then(resp => resp.json())
      .then(data => this.setState({"data": data}))
  }

  render() {
    if (!this.state.data) {
      return <div>Loading</div>
    }

    var wfid = this.props.match.params.wfid
    if (!wfid) {
      return <div>Workflow not found</div>
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
      <h3>Mortar Steps</h3>
      <table>
        <thead><tr>
          {header}
        </tr></thead>
        <tbody>{rows}</tbody>
      </table>
    </div>)
  }
}
const Cell = (rid, run) => {
  var cn = ""

  if (run.Total == run.Complete) {
    cn = "complete"
  }
  return (<td key={"run-" + rid} className={cn}>
    <Link to={"/run/" + rid}>{run.Complete} / {run.Total}</Link>
  </td>)
}

ReactDOM.render(
  (<Router>
    <div>
      <Route exact path="/" component={Home} />
      <Route path="/run/:rid" component={Run} />
      <Route path="/runs/:wfid" component={RunsForWorkflow} />
    </div>
  </Router>),
  document.getElementById('root')
)
