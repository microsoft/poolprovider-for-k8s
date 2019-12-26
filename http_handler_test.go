package main

import (
	"testing"
	"bytes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/runtime"
	"net/http"
	"net/http/httptest"
	"encoding/json"
	v1alpha1 "github.com/microsoft/k8s-poolprovider/pkg/apis/dev/v1alpha1"
	v1controller "github.com/microsoft/k8s-poolprovider/pkg/controller/azurepipelinespool"
	corev1 "k8s.io/api/core/v1"
	appsv1 "k8s.io/api/apps/v1"
	
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"context"

)

func TestAcquireHandlerShouldBeSuccessful(t *testing.T) {

	var response AgentProvisionResponse
	var jsonStr = []byte(`{"AgentId":"1"}`)
	
	SetupCustomResource()
	req, _ := http.NewRequest("POST", "/acquire", bytes.NewBuffer(jsonStr))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Azure-Signature", "4f6a97c5aa13477ed775dd20cdd7cf44477e310ba683545144d5112e77d88be967a43791da0696ded702a32ca0c190ab831dcd9204521b9a9ebe413066699ef9")

	resp := httptest.NewRecorder()
	http.HandlerFunc(AcquireAgentHandler).ServeHTTP(resp,req)

	if status := resp.Code; status != http.StatusCreated { //Must be 201
		t.Errorf("Status code differs. Expected %d. Got %d", http.StatusCreated, status)
	} else {

		cs := CreateClientSet()
		podClient := cs.clientset.CoreV1().Pods("azuredevops")
		pods, _ := podClient.List(metav1.ListOptions{LabelSelector: agentIdLabel + "=" + "1"})

		if pods == nil || len(pods.Items) == 0 {
			t.Errorf("Http Aquire Call failed")
		} else {

			json.Unmarshal([]byte(resp.Body.String()), &response)

			if response.Accepted != true {
				t.Errorf("Http Aquire Call failed")
			} else if response.ResponseType != "Success" {
				t.Errorf("Http Aquire Call failed")
			} else if response.ErrorMessage != "" {
				t.Errorf("Http Aquire Call failed")
			}
		}
	}
}

func TestAcquireHandlerShouldFailIfGetRequest(t *testing.T) {
	SetupCustomResource()

	var jsonStr = []byte(`{"AgentId":"1"}`)

	req, _ := http.NewRequest("GET", "/acquire", bytes.NewBuffer(jsonStr))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Azure-Signature", "4f6a97c5aa13477ed775dd20cdd7cf44477e310ba683545144d5112e77d88be967a43791da0696ded702a32ca0c190ab831dcd9204521b9a9ebe413066699ef9")

	resp := httptest.NewRecorder()
	http.HandlerFunc(AcquireAgentHandler).ServeHTTP(resp,req)

	if status := resp.Code; status != http.StatusMethodNotAllowed { //Must be 405
		t.Errorf("Status code differs. Expected %d. Got %d", http.StatusMethodNotAllowed, status)
	}
}

func TestAcquireHandlerShouldFailIfHmacNotValid(t *testing.T) {
	SetupCustomResource()

	var jsonStr = []byte(`{"AgentId":"12"}`)
	
	req, _ := http.NewRequest("POST", "/acquire", bytes.NewBuffer(jsonStr))
	req.Header.Add("Content-Type", "application/json")

	// wrong encoding
	req.Header.Add("X-Azure-Signature", "4f6a97c5aa13477ed775dd20cdd7cf44477e310ba683545144d5112e77d88be967a43791da0696ded702a32ca0c190ab831dcd9204521b9a9ebe413066699ef9")

	resp := httptest.NewRecorder()
	http.HandlerFunc(AcquireAgentHandler).ServeHTTP(resp,req)

	if status := resp.Code; status != http.StatusForbidden { //Must be 403
		t.Errorf("Status code differs. Expected %d. Got %d", http.StatusForbidden, status)
	}
}

