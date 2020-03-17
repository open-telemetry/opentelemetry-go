#!/bin/bash

set -e

help()
{
   echo ""
   echo "Usage: $0 -t tag -c commit-hash"
   echo "\t-t New tag that you would like to create"
   echo "\t-c Commit hash to associate with the new tag"
   exit 1 # Exit script after printing help
}

while getopts "t:c:" opt
do
   case "$opt" in
      t ) TAG="$OPTARG" ;;
      c ) COMMIT_HASH="$OPTARG" ;;
      ? ) help ;; # Print help
   esac
done

# Print help in case parameters are empty
if [ -z "$TAG" ] || [ -z "${COMMIT_HASH}" ]
then
   echo "Some or all of the parameters are missing";
   help
fi

cd $(dirname $0)

# Check if the commit-hash is valid
COMMIT_FOUND=`git log -50 --pretty=format:"%H" | grep ${COMMIT_HASH}`
if [ "${COMMIT_FOUND}" != "${COMMIT_HASH}" ] ; then
	echo "Commit ${COMMIT_HASH} not found"
	exit -1
fi

# Check if the tag doesn't alread exists.
TAG_FOUND=`git tag --list ${TAG}`
if [ "${TAG_FOUND}" = "${TAG}" ] ; then
	echo "Tag ${TAG} already exists"
	exit -1
fi

# Save most recent tag for generating logs
TAG_CURRENT=`git tag | grep '^v' | tail -1`

PACKAGE_DIRS=$(find . -mindepth 2 -type f -name 'go.mod' -exec dirname {} \; | egrep -v 'tools' | sed 's/^\.\///' | sort)

# Create tag for root module
git tag -a "${TAG}" -m "Version ${TAG}" ${COMMIT_HASH}

# Create tag for submodules
for dir in $PACKAGE_DIRS; do
	git tag -a "${dir}/${TAG}" -m "Version ${dir}/${TAG}" ${COMMIT_HASH}
done

# Generate commit logs
git log --pretty=oneline ${TAG_CURRENT}..${TAG}
echo -e "New tag ${TAG} created"
echo -e "\n\n\nChange log since previous tag ${TAG_CURRENT}"
echo -e "======================================\n"
git --no-pager log --pretty=oneline ${TAG_CURRENT}..${TAG}

