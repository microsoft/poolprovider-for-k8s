# k8s-poolprovider

## Introduction

When using multi-tenant hosted pools there are times the jobs remain in queued state because all the agents are occupied or we still hold the physical resources even when there are very few requests present to be addressed. This causes performance issues. To address the problems above, we have implemented Kubernetes based poolprovider which provides the elasticity of agent pools. 

> This feature is in private preview. We recommend that you <ins>**do not use this in production**.</ins> 

The k8s-poolprovider uses Kubernetes cluster as the build infrastructure.

This repository consists implementation of two major helm charts -
#### 1. k8s-poolprovidercrd :
This helm chart installs all the resources required for configuring Kubernetes poolprovider resources on the Kuberenetes cluster. It first installs the controller implemented using Operator-SDK. This is required for lifecycle management of poolprovider resources deployed in the cluster. As soon as user applies the custom resource yaml i.e. [azurepipelinespool_cr.yaml](https://github.com/microsoft/k8s-poolprovider/blob/prebansa-readme/helm/k8s-poolprovidercrd/azurepipelinescr/azurepipelinespool_cr.yaml); the controller instantiates multiple external resources like webserver deployment, service, buildkit pods etc. The controller handles the reinitialization and reconfiguration at runtime if any changes are observed in the configured instances.
  User can make changes to Custom resource file i.e. azurepipelinespool_cr.yaml as per requirements. In this file user can add modified controller container image, change the number of buildkit pods instances and add the customised agent container images, refer this [CRD](https://github.com/microsoft/k8s-poolprovider/blob/master/helm/k8s-poolprovidercrd/templates/azurepipelinespools_crd.yaml) specification.

#### 2. k8s-certmanager :
This helm chart installs different resources required for configuring the load balancer endpoint with https support.
  ##### Approach 1 - User provides the existing certificates and Key 
   In this helm chart installs the ingress resource to configure the rules that route traffic to internal webserver already installed as part of previous helm chart. Assuming user has already created a tls-secret with the existing certificate and key.
  ##### Approach 2 - Use Let's Encrypt to create a valid certificate and Key 
   In this helm chart installs the ClusterIsuer and Certificate along with ingress resource.
   
In order to set up your Kubernetes cluster as the build infrastructure, you need to
1. Configure the pool provider on Kuberentes cluster
2. Add the Agent pool configured as Kubernetes poolprovider
    
## 1. Configure the poolprovider on Kubernetes cluster

1. Install k8s-poolprovidercrd helm chart   
   `helm install k8s-poolprovidercrd --name-template k8spoolprovidercrd --set "azurepipelines.VSTS_SECRET=sharedsecretval" --set  "app.namespace=namespaceval"`   
   sharedsecretval - Value must be of atleast 16 characters    
   namespaceval - Namespace where all the poolprovider resources will be deployed 
2. Apply poolprovider custom resource yaml   
   `kubectl apply azurepipelinespool_cr.yaml`
3. Run helm install stable/nginx-ingress   
   `helm install stable/nginx-ingress --generate-name --namespace $namespaceval`
4. Execute commands to link the ingress service public ip with valid DNS name   
   For azure following set of commands are used -     
   ```
   kubectl get service -l app=nginx-ingress --namespace=namespaceval -o=jsonpath='{.items[0].status.loadBalancer.ingress[0].ip}'
   publicpid=$(az network public-ip list --query "[?ipAddress!=null]|[?contains(ipAddress, 'ingressip')].[id]" --output tsv) 
   
   az network public-ip update --ids $publicpid --dns-name dnsname 
   ```
    Note : You can learn more about the az network public-ip update command [here](https://docs.microsoft.com/en-us/cli/azure/network/public-ip?view=azure-cli-latest#az-network-public-ip-update)
5. Run helm install cert-manager if you want to use Let's Encrypt else execute    
   `kubectl create secret tls tls-secret --key keypath --cert certpath -n namespace`   
   keypath - Specify path for key    
   certpath - Specify path for certificate   
6. Install k8s-certmanager helm chart   
   `helm install k8s-certmanager --name-template k8spoolprovidercert --set "configvalues.dnsname=fqdn" --set "letsencryptcert.val=false"  --set "app.namespace=namespaceval"`
   
   fqdn - Fully qualified domain name for which the key and certificate are generated    
   namespaceval - Namespace where all the poolprovider resources will be deployed, this parameter is same as required in Step 1
   
### User can configure Azure Kubernetes Cluster using existing setup script - 
Note - If using an existing AKS cluster, user needs to have az login and get access credentials for a managed Kubernetes cluster using `az aks get-credentials` command. Refer [here](https://docs.microsoft.com/cli/azure/aks?view=azure-cli-latest#az-aks-get-credentials) for the command documentation.

Before running the script user need to have az login.
##### Approach 1 - User provides the existing certificates and Key
   ./setup.sh -s "sharedsecret" -d "dnsname" -u "useletsencrypt" -k "keypath" -c "certificate path"
##### Approach 2 - Use Let's Encrypt to create a valid certificate and Key 
   ./setup.sh -s "sharedsecret" -d "dnsname" -u "useletsencrypt"

##### Description of option arguments passed in the script
      
      -d : (string) dnsname (mandatory) ex: testdomainname
      -u : (bool - true|false) (mandatory) uses letsencrypt if set to true else pass the exiting certifacte path
      -k : (string) indicates existing key path; used when -u is set to false
      -c : (string) indicates existing certificate path; used when -u is set to false
      -n : (string) namespace (optional)
      -s : (string) sharedsecret (mandatory)
      -h : help

Note : As part of setup script we bind the public ip of ingress with the DNS name provided by user. Currently to perform this operation script is using az commands if you want to configure cluster other than AKS please change those commands.

## 2. Add Agent pool configured as Kubernetes poolprovider

1. Run the powershell script poolprovidersetup.ps1

	./poolprovidersetup.ps1 
  
  
   ##### Description of option arguments passed in the script
   
        URI : Account URI to be configured for poolprovider
        PATToken : PAT token for the account
        PoolName : AgentPool name to be configured as Kubernetes poolprovider
        DNSName : Same DNS name with which the key and secrets are generated
        Sharedsecret : Secret value having atleast 16 characters; needs to be xact same value as provided while configuring the cluster
        TargetSize : Target parallelism required in agent pool
