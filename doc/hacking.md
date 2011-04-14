# Hacking on Doozer

If you want to hack on doozer, the first step is to
discuss your plans on the [mailing list][mail].
Here are some instructions for building doozer from
source.

## Installing Go

I recommend not using apt, homebrew, or any other packaging system to install
Go. It's better to install straight from source.

Abridged version:

    $ GOROOT=$HOME/src/go # adjust this as you wish
    $ hg clone -r release https://go.googlecode.com/hg/ $GOROOT
    $ cd $GOROOT/src
    $ ./all.bash
    (put the value of $GOROOT/bin in your path)

Note: all.bash runs unit tests, and some of the tests
try to connect to Google servers over the internet.
If this causes problems, you can skip the network tests
by setting env var `DISABLE_NET_TESTS` before running
all.bash.

For full details, see <http://golang.org/doc/install.html>.

## Installing Dependencies

You need the `protoc` command
(from <http://code.google.com/p/protobuf/>):

    $ sudo apt-get install protobuf-compiler
    (or)
    $ brew install protobuf

If you want to run doozer's tests, install
<http://bmizerany.github.com/roundup/>.

## Building Doozer

    $ git clone https://github.com/bmizerany/doozer.git
    $ cd doozer/src
    $ ./all.sh

This will build the rest of the dependencies and
all doozer packages and commands,
copy the commands into `$GOROOT/bin`,
and run tests.

## Try It Out

    $ doozerd >/dev/null 2>&1 &
    $ open http://localhost:8080/

This will start up one doozer process and show a web view of its contents.

[mail]: https://groups.google.com/group/doozer
