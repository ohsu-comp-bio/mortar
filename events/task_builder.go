package events

import (
	"context"

	"github.com/ohsu-comp-bio/tes"
)

// TaskBuilder aggregates events into an in-memory Task object.
type TaskBuilder struct {
	*tes.Task
}

// WriteEvent updates the Task object.
func (tb TaskBuilder) WriteEvent(ctx context.Context, ev *Event) error {
	t := tb.Task
	t.Id = ev.Id
	attempt := int(ev.Attempt)
	index := int(ev.Index)

	switch ev.Type {
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
	}

	return nil
}