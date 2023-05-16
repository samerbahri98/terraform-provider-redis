#!/usr/bin/env bash

set -e

if [[ ! -f CHANGELOG.md ]]; then
  echo "ERROR: CHANGELOG.md not found in pwd."
  echo "Please run this from the root of the terraform provider repository"
  exit 1
fi

if [[ $(uname) == "Darwin" ]]; then
  echo "Using BSD sed"
  SED="sed -i.bak -E -e"
else
  echo "Using GNU sed"
  SED="sed -i.bak -r -e"
fi

PROVIDER_URL="https:\/\/github.com\/terraform-providers\/terraform-provider-digitalocean\/issues"

$SED "s/GH-([0-9]+)/\[#\1\]\($PROVIDER_URL\/\1\)/g" -e 's/\[\[#(.+)([0-9])\)]$/(\[#\1\2))/g' CHANGELOG.md

rm CHANGELOG.md.bak
