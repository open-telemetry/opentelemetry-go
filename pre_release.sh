#!/bin/sh

set -e

help()
{
   printf "\n"
   printf "Usage: $0 -t tag\n"
   printf "\t-t Unreleased tag. Update all go.mod with this tag.\n"
   exit 1 # Exit script after printing help
}

while getopts "t:" opt
do
   case "$opt" in
      t ) TAG="$OPTARG" ;;
      ? ) help ;; # Print help
   esac
done

# Print help in case parameters are empty
if [ -z "$TAG" ]
then
   printf "Tag is missing\n";
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

TAG_FOUND=`git tag --list ${TAG}`
if [ ${TAG_FOUND} = ${TAG} ] ; then
        printf "Tag ${TAG} already exists\n"
        exit -1
fi


cd $(dirname $0)

# Update go.mod
git checkout -b pre_release_${TAG} master
PACKAGE_DIRS=$(find . -mindepth 2 -type f -name 'go.mod' -exec dirname {} \; | egrep -v 'tools' | sed 's/^\.\///' | sort)

for dir in $PACKAGE_DIRS; do
	sed -i .bak "s/opentelemetry.io\/otel\([^ ]*\) v[0-9]*\.[0-9]*\.[0-9]/opentelemetry.io\/otel\1 ${TAG}/" ${dir}/go.mod
	rm -f ${dir}/go.mod.bak
done

# Run lint to update go.sum
make lint

# Add changes and commit.
git add .
make ci
git commit -m "Prepare for releasing $TAG"

printf "Now run following to verify the changes.\ngit diff master\n"
printf "\nThen push the changes to upstream\n"
