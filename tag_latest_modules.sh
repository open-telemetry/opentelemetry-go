#!/bin/sh

set -xe

cd $(dirname $0)

LATEST_TAG=$(git tag -l|grep '^v'|sort -V -r|head -n 1)
LATEST_REF=$(git rev-parse $LATEST_TAG)
PACKAGE_DIRS=$(find . -mindepth 2 -type f -name 'go.mod' -exec dirname {} \; | sort)

for dir in $PACKAGE_DIRS; do
	git tag -m "Submodule ${LATEST_TAG}" -s ${dir#./}/${LATEST_TAG} ${LATEST_REF}
done
