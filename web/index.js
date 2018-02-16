import React, { Component } from "react"
import ReactDOM from "react-dom"
import { BrowserRouter as Router, Route, Link } from 'react-router-dom'

/*
fetch("/data.json").then(function(resp) {
  resp.json().then(function(data) {
    console.log(data)
    var dat = JSON.stringify(data)


})})
*/

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

    for (var rid of data.RunIDs) {
      var run = data.Runs[rid]
      header.push(<th key={"run-th-" + rid}>{run.Name}</th>)
    }

    for (var wfid of data.WorkflowIDs) {
      var wf = data.Workflows[wfid]
      var cells = []
      cells.push(<td key={"wf-name-" + wfid}>{wf.Name}</td>)

      for (var rid of data.RunIDs) {
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

    for (var sid of stepIDs) {
      var step = run.Steps[sid]
      var tasks = run.StepTasks[sid]
      var task = run.Tasks[tasks[0]]

      rows.push(<tr key={"step-tr-" + sid}>
        <td>{step.gid}</td>
        <td>{task.state}</td>
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
    </div>
  </Router>),
  document.getElementById('root')
)
