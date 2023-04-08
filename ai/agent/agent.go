package agent

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/logger"
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
	IsRunning() bool
	GetRunningTasks() AgentTaskMap
	GetCompletedTasks() AgentTaskMap
	GetKilledTasks() AgentTaskMap
	AddTask(task IAgentTask) error
	addSequentialTask(task IAgentTask) error
	addStandardTask(task IAgentTask) error
	runTask(task IAgentTask) error
	runTaskInBackground(task IAgentTask) error
	executeTaskInBackground(task IAgentTask)
	AwaitAllTasks()
	runTaskLoop(wg *sync.WaitGroup)
	runSequentialTaskLoop(taskType AgentTaskType, wg *sync.WaitGroup)
	runSequentialTaskEventLoop(wg *sync.WaitGroup)
	Start()
	Stop()
	Kill()
}

type agentSequentialTaskEvent string

const (
	agentSequentialTaskEventNewType agentSequentialTaskEvent = "newType"
)

type AgentTaskMap map[string]*IAgentTask
type AgentSequentialTaskQueueMap map[AgentTaskType]chan IAgentTask
type agentSequentialTaskEvents chan struct {
	agentSequentialTaskEvent
	string
}
type Agent struct {
	IAgent
	ctx                  context.Context             // Agent context
	id                   string                      // Agent ID
	name                 string                      // Agent human-readable name
	_type                string                      // Agent type
	config               interface{}                 // Agent configuration
	isRunning            bool                        // Agent running state
	completedTasks       AgentTaskMap                // Agent completed tasks
	runningTasks         AgentTaskMap                // Agent running tasks
	killedTasks          AgentTaskMap                // Agent killed tasks
	taskQueue            chan IAgentTask             // Agent general task queue
	sequentialTaskQueues AgentSequentialTaskQueueMap // Agent queues for sequential tasks
	sequentialTaskEvents agentSequentialTaskEvents   // Agent events for sequential tasks
	routines             sync.WaitGroup              // Waitgroup for running tasks
	taskLoopRoutines     sync.WaitGroup              // Waitgroup for task loop
	taskLoopMutex        sync.Mutex                  // Mutex for task loop
	killChannel          chan bool                   // Channel to kill agent
}

