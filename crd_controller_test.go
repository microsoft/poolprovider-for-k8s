package main

import (
	"testing"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/runtime"

	v1alpha1 "github.com/microsoft/k8s-poolprovider/pkg/apis/dev/v1alpha1"
	v1controller "github.com/microsoft/k8s-poolprovider/pkg/controller/azurepipelinespool"
	corev1 "k8s.io/api/core/v1"
	appsv1 "k8s.io/api/apps/v1"
	
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"context"
	"time"

)

var azurepipelinepoolcr *v1alpha1.AzurePipelinesPool
var (
	name      = "azurepipelinepool-operator"
	namespace = "azuredevops"
)

func SetupCustomResource(){
	
	// create custom resource
	azurepipelinepoolcr = &v1alpha1.AzurePipelinesPool{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec:  v1alpha1.AzurePipelinesPoolSpec{
			ControllerName: "prebansa/webserverimage",
			BuildkitReplicaCount: 1,
			AgentPools: []v1alpha1.AgentPoolSpec {
				{
					PoolName: "linux",
					PoolSpec: &corev1.PodSpec {
						Containers: []corev1.Container {
							{
								Name:   "vsts-agent",
								Image:  "prebansa/myagent:v5.16",
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
	objs := []runtime.Object {
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
    for i:= 0; i < 5; i++ {
		res, err := r.Reconcile(req)
		if err != nil {
			t.Fatalf("reconcile: (%v)", err)
		}
		if res != (reconcile.Result{}) {
			t.Error("reconcile did not return an empty Result")
		}
	}
    

	// Check the pod is created
	expectedPod := v1controller.AddnewPodForCR(azurepipelinepoolcr)
	pod := &corev1.Pod{}
	err := cl.Get(context.TODO(), types.NamespacedName{Name: expectedPod.Name, Namespace: expectedPod.Namespace}, pod)
	if err != nil {
		t.Fatalf("get pod: (%v)", err)
	}

	expectedPod1 := v1controller.AddnewBuildkitPodForCR(azurepipelinepoolcr)
	pod1 := &appsv1.StatefulSet{}
	err = cl.Get(context.TODO(), types.NamespacedName{Name: expectedPod1.Name, Namespace: expectedPod1.Namespace}, pod1)
	if err != nil {
		t.Fatalf("get pod: (%v)", err)
	}

	expectedService := v1controller.AddnewServiceForCR(azurepipelinepoolcr)
	svc := &corev1.Service{}
	err = cl.Get(context.TODO(), types.NamespacedName{Name: expectedService.Name, Namespace: expectedService.Namespace}, svc)
	if err != nil {
		t.Fatalf("get pod: (%v)", err)
	}

	expectedService1 := v1controller.AddnewBuildkitServiceForCR(azurepipelinepoolcr)
	svc1 := &corev1.Service{}
	err = cl.Get(context.TODO(), types.NamespacedName{Name: expectedService1.Name, Namespace: expectedService.Namespace}, svc1)
	if err != nil {
		t.Fatalf("get pod: (%v)", err)
	}

	expectedMap := v1controller.AddnewConfigMapForCR(azurepipelinepoolcr)
	map1 := &corev1.ConfigMap{}
	err = cl.Get(context.TODO(), types.NamespacedName{Name: expectedMap.Name, Namespace: expectedMap.Namespace},map1)
	if err != nil {
		t.Fatalf("get pod: (%v)", err)
	}
}

func TestControllerMustRecreatePodIfDeleted(t *testing.T) {
	SetupCustomResource()
	objs := []runtime.Object {
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
    for i:= 0; i < 5; i++ {
		res, err := r.Reconcile(req)
		if err != nil {
			t.Fatalf("reconcile: (%v)", err)
		}
		if res != (reconcile.Result{}) {
			t.Error("reconcile did not return an empty Result")
		}
	}

	expectedPod := v1controller.AddnewPodForCR(azurepipelinepoolcr)
	pod := &corev1.Pod{}
	err := cl.Get(context.TODO(), types.NamespacedName{Name: expectedPod.Name, Namespace: expectedPod.Namespace}, pod)
	if err != nil {
		t.Fatalf("get pod: (%v)", err)
	}

	//prevpodUID := &pod.UID

	//delete pod
	errd := cl.Delete(context.TODO(), pod)
	if errd != nil {
		t.Fatalf("delete pod failed: (%v)", errd)
	}

	req1 := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      name,
			Namespace: namespace,
		},
	}

	res, err := r.Reconcile(req1)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}
	if res != (reconcile.Result{}) {
		t.Error("reconcile did not return an empty Result")
	}

	time.Sleep(time.Second*1)
	pod1 := &corev1.Pod{}
	err1 := cl.Get(context.TODO(), types.NamespacedName{Name: expectedPod.Name, Namespace: expectedPod.Namespace}, pod1)
	if err1 != nil {
		t.Fatalf("get pod not restarted: (%v)", err1)
	}
}

func TestControllerMustRecreateStatefulsetIfDeleted(t *testing.T) {
	SetupCustomResource()
	objs := []runtime.Object {
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
    for i:= 0; i < 5; i++ {
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