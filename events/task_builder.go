package events

import (
	"github.com/ohsu-comp-bio/tes"
)

// WriteEvent updates the Task object.
func WriteEvent(t *tes.Task, ev *Event) {

	attempt := int(ev.Attempt)
	index := int(ev.Index)

	switch ev.Type {
  case Type_TASK_CREATED:
    // TODO this isn't 100% correct. if the events come out of order,
    //      this should deep merge with the existing data.
    et := ev.GetTask()
    // TODO this is weird. Find a better way. Maybe return a task.
    *t = *et

	case Type_TASK_STATE:
		to := ev.GetState()
		t.State = to

	case Type_TASK_START_TIME:
		t.GetTaskLog(attempt).StartTime = ev.GetStartTime()

	case Type_TASK_END_TIME:
		t.GetTaskLog(attempt).EndTime = ev.GetEndTime()

	case Type_TASK_OUTPUTS:
		t.GetTaskLog(attempt).Outputs = ev.GetOutputs().Value

	case Type_TASK_METADATA:
		t.GetTaskLog(attempt).Metadata = ev.GetMetadata().Value

	case Type_EXECUTOR_START_TIME:
		t.GetExecLog(attempt, index).StartTime = ev.GetStartTime()

	case Type_EXECUTOR_END_TIME:
		t.GetExecLog(attempt, index).EndTime = ev.GetEndTime()

	case Type_EXECUTOR_EXIT_CODE:
		t.GetExecLog(attempt, index).ExitCode = ev.GetExitCode()

	case Type_EXECUTOR_STDOUT:
		t.GetExecLog(attempt, index).Stdout += ev.GetStdout()

	case Type_EXECUTOR_STDERR:
		t.GetExecLog(attempt, index).Stderr += ev.GetStderr()

  // TODO include system logs?
	}
}
