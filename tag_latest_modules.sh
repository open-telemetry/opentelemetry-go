#!/bin/sh

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

set -xe

cd $(dirname $0)

LATEST_TAG=$(git tag -l|grep '^v'|sort -V -r|head -n 1)
LATEST_REF=$(git rev-parse $LATEST_TAG)
PACKAGE_DIRS=$(find . -mindepth 2 -type f -name 'go.mod' -exec dirname {} \; | sort)

for dir in $PACKAGE_DIRS; do
	git tag -m "Submodule ${LATEST_TAG}" -s ${dir#./}/${LATEST_TAG} ${LATEST_REF}
done
