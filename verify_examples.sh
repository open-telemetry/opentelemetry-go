#!/bin/sh
# Copyright 2020, OpenTelemetry Authors
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

cd $(dirname $0)

if [ -z "${GOPATH}" ] ; then
	printf "GOPATH is not defined.\n"
	exit -1
fi

if [ ! -d "${GOPATH}" ] ; then
	printf "GOPATH ${GOPATH} is invalid \n"
	exit -1
fi

DIR_TMP="${GOPATH}/src/oteltmp/"
rm -rf $DIR_TMP
mkdir -p $DIR_TMP

printf "Copy examples to ${DIR_TMP}\n"
rsync -r ./example ${DIR_TMP}

# Update go.mod files
printf "Update go.mod: rename module and remove replace\n"

PACKAGE_DIRS=$(find . -mindepth 2 -type f -name 'go.mod' -exec dirname {} \; | egrep 'example' | sed 's/^\.\///' | sort)
for dir in $PACKAGE_DIRS; do
	printf "  Update go.mod for $dir\n"
	#rsync -r ${dir} ${DIR_TMP}/example
	(cd "${DIR_TMP}/${dir}" && \
         sed -i .bak "s/module go.opentelemetry.io\/otel/module oteltmp/" go.mod && \
         sed -i .bak "s/^.*=\>.*$//" go.mod && \
	 go mod tidy)
done
printf "Update done:\n\n"

# Build directories that contain main package. These directories are different than 
# directories that contain go.mod files.
# 
printf "Build examples:\n"
EXAMPLES=`./get_main_pkgs.sh ./example`
for ex in $EXAMPLES; do
	printf "  Build $ex in ${DIR_TMP}/${ex}\n"
	(cd "${DIR_TMP}/${ex}" && \
	 go build .)
done

# Cleanup
printf "Remove copied files.\n"
rm -rf $DIR_TMP
