#!/bin/sh

set -ex

go test -race -v ./...
