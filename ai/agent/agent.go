package agent

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
)

func generateUUID() string {
	uuid := uuid.New()
	return uuid.String()
}

type IAgent interface {
	GetID() string
	GetName() string
	GetType() string
	GetConfig() interface{}
	GetRunningTasks() AgentTaskMap
	GetCompletedTasks() AgentTaskMap
	GetKilledTasks() AgentTaskMap
	AddTask(task IAgentTask) (*IAgentTask, error)
	addSequentialTask(task IAgentTask) (*IAgentTask, error)
	addStandardTask(task IAgentTask) (*IAgentTask, error)
	runTask(task IAgentTask)
	runTaskInBackground(task IAgentTask)
	executeTaskInBackground(task IAgentTask)
	AwaitAllTasks()
	runTaskLoop(quit chan bool)
	Start()
	Stop()
	Kill()
}

type AgentTaskMap map[string]*IAgentTask
type AgentSequentialTaskQueueMap map[AgentTaskType]chan IAgentTask
type Agent struct {
	IAgent
	id                   string                      // Agent ID
	name                 string                      // Agent human-readable name
	_type                string                      // Agent type
	config               interface{}                 // Agent configuration
	completedTasks       AgentTaskMap                // Agent completed tasks
	runningTasks         AgentTaskMap                // Agent running tasks
	killedTasks          AgentTaskMap                // Agent killed tasks
	taskQueue            chan IAgentTask             // Agent general task queue
	sequentialTaskQueues AgentSequentialTaskQueueMap // Agent queues for sequential tasks
	routines             *sync.WaitGroup             // Waitgroup for running tasks
	taskLoopRoutines     *sync.WaitGroup             // Waitgroup for task loop
	killChannel          chan bool                   // Channel to kill agent
}

// NewAgent creates a new agent
//
// id: Unique ID for the agent
// name: Human-readable name for the agent
// agentType: Agent type
// config: Agent configuration
func NewAgent(name string, agentType string, config interface{}) *Agent {
	id := generateUUID()
	routines := &sync.WaitGroup{}
	taskLoopRoutines := &sync.WaitGroup{}
	completedTasks := make(AgentTaskMap)
	runningTasks := make(AgentTaskMap)
	killedTasks := make(AgentTaskMap)
	taskQueue := make(chan IAgentTask, 100)
	sequentialTaskQueues := make(AgentSequentialTaskQueueMap)
	kill := make(chan bool)
	return &Agent{
		id:                   id,
		name:                 name,
		_type:                agentType,
		config:               config,
		routines:             routines,
		taskLoopRoutines:     taskLoopRoutines,
		runningTasks:         runningTasks,
		completedTasks:       completedTasks,
		killedTasks:          killedTasks,
		taskQueue:            taskQueue,
		sequentialTaskQueues: sequentialTaskQueues,
		killChannel:          kill,
	}
}

func (a *Agent) GetID() string {
	return a.id
}

func (a *Agent) GetName() string {
	return a.name
}

func (a *Agent) GetType() string {
	return a._type
}

func (a *Agent) GetConfig() interface{} {
	return a.config
}

func (a *Agent) GetRunningTasks() AgentTaskMap {
	return a.runningTasks
}

func (a *Agent) GetCompletedTasks() AgentTaskMap {
	return a.completedTasks
}

func (a *Agent) GetKilledTasks() AgentTaskMap {
	return a.killedTasks
}

// AddTask adds a task for the agent to execute in its task loop
//
// task: Task to add
func (a *Agent) AddTask(task IAgentTask) error {
	if task.IsSequential() {
		return a.addSequentialTask(task)
	}
	return a.addStandardTask(task)
}

func (a *Agent) addSequentialTask(task IAgentTask) error {
	if !a.taskExists(task) {
		taskType := task.GetType()
		_, sequentialTaskQueueExists := a.sequentialTaskQueues[taskType]
		if !sequentialTaskQueueExists {
			a.sequentialTaskQueues[taskType] = make(chan IAgentTask, 100)
		}
		a.sequentialTaskQueues[taskType] <- task
		return nil
	}
	return fmt.Errorf("task already exists with <ID: %s>", task.GetID())
}

func (a *Agent) addStandardTask(task IAgentTask) error {
	if !a.taskExists(task) {
		a.taskQueue <- task
		return nil
	}
	return fmt.Errorf("task already exists with <ID: %s>", task.GetID())
}

