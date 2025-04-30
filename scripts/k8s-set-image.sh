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
    break
  else
    echo "Invalid selection. Try again."
  fi
done

read -p "Enter new image for container '$container': " NEW_IMAGE

echo "Updating image of container '$container' in deployment '$DEPLOYMENT' to '$NEW_IMAGE'..."

kubectl set image deployment/"$DEPLOYMENT" "$container"="$NEW_IMAGE" -n "$NAMESPACE"

if [ $? -eq 0 ]; then
  echo "✅ Image updated successfully."
else
  echo "❌ Failed to update image."
fi
