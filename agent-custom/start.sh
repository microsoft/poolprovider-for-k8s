#!/bin/bash

set -e

AGENT_FOLDER="agent"
AZP_AGENT_VERSION="$(cat /vsts/agent/.agentVersion)"
AGENTVERSION="$(curl -s "https://api.github.com/repos/microsoft/azure-pipelines-agent/releases/latest" | jq -r .tag_name[1:])"

echo "Latest agent release version $AGENTVERSION" 
echo "Agent Version from request $AZP_AGENT_VERSION"

if [ $AGENTVERSION != $AZP_AGENT_VERSION ]; then

  echo "Downloading the Azure Pipelines agent from agent request....."

  # Deleting previous downloaded folder
  rm -rf /azp/$AGENT_FOLDER

  # setting the agent folder name to agent version from request
  AGENT_FOLDER=AZP_AGENT_VERSION
  mkdir /azp/$AGENT_FOLDER
  cd /azp/$AGENT_FOLDER

  # Read download URL from the secret
  AZP_DOWNLOAD_URL="$(cat /vsts/agent/.url)"

  # Download the requested agent, else fall back to the version that was default when this was released
  if [ -z "$AZP_DOWNLOAD_URL" ]; then
    AZP_DOWNLOAD_URL="https://vstsagentpackage.azureedge.net/agent/2.158.0/vsts-agent-linux-x64-2.158.0.tar.gz"
  fi
  
  echo "Installing Azure Pipelines agent from the agent request..."

  curl -LsS $AZP_DOWNLOAD_URL | tar -xz & wait $!

else
  echo "Using Azure Pipelines agent downloaded from the agent custom image..."
  cd /azp/$AGENT_FOLDER
fi

source ./env.sh

trap 'cleanup; exit 130' INT
trap 'cleanup; exit 143' TERM

echo "Mounting the .agent and .credentials files from a different source"
# /azurepipelines/agent is the default path change below lines if something else is set as mount path in CRD.
ln -fs /vsts/agent/.agent /azp/$AGENT_FOLDER/.agent
ln -fs /vsts/agent/.credentials /azp/$AGENT_FOLDER/.credentials

echo "Running Azure Pipelines agent..."

# `exec` the node runtime so it's aware of TERM and INT signals
# AgentService.js understands how to handle agent self-update and restart
exec ./externals/node/bin/node ./bin/AgentService.js interactive