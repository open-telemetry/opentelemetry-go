#!/bin/sh

set -xe

help()
{
   echo ""
   echo "Usage: $0 -t tag"
   echo "\t-t Unreleased tag. Update all go.mod with this tag."
   exit 1 # Exit script after printing help
}

while getopts "t:c:" opt
do
   case "$opt" in
      t ) TAG="$OPTARG" ;;
      ? ) help ;; # Print help
   esac
done

# Print help in case parameters are empty
if [ -z "$TAG" ]
then
   echo "Tag is missing";
   help
fi

TAG_FOUND=`git tag --list ${TAG}`
echo "TAG found $TAG_FOUND"
if [ ${TAG_FOUND} = ${TAG} ] ; then
        echo "Tag ${TAG} already exists"
        exit -1
fi


cd $(dirname $0)

# Update go.mod
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

echo "Now run following to verify the changes.\ngit diff master"
echo "\nThen push the changes to upstream"
