#!/bin/sh
# Part of E-2567's verify suite. DO NOT EDIT unless you ARE E-2567.
# Builds a throwaway module that passes a plain string to NewHostname and
# NewIPAddress, and prints the compiler's complaints; then proves a member
# type does compile, so the failure is the constraint and not the setup.
set -eu
root=$(git rev-parse --show-toplevel)
tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT
printf 'module x\n\ngo 1.27.0\n\nrequire github.com/mikeschinkel/go-dt v0.0.0\n\nreplace github.com/mikeschinkel/go-dt => %s\n' "$root" >"$tmp/go.mod"
cd "$tmp"
export GOWORK=off GOFLAGS=-mod=mod
printf 'package main\n\nimport "github.com/mikeschinkel/go-dt"\n\nfunc main() {\n\ts := "example.com"\n\t_ = dt.NewHostname(s)\n\t_ = dt.NewIPAddress(s)\n}\n' >main.go
go build ./... 2>&1 | grep -c 'string does not satisfy' || true
printf 'package main\n\nimport "github.com/mikeschinkel/go-dt"\n\nfunc main() {\n\t_ = dt.NewHostname(dt.InternetDomain("example.com"))\n\t_ = dt.NewIPAddress(dt.IPv4Address("192.0.2.1"))\n}\n' >main.go
go build ./... && echo "member types compile"
