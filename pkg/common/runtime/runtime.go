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

package runtime

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

// GetCPUNum gets current os's cpu number
func GetCPUNum() int {
	if isContainer() {
		cpus, _ := numCPU()
		return cpus
	}
	return runtime.NumCPU()
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
		if strings.Contains(line, _dockerPath) ||
			strings.Contains(line, _kubepodsPath) {
			return true
		}
	}
	return false
}

// copied from https://github.com/containerd/cgroups/blob/318312a373405e5e91134d8063d04d59768a1bff/utils.go#L243
func readUint(path string) (uint64, error) {
	v, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}
	return parseUint(strings.TrimSpace(string(v)), 10, 64)
}

// copied from https://github.com/containerd/cgroups/blob/318312a373405e5e91134d8063d04d59768a1bff/utils.go#L251
func parseUint(s string, base, bitSize int) (uint64, error) {
	v, err := strconv.ParseUint(s, base, bitSize)
	if err != nil {
		intValue, intErr := strconv.ParseInt(s, base, bitSize)
		// 1. Handle negative values greater than MinInt64 (and)
		// 2. Handle negative values lesser than MinInt64
		if intErr == nil && intValue < 0 {
			return 0, nil
		} else if intErr != nil &&
			intErr.(*strconv.NumError).Err == strconv.ErrRange &&
			intValue < 0 {
			return 0, nil
		}
		return 0, err
	}
	return v, nil
}

// numCPU returns the CPU quota
func numCPU() (num int, err error) {
	if !isContainer() {
		return runtime.NumCPU(), nil
	}

	// If the container is running in a cgroup, we can use the cgroup cpu
	// quota to limit the number of CPUs.
	period, err := readUint(_cpuPeriodPath)
	if err != nil {
		return runtime.NumCPU(), err
	}
	quota, err := readUint(_cpuQuotaPath)
	if err != nil {
		return runtime.NumCPU(), err
	}

	// The number of CPUs is the quota divided by the period.
	// See https://www.kernel.org/doc/Documentation/scheduler/sched-bwc.txt
	if quota <= 0 || period <= 0 {
		return runtime.NumCPU(), err
	}

	return int(quota) / int(period), nil
}
