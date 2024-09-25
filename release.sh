#!/usr/bin/env bash

# make sure that the repo is clean
DIRTY=$(git status --porcelain)
if [ ! -z "$DIRTY" ]; then
  echo "Repo is not clean"
  exit 1
fi

git pull origin master

# get the latest tag
LATEST_TAG=$(git describe --abbrev=0 --tags)

# increment the minor version
NEXT_VERSION=$(echo $LATEST_TAG | awk -F. -v OFS=. 'NF==1{print ++$NF}; NF>1{if(length($NF+1)>length($NF))$(NF-1)++; $NF=sprintf("%0*d", length($NF), ($NF+1)%(10^length($NF))); print}')

echo "New version: $NEXT_VERSION"

git tag -a "$NEXT_VERSION" -m "Version $NEXT_VERSION"
git push origin "$NEXT_VERSION"

goreleaser release --clean