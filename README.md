# k8s-poolprovider
helm install k8s-poolprovider --name=myhelmchart --set "vsts.VSTS_ACCOUNT=accountname" --set "vsts.VSTS_POOL=poolname" --set "vsts.VSTS_TOKEN=pat token" --set "vsts.VSTS_SECRET=shared secret"

## For local testing
Use the kubernetes.yaml and buildkit.yaml under Manifests folder to set up the application on your cluster. This works if you want to test the containerised application on a kubernetes cluster.

When using VS Code, the launch.json has been appropriately modified for local debugging. If debugging locally, 
1. Make sure you can run kubectl commands from your machine, and a kubeconfig file is present on your machine.
2. Create a new namespace 'azuredevops' on your cluster.
3. Create an opaque secret in the azuredevops namespace with the required secrets set. `kubectl create secret generic vsts --from-literal=VSTS_TOKEN=<token> --from-literal=VSTS_ACCOUNT=<account> --from-literal=VSTS_POOL=<poolname> -n azuredevops`
4. Voila! You can start debugging directly from VS code. When the deployment happens, hit localhost:8082 with the correct APIs and see the app in action.

