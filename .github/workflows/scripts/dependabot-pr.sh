#!/bin/zsh -ex

# Copyright The OpenTelemetry Authors
# SPDX-License-Identifier: Apache-2.0

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
    if [[ $line != build* ]]; then
        continue
    fi

    module=$(echo $line | cut -f 3 -d " ")
    if [[ $module == go.opentelemetry.io/otel* ]]; then
        continue
    fi
    version=$(echo $line | cut -f 7 -d " ")

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

make go-mod-tidy golangci-lint-fix build

git add go.sum go.mod
git add "**/go.sum" "**/go.mod"
git commit -m "dependabot updates `date`
$message"
git push origin $BRANCH

gh pr create --title "[chore] dependabot updates `date`" --body "$message"
