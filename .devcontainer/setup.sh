#!/bin/bash

go install github.com/go-task/task/v3/cmd/task@latest
go install github.com/jesseduffield/lazygit@latest

sudo apt update
sudo apt install -y python3.10-venv

# Mkdocs Social Cards Native Dependencies
sudo apt install -y libcairo2-dev libfreetype6-dev libffi-dev libjpeg-dev libpng-dev libz-dev