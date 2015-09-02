#!/bin/bash -x

GOARCH=amd64 GOOS=linux go build $@
