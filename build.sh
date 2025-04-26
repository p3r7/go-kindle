#!/usr/bin/env bash

GOARCH=arm GOARM=7 CGO_ENABLED=0 go build
