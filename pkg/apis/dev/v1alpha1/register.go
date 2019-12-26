// NOTE: Boilerplate only.  Ignore this file.

// Package v1alpha1 contains API Schema definitions for the dev v1alpha1 API group
// +k8s:deepcopy-gen=package,register
// +groupName=dev.azure.com
package v1alpha1

import (
	"log"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/rest"
)

var (
	// SchemeGroupVersion is group version used to register these objects
	SchemeGroupVersion = schema.GroupVersion{Group: "dev.azure.com", Version: "v1alpha1"}

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuilder = &scheme.Builder{GroupVersion: SchemeGroupVersion}
)

var testingclient AzurePipelinesPoolV1Alpha1Client

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&AzurePipelinesPool{},
		&AzurePipelinesPoolList{},
	)
	meta_v1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}

func NewClient(cfg *rest.Config) (*AzurePipelinesPoolV1Alpha1Client, error) {
	scheme := runtime.NewScheme()
	SchemeBuilder := runtime.NewSchemeBuilder(addKnownTypes)
	if err := SchemeBuilder.AddToScheme(scheme); err != nil {
		return nil, err
	}
	config := *cfg
	config.GroupVersion = &SchemeGroupVersion
	config.APIPath = "/apis"
	config.ContentType = runtime.ContentTypeJSON
	config.NegotiatedSerializer = serializer.NewCodecFactory(scheme)

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}
	log.Println("rest client inside newclient", client)
	return &AzurePipelinesPoolV1Alpha1Client{RestClient: client}, nil
}

func NewClientTest() (*AzurePipelinesPoolV1Alpha1Client, error) {
	/*scheme := runtime.NewScheme()
	SchemeBuilder := runtime.NewSchemeBuilder(addKnownTypes)
	if err := SchemeBuilder.AddToScheme(scheme); err != nil {
		return nil, err
	}
	client, err := rest.RESTClientFor(&rest.Config{APIPath: "/apis", ContentConfig: rest.ContentConfig{GroupVersion: &SchemeGroupVersion, NegotiatedSerializer: serializer.NewCodecFactory(scheme)}})
	if err != nil {
		return nil, err
	}
	log.Println("rest client inside newclient", client)*/
	return &testingclient, nil
}

func SetClient(s *runtime.Scheme ) {
	client, err := rest.RESTClientFor(&rest.Config{APIPath: "/apis", ContentConfig: rest.ContentConfig{GroupVersion: &SchemeGroupVersion, NegotiatedSerializer: serializer.NewCodecFactory(s)}})
	if err != nil {
		
	}
	log.Println("rest client inside newclient", client)
	testingclient.RestClient = client
	//testingclient.RestClient =  s
}
