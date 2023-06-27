#!/usr/bin/env bash

# This script returns the id of the draft release

# standard bash error handling
set -o nounset  # treat unset variables as an error and exit immediately.
set -o errexit  # exit immediately when a command fails.
set -E          # needs to be set if we want the ERR trap
set -o pipefail # prevents errors in a pipeline from being masked
set +x

RELEASE_TAG=$1

REPOSITORY=${REPOSITORY:-barchw/istio}
GITHUB_URL=https://api.github.com/repos/${REPOSITORY}
GITHUB_AUTH_HEADER="Authorization: Bearer ${GITHUB_TOKEN}"
CHANGELOG_FILE=$(cat CHANGELOG.md)
RELEASE_NOTES_FILE=$(cat "docs/release-notes/${RELEASE_TAG}.md")
BODY="${RELEASE_NOTES_FILE}\\n${CHANGELOG_FILE}"

JSON_PAYLOAD=$(jq -n \
  --arg tag_name "$RELEASE_TAG" \
  --arg name "$RELEASE_TAG" \
  --arg body "$BODY" \
  '{
    "tag_name": $tag_name,
    "name": $name,
    "body": $body,
    "draft": true
  }')

CURL_RESPONSE=$(curl -L \
  -X POST \
  -H "Accept: application/vnd.github+json" \
  -H "${GITHUB_AUTH_HEADER}" \
  -H "X-GitHub-Api-Version: 2022-11-28" \
  "${GITHUB_URL}/releases" \
  -d "$JSON_PAYLOAD")

echo "$CURL_RESPONSE" | jq -r ".id"
