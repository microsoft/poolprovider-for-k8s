package main

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"

	v1alpha1 "github.com/microsoft/k8s-poolprovider/pkg/apis/dev/v1alpha1"
	v1controller "github.com/microsoft/k8s-poolprovider/pkg/controller/azurepipelinespool"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"

	"context"
	"time"

	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var azurepipelinepoolcr *v1alpha1.AzurePipelinesPool
var (
	name      = "azurepipelinepool-operator"
	namespace = "azuredevops"
)

func SetupCustomResource() {

	// create custom resource
	azurepipelinepoolcr = &v1alpha1.AzurePipelinesPool{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: v1alpha1.AzurePipelinesPoolSpec{
			ControllerName:       "prebansa/webserverimage",
			BuildkitReplicaCount: 1,
			AgentPools: []v1alpha1.AgentPoolSpec{
				{
					PoolName: "linux",
					PoolSpec: &corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  "vsts-agent",
								Image: "prebansa/myagent:v5.16",
							},
						},
					},
				},
			},
			Initialized: true,
		},
	}

	SetTestingEnvironmentVariables()

	s := scheme.Scheme
	s.AddKnownTypes(v1alpha1.SchemeGroupVersion, azurepipelinepoolcr)
	v1alpha1.SetClient(s)
}

func TestControllerMustCreateExternalResources(t *testing.T) {

	SetupCustomResource()
	objs := []runtime.Object{
		azurepipelinepoolcr,
	}

	s := scheme.Scheme

	// Create a fake client to mock API calls.
	cl := fake.NewFakeClient(objs...)
	v1alpha1.SetClient(s)

	r := &v1controller.ReconcileAzurePipelinesPool{Client: cl, Scheme: s}

	// Mock request to simulate Reconcile() being called on an event for a
	// watched resource .
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      name,
			Namespace: namespace,
		},
	}
	for i := 0; i < 5; i++ {
		res, err := r.Reconcile(req)
		if err != nil {
			t.Fatalf("reconcile: (%v)", err)
		}
		if res != (reconcile.Result{}) {
			t.Error("reconcile did not return an empty Result")
		}
	}

	// Check the deployment is created
	expectedDeployment := v1controller.AddnewDeploymentForCR(azurepipelinepoolcr)
	deployment := &appsv1.Deployment{}
	err := cl.Get(context.TODO(), types.NamespacedName{Name: expectedDeployment.Name, Namespace: expectedDeployment.Namespace}, deployment)
	if err != nil {
		t.Fatalf("get deployment failed: (%v)", err)
	}

	expectedBuildkitPod := v1controller.AddnewBuildkitPodForCR(azurepipelinepoolcr)
	buildkitpod := &appsv1.StatefulSet{}
	err = cl.Get(context.TODO(), types.NamespacedName{Name: expectedBuildkitPod.Name, Namespace: expectedBuildkitPod.Namespace}, buildkitpod)
	if err != nil {
		t.Fatalf("get buildkit pod failed: (%v)", err)
	}

	expectedService := v1controller.AddnewServiceForCR(azurepipelinepoolcr)
	svc := &corev1.Service{}
	err = cl.Get(context.TODO(), types.NamespacedName{Name: expectedService.Name, Namespace: expectedService.Namespace}, svc)
	if err != nil {
		t.Fatalf("get service failed: (%v)", err)
	}

	expectedBuildkitService := v1controller.AddnewBuildkitServiceForCR(azurepipelinepoolcr)
	buildkitsvc := &corev1.Service{}
	err = cl.Get(context.TODO(), types.NamespacedName{Name: expectedBuildkitService.Name, Namespace: expectedBuildkitService.Namespace}, buildkitsvc)
	if err != nil {
		t.Fatalf("get buildkit service failed: (%v)", err)
	}

	expectedMap := v1controller.AddnewConfigMapForCR(azurepipelinepoolcr)
	configmap := &corev1.ConfigMap{}
	err = cl.Get(context.TODO(), types.NamespacedName{Name: expectedMap.Name, Namespace: expectedMap.Namespace}, configmap)
	if err != nil {
		t.Fatalf("get pod: (%v)", err)
	}
}