func (a *Agent) taskExists(task IAgentTask) bool {
	_, completedTaskExists := a.completedTasks[task.GetID()]
	_, runningTaskExists := a.runningTasks[task.GetID()]
	return completedTaskExists || runningTaskExists
}

func (a *Agent) runTaskInBackground(task IAgentTask) (*IAgentTask, error) {
	if !a.taskExists(task) {
		a.runningTasks[task.GetID()] = &task
		task := a.executeTaskInBackground(task)
		return task, nil
	}
	return nil, fmt.Errorf("task already exists with <ID: %s>", task.GetID())
}

func (a *Agent) executeTaskInBackground(task IAgentTask) *IAgentTask {
	a.runningTasks[task.GetID()] = &task
	a.routines.Add(1)
	go func() {
		task.Execute(func() {
			a.routines.Done()
			delete(a.runningTasks, task.GetID())
			if task.WasKilled() {
				a.completedTasks[task.GetID()] = &task
			} else {
				a.killedTasks[task.GetID()] = &task
			}
		})
	}()
	return &task
}

func (a *Agent) AwaitAllTasks() {
	a.routines.Wait()
}

func (a *Agent) runTaskLoop(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case task := <-a.taskQueue:
			a.runTaskInBackground(task)
		case <-a.killChannel:
			close(a.taskQueue)
			return
		}
	}
}

func (a *Agent) runSequentialTaskLoop(wg *sync.WaitGroup, taskType AgentTaskType) {
	defer wg.Done()
	for {
		select {
		case task := <-a.sequentialTaskQueues[taskType]:
			a.runTaskInBackground(task)
		case <-a.killChannel:
			close(a.sequentialTaskQueues[taskType])
			return
		}
	}
}

func (a *Agent) Start() {
	var wg sync.WaitGroup
	wg.Add(1)
	go a.runTaskLoop(&wg)
	for taskType := range a.sequentialTaskQueues {
		wg.Add(1)
		go a.runSequentialTaskLoop(&wg, taskType)
	}
	a.taskLoopRoutines = &wg
}

func (a *Agent) Stop() {
	close(a.taskQueue)
	for taskType := range a.sequentialTaskQueues {
		close(a.sequentialTaskQueues[taskType])
	}
	a.taskLoopRoutines.Wait()
}

func (a *Agent) Kill() {
	a.killChannel <- true
	for task := range a.runningTasks {
		(*a.runningTasks[task]).Kill()
		a.killedTasks[task] = a.runningTasks[task]
	}
}

type AgentTaskType struct {
	Type         string
	IsSequential bool
}

type IAgentTask interface {
	GetID() string
	GetName() string
	GetType() AgentTaskType
	IsSequential() bool
	WasKilled() bool
	GetResult() interface{}
	Execute(callback func())
	Kill()
	AwaitCompletion() interface{}
	IsCompleted() bool
}

type handlerFunction func(kill chan bool) interface{}
type AgentTask struct {
	IAgentTask
	id                string
	name              string
	_type             AgentTaskType
	completionChannel chan interface{}
	result            interface{}
	isCompleted       bool
	wasKilled         bool
	handler           handlerFunction
	kill              chan bool
}

func NewAgentTask(name string, taskType AgentTaskType, handler handlerFunction) *AgentTask {
	id := generateUUID()
	var result interface{}
	channel := make(chan interface{}, 1)
	kill := make(chan bool, 1)
	isCompleted := false
	wasKilled := false
	return &AgentTask{
		id:                id,
		name:              name,
		_type:             taskType,
		completionChannel: channel,
		result:            result,
		kill:              kill,
		isCompleted:       isCompleted,
		wasKilled:         wasKilled,
		handler:           handler,
	}
}

func (t *AgentTask) GetID() string {
	return t.id
}

func (t *AgentTask) GetName() string {
	return t.name
}

func (t *AgentTask) GetType() AgentTaskType {
	return t._type
}

func (t *AgentTask) IsSequential() bool {
	return t._type.IsSequential
}

func (t *AgentTask) GetResult() interface{} {
	return t.result
}

func (t *AgentTask) Execute(callback func()) {
	result := t.handler(t.kill)
	t.result = result
	callback()
	t.isCompleted = true
	t.completionChannel <- result
}

func (t *AgentTask) Kill() {
	t.kill <- true
	t.wasKilled = true
}

func (t *AgentTask) AwaitCompletion() interface{} {
	return <-t.completionChannel
}

func (t *AgentTask) IsCompleted() bool {
	return t.isCompleted
}

func (t *AgentTask) WasKilled() bool {
	return t.wasKilled
}
