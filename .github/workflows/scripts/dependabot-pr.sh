#!/bin/zsh -ex

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

git config user.name opentelemetrybot
git config user.email 107717825+opentelemetrybot@users.noreply.github.com

BRANCH=dependabot/dependabot-prs/`date +'%Y-%m-%dT%H%M%S'`
git checkout -b $BRANCH

IFS=$'\n'
requests=($( gh pr list --search "author:app/dependabot" --json title --jq '.[].title' ))
message=""
dirs=(`find . -type f -name "go.mod" -exec dirname {} \; | sort | egrep  '^./'`)

declare -A mods

for line in $requests; do
    echo $line
    if [[ $line != Bump* ]]; then
        continue
    fi

    module=$(echo $line | cut -f 2 -d " ")
    if [[ $module == go.opentelemetry.io/otel* ]]; then
        continue
    fi
    version=$(echo $line | cut -f 6 -d " ")

    mods[$module]=$version
    message+=$line
    message+=$'\n'
done

for module version in ${(kv)mods}; do
    topdir=`pwd`
    for dir in $dirs; do
        echo "checking $dir"
        cd $dir && if grep -q "$module " go.mod; then go get "$module"@v"$version"; fi
        cd $topdir
    done
done

make go-mod-tidy
make build

git add go.sum go.mod
git add "**/go.sum" "**/go.mod"
git commit -m "dependabot updates `date`
$message"
git push origin $BRANCH

gh pr create --title "[chore] dependabot updates `date`" --body "$message"