func SetupCustomResource(){
	var (
		//name      = "azurepipelinepool-operator"
		namespace = "azuredevops"
	)
	// create custom resource
	azurepipelinepoolcr := &v1alpha1.AzurePipelinesPool{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "azurepipelinepool-operator",
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

func TestK8PoolProviderCreate(t *testing.T) {

	var (
		name      = "azurepipelinepool-operator"
		namespace = "azuredevops"
	)
	// create custom resource
	azurepipelinepoolcr := &v1alpha1.AzurePipelinesPool{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "azurepipelinepool-operator",
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
	objs := []runtime.Object {
		azurepipelinepoolcr,
	}

	s := scheme.Scheme
	s.AddKnownTypes(v1alpha1.SchemeGroupVersion, azurepipelinepoolcr)
	// Create a fake client to mock API calls.
	cl := fake.NewFakeClient(objs...)
	v1alpha1.SetClient(s)
//	SetClientSet(s)
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
	t.Log("Pod-----------")
	t.Log(pod)

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
	t.Log("Called from here.............")
	
	//AcquireHandlerShouldBeSuccessful(t)
	//AcquireHandlerShouldFailIfGetRequest(t)
	//AcquireHandlerShouldFailIfHmacNotValid(t)
	//ReleaseHandlerShouldBeSuccessful(t)
	//(t)
	//return nil
}

func TestReleaseHandlerShouldBeSuccessful(t *testing.T) {
	SetupCustomResource()
	var agentrequest AgentRequest
	agentrequest.AgentId = "1"
	testPod := CreatePod(agentrequest, "azuredevops")

	if (testPod.Accepted != true){
		t.Errorf("Pod creation failed")
	}

	var response PodResponse
	var jsonStr = []byte(`{"AgentId":"1"}`)
	
	req, _ := http.NewRequest("POST", "/release", bytes.NewBuffer(jsonStr))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Azure-Signature", "4f6a97c5aa13477ed775dd20cdd7cf44477e310ba683545144d5112e77d88be967a43791da0696ded702a32ca0c190ab831dcd9204521b9a9ebe413066699ef9")

	resp := httptest.NewRecorder()
	http.HandlerFunc(ReleaseAgentHandler).ServeHTTP(resp,req)

	if status := resp.Code; status != http.StatusCreated { //Must be 201
		t.Errorf("Status code differs. Expected %d. Got %d", http.StatusCreated, status)
	} else {

		// Now check the pod is deleted or not
		cs := CreateClientSet()
		podClient := cs.clientset.CoreV1().Pods("azuredevops")
		pods, _ := podClient.List(metav1.ListOptions{LabelSelector: agentIdLabel + "=" + "1"})

		if pods == nil || len(pods.Items) == 1 {
			t.Errorf("Http Release Call failed")
		} else {
			json.Unmarshal([]byte(resp.Body.String()), &response)

			if response.Status != "success" {
				t.Errorf("Http Release Call failed")
			} else if response.Message == "" {
				t.Errorf("Http Release Call failed")
			}
		}
	}
}

func TestGetBuildPodHandlerShouldBeSuccessful(t *testing.T) {
	SetTestingEnvironmentVariables()
	CreateDummyBuildKitPod()

	var response PodResponse
	var jsonStr = []byte("")
	req, _ := http.NewRequest("GET", "/buildPod", bytes.NewBuffer(jsonStr))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Azure-Signature", "4f6a97c5aa13477ed775dd20cdd7cf44477e310ba683545144d5112e77d88be967a43791da0696ded702a32ca0c190ab831dcd9204521b9a9ebe413066699ef9")

	resp := httptest.NewRecorder()
	http.HandlerFunc(GetBuildPodHandler).ServeHTTP(resp,req)

	if status := resp.Code; status != http.StatusCreated { //Must be 201
		t.Errorf("Status code differs. Expected %d. Got %d", http.StatusCreated, status)
	} else {
			json.Unmarshal([]byte(resp.Body.String()), &response)

			if response.Status != "success" {
				t.Errorf("Http buildPod Call failed")
			} else if response.Message != "buildkitd-0" {
				t.Errorf("Http buildPod Call failed")
			}
	}
}