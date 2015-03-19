## Golang Environment Setup

* Version 1.3
* Install from source (http://golang.org/doc/install/source)
  1. http://golang.org/dl/go1.3.src.tar.gz
  2. unzip and cd to go/src
  3. run ./all.bash

* Setup cross compilation for linux
  1. set environment variable GOOS=linux
  2. set environment variable GOARCH=amd64
  3. Rebuild with ./all.bash

