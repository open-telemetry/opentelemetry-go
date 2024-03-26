#!/bin/zsh -ex

# Copyright The OpenTelemetry Authors
# SPDX-License-Identifier: Apache-2.0

dirs=(`find . -type d -not -path "*/internal*" -not -path "*/test*" -not -path "*/example*" -not -path "*/.*" | sort`)
topdir=`pwd`

for dir in $dirs; do
	echo "checking $dir"

	cd $dir
	pwd
	if [ ! -f "README.md" ]; then
		echo "couldn't find README.md for $dir"
		exit 1
	fi
	cd $topdir
done
