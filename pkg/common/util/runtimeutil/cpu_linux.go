//go:build linux
// +build linux

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

package runtimeutil

import (
	"bufio"
	"os"
	"runtime"
	"strconv"
	"strings"
)

const (
	_cgroupPath    = "/proc/self/cgroup"
	_dockerPath    = "/docker"
	_kubepodsPath  = "/kubepods"
	_cpuPeriodPath = "/sys/fs/cgroup/cpu/cpu.cfs_period_us"
	_cpuQuotaPath  = "/sys/fs/cgroup/cpu/cpu.cfs_quota_us"
)

// NumCPU returns the CPU quota
func NumCPU() int {
	if !isContainer() {
		return runtime.NumCPU()
	}

	// If the container is running in a cgroup, we can use the cgroup cpu
	// quota to limit the number of CPUs.
	periodLines := readLinesFromFile(_cpuPeriodPath)
	quotaLines := readLinesFromFile(_cpuQuotaPath)
	if len(periodLines) == 0 || len(quotaLines) == 0 {
		return runtime.NumCPU()
	}

	period, err := strconv.ParseInt(periodLines[0], 10, 64)
	if err != nil {
		return runtime.NumCPU()
	}
	quota, err := strconv.ParseInt(quotaLines[0], 10, 64)
	if err != nil {
		return runtime.NumCPU()
	}

	// The number of CPUs is the quota divided by the period.
	// See https://www.kernel.org/doc/Documentation/scheduler/sched-bwc.txt
	if quota <= 0 || period < quota {
		return runtime.NumCPU()
	}

	return int(quota) / int(period)
}

// readLinesFromFile reads the lines from a file.
func readLinesFromFile(filepath string) []string {
	res := make([]string, 0)
	f, err := os.Open(filepath)
	if err != nil {
		return res
	}
	defer f.Close()
	buff := bufio.NewReader(f)
	for {
		line, _, err := buff.ReadLine()
		if err != nil {
			return res
		}
		res = append(res, string(line))
	}
}

// isContainer returns true if the process is running in a container.
func isContainer() bool {
	lines := readLinesFromFile(_cgroupPath)
	for _, line := range lines {
		if strings.HasPrefix(line, _dockerPath) ||
			strings.HasPrefix(line, _kubepodsPath) {
			return true
		}
	}
	return false
}
