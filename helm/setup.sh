#!/bin/bash

# command to run script and generate certificate on-fly ./setup.sh azpipelinespooltestcert true
# command to run script using existing certificate ./setup.sh azpipelinespooltest false privateKey.key certificate.crt
if [ "$#" -lt "2" ]
then
    echo "Pass atleast two arguments - DNS Name, Bool value to indicate if existing certificate used"
    exit 0
fi

if [ "$2" = false -a "$#" -lt "4" ]
then
    echo "Provide the files with key and certificate"
    exit 0
fi

dnsname=$1
echo "1. DNS name set"

useletsencrypt=$2
echo "2. Use letsencrypt variable set"

if [ "$useletsencrypt" = false ]
then
    keypath=$3
    certpath=$4
fi

helm install k8s-poolprovidercrd --name-template k8spoolprovidercrd --set "azurepipelines.VSTS_SECRET=sharedsecret1234"
echo "3. k8s-poolprovidercrd helm chart installed"

kubectl apply -f k8s-poolprovidercrd/azurepipelinescr/azurepipelinespool_cr.yaml
echo "4. Custom resource yaml applied"

helm repo add stable https://kubernetes-charts.storage.googleapis.com 
echo "5. Stable repo added"

helm repo update
echo "6. Helm repo updated"

helm install stable/nginx-ingress --generate-name --namespace azuredevops 
echo "7. Installed nginx-ingress"

while true; do

  ingressip=$(kubectl get service -l app=nginx-ingress --namespace=azuredevops -o=jsonpath='{.items[0].status.loadBalancer.ingress[0].ip}')

  if [ -n "$ingressip" ] 
  then
    echo "8. Found ingressip :" $ingressip
    break
  fi
done

publicpid=$(az network public-ip list --query "[?ipAddress!=null]|[?contains(ipAddress, '$ingressip')].[id]" --output tsv)
echo "9. Fetched resource id"

# Update public ip address with DNS name
response=$(az network public-ip update --ids $publicpid --dns-name $dnsname )
echo "10. Assigned DnsName with ip address"

fqdn=`echo $response | jq '.dnsSettings.fqdn'`

echo "11. Fetched fully qualified domain name: " $fqdn

if [ "$useletsencrypt" = true ]
then
    kubectl apply -f https://raw.githubusercontent.com/jetstack/cert-manager/release-0.8/deploy/manifests/00-crds.yaml
    echo "12. Installed cert-manager CRD"

    kubectl create namespace cert-manager
    echo "13. Created cert-manager namespace"

    # Label the cert-manager namespace to disable resource validation
    kubectl label namespace cert-manager certmanager.k8s.io/disable-validation=true
    echo "14. Labeled cert-manager namespace to disable validation"

    # Add the Jetstack Helm repository
    helm repo add jetstack https://charts.jetstack.io
    echo "15. Added jetstack repo"

    helm repo update
    echo "16. Updated helm repo"

    # Install the cert-manager Helm chart
    helm install --name-template cert-manager --namespace cert-manager --version v0.8.0 jetstack/cert-manager
    echo "17. Installed helm repo for cert-manager"

    sleep 70
else
    kubectl create secret tls tls-secret --key $keypath --cert $certpath -n azuredevops
    echo "12. tls-secret created"
fi

helm install k8s-certmanager --name-template k8spoolprovidercert --set "configvalues.dnsname=$fqdn" --set "letsencryptcert.val=$useletsencrypt"
echo "---- Cluster configuration successfully done. ----"