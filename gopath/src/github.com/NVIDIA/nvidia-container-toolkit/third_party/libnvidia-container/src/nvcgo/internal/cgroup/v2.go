/*
 * Copyright (c) 2021, NVIDIA CORPORATION.  All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cgroup

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cilium/ebpf"
	"golang.org/x/sys/unix"
	internal "nvcgo/internal/ebpf"
)

// GetDeviceCGroupMountPath returns the mount path for the device cgroup controller associated with pid
func (c *cgroupv2) GetDeviceCGroupMountPath(procRootPath string, pid int) (string, error) {
	// Open the pid's mountinfo file in /proc.
	path := fmt.Sprintf(filepath.Join(procRootPath, "proc", "%v", "mountinfo"), pid)
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Create a scanner to loop through the file's contents.
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	// Loop through the file looking for a subsystem of '' (i.e. unified) entry.
	for scanner.Scan() {
		// Split each entry by '[space]'
		parts := strings.Split(scanner.Text(), " ")
		if len(parts) < 5 {
			return "", fmt.Errorf("malformed mountinfo entry: %v", scanner.Text())
		}
		// Look for an entry with cgroup2 as the mount type.
		if parts[len(parts)-3] != "cgroup2" {
			continue
		}
		// Return the 4th element as the moint point of the devices cgroup.
		return parts[4], nil
	}

	return "", fmt.Errorf("no cgroup2 filesystem in mountinfo file")
}

// GetDeviceCGroupMountPath returns the root path for the device cgroup controller associated with pid
func (c *cgroupv2) GetDeviceCGroupRootPath(procRootPath string, pid int) (string, error) {
	// Open the pid's cgroup file in /proc.
	path := fmt.Sprintf(filepath.Join(procRootPath, "proc", "%v", "cgroup"), pid)
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Create a scanner to loop through the file's contents.
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	// Loop through the file looking for either a '' (i.e. unified) entry.
	for scanner.Scan() {
		// Split each entry by ':'
		parts := strings.SplitN(scanner.Text(), ":", 3)
		if len(parts) != 3 {
			return "", fmt.Errorf("malformed cgroup entry: %v", scanner.Text())
		}
		if parts[1] != "" {
			return "", fmt.Errorf("malformed cgroup entry: %v", scanner.Text())
		}
		// Return the cgroup root from the 2nd element.
		return parts[2], nil
	}

	return "", fmt.Errorf("no lines in cgroup file")
}

// AddDeviceRules adds a set of device rules for the device cgroup at cgroupPath
func (c *cgroupv2) AddDeviceRules(cgroupPath string, rules []DeviceRule) error {
	// Open the cgroup path.
	dirFD, err := unix.Open(cgroupPath, unix.O_DIRECTORY|unix.O_RDONLY, 0600)
	if err != nil {
		return fmt.Errorf("unable to open the cgroup path: %v", err)
	}
	defer unix.Close(dirFD)

	// Find any existing eBPF device filter programs attached to this cgroup.
	oldProgs, err := FindAttachedCgroupDeviceFilters(dirFD)
	if err != nil {
		return fmt.Errorf("unable to find any existing device filters attached to the cgroup: %v", err)
	}

	// For each existing program that we found...
	for _, oldProg := range oldProgs {
		// Retreive the instructions from the original program.
		oldInsts, oldLicense, err := (&internal.Program{oldProg}).GetInstructions()
		if err != nil {
			return fmt.Errorf("unable to get the instructions of the original device filters program: %v", err)
		}

		// Prepend instructions for the new devices to the original set of instructions.
		newInsts, err := PrependDeviceFilter(rules, oldInsts)
		if err != nil {
			return fmt.Errorf("unable to prepend new device filters to the original device filters program: %v", err)
		}

		// Generate new eBPF program for the merged device filter instructions.
		spec := &ebpf.ProgramSpec{
			Type:         oldProg.Type(),
			Instructions: newInsts,
			License:      oldLicense,
		}
		newProg, err := ebpf.NewProgram(spec)
		if err != nil {
			return fmt.Errorf("unable to create new device filters program: %v", err)
		}

		// Replace the eBPF program with the merged device filter.
		_, err = ReplaceCgroupDeviceFilter(newProg, oldProg, dirFD)
		if err != nil {
			return fmt.Errorf("unable to replace original device filters program with new one: %v", err)
		}
	}

	return nil
}
