package agent

import (
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
	assert.Equal(t, false, task.IsSequential())
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

func TestAddSequentialTasks(t *testing.T) {
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
	assert.Equal(t, true, task.IsCompleted())
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
	assert.Equal(t, true, task.IsCompleted())
}

func TestTaskWithKill(t *testing.T) {
	agent := NewAgent("testName", "testAgentType", nil)
	agent.Start()
	defer agent.Kill()
	task := NewAgentTask("test", testTaskType, testTaskHandlerWaitKill)
	err := agent.AddTask(task)
	assert.Nil(t, err)
	task.Kill()
	assert.Equal(t, true, task.WasKilled())
}

func TestSequentialTaskWithKill(t *testing.T) {
	agent := NewAgent("testName", "testAgentType", nil)
	agent.Start()
	defer agent.Kill()
	task := NewAgentTask("test", testSequentialTaskType, testTaskHandlerWaitKill)
	err := agent.AddTask(task)
	assert.Nil(t, err)
	task.Kill()
	assert.Equal(t, true, task.WasKilled())
}

func TestMischeviousTaskWithKill(t *testing.T) {
	agent := NewAgent("testName", "testAgentType", nil)
	agent.Start()
	defer agent.Kill()
	task := NewAgentTask("test", testTaskType, testTaskHandlerWaitForever)
	err := agent.AddTask(task)
	assert.Nil(t, err)
	task.Kill()
	assert.Equal(t, true, task.WasKilled())
}
