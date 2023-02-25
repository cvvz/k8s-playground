#!/bin/bash

set -x

export CLUSTER_NAME="weizhichen-wi"
export CLUSTER_RESOURCE_GROUP="weizhichen-wi"
export LOCATION="westus"


# 1. 创建AKS集群

# az cloud set -n AzureCloud && \
# az account set --subscription 8ecadfc9-d1a3-4ea4-b844-0d9f87e4d7c8 && \
# az group create -l ${LOCATION} -n ${CLUSTER_RESOURCE_GROUP} && \
# az aks create -n ${CLUSTER_NAME} -g ${CLUSTER_RESOURCE_GROUP} --enable-oidc-issuer --enable-workload-identity


export SERVICE_ACCOUNT_ISSUER=$(az aks show -n $CLUSTER_NAME -g $CLUSTER_RESOURCE_GROUP --query "oidcIssuerProfile.issuerUrl" -otsv)
export USER_ASSIGNED_IDENTITY_NAME="weizhichen-wi-test-msi"
export RESOURCE_GROUP="weizhichen-wi-test-rg"
export SERVICE_ACCOUNT_NAMESPACE="default"
export SERVICE_ACCOUNT_NAME="weizhichen-wi-test-sa"

export KEYVAULT_NAME="azwi-kv-$(openssl rand -hex 2)"
export KEYVAULT_SECRET_NAME="my-secret"


az group create --name "${RESOURCE_GROUP}" --location "${LOCATION}"


# 2 创建KeyVault
az keyvault create --resource-group "${RESOURCE_GROUP}" \
   --location "${LOCATION}" \
   --name "${KEYVAULT_NAME}"

az keyvault secret set --vault-name "${KEYVAULT_NAME}" \
   --name "${KEYVAULT_SECRET_NAME}" \
   --value "Hello\!"

# 3 创建MSI并赋予权限
az identity create --name "${USER_ASSIGNED_IDENTITY_NAME}" --resource-group "${RESOURCE_GROUP}"

export USER_ASSIGNED_IDENTITY_CLIENT_ID="$(az identity show --name "${USER_ASSIGNED_IDENTITY_NAME}" --resource-group "${RESOURCE_GROUP}" --query 'clientId' -otsv)"
export USER_ASSIGNED_IDENTITY_OBJECT_ID="$(az identity show --name "${USER_ASSIGNED_IDENTITY_NAME}" --resource-group "${RESOURCE_GROUP}" --query 'principalId' -otsv)"
az keyvault set-policy --name "${KEYVAULT_NAME}" \
  --secret-permissions get \
  --object-id "${USER_ASSIGNED_IDENTITY_OBJECT_ID}"


# 4 MSI和sa绑定
az identity federated-credential create \
  --name "kubernetes-federated-credential" \
  --identity-name "${USER_ASSIGNED_IDENTITY_NAME}" \
  --resource-group "${RESOURCE_GROUP}" \
  --issuer "${SERVICE_ACCOUNT_ISSUER}" \
  --subject "system:serviceaccount:${SERVICE_ACCOUNT_NAMESPACE}:${SERVICE_ACCOUNT_NAME}"


# 5 创建sa和pod
export KEYVAULT_URL="$(az keyvault show -g ${RESOURCE_GROUP} -n ${KEYVAULT_NAME} --query properties.vaultUri -o tsv)"
cat <<EOF | kubectl create -f -
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name:  ${SERVICE_ACCOUNT_NAME}
  namespace: ${SERVICE_ACCOUNT_NAMESPACE}
  annotations: 
    azure.workload.identity/client-id: "${USER_ASSIGNED_IDENTITY_CLIENT_ID}" 
  labels: 
    azure.workload.identity/use: "true"
---
apiVersion: v1
kind: Pod
metadata:
  name: quick-start
  namespace: ${SERVICE_ACCOUNT_NAMESPACE}
  labels:
    azure.workload.identity/use: "true"
spec:
  serviceAccountName: ${SERVICE_ACCOUNT_NAME}
  containers:
    - image: ghcr.io/azure/azure-workload-identity/msal-go
      name: oidc
      env:
      - name: KEYVAULT_URL
        value: ${KEYVAULT_URL}
      - name: SECRET_NAME
        value: ${KEYVAULT_SECRET_NAME}
  nodeSelector:
    kubernetes.io/os: linux
EOF

#kubectl logs quick-start -f

# 6 clenaup

#kubectl delete pod quick-start
#kubectl delete sa "${SERVICE_ACCOUNT_NAME}" --namespace "${SERVICE_ACCOUNT_NAMESPACE}"

#az group delete --name "${RESOURCE_GROUP}"