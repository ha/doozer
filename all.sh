#!/bin/sh
set -e
if [ -f env.sh ]
then . ./env.sh
else
    echo 1>&2 "! $0 must be run from the root directory"
    exit 1
fi

xcd() {
    echo
    cd $1
    echo --- cd $1
}

mk() {
    set -e
    xcd $1
    shift
    $*
}

make install
for cmd in $CMDS
do (mk cmd/$cmd make install)
done
