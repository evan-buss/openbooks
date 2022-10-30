#!/bin/bash

go install github.com/go-task/task/v3/cmd/task@latest

sudo apt update
sudo apt install python3.10-venv