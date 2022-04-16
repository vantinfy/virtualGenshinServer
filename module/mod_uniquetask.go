package module

type MyTask struct {
	TaskId int
	State  int
}

type ModUniqueTask struct {
	MyTaskInfo map[int]*MyTask
}

func (p *ModUniqueTask) IsTaskFinish(taskId int) bool {
	task, ok := p.MyTaskInfo[taskId]
	if !ok {
		return false
	}
	return task.State == TASK_STATE_FINISH
}
