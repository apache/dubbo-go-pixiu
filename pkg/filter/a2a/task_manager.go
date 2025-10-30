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
	"fmt"
	"sync"
	"time"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

// TaskManager manages tasks in memory (simplified implementation for phase 1)
type TaskManager struct {
	config        *TaskConfig
	tasks         map[string]*Task
	mutex         sync.RWMutex
	idCounter     int64
	stopCh        chan struct{}
	cleanupTicker *time.Ticker
}

// NewTaskManager creates a new task manager instance
func NewTaskManager(config *TaskConfig) *TaskManager {
	tm := &TaskManager{
		config:    config,
		tasks:     make(map[string]*Task),
		idCounter: 1,
		stopCh:    make(chan struct{}),
	}

	// Start cleanup routine if cleanup interval is configured
	if config.CleanupInterval > 0 {
		tm.cleanupTicker = time.NewTicker(time.Duration(config.CleanupInterval) * time.Millisecond)
		go tm.cleanupRoutine()
		logger.Infof("[dubbo-go-pixiu] A2A TaskManager cleanup routine started with interval %d ms", config.CleanupInterval)
	}

	return tm
}

// CreateTask creates a new task
func (tm *TaskManager) CreateTask(from, to, taskType string, content map[string]any, timeout int64) (*Task, error) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Check concurrent task limit
	if len(tm.tasks) >= tm.config.MaxConcurrentTasks {
		return nil, fmt.Errorf("maximum concurrent tasks limit reached: %d", tm.config.MaxConcurrentTasks)
	}

	// Generate task ID
	taskID := fmt.Sprintf("task_%d_%d", time.Now().UnixMilli(), tm.idCounter)
	tm.idCounter++

	// Use default timeout if not specified
	if timeout <= 0 {
		timeout = tm.config.DefaultTimeout
	}

	now := time.Now().UnixMilli()
	task := &Task{
		TaskID:    taskID,
		From:      from,
		To:        to,
		Type:      taskType,
		Content:   content,
		Status:    TaskPending,
		CreatedAt: now,
		UpdatedAt: now,
		Timeout:   timeout,
	}

	tm.tasks[taskID] = task
	logger.Infof("[dubbo-go-pixiu] A2A TaskManager created task %s: %s -> %s (%s)", taskID, from, to, taskType)

	return task, nil
}

// GetTask retrieves a task by ID
func (tm *TaskManager) GetTask(taskID string) (*Task, error) {
	tm.mutex.RLock()
	defer tm.mutex.RUnlock()

	task, exists := tm.tasks[taskID]
	if !exists {
		return nil, fmt.Errorf("task not found: %s", taskID)
	}

	// Return a copy to prevent external modification
	taskCopy := *task
	return &taskCopy, nil
}

// UpdateTask updates task status and result
func (tm *TaskManager) UpdateTask(taskID string, status TaskStatus, result map[string]any, errorMsg string) error {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	task, exists := tm.tasks[taskID]
	if !exists {
		return fmt.Errorf("task not found: %s", taskID)
	}

	// Update task fields
	task.Status = status
	task.UpdatedAt = time.Now().UnixMilli()

	if result != nil {
		task.Result = result
	}

	if errorMsg != "" {
		task.Error = errorMsg
	}

	logger.Infof("[dubbo-go-pixiu] A2A TaskManager updated task %s status to %s", taskID, status)
	return nil
}

// CancelTask cancels a task
func (tm *TaskManager) CancelTask(taskID string) error {
	return tm.UpdateTask(taskID, TaskCancelled, nil, "Task canceled")
}

// ListTasks returns all tasks (for debugging purposes)
func (tm *TaskManager) ListTasks() map[string]*Task {
	tm.mutex.RLock()
	defer tm.mutex.RUnlock()

	// Return copies to prevent external modification
	result := make(map[string]*Task, len(tm.tasks))
	for id, task := range tm.tasks {
		taskCopy := *task
		result[id] = &taskCopy
	}

	return result
}

// GetTaskCount returns the current number of tasks
func (tm *TaskManager) GetTaskCount() int {
	tm.mutex.RLock()
	defer tm.mutex.RUnlock()
	return len(tm.tasks)
}

// Stop stops the task manager and cleanup routines
func (tm *TaskManager) Stop() {
	if tm.cleanupTicker != nil {
		tm.cleanupTicker.Stop()
	}
	close(tm.stopCh)
	logger.Infof("[dubbo-go-pixiu] A2A TaskManager stopped")
}

// cleanupRoutine periodically cleans up completed and expired tasks
func (tm *TaskManager) cleanupRoutine() {
	for {
		select {
		case <-tm.cleanupTicker.C:
			tm.cleanupExpiredTasks()
		case <-tm.stopCh:
			return
		}
	}
}

// cleanupExpiredTasks removes expired and old completed tasks
func (tm *TaskManager) cleanupExpiredTasks() {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	now := time.Now().UnixMilli()
	var toDelete []string

	for taskID, task := range tm.tasks {
		shouldDelete := false

		// Check if task is expired (timeout exceeded)
		if task.Status == TaskPending || task.Status == TaskRunning {
			if now-task.CreatedAt > task.Timeout {
				task.Status = TaskFailed
				task.Error = "Task timeout"
				task.UpdatedAt = now
				logger.Warnf("[dubbo-go-pixiu] A2A TaskManager task %s timed out", taskID)
			}
		}

		// Check if completed task should be cleaned up
		if task.Status == TaskCompleted || task.Status == TaskFailed || task.Status == TaskCancelled {
			if now-task.UpdatedAt > tm.config.TaskRetention {
				shouldDelete = true
			}
		}

		if shouldDelete {
			toDelete = append(toDelete, taskID)
		}
	}

	// Delete expired tasks
	for _, taskID := range toDelete {
		delete(tm.tasks, taskID)
		logger.Debugf("[dubbo-go-pixiu] A2A TaskManager cleaned up task %s", taskID)
	}

	if len(toDelete) > 0 {
		logger.Infof("[dubbo-go-pixiu] A2A TaskManager cleaned up %d expired tasks", len(toDelete))
	}
}
