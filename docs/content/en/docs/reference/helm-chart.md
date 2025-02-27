---
title: "Helm chart"
description: "This reference is generated from the Tetragon Helm chart values."
---

{{< comment >}}
This page was generated with github.io/cilium/tetragon/install/kubernetes/export-doc.sh,
please do not edit directly.
{{< /comment >}}

The Tetragon Helm chart source is available under 
[github.io/cilium/tetragon/install/kubernetes](https://github.com/cilium/tetragon/tree/main/install/kubernetes)
and is distributed from the Cilium helm charts repository [helm.cilium.io](https://helm.cilium.io).

To deploy Tetragon using this Helm chart you can run the following commands:
```shell-session
helm repo add cilium https://helm.cilium.io
helm repo update
helm install tetragon cilium/tetragon -n kube-system
```

To use [the values available](#values), with `helm install` or `helm upgrade`, use `--set key=value`.

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| affinity | object | `{}` |  |
| daemonSetAnnotations | object | `{}` |  |
| daemonSetLabelsOverride | object | `{}` |  |
| dnsPolicy | string | `"Default"` |  |
| enabled | bool | `true` |  |
| export.filenames[0] | string | `"tetragon.log"` |  |
| export.mode | string | `"stdout"` |  |
| export.resources | object | `{}` |  |
| export.securityContext | object | `{}` |  |
| export.stdout.argsOverride | list | `[]` |  |
| export.stdout.commandOverride | list | `[]` |  |
| export.stdout.enabledArgs | bool | `true` |  |
| export.stdout.enabledCommand | bool | `true` |  |
| export.stdout.extraEnv | list | `[]` |  |
| export.stdout.extraVolumeMounts | list | `[]` |  |
| export.stdout.image.override | string | `nil` |  |
| export.stdout.image.repository | string | `"quay.io/cilium/hubble-export-stdout"` |  |
| export.stdout.image.tag | string | `"v1.0.3"` |  |
| exportDirectory | string | `"/var/run/cilium/tetragon"` |  |
| exportFileCreationInterval | string | `"120s"` |  |
| extraConfigmapMounts | list | `[]` |  |
| extraHostPathMounts | list | `[]` |  |
| extraVolumes | list | `[]` |  |
| hostNetwork | bool | `true` |  |
| imagePullPolicy | string | `"IfNotPresent"` |  |
| imagePullSecrets | list | `[]` |  |
| nodeSelector | object | `{}` |  |
| podAnnotations | object | `{}` |  |
| podLabelsOverride | object | `{}` |  |
| podSecurityContext | object | `{}` |  |
| podWatcher.enabled | bool | `false` |  |
| selectorLabelsOverride | object | `{}` |  |
| serviceAccount.annotations | object | `{}` |  |
| serviceAccount.create | bool | `true` |  |
| serviceAccount.name | string | `""` |  |
| serviceLabelsOverride | object | `{}` |  |
| tetragon.argsOverride | list | `[]` |  |
| tetragon.btf | string | `""` |  |
| tetragon.commandOverride | list | `[]` |  |
| tetragon.enableCiliumAPI | bool | `false` |  |
| tetragon.enableK8sAPI | bool | `true` |  |
| tetragon.enableMsgHandlingLatency | bool | `false` |  |
| tetragon.enablePolicyFilter | bool | `false` |  |
| tetragon.enablePolicyFilterDebug | bool | `false` |  |
| tetragon.enableProcessCred | bool | `false` |  |
| tetragon.enableProcessNs | bool | `false` |  |
| tetragon.enabled | bool | `true` |  |
| tetragon.exportAllowList | string | `"{\"event_set\":[\"PROCESS_EXEC\", \"PROCESS_EXIT\", \"PROCESS_KPROBE\", \"PROCESS_UPROBE\"]}"` |  |
| tetragon.exportDenyList | string | `"{\"health_check\":true}\n{\"namespace\":[\"\", \"cilium\", \"kube-system\"]}"` |  |
| tetragon.exportFileCompress | bool | `false` |  |
| tetragon.exportFileMaxBackups | int | `5` |  |
| tetragon.exportFileMaxSizeMB | int | `10` |  |
| tetragon.exportFilename | string | `"tetragon.log"` |  |
| tetragon.exportRateLimit | int | `-1` |  |
| tetragon.extraArgs | object | `{}` |  |
| tetragon.extraEnv | list | `[]` |  |
| tetragon.extraVolumeMounts | list | `[]` |  |
| tetragon.fieldFilters | string | `"{}"` |  |
| tetragon.gops.address | string | `"localhost"` | The address at which to expose gops. |
| tetragon.gops.port | int | `8118` | The port at which to expose gops. |
| tetragon.grpc.address | string | `"localhost:54321"` | The address at which to expose gRPC. Examples: localhost:54321, unix:///var/run/tetragon/tetragon.sock |
| tetragon.grpc.enabled | bool | `true` | Whether to enable exposing Tetragon gRPC. |
| tetragon.image.override | string | `nil` |  |
| tetragon.image.repository | string | `"quay.io/cilium/tetragon"` |  |
| tetragon.image.tag | string | `"v0.11.0"` |  |
| tetragon.processCacheSize | int | `65536` |  |
| tetragon.prometheus.address | string | `""` | The address at which to expose metrics. Set it to "" to expose on all available interfaces. |
| tetragon.prometheus.enabled | bool | `true` | Whether to enable exposing Tetragon metrics. |
| tetragon.prometheus.port | int | `2112` | The port at which to expose metrics. |
| tetragon.prometheus.serviceMonitor.enabled | bool | `false` | Whether to create a 'ServiceMonitor' resource targeting the 'tetragon' pods. |
| tetragon.prometheus.serviceMonitor.labelsOverride | object | `{}` | The set of labels to place on the 'ServiceMonitor' resource. |
| tetragon.resources | object | `{}` |  |
| tetragon.securityContext.privileged | bool | `true` |  |
| tetragonOperator.enabled | bool | `true` | Enable the tetragon-operator component (required). |
| tetragonOperator.image | object | `{"override":null,"repository":"quay.io/cilium/tetragon-operator","suffix":"","tag":"v0.11.0"}` | tetragon-operator image. |
| tetragonOperator.skipCRDCreation | bool | `false` |  |
| tolerations[0].operator | string | `"Exists"` |  |
| updateStrategy | object | `{}` |  |

