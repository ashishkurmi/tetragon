// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Tetragon

package syscallmetrics

import (
	"github.com/cilium/tetragon/api/v1/tetragon"
	"github.com/cilium/tetragon/pkg/metrics"
	"github.com/cilium/tetragon/pkg/metrics/consts"
	"github.com/cilium/tetragon/pkg/syscallinfo"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	syscallStats = metrics.NewCounterVecWithPod(prometheus.CounterOpts{
		Namespace:   consts.MetricsNamespace,
		Name:        "syscalls_total",
		Help:        "System calls observed.",
		ConstLabels: nil,
	}, []string{"syscall", "namespace", "workload", "pod", "binary"})
)

func InitMetrics(registry *prometheus.Registry) {
	registry.MustRegister(syscallStats)
}

func Handle(event interface{}) {
	ev, ok := event.(*tetragon.GetEventsResponse)
	if !ok {
		return
	}

	var syscall, namespace, workload, pod, binary string
	if tpEvent := ev.GetProcessTracepoint(); tpEvent != nil {
		if tpEvent.Subsys == "raw_syscalls" && tpEvent.Event == "sys_enter" {
			syscall = rawSyscallName(tpEvent)
			if tpEvent.Process != nil {
				if tpEvent.Process.Pod != nil {
					namespace = tpEvent.Process.Pod.Namespace
					workload = tpEvent.Process.Pod.Workload
					pod = tpEvent.Process.Pod.Name
				}
				binary = tpEvent.Process.Binary
			}
		}
	}

	if syscall != "" {
		syscallStats.WithLabelValues(syscall, namespace, workload, pod, binary).Inc()
	}
}

func rawSyscallName(tp *tetragon.ProcessTracepoint) string {
	sysID := int64(-1)
	if len(tp.Args) > 0 && tp.Args[0] != nil {
		if x, ok := tp.Args[0].GetArg().(*tetragon.KprobeArgument_LongArg); ok {
			sysID = x.LongArg
		}
	}
	if sysID == -1 {
		return ""
	}
	return syscallinfo.GetSyscallName(int(sysID))
}
