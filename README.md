# Doozer Client

**Note:** [doozerd](/ha/doozerd) is the server.
This is the Go client driver for doozer.

## How to use

To install:

    $ goinstall github.com/ha/doozer

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

## License and Authors

Doozer is distributed under the terms of the MIT
License. See [LICENSE][] for details.

Doozer was created by Blake Mizerany and Keith Rarick.
Type `git shortlog -s` for a full list of contributors.

[mail]: https://groups.google.com/group/doozer
[LICENSE]: /ha/doozer/blob/master/LICENSE
