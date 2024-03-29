/**
# Copyright (c) NVIDIA CORPORATION.  All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
**/

package devchar

import (
	"path/filepath"
	"strings"

	"github.com/NVIDIA/nvidia-container-toolkit/internal/lookup"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
)

type nodeLister interface {
	DeviceNodes() ([]deviceNode, error)
}

type existing struct {
	logger     *logrus.Logger
	driverRoot string
}

// DeviceNodes returns a list of NVIDIA device nodes in the specified root.
// The nvidia-nvswitch* and nvidia-nvlink devices are excluded.
func (m existing) DeviceNodes() ([]deviceNode, error) {
	locator := lookup.NewCharDeviceLocator(
		lookup.WithLogger(m.logger),
		lookup.WithRoot(m.driverRoot),
		lookup.WithOptional(true),
	)

	devices, err := locator.Locate("/dev/nvidia*")
	if err != nil {
		m.logger.Warnf("Error while locating device: %v", err)
	}

	capDevices, err := locator.Locate("/dev/nvidia-caps/nvidia-*")
	if err != nil {
		m.logger.Warnf("Error while locating caps device: %v", err)
	}

	if len(devices) == 0 && len(capDevices) == 0 {
		m.logger.Infof("No NVIDIA devices found in %s", m.driverRoot)
		return nil, nil
	}

	var deviceNodes []deviceNode
	for _, d := range append(devices, capDevices...) {
		if m.nodeIsBlocked(d) {
			continue
		}

		var stat unix.Stat_t
		err := unix.Stat(d, &stat)
		if err != nil {
			m.logger.Warnf("Could not stat device: %v", err)
			continue
		}
		deviceNode := deviceNode{
			path:  d,
			major: unix.Major(uint64(stat.Rdev)),
			minor: unix.Minor(uint64(stat.Rdev)),
		}

		deviceNodes = append(deviceNodes, deviceNode)
	}

	return deviceNodes, nil
}

// nodeIsBlocked returns true if the specified device node should be ignored.
func (m existing) nodeIsBlocked(path string) bool {
	blockedPrefixes := []string{"nvidia-fs", "nvidia-nvswitch", "nvidia-nvlink"}
	nodeName := filepath.Base(path)
	for _, prefix := range blockedPrefixes {
		if strings.HasPrefix(nodeName, prefix) {
			return true
		}
	}
	return false
}
