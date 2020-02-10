# k8s-poolprovider
Kubernetes based pool provider implementation for Azure Pipelines. It uses Kubernetes cluster as the build infrastructure.

This repository consists implementation of two major helm charts -
#### 1. k8s-poolprovidercrd :
This helm chart installs all the resources required for configuring Kubernetes cluster. It first installs the controller implemented using Operator-SDK. This is required for lifecycle management of external resources deployed in the cluster. As soon as user applies the custom resource yaml i.e. azurepipelinespool_cr.yaml; the controller instantiates multiple external resources like webserver deployment, service, buildkit pods etc. The controller handles the reinitialization and reconfiguration at runtime if any changes are observed in the configured instances.
  User can make changes to Custom reosurce file i.e. azurepipelinespool_cr.yaml as per his requirement. In this file user can add modified controller container image, change the number of buildkit pods instances and add the customised agent container images.

#### 2. k8s-certmanager :
This helm chart installs different resources required for configuring the load balancer endpoint with https support.
  ##### Approach 1 - User provides the existing certificates and Key 
   In this helm chart installs the ingress resource to configure the rules that route traffic to internal webserver already installed as part of previous helm chart. Assuming user has alreday created a tls-secret with the existing certificate and key.
  ##### Approach 2 - Use Let's Encrypt to create a valid certificate and Key 
   In this helm chart installs the ClusterIsuer and certificate along with ingress resource. Other two resources are required for configuring valid certificate created at runtime.
    
## Steps to configure the Kubernetes Cluster

1. Install k8s-poolprovidercrd helm chart
2. Run kubectl apply azurepipelinespool_cr.yaml
3. Run helm install stable/nginx-ingress
4. Execute commands to link the ingress service public ip with valid DNS name
5. Fetch the fully qualified domain name 
6. Run helm install cert-manager if you want to use Let's Encrypt else execute 
   kubectl create secret tls tls-secret --key $keypath --cert $certpath -n $namespace
7. Install k8s-certmanager helm chart
    
User can configure Azure Kubernetes Cluster using existing setup script - 
Before running the script user need to have az login.
##### Approach 1 - User provides the existing certificates and Key
   ./setup.sh -s "sharedsecret" -d "dnsname" -u "useletsencrypt" -k "keypath" -c "certificate path"
##### Approach 2 - Use Let's Encrypt to create a valid certificate and Key 
   ./setup.sh -s "sharedsecret" -d "dnsname" -u "useletsencrypt"

Description of option arguments passed in the script
      
      -d : (string) dnsname (mandatory) ex: testdomainname
      -u : (bool - true|false) (mandatory) uses letsencrypt if set to true else pass the exiting certifacte path
      -k : (string) indicates existing key path; used when -u is set to false
      -c : (string) indicates existing certificate path; used when -u is set to false
      -n : (string) namespace (optional)
      -s : (string) sharedsecret (mandatory)
      -h : help

Note : As part of setup script we bind the public ip of ingress with the DNS name provided by user. Currently to perform this operation script is using az commands if you want to configure cluster other than AKS please change those commands.