// NewAgent creates a new agent
//
// id: Unique ID for the agent
// name: Human-readable name for the agent
// agentType: Agent type
// config: Agent configuration
func NewAgent(name string, agentType string, config interface{}) *Agent {
	ctx := context.Background()
	id := generateUUID()
	isRunning := false
	completedTasks := make(AgentTaskMap)
	runningTasks := make(AgentTaskMap)
	killedTasks := make(AgentTaskMap)
	taskQueue := make(chan IAgentTask, 100)
	sequentialTaskQueues := make(AgentSequentialTaskQueueMap)
	sequentialTaskEvents := make(agentSequentialTaskEvents)
	kill := make(chan bool)
	return &Agent{
		ctx:                  ctx,
		id:                   id,
		name:                 name,
		_type:                agentType,
		config:               config,
		isRunning:            isRunning,
		routines:             sync.WaitGroup{},
		taskLoopRoutines:     sync.WaitGroup{},
		runningTasks:         runningTasks,
		completedTasks:       completedTasks,
		killedTasks:          killedTasks,
		taskQueue:            taskQueue,
		sequentialTaskQueues: sequentialTaskQueues,
		taskLoopMutex:        sync.Mutex{},
		sequentialTaskEvents: sequentialTaskEvents,
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

func (a *Agent) IsRunning() bool {
	return a.isRunning
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
		logger.Infof("Adding sequential task <ID: %s, Name: %s> to agent <ID: %s, Name: %s>", task.GetID(), task.GetName(), a.GetID(), a.GetName())
		taskType := task.GetType()
		_, sequentialTaskQueueExists := a.sequentialTaskQueues[taskType]
		if !sequentialTaskQueueExists {
			a.sequentialTaskQueues[taskType] = make(chan IAgentTask, 100)
			a.sequentialTaskEvents <- struct {
				agentSequentialTaskEvent
				string
			}{agentSequentialTaskEventNewType, string(taskType.Type)}
		}
		a.sequentialTaskQueues[taskType] <- task
		return nil
	}
	return fmt.Errorf("task already exists with <ID: %s> on agent <ID: %s>", task.GetID(), a.GetID())
}

func (a *Agent) addStandardTask(task IAgentTask) error {
	if !a.taskExists(task) {
		logger.Infof("Adding standard task <ID: %s, Name: %s> to agent <ID: %s, Name: %s>", task.GetID(), task.GetName(), a.GetID(), a.GetName())
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

func (a *Agent) runTaskInBackground(task IAgentTask) error {
	if !a.taskExists(task) {
		a.runningTasks[task.GetID()] = &task
		a.executeTaskInBackground(task)
		return nil
	}
	return fmt.Errorf("task already exists with <ID: %s>", task.GetID())
}

func (a *Agent) runTask(task IAgentTask) error {
	err := a.runTaskInBackground(task)
	if err != nil {
		return err
	}
	task.AwaitCompletion()
	return nil
}

func (a *Agent) executeTaskInBackground(task IAgentTask) {
	a.runningTasks[task.GetID()] = &task
	a.incrementTaskRoutines()
	go func() {
		logger.Infof("Executing task <ID: %s, Name: %s> on agent <ID: %s, Name: %s>", task.GetID(), task.GetName(), a.GetID(), a.GetName())
		task.Execute(func() {
			a.decrementTaskRoutines()
			delete(a.runningTasks, task.GetID())
			if task.WasKilled() {
				a.killedTasks[task.GetID()] = &task
			} else {
				a.completedTasks[task.GetID()] = &task
			}
		})
		logger.Infof("Finished executing task <ID: %s, Name: %s> on agent <ID: %s, Name: %s>", task.GetID(), task.GetName(), a.GetID(), a.GetName())
	}()
}

func (a *Agent) AwaitAllTasks() {
	a.routines.Wait()
}

func (a *Agent) lockTaskLoop() {
	a.taskLoopMutex.Lock()
}

func (a *Agent) unlockTaskLoop() {
	a.taskLoopMutex.Unlock()
}

func (a *Agent) incrementTaskLoopRoutines() {
	a.lockTaskLoop()
	a.taskLoopRoutines.Add(1)
	a.unlockTaskLoop()
}

func (a *Agent) decrementTaskLoopRoutines() {
	a.lockTaskLoop()
	a.taskLoopRoutines.Done()
	a.unlockTaskLoop()
}

func (a *Agent) incrementTaskRoutines() {
	a.lockTaskLoop()
	a.routines.Add(1)
	a.unlockTaskLoop()
}

func (a *Agent) decrementTaskRoutines() {
	a.lockTaskLoop()
	a.routines.Done()
	a.unlockTaskLoop()
}

func (a *Agent) runTaskLoop() {
	defer a.decrementTaskLoopRoutines()
	for {
		select {
		case task, ok := <-a.taskQueue:
			if !ok {
				return
			}
			logger.Infof("Running standard task <ID: %s, Name: %s> on agent <ID: %s, Name: %s>", task.GetID(), task.GetName(), a.GetID(), a.GetName())
			a.runTaskInBackground(task)
		case <-a.killChannel:
			close(a.taskQueue)
			return
		}
	}
}
func (a *Agent) runSequentialTaskLoop(taskType AgentTaskType) {
	defer a.decrementTaskLoopRoutines()
	for {
		select {
		case task, ok := <-a.sequentialTaskQueues[taskType]:
			if !ok {
				return
			}
			logger.Infof("Running sequential task <ID: %s, Name: %s> on agent <ID: %s, Name: %s>", task.GetID(), task.GetName(), a.GetID(), a.GetName())
			a.runTask(task) // Runs the task blocking the task loop
			logger.Infof("Finished sequential task <ID: %s, Name: %s> on agent <ID: %s, Name: %s>", task.GetID(), task.GetName(), a.GetID(), a.GetName())
		case <-a.killChannel:
			close(a.sequentialTaskQueues[taskType])
			return
		}
	}
}

func (a *Agent) runSequentialTaskEventLoop() {
	defer a.decrementTaskLoopRoutines()
	for {
		select {
		case eventType, ok := <-a.sequentialTaskEvents:
			if !ok {
				return
			}
			switch eventType.agentSequentialTaskEvent {
			case agentSequentialTaskEventNewType:
				a.incrementTaskLoopRoutines()
				taskType := AgentTaskType{
					Type:         eventType.string,
					IsSequential: true,
				}
				go a.runSequentialTaskLoop(taskType)
			}
		case <-a.killChannel:
			for taskType := range a.sequentialTaskQueues {
				close(a.sequentialTaskQueues[taskType])
			}
			close(a.sequentialTaskEvents)
		}
	}
}

func (a *Agent) Start() {
	logger.Infof("Starting agent <ID: %s, Name: %s>", a.GetID(), a.GetName())
	if a.isRunning {
		logger.Infof("Agent <ID: %s, Name: %s> is already running. Start canceled.", a.GetID(), a.GetName())
		return
	}
	a.isRunning = true
	a.incrementTaskLoopRoutines()
	go a.runTaskLoop()
	a.incrementTaskLoopRoutines()
	go a.runSequentialTaskEventLoop()
	logger.Infof("Started agent <ID: %s, Name: %s>", a.GetID(), a.GetName())
}

func (a *Agent) Stop() {
	logger.Infof("Stopping agent <ID: %s, Name: %s>", a.GetID(), a.GetName())
	if !a.isRunning {
		logger.Infof("Agent <ID: %s, Name: %s> is not running. Stop canceled.", a.GetID(), a.GetName())
		return
	}
	close(a.taskQueue)
	if a.sequentialTaskQueues != nil {
		for taskType := range a.sequentialTaskQueues {
			close(a.sequentialTaskQueues[taskType])
		}
	}
	close(a.sequentialTaskEvents)
	a.taskLoopRoutines.Wait()
	a.isRunning = false
	logger.Infof("Stopped agent <ID: %s, Name: %s>", a.GetID(), a.GetName())
}

func (a *Agent) Kill() {
	logger.Infof("Killing agent <ID: %s, Name: %s>", a.GetID(), a.GetName())
	if !a.isRunning {
		logger.Infof("Agent <ID: %s, Name: %s> is not running. Kill canceled.", a.GetID(), a.GetName())
		return
	}
	a.killChannel <- true
	for task := range a.runningTasks {
		(*a.runningTasks[task]).Kill()
		a.killedTasks[task] = a.runningTasks[task]
	}
	a.isRunning = false
	logger.Infof("Killed agent <ID: %s, Name: %s>", a.GetID(), a.GetName())
}

type AgentTaskType struct {
	Type         string
	IsSequential bool
}

func NewAgentTaskType(taskType string, isSequential bool) AgentTaskType {
	return AgentTaskType{
		Type:         taskType,
		IsSequential: isSequential,
	}
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

type HandlerFunction func(kill chan bool) interface{}
type AgentTask struct {
	IAgentTask
	id                string
	name              string
	_type             AgentTaskType
	completionChannel chan interface{}
	result            interface{}
	isCompleted       bool
	wasKilled         bool
	handler           HandlerFunction
	killChannel       chan bool
}

func NewAgentTask(name string, taskType AgentTaskType, handler HandlerFunction) *AgentTask {
	id := generateUUID()
	var result interface{}
	completionChannel := make(chan interface{}, 1)
	killChannel := make(chan bool, 1)
	isCompleted := false
	wasKilled := false
	return &AgentTask{
		id:                id,
		name:              name,
		_type:             taskType,
		completionChannel: completionChannel,
		result:            result,
		killChannel:       killChannel,
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
	result := t.handler(t.killChannel)
	t.result = result
	callback()
	t.isCompleted = true
}

func (t *AgentTask) Kill() {
	for len(t.killChannel) > 0 {
		<-t.killChannel
	}
	t.killChannel <- true
	t.wasKilled = true
}

func (t *AgentTask) AwaitCompletion() interface{} {
	for !t.IsCompleted() {
	}
	return t.GetResult()
}

func (t *AgentTask) IsCompleted() bool {
	return t.isCompleted
}

func (t *AgentTask) WasKilled() bool {
	return t.wasKilled
}
