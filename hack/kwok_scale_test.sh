#!/usr/bin/env bash

## KWOK repository
#KWOK_REPO=kubernetes-sigs/kwok
## Get latest
#KWOK_LATEST_RELEASE=$(curl "https://api.github.com/repos/${KWOK_REPO}/releases/latest" | jq -r '.tag_name')
#
#kubectl apply -f "https://github.com/${KWOK_REPO}/releases/download/${KWOK_LATEST_RELEASE}/kwok.yaml"
#
#kubectl apply -f "https://github.com/${KWOK_REPO}/releases/download/${KWOK_LATEST_RELEASE}/stage-fast.yaml"
#kubectl apply -f "https://github.com/${KWOK_REPO}/releases/download/${KWOK_LATEST_RELEASE}/metrics-usage.yaml"

echo "How many fake Nodes do you want to create?"
read -r nodeNumber

for i in $(seq 0 $nodeNumber); do
kubectl apply -f - <<EOF
apiVersion: v1
kind: Node
metadata:
  annotations:
    node.alpha.kubernetes.io/ttl: "0"
    kwok.x-k8s.io/node: fake
  labels:
    beta.kubernetes.io/arch: amd64
    beta.kubernetes.io/os: linux
    kubernetes.io/arch: amd64
    kubernetes.io/hostname: kwok-node-0
    kubernetes.io/os: linux
    kubernetes.io/role: agent
    node-role.kubernetes.io/agent: ""
    type: kwok
  name: kwok-node-${i}
spec:
  taints: # Avoid scheduling actual running pods to fake Node
  - effect: NoSchedule
    key: kwok.x-k8s.io/node
    value: fake
status:
  allocatable:
    cpu: 32
    memory: 256Gi
    pods: 110
  capacity:
    cpu: 32
    memory: 256Gi
    pods: 110
  nodeInfo:
    architecture: amd64
    bootID: ""
    containerRuntimeVersion: ""
    kernelVersion: ""
    kubeProxyVersion: v1.33.3+k3s1
    kubeletVersion: v1.33.3+k3s1
    machineID: ""
    operatingSystem: linux
    osImage: ""
    systemUUID: ""
  phase: Running
EOF
kubectl create namespace "fake-namespace-${i}" --dry-run=client -o yaml | kubectl apply -f -
kubectl label namespace "fake-namespace-${i}" istio-injection=enabled

done

echo "Proceed with deploying fake pods? (y/n)"
read -r answer
if [[ $answer == "n" || $answer == "N" ]]
then
  echo "Aborting redeploying k3d cluster"
  return
fi

for i in $(seq 0 $nodeNumber); do
  for j in $(seq 0 9); do
kubectl apply -f - <<EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  name: fake-pod-${i}-${j}
  namespace: fake-namespace-${i}
spec:
  replicas: 10
  selector:
    matchLabels:
      app: fake-pod-${i}-${j}
  template:
    metadata:
      labels:
        app: fake-pod-${i}-${j}
    spec:
      affinity:
        nodeAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
            nodeSelectorTerms:
            - matchExpressions:
              - key: type
                operator: In
                values:
                - kwok
      # A taints was added to an automatically created Node.
      # You can remove taints of Node or add this tolerations.
      tolerations:
      - key: "kwok.x-k8s.io/node"
        operator: "Exists"
        effect: "NoSchedule"
      containers:
      - name: fake-container
        image: fake-image
EOF
  done
done

for i in $(seq 0 $nodeNumber); do
  for j in $(seq 0 9); do
