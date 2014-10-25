#!/bin/sh

go test -coverprofile=coverage.out -coverpkg encoding/json
go tool cover -html=coverage.out
