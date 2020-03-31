#!/bin/bash

# Copyright The OpenTelemetry Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -e

help()
{
   printf "\n"
   printf "Usage: $0 -t tag -c commit-hash\n"
   printf "\t-t New tag that you would like to create\n"
   printf "\t-c Commit hash to associate with the new tag\n"
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
   printf "Some or all of the parameters are missing\n";
   help
fi

# Validate semver
SEMVER_REGEX="^v(0|[1-9][0-9]*)\\.(0|[1-9][0-9]*)\\.(0|[1-9][0-9]*)(\\-[0-9A-Za-z-]+(\\.[0-9A-Za-z-]+)*)?(\\+[0-9A-Za-z-]+(\\.[0-9A-Za-z-]+)*)?$"


if [[ "${TAG}" =~ ${SEMVER_REGEX} ]]; then
        printf "${TAG} is valid semver tag.\n"
else
        printf "${TAG} is not a valid semver tag.\n"
        exit -1
fi

cd $(dirname $0)

# Check if the commit-hash is valid
COMMIT_FOUND=`git log -50 --pretty=format:"%H" | grep ${COMMIT_HASH}`
if [ "${COMMIT_FOUND}" != "${COMMIT_HASH}" ] ; then
	printf "Commit ${COMMIT_HASH} not found\n"
	exit -1
fi

# Check if the tag doesn't alread exists.
TAG_FOUND=`git tag --list ${TAG}`
if [ "${TAG_FOUND}" = "${TAG}" ] ; then
	printf "Tag ${TAG} already exists\n"
	exit -1
fi

# Save most recent tag for generating logs
TAG_CURRENT=`git tag | grep '^v' | tail -1`

PACKAGE_DIRS=$(find . -mindepth 2 -type f -name 'go.mod' -exec dirname {} \; | egrep -v 'tools' | sed 's/^\.\///' | sort)

# Create tag for root module
git tag -a "${TAG}" -s -m "Version ${TAG}" ${COMMIT_HASH}

# Create tag for submodules
for dir in $PACKAGE_DIRS; do
	git tag -a "${dir}/${TAG}" -s -m "Version ${dir}/${TAG}" ${COMMIT_HASH}
done

# Generate commit logs
printf "New tag ${TAG} created.\n"
printf "\n\n\nChange log since previous tag ${TAG_CURRENT}\n"
printf "======================================\n"
git --no-pager log --pretty=oneline ${TAG_CURRENT}..${TAG}

