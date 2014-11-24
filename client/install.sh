#!/bin/bash
export GOPATH=$PWD

chmod +x hooks/pre-commit
ln -s ../../client/hooks/pre-commit ../.git/hooks/
