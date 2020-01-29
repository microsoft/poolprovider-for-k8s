#!/bin/bash

# command to run script and generate certificate on-fly ./setup.sh azpipelinespooltestcert true
# command to run script using existing certificate ./setup.sh azpipelinespooltest false privateKey.key certificate.crt

usage() {
    echo "Usage :"
    echo "./setup.sh -d <dnsname> -u <useletsencrypt> or"
    echo "./setup.sh -d <dnsname> -u <useletsencrypt> -k <keypath> -c <certificate path>"
    echo "-d : (string) dnsname ex: testdomainname"
    echo "-u : (bool - true|false) uses letsencrypt if set to true else pass the exiting certifacte path"
    echo "-k : (string) indicates existing key path; used when -u is set to false"
    echo "-c : (string) indicates existing certificate path; used when -u is set to false"
    echo "-n : (string) namespace"
    echo "-h : help"
    exit 0;
    }

namespaceval="azuredevops"

if [ "$#" -lt "2" ]
then
    usage
fi

while getopts ":d:u:k:c:n:h" o;
do

  case "${o}" in
    d)
        echo "Dns name set"
        dnsname=${OPTARG}
        ;;
    u)
        echo "Use letsencrypt variable"
        useletsencrypt=${OPTARG}
        ;;
    k)
        echo "Keypath set"
        keypath=${OPTARG}
        ;;
    c)
        echo "Certificate path set"
        certpath=${OPTARG}
        ;;
    n)
        echo "namespace set"
        namespaceval=${OPTARG}
        ;;
    *)
        usage
  esac
done

if [ "$useletsencrypt" = false ]
then
    if [ -z "$keypath" -o -z "$certpath" ]
    then
        echo "If using existing certificate keypath and certificate path are mandatory"
        usage
    fi
fi

helm install k8s-poolprovidercrd --name-template k8spoolprovidercrd --set "azurepipelines.VSTS_SECRET=sharedsecret1234" --set "app.namespace=$namespaceval"
echo "K8s-poolprovidercrd helm chart installed"

sed -i 's/\(.*namespace:.*\)/  namespace: '$namespaceval'/g' k8s-poolprovidercrd/azurepipelinescr/azurepipelinespool_cr.yaml

kubectl apply -f k8s-poolprovidercrd/azurepipelinescr/azurepipelinespool_cr.yaml
echo "Custom resource yaml applied"

helm repo add stable https://kubernetes-charts.storage.googleapis.com 
echo "Stable repo added"

helm repo update
echo "Helm repo updated"

helm install stable/nginx-ingress --generate-name --namespace $namespaceval
echo "Installed nginx-ingress"

cnt=0

while [ $cnt -lt 100 ]
do

  ingressip=$(kubectl get service -l app=nginx-ingress --namespace=$namespaceval -o=jsonpath='{.items[0].status.loadBalancer.ingress[0].ip}')

  if [ -n "$ingressip" ] 
  then
    echo "Found ingressip :" $ingressip
    break
  fi
  cnt=`expr $cnt + 1`
  sleep 2
  echo "Waiting for ingressip to be available...."

done

publicpid=$(az network public-ip list --query "[?ipAddress!=null]|[?contains(ipAddress, '$ingressip')].[id]" --output tsv)
echo "Fetched resource id"

# Update public ip address with DNS name
response=$(az network public-ip update --ids $publicpid --dns-name $dnsname )
echo "Assigned DnsName with ip address"

fqdn=`echo $response | jq '.dnsSettings.fqdn'`

echo "Fetched fully qualified domain name: " $fqdn

if [ "$useletsencrypt" = true ]
then
    kubectl apply -f https://raw.githubusercontent.com/jetstack/cert-manager/release-0.8/deploy/manifests/00-crds.yaml
    echo "Installed cert-manager CRD"

    kubectl create namespace cert-manager
    echo "Created cert-manager namespace"

    # Label the cert-manager namespace to disable resource validation
    kubectl label namespace cert-manager certmanager.k8s.io/disable-validation=true
    echo "Labeled cert-manager namespace to disable validation"

    # Add the Jetstack Helm repository
    helm repo add jetstack https://charts.jetstack.io
    echo "Added jetstack repo"

    helm repo update
    echo "Updated helm repo"

    # Install the cert-manager Helm chart
    helm install --name-template cert-manager --namespace cert-manager --version v0.8.0 jetstack/cert-manager
    echo "Installed helm repo for cert-manager"

    sleep 70
else
    kubectl create secret tls tls-secret --key $keypath --cert $certpath -n $namespaceval
    echo "tls-secret created"
fi

helm install k8s-certmanager --name-template k8spoolprovidercert --set "configvalues.dnsname=$fqdn" --set "letsencryptcert.val=$useletsencrypt"  --set "app.namespace=$namespaceval"
echo "---- Cluster configuration successfully done. ----"