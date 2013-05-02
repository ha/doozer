# Doozer Client

**Note:** [doozerd](https://github.com/ha/doozerd) is the server.
This is the Go client driver for doozer.

## How to use

To install the Go client:

    $ go get github.com/ha/doozer

To install the CLI client:

    $ go get github.com/ha/doozer/src/util/doozer

To use:

    package main

    import (
    	"github.com/ha/doozer"
    	"os"
    )

    func main() {
    	doozer, err := doozer.Dial("localhost:8046")

    	if err != nil {
    		panic(err)
    	}

    	myfile, _, _ := doozer.Get("/myfile", nil)
    	os.Stdout.Write(myfile)
    }


## Hacking

Building doozer utility from git:

    $ make

## License and Authors

Doozer is distributed under the terms of the MIT
License. See [LICENSE](LICENSE) for details.

Doozer was created by Blake Mizerany and Keith Rarick.
Type `git shortlog -s` for a full list of contributors.
