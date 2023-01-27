#!/bin/sh

bail() {
  echo 'Error executing command, exiting'
  exit 1
}

exec_cmd_nobail() {
  echo "+ $1"
  sh -c "$1"
}

exec_cmd() {
  exec_cmd_nobail "$1" || bail
}

rm -vrf ./bin
mkdir -vp ./bin
PKGS=$(go list -find -f '{{.ImportPath}}' .)
for pkg in $PKGS; do
  exec_cmd "go build -race -ldflags '-extldflags "'"-static"'"' -o ./bin/ $pkg"
done
