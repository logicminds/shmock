#!/usr/bin/env bash

source ./preexec.bash
FAKE_BINS=bins/
if [! -e ./fake_bins]; then
	mkdir fake_bins
fi
PATH=bins/:$PATH

# function preexec () {
#     // Create binary file
       ln -s ./smock_client "${BASH_COMMAND}"
#      ./smock_client get "${BASH_COMMAND}" 

# }

# Default do-nothing implementation of precmd.
# function precmd () {
#   echo "SMOCK" 
# }

command_not_found_handle() {
	go run ./smock_client --help
	exit 128
}

preexec_install