func TestControllerMustRecreatePodIfDeleted(t *testing.T) {
	SetupCustomResource()
	objs := []runtime.Object{
		azurepipelinepoolcr,
	}

	s := scheme.Scheme

	// Create a fake client to mock API calls.
	cl := fake.NewFakeClient(objs...)
	v1alpha1.SetClient(s)

	r := &v1controller.ReconcileAzurePipelinesPool{Client: cl, Scheme: s}

	// Mock request to simulate Reconcile() being called on an event for a
	// watched resource .
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      name,
			Namespace: namespace,
		},
	}
	for i := 0; i < 5; i++ {
		res, err := r.Reconcile(req)
		if err != nil {
			t.Fatalf("reconcile: (%v)", err)
		}
		if res != (reconcile.Result{}) {
			t.Error("reconcile did not return an empty Result")
		}
	}

	expectedDeployment := v1controller.AddnewDeploymentForCR(azurepipelinepoolcr)
	deployment := &appsv1.Deployment{}
	err := cl.Get(context.TODO(), types.NamespacedName{Name: expectedDeployment.Name, Namespace: expectedDeployment.Namespace}, deployment)
	if err != nil {
		t.Fatalf("get deployment: (%v)", err)
	}

	//delete deployment
	errd := cl.Delete(context.TODO(), deployment)
	if errd != nil {
		t.Fatalf("delete pod failed: (%v)", errd)
	}

	deletereq := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      name,
			Namespace: namespace,
		},
	}

	res, err := r.Reconcile(deletereq)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}
	if res != (reconcile.Result{}) {
		t.Error("reconcile did not return an empty Result")
	}

	time.Sleep(time.Second * 1)
	recreateddeployment := &appsv1.Deployment{}
	errrecreate := cl.Get(context.TODO(), types.NamespacedName{Name: expectedDeployment.Name, Namespace: expectedDeployment.Namespace}, recreateddeployment)
	if errrecreate != nil {
		t.Fatalf("get pod not restarted: (%v)", errrecreate)
	}
}

func TestControllerMustRecreateStatefulsetIfDeleted(t *testing.T) {
	SetupCustomResource()
	objs := []runtime.Object{
		azurepipelinepoolcr,
	}

	s := scheme.Scheme

	// Create a fake client to mock API calls.
	cl := fake.NewFakeClient(objs...)
	v1alpha1.SetClient(s)

	r := &v1controller.ReconcileAzurePipelinesPool{Client: cl, Scheme: s}

	// Mock request to simulate Reconcile() being called on an event for a
	// watched resource .
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      name,
			Namespace: namespace,
		},
	}
	for i := 0; i < 5; i++ {
		res, err := r.Reconcile(req)
		if err != nil {
			t.Fatalf("reconcile: (%v)", err)
		}
		if res != (reconcile.Result{}) {
			t.Error("reconcile did not return an empty Result")
		}
	}

	expectedPod1 := v1controller.AddnewBuildkitPodForCR(azurepipelinepoolcr)
	statefulset := &appsv1.StatefulSet{}
	err := cl.Get(context.TODO(), types.NamespacedName{Name: expectedPod1.Name, Namespace: expectedPod1.Namespace}, statefulset)
	if err != nil {
		t.Fatalf("get pod: (%v)", err)
	}

	//delete statefulset
	errd := cl.Delete(context.TODO(), statefulset)
	if errd != nil {
		t.Fatalf("delete pod failed: (%v)", errd)
	}

	res, err := r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}
	if res != (reconcile.Result{}) {
		t.Error("reconcile did not return an empty Result")
	}

	newstatefulset := &appsv1.StatefulSet{}
	err1 := cl.Get(context.TODO(), types.NamespacedName{Name: expectedPod1.Name, Namespace: expectedPod1.Namespace}, newstatefulset)
	if err1 != nil {
		t.Fatalf("get pod not restarted: (%v)", err1)
	}
}
