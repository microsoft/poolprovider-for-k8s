#!/bin/bash
 
read -p "Enter URI of the account to be configured for poolprovider: "  URI

if [ -z "$URI" ]
then
    echo "Provide a valid URI value!!"
    exit 0;
fi
 
read -p "Enter PAT Token: " PATToken
 
if [ -z "$PATToken" ]
then
    echo "Provide a valid PAT token value!!"
    exit 0;
fi
 
read -p "Enter poolname: " poolname
 
if [ -z "$poolname" ]
then
    echo "Provide a valid poolname!!"
    exit 0;
fi
 
read -p "Enter DNSName (same as the ingress host name configured in k8s cluster): " dnsname
 
if [ -z "$dnsname" ]
then
    echo "Provide a valid dnsname value!!"
    exit 0;
fi
 
read -p "Enter shared secret (atleast 16 characters required; needs to be exactly same as set while configuring k8s cluster): " sharedSecret
 
if [ -z "$sharedSecret" ]
then
    echo "Provide a valid sharedsecret value!!"
    exit 0;
fi
 
if [ ${#sharedSecret} -lt 16 ]
then
    echo "Provide a valid 16 length sharedsecret value!!"
    exit 0;
fi
 
read -p "Enter target size required for parallelism: " targetSize
 
if [ -z "$targetSize" ]
then
    echo "Provide a valid targetsize value!!"
    exit 0;
fi
 
# generate the base64 encoded value
hextoken=$(echo -n ':'$PATToken | base64)
 
echo 'hextoken is' . ${hextoken}
 
agentcloudid=$(curl -v --header "Content-Type: application/json" --header "Accept: application/json" --header "Authorization: Basic ${hextoken}" -d "{\"name\":\"$poolname\", \"type\":\"Ignore\", \"acquireAgentEndpoint\":\"https://$dnsname/acquire\", \"releaseAgentEndpoint\":\"https://$dnsname/release\", \"sharedSecret\":\"$sharedSecret\"}" $URI/_apis/distributedtask/agentclouds?api-version=5.0-preview  |  jq -r  '.agentCloudId')
 
echo "agentcloud is " .${agentcloudid}
 
agentpoolresponse=$(curl -v --header "Content-Type: application/json" --header "Accept: application/json" --header "Authorization: Basic ${hextoken}" -d "{\"name\":\"$poolname\", \"agentCloudId\":\"$agentcloudid\", \"targetSize\":\"$targetSize\"" $URI/_apis/distributedtask/pools?api-version=5.0-preview)
 
 
echo ${agentpoolresponse}
