package agent

import (
	"fmt"
	"sync"
)

type IAgent interface {
	GetID() string
	GetName() string
	GetType() string
	GetConfig() interface{}
	GetRunningTasks() AgentTaskMap
	GetCompletedTasks() AgentTaskMap
	AddTask(task IAgentTask) (*IAgentTask, error)
	addSequentialTask(task IAgentTask) (*IAgentTask, error)
	addStandardTask(task IAgentTask) (*IAgentTask, error)
	runTask(task IAgentTask)
	runTaskInBackground(task IAgentTask)
	executeTaskInBackground(task IAgentTask)
	AwaitAllTasks()
	AwaitTask(task IAgentTask)
	RunTaskLoop(quit chan bool)
	Start()
	Stop()
	Kill()
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
	GetChannel() chan interface{}
	GetResult() interface{}
	Execute(callback func())
	Kill()
	AwaitCompletion() interface{}
	IsCompleted() bool
}

type AgentTaskMap map[string]*IAgentTask
type AgentSequentialTaskQueueMap map[AgentTaskType]chan IAgentTask
type Agent struct {
	id                   string                      // Agent ID
	name                 string                      // Agent human-readable name
	_type                string                      // Agent type
	config               interface{}                 // Agent configuration
	completedTasks       AgentTaskMap                // Agent completed tasks
	runningTasks         AgentTaskMap                // Agent running tasks
	taskQueue            chan IAgentTask             // Agent general task queue
	sequentialTaskQueues AgentSequentialTaskQueueMap // Agent queues for sequential tasks
	routines             *sync.WaitGroup             // Waitgroup for running tasks
	taskLoopRoutines     *sync.WaitGroup             // Waitgroup for task loop
	kill                 chan bool                   // Channel to kill agent
}

// NewAgent creates a new agent
//
// id: Unique ID for the agent
// name: Human-readable name for the agent
// agentType: Agent type
// config: Agent configuration
func NewAgent(id string, name string, agentType string, config interface{}) *Agent {
	routines := &sync.WaitGroup{}
	taskLoopRoutines := &sync.WaitGroup{}
	completedTasks := make(AgentTaskMap)
	runningTasks := make(AgentTaskMap)
	taskQueue := make(chan IAgentTask)
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
		taskQueue:            taskQueue,
		sequentialTaskQueues: sequentialTaskQueues,
		kill:                 kill,
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

func (a *Agent) AddTask(task IAgentTask) (*IAgentTask, error) {
	if task.IsSequential() {
		return a.addSequentialTask(task)
	}
	return a.addStandardTask(task)
}

func (a *Agent) addSequentialTask(task IAgentTask) (*IAgentTask, error) {
	if !a.taskExists(task) {
		taskType := task.GetType()
		_, sequentialTaskQueueExists := a.sequentialTaskQueues[taskType]
		if !sequentialTaskQueueExists {
			a.sequentialTaskQueues[taskType] = make(chan IAgentTask)
		}
		a.sequentialTaskQueues[taskType] <- task
	}
	return nil, fmt.Errorf("task already exists with <ID: %s>", task.GetID())
}

func (a *Agent) addStandardTask(task IAgentTask) (*IAgentTask, error) {
	if !a.taskExists(task) {
		a.taskQueue <- task
	}
	return nil, fmt.Errorf("task already exists with <ID: %s>", task.GetID())
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
			a.completedTasks[task.GetID()] = &task
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
		case <-a.kill:
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
		case <-a.kill:
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
	a.kill <- true
	for task := range a.runningTasks {
		(*a.runningTasks[task]).Kill()
	}
}
