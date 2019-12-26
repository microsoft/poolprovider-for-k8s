package azurepipelinespool

import (
	"context"

	devv1alpha1 "github.com/microsoft/k8s-poolprovider/pkg/apis/dev/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

const controllerName = "finalizer_azurepipelinespool"
var log = logf.Log.WithName("controller_azurepipelinespool")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new AzurePipelinesPool Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileAzurePipelinesPool{Client: mgr.GetClient(), Scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("azurepipelinespool-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource AzurePipelinesPool
	err = c.Watch(&source.Kind{Type: &devv1alpha1.AzurePipelinesPool{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &devv1alpha1.AzurePipelinesPool{},
	})
	if err != nil {
		return err
	}

    err = c.Watch(&source.Kind{Type: &corev1.ConfigMap{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &devv1alpha1.AzurePipelinesPool{},
	})
	if err != nil {
		return err
	}

    err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &devv1alpha1.AzurePipelinesPool{},
	})
	if err != nil {
		return err
	}

    err = c.Watch(&source.Kind{Type: &appsv1.StatefulSet{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &devv1alpha1.AzurePipelinesPool{},
	})
	if err != nil {
		return err
	}


	return nil
}

// blank assignment to verify that ReconcileAzurePipelinesPool implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileAzurePipelinesPool{}

// ReconcileAzurePipelinesPool reconciles a AzurePipelinesPool object
type ReconcileAzurePipelinesPool struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	Client client.Client
	Scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a AzurePipelinesPool object and makes changes based on the state read
// and what is in the AzurePipelinesPool.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileAzurePipelinesPool) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling AzurePipelinesPool")

	// Fetch the AzurePipelinesPool instance
	instance := &devv1alpha1.AzurePipelinesPool{}
	err := r.Client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	if ok := IsInitialized(instance); !ok {
		err := r.Client.Update(context.TODO(), instance)
		if err != nil {
			log.Error(err, "unable to update instance", "instance", instance)
		}
		return reconcile.Result{}, nil
	}

    if isBeingDeleted(instance) {
		if !hasFinalizer(instance, controllerName) {
			return reconcile.Result{}, nil
		}
		manageCleanUpLogic(instance)
		
		removeFinalizer(instance, controllerName)
		err = r.Client.Update(context.TODO(), instance)
		if err != nil {
			log.Error(err,"unable to update instance")
		}
		return reconcile.Result{}, nil
	}

	// Define a new Pod object
	pod := AddnewPodForCR(instance)

	// Set AzurePipelinePool instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, pod, r.Scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if this Pod already exists
	found := &corev1.Pod{}
	err = r.Client.Get(context.TODO(), types.NamespacedName{Name: pod.Name, Namespace: pod.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Pod", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)
		err = r.Client.Create(context.TODO(), pod)
		if err != nil {
			return reconcile.Result{}, err
		}

		// Pod created successfully - don't requeue
		return reconcile.Result{}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// Pod already exists - don't requeue
	reqLogger.Info("Skip reconcile: Pod already exists", "Pod.Namespace", found.Namespace, "Pod.Name", found.Name)

    // Define a new ConfigMapobject
	configMap := AddnewConfigMapForCR(instance)

	// Set AzurePipelinePool instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, configMap, r.Scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if this ConfigMap already exists
	found1 := &corev1.ConfigMap{}
	err = r.Client.Get(context.TODO(), types.NamespacedName{Name: configMap.Name, Namespace: configMap.Namespace}, found1)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new ConfigMap", "ConfigMap.Namespace", configMap.Namespace, "ConfigMap.Name", configMap.Name)
		err = r.Client.Create(context.TODO(), configMap)
		if err != nil {
			return reconcile.Result{}, err
		}

		// ConfigMap created successfully - don't requeue
		return reconcile.Result{}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// ConfigMap already exists - don't requeue
	reqLogger.Info("Skip reconcile: ConfigMap already exists", "ConfigMap.Namespace", found1.Namespace, "ConfigMap.Name", found1.Name)

	service := AddnewServiceForCR(instance)

	// Set AzurePipelinePool instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, service, r.Scheme); err != nil {
		return reconcile.Result{}, err
	}

	found2 := &corev1.Service{}
	err = r.Client.Get(context.TODO(), types.NamespacedName{Name: service.Name, Namespace: service.Namespace}, found2)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Service", "Service.Namespace", service.Namespace, "Service.Name", service.Name)
		err = r.Client.Create(context.TODO(), service)
		if err != nil {
			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	reqLogger.Info("Skip reconcile: Service already exists", "Service.Namespace", found2.Namespace, "Service.Name", found2.Name)

	buildkitPod := AddnewBuildkitPodForCR(instance)

	// Set AzurePipelinePool instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, buildkitPod, r.Scheme); err != nil {
		return reconcile.Result{}, err
	}

	found3 := &appsv1.StatefulSet{}
	err = r.Client.Get(context.TODO(), types.NamespacedName{Name: buildkitPod.Name, Namespace: buildkitPod.Namespace}, found3)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Buildkit Pod", "BuildKitPod.Namespace", buildkitPod.Namespace, "BuildKitPod.Name", buildkitPod.Name)
		err = r.Client.Create(context.TODO(), buildkitPod)
		if err != nil {
			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	reqLogger.Info("Skip reconcile: Buildkit pod already exists", "BuildkitPod.Namespace", found3.Namespace, "BuildKitPod.Name", found3.Name)

	buildkitService := AddnewBuildkitServiceForCR(instance)

	// Set AzurePipelinePool instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, buildkitService, r.Scheme); err != nil {
		return reconcile.Result{}, err
	}

	found4 := &corev1.Service{}
	err = r.Client.Get(context.TODO(), types.NamespacedName{Name: buildkitService.Name, Namespace: buildkitService.Namespace}, found4)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Buildkit Service", "BuildKitService.Namespace", buildkitService.Namespace, "BuildKitService.Name", buildkitService.Name)
		err = r.Client.Create(context.TODO(), buildkitService)
		if err != nil {
			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	reqLogger.Info("Skip reconcile: Buildkit Service already exists", "BuildkitService.Namespace", found4.Namespace, "BuildKitService.Name", found4.Name)
	return reconcile.Result{}, nil
}

func IsInitialized(obj metav1.Object) bool {
	azurepipelinepoolobj, ok := obj.(*devv1alpha1.AzurePipelinesPool)
	if !ok {
		return false
	}
	if azurepipelinepoolobj.Spec.Initialized {
		return true
	}
	addFinalizer(azurepipelinepoolobj, controllerName)
	azurepipelinepoolobj.Spec.Initialized = true
	return false

}

func isBeingDeleted(obj metav1.Object) bool {
	return !obj.GetDeletionTimestamp().IsZero()
}

func addFinalizer(obj metav1.Object, finalizer string) {
	if !hasFinalizer(obj, finalizer) {
		  obj.SetFinalizers(append(obj.GetFinalizers(), finalizer))
	}
}

func hasFinalizer(obj metav1.Object, finalizer string) bool {
	for _, fin := range obj.GetFinalizers() {
			 if fin == finalizer {
				  return true
			 }
	}
	return false
}

func removeFinalizer(obj metav1.Object, finalizer string) {
	for i, fin := range obj.GetFinalizers() {
		   if fin == finalizer {
				 finalizers := obj.GetFinalizers()
				 finalizers[i] = finalizers[len(finalizers)-1]
				 obj.SetFinalizers(finalizers[:len(finalizers)-1])
				 return
		   }
	}	
}

func manageCleanUpLogic(cr *devv1alpha1.AzurePipelinesPool) error{
  // perform additional cleanup here
  return nil
}

func AddnewPodForCR(cr *devv1alpha1.AzurePipelinesPool) *corev1.Pod {
  labels := map[string]string{
	  "app": cr.Name,
	  "tier":"frontend",
  }
  return &corev1.Pod {
	  ObjectMeta: metav1.ObjectMeta {
		  Name:      "azurepipelinepod",
		  Namespace: cr.Namespace,
		  Labels:    labels,
	  },
	  Spec: corev1.PodSpec{
		  Containers: []corev1.Container {
			  {
				  Name:   cr.Name,
				  Image:  cr.Spec.ControllerName,
				  Env:    []corev1.EnvVar {
					  {
						  Name: "VSTS_SECRET",
						  ValueFrom: &corev1.EnvVarSource {
							  SecretKeyRef: &corev1.SecretKeySelector {
								  LocalObjectReference: corev1.LocalObjectReference{Name: "vsts"},
								  Key: "VSTS_SECRET",
							  },
						  },
					  },
					  {
						  Name: "POD_NAMESPACE",
                          Value: cr.Namespace,
					  },
				  },
				  Ports: []corev1.ContainerPort {
					  {
						  ContainerPort: 8080,
					  },
				  },
			  },
		  },
	  },
  }
}

func AddnewBuildkitPodForCR(cr *devv1alpha1.AzurePipelinesPool) *appsv1.StatefulSet {
  labels := map[string]string {
	  "app": cr.Name,
  }

  labels1 := map[string]string {
	  "app":  cr.Name,
	  "role": "buildkit",
  }

  annotations := map[string]string {
	  "container.apparmor.security.beta.kubernetes.io/buildkitd": "unconfined",
	  "container.seccomp.security.alpha.kubernetes.io/buildkitd": "unconfined",
  }

  return &appsv1.StatefulSet {
	  TypeMeta: metav1.TypeMeta {
		  Kind:       "StatefulSet",
		  APIVersion: "apps/v1",
	  },
	  ObjectMeta: metav1.ObjectMeta {
		  Name:      "buildkitd",
		  Namespace: cr.Namespace,
		  Labels:    labels,
	  },
	  Spec: appsv1.StatefulSetSpec {
		  Selector:    &metav1.LabelSelector{MatchLabels:labels},
		  ServiceName: "buildkitd",
		  Replicas:    &cr.Spec.BuildkitReplicaCount,
		  Template: corev1.PodTemplateSpec {
			  ObjectMeta: metav1.ObjectMeta {
				  Labels: labels1,
				  Annotations: annotations,
			  },
			  Spec: corev1.PodSpec {
				  Containers: []corev1.Container {
					  corev1.Container {
						  Name:  "buildkitd",
						  Image: "moby/buildkit:master-rootless",
						  Args:  []string{"--addr","unix:///run/user/1000/buildkit/buildkitd.sock","--addr","tcp://0.0.0.0:1234","--oci-worker-no-process-sandbox"},
						  Ports: []corev1.ContainerPort {
							  {
								  ContainerPort: 1234,
							  },
						  },
					  },
				  },
			  },
		  },
	  },
  }
}

func AddnewBuildkitServiceForCR(cr *devv1alpha1.AzurePipelinesPool) *corev1.Service {
	 labels := map[string]string {
	  "app": cr.Name,
  }
  return &corev1.Service {
	  ObjectMeta: metav1.ObjectMeta {
		  Namespace: cr.Namespace,
		  Name:      "buildkitd",
		  Labels:    labels,
	  },
	  Spec: corev1.ServiceSpec {
		  Selector: labels,
		  Ports:    []corev1.ServicePort {
			  {
				  Port:       1234,
				  Protocol:   "TCP",
			  },
		  },
	  },
  }
}

func AddnewConfigMapForCR(cr *devv1alpha1.AzurePipelinesPool) *corev1.ConfigMap {
  return &corev1.ConfigMap {
	  ObjectMeta: metav1.ObjectMeta {
		  Name:      "kubernetes-config",
		  Namespace: cr.Namespace,
	  },
	  Data: map[string]string {
		  "type": "KUBERNETES",
	  },
  }
}

func AddnewServiceForCR(cr *devv1alpha1.AzurePipelinesPool) *corev1.Service {
  labels := map[string]string {
	  "app": cr.Name,
	  "tier":"frontend",
  }
  return &corev1.Service {
	  ObjectMeta: metav1.ObjectMeta {
		  Namespace: cr.Namespace,
		  Name:      "azurepipelinespoolservice",
		  Labels:    labels,
	  },
	  Spec: corev1.ServiceSpec {
		  Selector: labels,
		  Type:     corev1.ServiceTypeLoadBalancer,
		  Ports:    []corev1.ServicePort {
			  {
				  Port: 8080,
			  },
		  },
	  },
  }
}