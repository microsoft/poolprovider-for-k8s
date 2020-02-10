# k8s-poolprovider

You can configure your kubernetes cluster using either of below listed approaches - 
1. Want to use an existing certificate 
   ./setup.sh -s <sharedsecret> -d <dnsname> -u <useletsencrypt> -k <keypath> -c <certificate path>
2. Want to create new certificate using letsencrypt
   ./setup.sh -s <sharedsecret> -d <dnsname> -u <useletsencrypt>

Before running the script user need to have az login.

This script installs below helm charts - 
1. k8s-poolprovidercrd - This helm chart installs all the resources required for configuring cluster like webserver deployment, internal service, buildkit pods etc.
  User can make changes to Custom reosurce file i.e. azurepipelinespool_cr.yaml. In this file he can add modified controller container image, change the count of buildkit pods and add the customised agent images.

2. nginx/ingress - This helm chart installs the nginx loadbalancer which supports https communication.

3. cert-manager - This helm chart is installed if using second approach where certificates are installed using letsencrypt. It is basically Kubernetes certificate management controller which helps in issueing of new valid certificates and even attempts to renew certificates before the expiration timestamp.

4. k8spoolprovidercert -  This helm chart installs different resources based on the approach being used.
  # Approach 1 - 
    Installs the ingress which creates rules which helps in diverting the incoming https requests to the internal service being installed
  # Approach 2 -
    Installs the ClusterIsuer and certificate along with ingress. Other two resources are required for certificate creation.
  
Note : As part of setup script we bind the public ip of ingress with the DNS name provided by user. Currently to perform this operation script is using az commands if you want to configure cluster other than AKS please change those commands
