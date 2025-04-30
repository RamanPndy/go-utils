#!/bin/bash

read -p "Enter namespace: " NAMESPACE

# Get deployments in the namespace
deployments=($(kubectl get deployments -n "$NAMESPACE" -o jsonpath="{.items[*].metadata.name}"))

if [ ${#deployments[@]} -eq 0 ]; then
  echo "No deployments found in namespace '$NAMESPACE'."
  exit 1
fi

echo "Select a deployment:"
select DEPLOYMENT in "${deployments[@]}"; do
  if [[ -n "$DEPLOYMENT" ]]; then
    break
  else
    echo "Invalid selection. Try again."
  fi
done

# Get containers in the selected deployment
containers=($(kubectl get deployment "$DEPLOYMENT" -n "$NAMESPACE" -o jsonpath="{.spec.template.spec.containers[*].name}"))

if [ ${#containers[@]} -eq 0 ]; then
  echo "No containers found in deployment '$DEPLOYMENT'."
  exit 1
fi

echo "Select a container:"
select container in "${containers[@]}"; do
  if [[ -n "$container" ]]; then
    image=$(kubectl get deployment "$DEPLOYMENT" -n "$NAMESPACE" -o jsonpath="{.spec.template.spec.containers[?(@.name==\"$container\")].image}")
    echo "Image used by container '$container': $image"
    break
  else
    echo "Invalid selection. Try again."
  fi
done
