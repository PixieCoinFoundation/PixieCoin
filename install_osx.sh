/bin/bash install_setup.sh

CURDIR=`pwd`
OLDGOPATH="$GOPATH"
export GOPATH="$CURDIR"
export GODEBUG=gctrace=0

# echo "gopath:$GOPATH"

if [ "$1" != "simple" ]; then
   	echo 'format code'
	gofmt -w src
fi

echo 'install...'
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install launcher
echo 'mv bin...'
mv bin/linux_amd64/launcher bin/
rm -rf bin/linux_amd64
rm -rf bin/.git
echo 'tar...'
tar -czf pixieserver.tar.gz bin

# echo 'copy to other bin..'

# rm -rf ~/Work/GFServer_all_bin/*
# cp -r bin/* ~/Work/GFServer_all_bin/
# cp -r ~/Work/GFServer_all_bin/odpscmd/bin/odpscmd_biliall ~/Work/GFServer_all_bin/odpscmd/bin/odpscmd

# rm -rf ~/Work/GFServer_iosts_bin/*
# cp -r bin/* ~/Work/GFServer_iosts_bin/

export GOPATH="$OLDGOPATH"

echo 'finished'
