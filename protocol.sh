#!/bin/bash

#protoc -I /home/night/lang/go/sp /home/night/lang/go/sp/protocol.proto --go_out=/home/night/lang/go/sp/proto/
protoc --proto_path=./ --go_out=. ./protocol.proto

#protoc-gen-go: program not found or is not executable
#Please specify a program using absolute path or make sure the program is available in your PATH system variable
#--go_out: protoc-gen-go: Plugin failed with status code 1.