## Smock
A system level binary mocking subsystem

### Why
System testing can be difficult and very time consuming.  It would be great if 
we could mock system level commands so that we can create better code that utilizes 
the real system with canned mocked responses.

### Goals

* mock any binary on any system for through system testing
* Written in go language so install is easy
* Should stub the PATH env variable so that when any binary is called
  it would have to check with smock to see if its being mocked.
* Should have a pluggable system so writting "fixtures" or common responses can be shared among others easily and be downloadable.
* Shim the path like RVM or rbenv
* Should also have a "Gemfile" that specifies which binaries are to be mocked and which set of fixtures they should use.
* Should be able to switch fixtures easily (maybe use ENV Variable)
* Should have a web interface to show which things are being smocked.
* This shims should make rest calls to the REST API to retrieve the response
* There should be a Web Server that takes binary requests and outputs a response
* The responses should render templates for the response
* Web server should fetch from local and then remote central library of responses.
* Should be able to turn on/off the mocking system with one command/API call
* Each fixture should be split into two types.  
   1. with Smockfile (entire list of OS binaries to mock)
   2. without Smockfile  (would be included in a Smockfile)

* The goal of the Smockfile is to easily sync changes from upstream fixtures.  The smock binary will download new changes
* Returned data of rendered fixture should fail gracefully and return smock fixture failed to render

* Stub system calls using bash_prompt and setting path to include shims
  The PROMPT_COMMAND will execute the main binary that then creates a shim specifically for the system call.
  We can then changes the PROMPT to show its a mocked response.



### Similar tools

* https://bitheap.org/cram/
* http://pbrisbin.com/posts/mocking_bash/
