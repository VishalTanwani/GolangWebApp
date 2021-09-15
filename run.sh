#!/bin/bash

go build -o GolangWebApp cmd/web/*.go 
./GolangWebApp -dbname=GolangWebApp -dbuser=vishal -cache=false -production=false