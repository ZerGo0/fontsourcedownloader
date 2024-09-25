#!/usr/bin/env bash

# make sure that the GITHUB_TOKEN env var is set
if [ -z "$GITHUB_TOKEN" ]; then
  echo "GITHUB_TOKEN is not set"
  exit 1
fi

# make sure that the repo is clean
DIRTY=$(git status --porcelain)
if [ ! -z "$DIRTY" ]; then
  echo "Repo is not clean"
  exit 1
fi

git pull origin master

LATEST_TAG=$(git describe --abbrev=0 --tags)
NEXT_VERSION=$(echo $LATEST_TAG | awk -F. -v OFS=. 'NF==1{print ++$NF}; NF>1{if(length($NF+1)>length($NF))$(NF-1)++; $NF=sprintf("%0*d", length($NF), ($NF+1)%(10^length($NF))); print}')

echo "New version: $NEXT_VERSION"

# ask user for confirmation
read -p "Are you sure? " -n 1 -r

git tag -a "$NEXT_VERSION" -m "Version $NEXT_VERSION"
git push origin "$NEXT_VERSION"

goreleaser release --clean