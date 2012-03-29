#!/bin/bash
set -e
./clean.sh
./all.sh
base=`./cmd/doozer/doozer -v|tr ' ' -`
trap "rm -rf $base" 0
mkdir $base
cp cmd/doozer/doozer $base
cat <<end >$base/README
This is the command line client for Doozer,
a consistent, fault-tolerant data store.

See http://github.com/ha/doozer
and http://github.com/ha/doozerd
end
file=$base-$GOOS-$GOARCH.tar
tar cf $file $base
gzip $file
