#!/bin/bash

cd ..
go test -v -coverprofile cp.out ./...
go tool cover -html=cp.out
