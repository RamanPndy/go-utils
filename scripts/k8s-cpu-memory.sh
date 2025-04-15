#!/bin/bash

# Check if namespace is passed
if [ -z "$1" ]; then
  echo "Usage: $0 <namespace>"
  exit 1
fi

NAMESPACE=$1

# Print Header
echo -e "DEPLOYMENT\tCONTAINER\tREQ_CPU\tREQ_MEM\tLIMIT_CPU\tLIMIT_MEM"

# Get and format deployment resource data
kubectl get deployments -n "$NAMESPACE" -o json | jq -r '
  .items[] |
  .metadata.name as $dname |
  .spec.template.spec.containers[] |
  [$dname, .name, (.resources.requests.cpu // "none"), (.resources.requests.memory // "none"), (.resources.limits.cpu // "none"), (.resources.limits.memory // "none")] |
  @tsv
' | column -t -s $'\t'
