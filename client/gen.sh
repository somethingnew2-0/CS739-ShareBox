#!/bin/bash
cd src/client/thrift
thrift -out . -r --gen go tutorial.thrift
