#!/usr/bin/env bash

echo "How many fake Nodes were created? (used to determine cleanup range)"
read -r nodeNumber

echo "Deleting Deployments in fake namespaces..."
for i in $(seq 0 $nodeNumber); do
  for j in $(seq 0 9); do
    kubectl delete deployment "fake-pod-${i}-${j}" -n "fake-namespace-${i}" --ignore-not-found
  done
done

echo "Deleting fake namespaces..."
for i in $(seq 0 $nodeNumber); do
  kubectl delete namespace "fake-namespace-${i}" --ignore-not-found
done

echo "Deleting fake KWOK nodes..."
for i in $(seq 0 $nodeNumber); do
  kubectl delete node "kwok-node-${i}" --ignore-not-found
done

echo "Cleanup complete."

