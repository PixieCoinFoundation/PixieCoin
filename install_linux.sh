/bin/bash install_setup.sh linux

CURDIR=`pwd`
OLDGOPATH="$GOPATH"
export GOPATH="$CURDIR"
export GODEBUG=gctrace=0

gofmt -w src
go install launcher
#go install test

export GOPATH="$OLDGOPATH"

echo 'finished'
