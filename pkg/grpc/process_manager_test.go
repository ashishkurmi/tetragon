// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package grpc

import (
	"context"
	"encoding/base64"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/cilium/tetragon/pkg/grpc/exec"
	"github.com/cilium/tetragon/pkg/option"

	"github.com/cilium/tetragon/api/v1/tetragon"
	"github.com/cilium/tetragon/pkg/api/processapi"
	"github.com/cilium/tetragon/pkg/cilium"
	"github.com/cilium/tetragon/pkg/process"
	"github.com/cilium/tetragon/pkg/reader/node"
	"github.com/cilium/tetragon/pkg/rthooks"
	"github.com/cilium/tetragon/pkg/watcher"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestProcessManager_getPodInfo(t *testing.T) {
	controller := true
	podA := corev1.Pod{
		ObjectMeta: v1.ObjectMeta{
			Name:         "pod-a",
			Namespace:    "namespace-a",
			GenerateName: "test-workload-",
			OwnerReferences: []v1.OwnerReference{
				{
					Name:       "test-workload",
					Controller: &controller,
				},
			},
		},
		Status: corev1.PodStatus{
			ContainerStatuses: []corev1.ContainerStatus{
				{
					Name:        "pod-a-container-a-name",
					Image:       "pod-a-image-a-name",
					ImageID:     "pod-a-image-a-id",
					ContainerID: "docker://aaaaaaaaaaaaaaa",
					State: corev1.ContainerState{
						Running: &corev1.ContainerStateRunning{
							StartedAt: v1.Time{
								Time: time.Unix(1, 2),
							},
						},
					},
				},
			},
		},
	}

	_, err := cilium.InitCiliumState(context.Background(), false)
	assert.NoError(t, err)
	pods := []interface{}{&podA}
	err = process.InitCache(watcher.NewFakeK8sWatcher(pods), 10)
	assert.NoError(t, err)
	defer process.FreeCache()
	pod := process.GetPodInfo("container-id-not-found", "", "", 0)
	assert.Nil(t, pod)
	pod = process.GetPodInfo("aaaaaaa", "", "", 1234)
	assert.Equal(t,
		&tetragon.Pod{
			Namespace: podA.Namespace,
			Workload:  podA.OwnerReferences[0].Name,
			Name:      podA.Name,
			Container: &tetragon.Container{
				Id:   podA.Status.ContainerStatuses[0].ContainerID,
				Name: podA.Status.ContainerStatuses[0].Name,
				Image: &tetragon.Image{
					Id:   podA.Status.ContainerStatuses[0].ImageID,
					Name: podA.Status.ContainerStatuses[0].Image,
				},
				StartTime: &timestamppb.Timestamp{
					Seconds: int64(podA.Status.ContainerStatuses[0].State.Running.StartedAt.Second()),
					Nanos:   int32(podA.Status.ContainerStatuses[0].State.Running.StartedAt.Nanosecond()),
				},
				Pid: &wrapperspb.UInt32Value{Value: 1234},
			},
		}, pod)
}

func TestProcessManager_getPodInfoMaybeExecProbe(t *testing.T) {
	controller := true
	var podA = corev1.Pod{
		ObjectMeta: v1.ObjectMeta{
			Name:         "pod-a",
			Namespace:    "namespace-a",
			GenerateName: "test-workload-",
			OwnerReferences: []v1.OwnerReference{
				{
					Name:       "test-workload",
					Controller: &controller,
				},
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name: "pod-a-container-a-name",
					LivenessProbe: &corev1.Probe{
						ProbeHandler: corev1.ProbeHandler{
							Exec: &corev1.ExecAction{
								Command: []string{"command", "arg-a", "arg-b"},
							},
						},
					},
				},
			},
		},
		Status: corev1.PodStatus{
			ContainerStatuses: []corev1.ContainerStatus{
				{
					Name:        "pod-a-container-a-name",
					ContainerID: "docker://aaaaaaaaaaaaaaa",
				},
			},
		},
	}
	_, err := cilium.InitCiliumState(context.Background(), false)
	assert.NoError(t, err)
	pods := []interface{}{&podA}
	err = process.InitCache(watcher.NewFakeK8sWatcher(pods), 10)
	assert.NoError(t, err)
	defer process.FreeCache()
	pod := process.GetPodInfo("aaaaaaa", "/bin/command", "arg-a arg-b", 1234)
	assert.Equal(t,
		&tetragon.Pod{
			Namespace: podA.Namespace,
			Workload:  podA.OwnerReferences[0].Name,
			Name:      podA.Name,
			Container: &tetragon.Container{
				Id:             podA.Status.ContainerStatuses[0].ContainerID,
				Name:           podA.Status.ContainerStatuses[0].Name,
				Image:          &tetragon.Image{},
				Pid:            &wrapperspb.UInt32Value{Value: 1234},
				MaybeExecProbe: true,
			},
		}, pod)
}

