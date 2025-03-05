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

package controller

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	opensearchv1alpha1 "github.com/vanchonlee/oscale/api/v1alpha1"
	"github.com/vanchonlee/oscale/internal/pkg/duration"
	"github.com/vanchonlee/oscale/internal/pkg/schedule"
)

var _ = Describe("ClusterAutoscaler Controller", func() {
	Context("When reconciling a resource", func() {
		const resourceName = "test-resource"

		ctx := context.Background()

		typeNamespacedName := types.NamespacedName{
			Name:      resourceName,
			Namespace: "default", // TODO(user):Modify as needed
		}
		clusterautoscaler := &opensearchv1alpha1.ClusterAutoscaler{}

		BeforeEach(func() {
			By("creating the custom resource for the Kind ClusterAutoscaler")
			err := k8sClient.Get(ctx, typeNamespacedName, clusterautoscaler)
			if err != nil && errors.IsNotFound(err) {
				truePtr := true
				resource := &opensearchv1alpha1.ClusterAutoscaler{
					ObjectMeta: metav1.ObjectMeta{
						Name:      resourceName,
						Namespace: "default",
					},
					Spec: opensearchv1alpha1.ClusterAutoscalerSpec{
						Provider:             "aws",
						DomainName:           "test",
						TargetCPUUtilization: 80,
						ScaleUpStep:          1,
						ScaleDownStep:        1,
						ScalingEnabled:       &truePtr,
						Interval:             duration.Duration{DurationStr: "1m"},
						UpscaleStabilizationWindow: duration.Duration{
							DurationStr: "10m",
						},
						DownscaleStabilizationWindow: duration.Duration{
							DurationStr: "10m",
						},
						EvenOnly:     &truePtr,
						MinDataNodes: 1,
						MaxDataNodes: 10,
						MinDataNodesSchedule: schedule.Schedule{
							Entities: []schedule.Entity{
								{
									CronStart: "0 0 * * *",
									CronEnd:   "0 0 * * *",
									Count:     1,
								},
							},
						},
					},
				}
				Expect(k8sClient.Create(ctx, resource)).To(Succeed())
			}
		})

		AfterEach(func() {
			// TODO(user): Cleanup logic after each test, like removing the resource instance.
			resource := &opensearchv1alpha1.ClusterAutoscaler{}
			err := k8sClient.Get(ctx, typeNamespacedName, resource)
			Expect(err).NotTo(HaveOccurred())

			By("Cleanup the specific resource instance ClusterAutoscaler")
			Expect(k8sClient.Delete(ctx, resource)).To(Succeed())
		})
		It("should successfully reconcile the resource", func() {
			By("Reconciling the created resource")
			controllerReconciler := &ClusterAutoscalerReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
			}

			_, err := controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())
			// TODO(user): Add more specific assertions depending on your controller's reconciliation logic.
			// Example: If you expect a certain status condition after reconciliation, verify it here.
		})
	})
})
