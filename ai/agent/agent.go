package agent

type IAgent interface {
	GetID() string
	GetName() string
	GetType() string
	GetConfig() interface{}
	GetTasks() []IAgentTask
	AddTask(task IAgentTask)
	PopTask() IAgentTask
	GetLastTask() IAgentTask
}

type IAgentTask interface {
	GetID() string
	GetName() string
	GetType() string
	GetChannel() chan interface{}
	IsCompleted() bool
}

type Agent struct {
	ID     string
	Name   string
	Type   string
	Tasks  []IAgentTask
	Config interface{}
}

func (a *Agent) NewAgent(id string, name string, agentType string, config interface{}) *Agent {
	return &Agent{
		ID:     id,
		Name:   name,
		Type:   agentType,
		Config: config,
	}
}

func (a *Agent) GetID() string {
	return a.ID
}

func (a *Agent) GetName() string {
	return a.Name
}

func (a *Agent) GetType() string {
	return a.Type
}

func (a *Agent) GetConfig() interface{} {
	return a.Config
}

func (a *Agent) GetTasks() []IAgentTask {
	return a.Tasks
}

func (a *Agent) AddTask(task IAgentTask) {
	a.Tasks = append(a.Tasks, task)
}

func (a *Agent) PopTask() IAgentTask {
	if len(a.Tasks) == 0 {
		return nil
	}
	task := a.Tasks[0]
	a.Tasks = a.Tasks[1:]
	return task
}

func (a *Agent) GetLastTask() IAgentTask {
	if len(a.Tasks) == 0 {
		return nil
	}
	return a.Tasks[len(a.Tasks)-1]
}
