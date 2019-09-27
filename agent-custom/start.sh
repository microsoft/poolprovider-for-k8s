#!/bin/bash
set -e

rm -rf /azp/agent
mkdir /azp/agent
cd /azp/agent

# Read download URL from the secret
AZP_DOWNLOAD_URL="$(cat /vsts/agent/.url)"

# Download the requested agent, else fall back to the version that was default when this was released
if [ -z "$AZP_DOWNLOAD_URL" ]; then
  AZP_DOWNLOAD_URL="https://vstsagentpackage.azureedge.net/agent/2.158.0/vsts-agent-linux-x64-2.158.0.tar.gz"
fi

export AGENT_ALLOW_RUNASROOT="1"

print_header() {
  lightcyan='\033[1;36m'
  nocolor='\033[0m'
  echo -e "${lightcyan}$1${nocolor}"
}

# Let the agent ignore the token env variables

print_header "1. Downloading and installing Azure Pipelines agent from the agent request..."

curl -LsS $AZP_DOWNLOAD_URL | tar -xz & wait $!

source ./env.sh

trap 'cleanup; exit 130' INT
trap 'cleanup; exit 143' TERM

print_header "2. Mounting the .agent and .credentials files from a different source"
ln -fs /vsts/agent/.agent /azp/agent/.agent
ln -fs /vsts/agent/.credentials /azp/agent/.credentials
# ln -fs /vsts/agent/.credentials_rsaparams /azp/agent/.credentials_rsaparams

print_header "3. Running Azure Pipelines agent..."

# `exec` the node runtime so it's aware of TERM and INT signals
# AgentService.js understands how to handle agent self-update and restart
exec ./externals/node/bin/node ./bin/AgentService.js interactive