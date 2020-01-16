#!/bin/bash

helm install k8s-poolprovidercrd --generate-name --set "azurepipelines.VSTS_SECRET=sharedsecret1234"

kubectl apply -f k8s-poolprovidercrd/azurepipelinescr/azurepipelinespool_cr.yaml

helm repo add stable https://kubernetes-charts.storage.googleapis.com 

helm repo update

helm install stable/nginx-ingress --generate-name --namespace azuredevops 

sleep 70


ingressip=$(kubectl get service -l app=nginx-ingress --namespace=azuredevops -o=jsonpath='{.items[0].status.loadBalancer.ingress[0].ip}')
echo "nginx ingress ip is " $ingressip

dnsname="azurepipelinespool"

publicpid=$(az network public-ip list --query "[?ipAddress!=null]|[?contains(ipAddress, '$ingressip')].[id]" --output tsv)

# Update public ip address with DNS name
response=$(az network public-ip update --ids $publicpid --dns-name $dnsname )

fqdn=`echo $response | jq '.dnsSettings.fqdn'`

echo $fqdn

# Install the CustomResourceDefinition resources separately
kubectl apply -f https://raw.githubusercontent.com/jetstack/cert-manager/release-0.8/deploy/manifests/00-crds.yaml

# Create the namespace for cert-manager
kubectl create namespace cert-manager

# Label the cert-manager namespace to disable resource validation
kubectl label namespace cert-manager certmanager.k8s.io/disable-validation=true

# Add the Jetstack Helm repository
helm repo add jetstack https://charts.jetstack.io

# Update your local Helm chart repository cache
helm repo update

# Install the cert-manager Helm chart
helm install --name-template cert-manager --namespace cert-manager --version v0.8.0 jetstack/cert-manager

helm install k8s-certmanager --generate-name --set "configvalues.dnsname=$fqdn"