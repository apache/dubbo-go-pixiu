/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package a2a

import (
	"testing"
	"time"
)

import (
	"github.com/stretchr/testify/assert"
)

// TestTaskManager_CreateAndGet tests basic task creation and retrieval
func TestTaskManager_CreateAndGet(t *testing.T) {
	config := &TaskConfig{
		DefaultTimeout:     5000,
		MaxConcurrentTasks: 10,
		CleanupInterval:    0, // Disable cleanup for this test
	}

	tm := NewTaskManager(config)
	defer tm.Stop()

	// Create a task
	task, err := tm.CreateTask("agent-1", "agent-2", "test_task", map[string]interface{}{
		"message": "test",
	}, 3000)

	assert.NoError(t, err)
	assert.NotNil(t, task)
	assert.Equal(t, "agent-1", task.From)
	assert.Equal(t, "agent-2", task.To)
	assert.Equal(t, "test_task", task.Type)
	assert.Equal(t, TaskPending, task.Status)

	// Get the task
	retrievedTask, err := tm.GetTask(task.TaskID)
	assert.NoError(t, err)
	assert.Equal(t, task.TaskID, retrievedTask.TaskID)
	assert.Equal(t, task.From, retrievedTask.From)
}

// TestTaskManager_UpdateTask tests task updates
func TestTaskManager_UpdateTask(t *testing.T) {
	config := &TaskConfig{
		DefaultTimeout:     5000,
		MaxConcurrentTasks: 10,
		CleanupInterval:    0,
	}

	tm := NewTaskManager(config)
	defer tm.Stop()

	// Create a task
	task, err := tm.CreateTask("agent-1", "agent-2", "test_task", map[string]interface{}{
		"data": "test",
	}, 0)
	assert.NoError(t, err)

	// Update the task
	result := map[string]interface{}{
		"output": "success",
	}
	err = tm.UpdateTask(task.TaskID, TaskCompleted, result, "")
	assert.NoError(t, err)

	// Verify update
	updatedTask, err := tm.GetTask(task.TaskID)
	assert.NoError(t, err)
	assert.Equal(t, TaskCompleted, updatedTask.Status)
	assert.NotNil(t, updatedTask.Result)
	assert.Equal(t, "success", updatedTask.Result["output"])
}

// TestTaskManager_CancelTask tests task cancellation
func TestTaskManager_CancelTask(t *testing.T) {
	config := &TaskConfig{
		DefaultTimeout:     5000,
		MaxConcurrentTasks: 10,
		CleanupInterval:    0,
	}

	tm := NewTaskManager(config)
	defer tm.Stop()

	// Create a task
	task, err := tm.CreateTask("agent-1", "agent-2", "test_task", map[string]interface{}{}, 0)
	assert.NoError(t, err)

	// Cancel the task
	err = tm.CancelTask(task.TaskID)
	assert.NoError(t, err)

	// Verify cancellation
	cancelledTask, err := tm.GetTask(task.TaskID)
	assert.NoError(t, err)
	assert.Equal(t, TaskCancelled, cancelledTask.Status)
}

// TestTaskManager_ConcurrentLimit tests the concurrent task limit
func TestTaskManager_ConcurrentLimit(t *testing.T) {
	config := &TaskConfig{
		DefaultTimeout:     5000,
		MaxConcurrentTasks: 3, // Small limit for testing
		CleanupInterval:    0,
	}

	tm := NewTaskManager(config)
	defer tm.Stop()

	// Create tasks up to the limit
	for i := 0; i < 3; i++ {
		_, err := tm.CreateTask("agent-1", "agent-2", "test_task", map[string]interface{}{}, 0)
		assert.NoError(t, err)
	}

	// Try to create one more task (should fail)
	_, err := tm.CreateTask("agent-1", "agent-2", "test_task", map[string]interface{}{}, 0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "maximum concurrent tasks limit reached")
}

// TestTaskManager_GetTaskNotFound tests getting a non-existent task
func TestTaskManager_GetTaskNotFound(t *testing.T) {
	config := &TaskConfig{
		DefaultTimeout:     5000,
		MaxConcurrentTasks: 10,
		CleanupInterval:    0,
	}

	tm := NewTaskManager(config)
	defer tm.Stop()

	// Try to get a non-existent task
	_, err := tm.GetTask("non-existent-task-id")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "task not found")
}

// TestTaskManager_ListTasks tests listing all tasks
func TestTaskManager_ListTasks(t *testing.T) {
	config := &TaskConfig{
		DefaultTimeout:     5000,
		MaxConcurrentTasks: 10,
		CleanupInterval:    0,
	}

	tm := NewTaskManager(config)
	defer tm.Stop()

	// Create multiple tasks
	task1, _ := tm.CreateTask("agent-1", "agent-2", "task1", map[string]interface{}{}, 0)
	task2, _ := tm.CreateTask("agent-1", "agent-3", "task2", map[string]interface{}{}, 0)

	// List all tasks
	tasks := tm.ListTasks()
	assert.Equal(t, 2, len(tasks))
	assert.NotNil(t, tasks[task1.TaskID])
	assert.NotNil(t, tasks[task2.TaskID])
}

// TestTaskManager_Cleanup tests the cleanup mechanism
func TestTaskManager_Cleanup(t *testing.T) {
	config := &TaskConfig{
		DefaultTimeout:     100, // Very short timeout
		MaxConcurrentTasks: 10,
		CleanupInterval:    50,  // Frequent cleanup
		TaskRetention:      100, // Short retention
	}

	tm := NewTaskManager(config)
	defer tm.Stop()

	// Create a task
	task, err := tm.CreateTask("agent-1", "agent-2", "test_task", map[string]interface{}{}, 100)
	assert.NoError(t, err)

	// Wait for timeout
	time.Sleep(150 * time.Millisecond)

	// Manually trigger cleanup
	tm.cleanupExpiredTasks()

	// Task should be marked as failed due to timeout
	timedOutTask, err := tm.GetTask(task.TaskID)
	assert.NoError(t, err)
	assert.Equal(t, TaskFailed, timedOutTask.Status)
	assert.Contains(t, timedOutTask.Error, "timeout")
}

// TestTaskManager_GetTaskCount tests the task count functionality
func TestTaskManager_GetTaskCount(t *testing.T) {
	config := &TaskConfig{
		DefaultTimeout:     5000,
		MaxConcurrentTasks: 10,
		CleanupInterval:    0,
	}

	tm := NewTaskManager(config)
	defer tm.Stop()

	assert.Equal(t, 0, tm.GetTaskCount())

	// Create tasks
	tm.CreateTask("agent-1", "agent-2", "task1", map[string]interface{}{}, 0)
	assert.Equal(t, 1, tm.GetTaskCount())

	tm.CreateTask("agent-1", "agent-3", "task2", map[string]interface{}{}, 0)
	assert.Equal(t, 2, tm.GetTaskCount())
}
