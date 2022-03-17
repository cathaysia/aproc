#!/bin/bash

go mod tidy
go build
sudo cp ./aproc.service /usr/lib/systemd/system
sudo cp ./aproc /usr/bin
