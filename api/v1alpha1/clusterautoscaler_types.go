/*
Copyright 2025.

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

package v1alpha1

import (
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/vanchonlee/oscale/internal/pkg/schedule"
	"github.com/vanchonlee/oscale/internal/pkg/support"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ClusterAutoscalerSpec defines the desired state of ClusterAutoscaler.
type ClusterAutoscalerSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of ClusterAutoscaler. Edit clusterautoscaler_types.go to remove/update
	// DatajetCluster specify the Datajet cluster name, which is used to identify the cluster
	// provider is one of aws, k8s
	Provider string `json:"provider"`
	// possible values is "use1", "apse1", "apse2", "euw1", "staging"
	DatajetCluster string `json:"datajetCluster,omitempty"`
	// DomainName specify the OpenSearch domain name from AWS panel, like "feeds"
	DomainName string `json:"domainName"`
	// NameSpace specify the namespace of the OpenSearch cluster, if provider is k8s
	NameSpace string `json:"nameSpace,omitempty"`
	// ScalingEnabled specify if the scaling is enabled or not, if false -- will be dry run
	ScalingEnabled *bool `json:"scalingEnabled"`

	// Interval for calculating CPU utilization
	Interval duration.Duration `json:"interval"`
	// TargetCPUUtilization specify the target CPU utilization for the OpenSearch cluster
	TargetCPUUtilization int32 `json:"targetCPUUtilization"`

	// ScaleUpStep specify the minimum number of nodes to scale up. For example if step is 4,
	// if desired number of nodes is increased by 2, then the number of nodes will be increased by 4,
	// if desired number of nodes is increased by 5, then the number of nodes will be increased by 8
	ScaleUpStep int32 `json:"scaleUpStep"`
	// ScaleDownStep specify the minimum number of nodes to scale down. For example if step is 4,
	// if desired number of nodes is decreased by 2, then the number of nodes will be decreased by 0,
	// if desired number of nodes is decreased by 5, then the number of nodes will be decreased by 4,
	// if desired number of nodes is decreased by 9, then the number of nodes will be decreased by 8
	ScaleDownStep int32 `json:"scaleDownStep"`
	// UpscaleStabilizationWindow specify the duration for which the cluster should be stable before scaling up
	UpscaleStabilizationWindow duration.Duration `json:"upscaleStabilizationWindow,omitempty"`
	// DownscaleStabilizationWindow specify the duration for which the cluster should be stable before scaling down
	DownscaleStabilizationWindow duration.Duration `json:"downscaleStabilizationWindow,omitempty"`
	// EvenOnly specify if the number of nodes should be only even
	EvenOnly *bool `json:"evenOnly"`
	// MinDataNodes specify the minimum number of data nodes in the OpenSearch cluster
	MinDataNodes int32 `json:"minDataNodes"`
	// MaxDataNodes specify the maximum number of data nodes in the OpenSearch cluster
	MaxDataNodes int32 `json:"maxDataNodes"`
	// MinDataNodesSchedule is a schedule for the minimum number of data nodes
	MinDataNodesSchedule schedule.Schedule `json:"minDataNodesSchedule,omitempty"`
	// Supports allows to calculate the minimum number of nodes depends on historical usage
	Supports support.Support `json:"support,omitempty"`
}

// ClusterAutoscalerStatus defines the observed state of ClusterAutoscaler.
type ClusterAutoscalerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// ClusterAutoscaler is the Schema for the clusterautoscalers API.
type ClusterAutoscaler struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterAutoscalerSpec   `json:"spec,omitempty"`
	Status ClusterAutoscalerStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ClusterAutoscalerList contains a list of ClusterAutoscaler.
type ClusterAutoscalerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterAutoscaler `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ClusterAutoscaler{}, &ClusterAutoscalerList{})
}
