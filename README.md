## Shmock
A system level mocking subsystem.

## Why
System testing can be difficult and very time consuming.  Sometimes when your run integration tests you need to test
agains real hardware as such as case with IPMI.  It would be great if
we could mock system level commands so that we can create better code that utilizes the real system with canned mocked responses.

As a developer I am tired of mocking system calls and wish there was a way for my code to get a real result without
having to mock objects and rely on fixtures to supply the data.  I want a real response and real values that can be attained
easily.  There should be a standard way of sharing system command fixtures.

## Use Cases
* Projects that make system calls need real responses in a variety of conditions in order to make the code base solid.  Sometimes
  unexpected results occur outside the development environment.  Having the ability to use responses that are created outside
  the development environment can be extremely valuable during integration tests.  Additional, some system commands require
  real devices that not everyone has.  So having the ability to mock these real devices with real output from many types of systems
  makes testing easier.

* Some products are written entirely for the CLI and thus should allow the user to practice the commands in a fake environment
  to test out the app or library.

## Features

### List Commands
The shmock client has the ability to list the commands that are currently being "shmocked".  Additional if you specify the
namespace you will get a list of just the commands under that namespace.

### System Response Fixtures (a.k.a Shmocks)
A shmock fixture is basically a plain text file encoded in JSON format.  The use of JSON files allows any person to easily
tweak the output of the command.  Additionally these files can be grouped together and shared among other shmock users through
traditional file sharing methods.

Example Shmock response file

```json
{
   "0b560112b718db4dfd5f67a687035afe73f33a69": {
         "Stdout": "",
         "Stderr": "",
         "Exitcode": 0,
         "Delay": 0,
         "Stdin": "/usr/local/sbin/ipmi-chassis --hostname=192.168.1.16 --chassis-identify=5 --config-file=/var/folders/h6/v6nv76td37s7vqj902_z59kh0000gn/T/20150323-59411-f3rtjv --driver-type=LAN_2_0"
      },
 }
```

Example Shmocks files struture

```shell
   ├── shmocks
   │   ├── biz
   │   │   └── logicminds
   │   │       └── rubyipmi
   │   │           ├── bin
   │   │           │   └── echo.tmpl
   │   │           ├── usr
   │   │           │   └── local
   │   │           │       └── sbin
   │   │           │           └── ipmi-chassis.tmpl
   │   │           └── which.tmpl
   │   ├── internal_responses.json
   │   └── which.tmpl

```

### Sharable Shmock Fixtures
Shmock fixtures are easily sharable since they are just json files that contain the output of a system command.  These files
can be shared via traditional methods through VCS repos, tar files, zip files.  Each set of fixtures can be namespaced
to allow for different kinds of responses.

### Load Command
The shmock client has the ability to load new shmocks from a remote/local destination.  This allows the user to pull
down new shmocks manually.

### Shmockfile
The Shmockfile provides the shmock client or shmock server with a automated way to sync shmocks.  This is analogous to using a Gemfile to pull down Gems.
Underneath the Shmockfile we are just using the load command to pull down remote / local shmocks.

### Shmock Client
The shmock client is a CLI app in Golang that can be run on the command line to get shmocked commands or capture existing
commands that have not been shmocked yet.

### Shmock Capture
The shmock capture CLI app is a subset of the shmock client but is designed to capture commands as the default behavior
and can be used inside your test suite.

### Similar tools
* https://bitheap.org/cram/
* http://pbrisbin.com/posts/mocking_bash/

### Issues
