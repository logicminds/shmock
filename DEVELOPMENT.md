### Goals

* mock any binary on any system for rigorous system testing
* Should have a pluggable system so writting "fixtures" or common responses can be shared among others easily and be downloadable.
* Should also have a "Gemfile" that specifies which binaries are to be mocked and which set of fixtures they should use.
* Should be able to switch fixtures easily (maybe use ENV Variable)
* Should have a web interface to show which things are being smocked.
* Provide a way to shim the environment to use "smocked" calls to the system by prefixing the PATH with location of shims
* This shims should make rest calls to the REST API to retrieve the response
* There should be a Web Server that takes smock requests and outputs a response
* The responses should render templates for the response
* Web server should fetch from local and then remote central library of responses.
* Should be able to turn on/off the mocking system with one command/API call
* Each fixture should be split into two types.
   1. with Smockfile (entire list of OS binaries to mock)
   2. without Smockfile  (would be included in a Smockfile)

* The goal of the Smockfile is to easily sync changes from external fixtures.  The smock binary will download new changes
* Returned data of rendered fixture should fail gracefully and return smock fixture failed to render

* Stub system calls using bash_prompt and setting path to include shims
  The PROMPT_COMMAND will execute the main binary that then creates a shim specifically for the system call.
  We can then changes the PROMPT to show its a mocked response.

* Using in testing code
   Allow the user to set the context, like environment variables, and response namespaces which would be passed in the body
   or head to allow for different types of responses for a single command.

* Pass in the Environment using os.Environ in the headers
* Pass in the namespace of the command in the header

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