func TestProcessManager_GetProcessExec(t *testing.T) {
	_, err := cilium.InitCiliumState(context.Background(), false)
	assert.NoError(t, err)
	err = process.InitCache(watcher.NewFakeK8sWatcher(nil), 10)
	assert.NoError(t, err)
	defer process.FreeCache()
	var wg sync.WaitGroup

	option.Config.EnableProcessNs = false
	option.Config.EnableProcessCred = false
	option.Config.EnableCilium = false
	_, err = NewProcessManager(
		context.Background(),
		&wg,
		nil,
		&rthooks.Runner{})
	assert.NoError(t, err)
	pi := &exec.MsgExecveEventUnix{MsgExecveEventUnix: processapi.MsgExecveEventUnix{
		Common: processapi.MsgCommon{
			Ktime: 1234,
		},
		Capabilities: processapi.MsgCapabilities{
			Permitted:   1,
			Effective:   1,
			Inheritable: 1,
		},
		Process: processapi.MsgProcess{
			PID: 5678,
		},
	}}

	assert.Nil(t, exec.GetProcessExec(pi, false).Process.Cap)

	// cap field should be set with enable-process-cred flag.
	option.Config.EnableProcessCred = true
	assert.Equal(t,
		&tetragon.Capabilities{
			Permitted:   []tetragon.CapabilitiesType{tetragon.CapabilitiesType_CAP_CHOWN},
			Effective:   []tetragon.CapabilitiesType{tetragon.CapabilitiesType_CAP_CHOWN},
			Inheritable: []tetragon.CapabilitiesType{tetragon.CapabilitiesType_CAP_CHOWN},
		},
		exec.GetProcessExec(pi, false).Process.Cap)
}

func Test_getNodeNameForExport(t *testing.T) {
	assert.Equal(t, "", node.GetNodeNameForExport())
	assert.NoError(t, os.Setenv("NODE_NAME", "from-node-name"))
	assert.Equal(t, "from-node-name", node.GetNodeNameForExport())
	assert.NoError(t, os.Setenv("HUBBLE_NODE_NAME", "from-hubble-node-name"))
	assert.Equal(t, "from-hubble-node-name", node.GetNodeNameForExport())
	assert.NoError(t, os.Unsetenv("NODE_NAME"))
	assert.NoError(t, os.Unsetenv("HUBBLE_NODE_NAME"))
}

func TestProcessManager_GetProcessID(t *testing.T) {
	assert.NoError(t, os.Setenv("NODE_NAME", "my-node"))

	_, err := cilium.InitCiliumState(context.Background(), false)
	assert.NoError(t, err)

	err = process.InitCache(watcher.NewFakeK8sWatcher([]interface{}{}), 10)
	assert.NoError(t, err)
	defer process.FreeCache()
	id := process.GetProcessID(1, 2)
	decoded, err := base64.StdEncoding.DecodeString(id)
	assert.NoError(t, err)
	assert.Equal(t, "my-node:2:1", string(decoded))
	assert.NoError(t, os.Unsetenv("NODE_NAME"))
}
