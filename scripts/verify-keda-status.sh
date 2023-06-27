#!/usr/bin/env bash



echo "Checking status of POST Jobs for istio"

fullstatus=`curl -L -H "Accept: application/vnd.github+json" -H "X-GitHub-Api-Version: 2022-11-28" https://api.github.com/repos/barchw/istio/commits/main/status | head -n 2 `

sleep 10
echo $fullstatus

if [[ "$fullstatus" == *"success"* ]]; then
  echo "All jobs succeeded"
else
  echo "Jobs failed or pending - Check Prow status"
  exit 1
fi