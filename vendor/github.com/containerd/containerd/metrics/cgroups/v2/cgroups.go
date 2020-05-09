// +build linux

/*
   Copyright The containerd Authors.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package v2

import (
	"context"

	"github.com/containerd/containerd/errdefs"
	"github.com/containerd/containerd/events"
	"github.com/containerd/containerd/runtime"
	"github.com/containerd/containerd/runtime/v1/linux"
	metrics "github.com/docker/go-metrics"
)

// NewTaskMonitor returns a new cgroups monitor
func NewTaskMonitor(ctx context.Context, publisher events.Publisher, ns *metrics.Namespace) (runtime.TaskMonitor, error) {
	collector := NewCollector(ns)
	return &cgroupsMonitor{
		collector: collector,
		context:   ctx,
		publisher: publisher,
	}, nil
}

type cgroupsMonitor struct {
	collector *Collector
	context   context.Context
	publisher events.Publisher
}

func (m *cgroupsMonitor) Monitor(c runtime.Task) error {
	if err := m.collector.Add(c); err != nil {
		return err
	}
	t, ok := c.(*linux.Task)
	if !ok {
		return nil
	}
	cg, err := t.Cgroup()
	if err != nil {
		if errdefs.IsNotFound(err) {
			return nil
		}
		return err
	}
	// OOM handler is not implemented yet
	_ = cg
	return nil
}

func (m *cgroupsMonitor) Stop(c runtime.Task) error {
	m.collector.Remove(c)
	return nil
}
