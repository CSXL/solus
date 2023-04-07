package agent

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	// Test task handler functions
	testTaskHandlerWaitKill = func(kill chan bool) interface{} {
		<-kill
		return "CSX Labs, Launching ideas into cyberspace. ;)"
	}
	testTaskHandlerWithResult = func(kill chan bool) interface{} {
		return "test"
	}
	testTaskHandlerWaitForever = func(kill chan bool) interface{} {
		mischeviousChannel := make(chan bool)
		<-mischeviousChannel
		return "This will never be returned :("
	}
	// Test task types
	testTaskType = AgentTaskType{
		Type:         "test",
		IsSequential: false,
	}
	testSequentialTaskType = AgentTaskType{
		Type:         "testSequential",
		IsSequential: true,
	}
)

func TestNewAgent(t *testing.T) {
	agent := NewAgent("testName", "testType", nil)
	assert.Equal(t, "testName", agent.GetName())
	assert.Equal(t, "testType", agent.GetType())
	assert.Equal(t, 0, len(agent.GetRunningTasks()))
	assert.Equal(t, 0, len(agent.GetCompletedTasks()))
}

func TestNewAgentTask(t *testing.T) {
	task := NewAgentTask("test", testTaskType, testTaskHandlerWithResult)
	assert.Equal(t, "test", task.GetName())
	assert.Equal(t, testTaskType, task.GetType())
	assert.False(t, task.IsSequential())
}

func TestAddTask(t *testing.T) {
	agent := NewAgent("testName", "testAgentType", nil)
	agent.Start()
	defer agent.Kill()
	task := NewAgentTask("test", testTaskType, testTaskHandlerWithResult)
	err := agent.AddTask(task)
	assert.Nil(t, err)
}

func TestAddMultipleTasks(t *testing.T) {
	agent := NewAgent("testName", "testAgentType", nil)
	agent.Start()
	defer agent.Kill()
	task := NewAgentTask("test", testTaskType, testTaskHandlerWithResult)
	for i := 0; i < 10; i++ {
		err := agent.AddTask(task)
		assert.Nil(t, err)
	}
}

func TestAddSequentialTask(t *testing.T) {
	agent := NewAgent("testName", "testAgentType", nil)
	agent.Start()
	defer agent.Kill()
	task := NewAgentTask("test", testSequentialTaskType, testTaskHandlerWithResult)
	err := agent.AddTask(task)
	assert.Nil(t, err)
}

func TestAddMultipleSequentialTasks(t *testing.T) {
	agent := NewAgent("testName", "testAgentType", nil)
	agent.Start()
	defer agent.Kill()
	task := NewAgentTask("test", testSequentialTaskType, testTaskHandlerWithResult)
	for i := 0; i < 10; i++ {
		err := agent.AddTask(task)
		assert.Nil(t, err)
	}
}

func TestTaskWithResult(t *testing.T) {
	agent := NewAgent("testName", "testAgentType", nil)
	agent.Start()
	defer agent.Kill()
	task := NewAgentTask("test", testTaskType, testTaskHandlerWithResult)
	err := agent.AddTask(task)
	assert.Nil(t, err)
	result := task.AwaitCompletion()
	assert.Equal(t, "test", result)
	assert.True(t, task.IsCompleted())
}

func TestSequentialTaskWithResult(t *testing.T) {
	agent := NewAgent("testName", "testAgentType", nil)
	agent.Start()
	defer agent.Kill()
	task := NewAgentTask("test", testSequentialTaskType, testTaskHandlerWithResult)
	err := agent.AddTask(task)
	assert.Nil(t, err)
	result := task.AwaitCompletion()
	assert.Equal(t, "test", result)
	assert.True(t, task.IsCompleted())
}

func TestTaskWithKill(t *testing.T) {
	agent := NewAgent("testName", "testAgentType", nil)
	agent.Start()
	defer agent.Kill()
	task := NewAgentTask("test", testTaskType, testTaskHandlerWaitKill)
	err := agent.AddTask(task)
	assert.Nil(t, err)
	task.Kill()
	assert.True(t, task.WasKilled())
}

func TestSequentialTaskWithKill(t *testing.T) {
	agent := NewAgent("testName", "testAgentType", nil)
	agent.Start()
	defer agent.Kill()
	task := NewAgentTask("test", testSequentialTaskType, testTaskHandlerWaitKill)
	err := agent.AddTask(task)
	assert.Nil(t, err)
	task.Kill()
	assert.True(t, task.WasKilled())
}

func TestMischeviousTaskWithKill(t *testing.T) {
	agent := NewAgent("testName", "testAgentType", nil)
	agent.Start()
	defer agent.Kill()
	task := NewAgentTask("test", testTaskType, testTaskHandlerWaitForever)
	err := agent.AddTask(task)
	assert.Nil(t, err)
	task.Kill()
	assert.True(t, task.WasKilled())
}

func TestMischeviousSequentialTaskWithKill(t *testing.T) {
	agent := NewAgent("testName", "testAgentType", nil)
	agent.Start()
	defer agent.Kill()
	task := NewAgentTask("test", testSequentialTaskType, testTaskHandlerWaitForever)
	err := agent.AddTask(task)
	assert.Nil(t, err)
	task.Kill()
	assert.True(t, task.WasKilled())
}

func TestKill(t *testing.T) {
	agent := NewAgent("testName", "testAgentType", nil)
	agent.Start()
	task := NewAgentTask("test", testTaskType, testTaskHandlerWaitKill)
	err := agent.AddTask(task)
	assert.Nil(t, err)
	agent.Kill()
	assert.False(t, agent.IsRunning())
}

func TestKillWithTasks(t *testing.T) {
	agent := NewAgent("testName", "testAgentType", nil)
	agent.Start()
	task := NewAgentTask("test", testTaskType, testTaskHandlerWaitKill)
	err := agent.AddTask(task)
	assert.Nil(t, err)
	agent.Kill()
	assert.False(t, agent.IsRunning())
}

func TestKillWithSequentialTasks(t *testing.T) {
	agent := NewAgent("testName", "testAgentType", nil)
	agent.Start()
	task := NewAgentTask("test", testSequentialTaskType, testTaskHandlerWaitKill)
	err := agent.AddTask(task)
	assert.Nil(t, err)
	agent.Kill()
	assert.False(t, agent.IsRunning())
}

func TestSequentialTaskCompletionOrder(t *testing.T) {
	agent := NewAgent("testAgent", "testAgentType", nil)
	agent.Start()
	defer agent.Kill()
	handlerSignal := make(chan bool)
	testerSignal := make(chan bool)
	testHandler := func(kill chan bool) interface{} {
		handlerSignal <- true
		<-testerSignal
		return nil
	}
	var taskList []*AgentTask
	for i := 0; i < 10; i++ {
		taskName := fmt.Sprintf("task#%d", i)
		task := NewAgentTask(taskName, testSequentialTaskType, testHandler)
		err := agent.AddTask(task)
		assert.Nil(t, err)
		taskList = append(taskList, task)
	}
	for _, task := range taskList {
		<-handlerSignal
		_, inRunningTasks := agent.GetRunningTasks()[task.GetID()]
		assert.Truef(t, inRunningTasks, "Task <ID: %s, Name: %s> not in running tasks <%s>", task.GetID(), task.GetName(), fmt.Sprint(agent.GetRunningTasks()))
		testerSignal <- true
		task.AwaitCompletion()
		assert.True(t, task.IsCompleted())
		_, inCompletedTasks := agent.GetCompletedTasks()[task.GetID()]
		assert.True(t, inCompletedTasks)
	}
}
