#!/bin/bash

set -e

# This will have to be updated every so often to the latest version. You can
# find the latest version here: https://sqlite.org/download.html
export SQLITE3_FILE=sqlite-amalgamation-3190300

# Download Sqlite3 amalgamated source.
curl https://sqlite.org/2017/$SQLITE3_FILE.zip > /tmp/$SQLITE3_FILE.zip
unzip /tmp/$SQLITE3_FILE.zip -d /tmp

# Build c2go and transpile the SQLite3 files.
go build
./c2go transpile /tmp/sqlite-amalgamation-3190300/shell.c
./c2go transpile /tmp/sqlite-amalgamation-3190300/sqlite3.c
