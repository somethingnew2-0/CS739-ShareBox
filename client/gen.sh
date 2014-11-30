#!/bin/bash
cd thrift
thrift -out ../src/thrift -r --gen go tutorial.thrift
