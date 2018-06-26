#!/bin/bash

set -e

# These steps are from the README to verify it can be installed and run as
# documented.
go build


# Variable for location of temp sqlite files
export TRIANGLE_TEMP_FOLDER="/tmp/TRIANGLE"
mkdir -p $TRIANGLE_TEMP_FOLDER

export TRIANGLE_FILE="triangle"

# Download/unpack SQLite if required.
if [ ! -e $TRIANGLE_TEMP_FOLDER/$TRIANGLE_FILE.zip ]; then
    curl http://www.netlib.org/voronoi/$TRIANGLE_FILE.zip > $TRIANGLE_TEMP_FOLDER/$TRIANGLE_FILE.zip
    unzip $TRIANGLE_TEMP_FOLDER/$TRIANGLE_FILE.zip -d $TRIANGLE_TEMP_FOLDER
fi

# Clean generated files. This should not be required, but it's polite.
rm -f $TRIANGLE_TEMP_FOLDER/$TRIANGLE_FILE.go

# Transpile files.
echo "Transpiling $TRIANGLE_FILE.c..."
./c2go transpile -o=$TRIANGLE_TEMP_FOLDER/$TRIANGLE_FILE.go $TRIANGLE_TEMP_FOLDER/$TRIANGLE_FILE.c

# Show amount "Warning" in Go codes
TRIANGLE_WARNINGS=`cat $TRIANGLE_TEMP_FOLDER/$TRIANGLE_FILE.go | grep "// Warning" | wc -l`
echo "In file $TRIANGLE_FILE summary : $TRIANGLE_WARNINGS warnings."

# Show amount error from `go build`:
TRIANGLE_WARNINGS_GO=`go build -gcflags="-e" $TRIANGLE_TEMP_FOLDER/$TRIANGLE_FILE.go 2>&1 | wc -l`
echo "In file $TRIANGLE_FILE summary : $TRIANGLE_WARNINGS_GO warnings in go build."

# Amount warning from gometalinter
echo "Calculation warnings by gometalinter"
TRIANGLE_GOMETALINTER_WARNINGS=`$GOPATH/bin/gometalinter $TRIANGLE_TEMP_FOLDER/$TRIANGLE_FILE.go 2>&1 | wc -l`
echo "Amount found warnings by gometalinter at 30 second : $TRIANGLE_GOMETALINTER_WARNINGS warnings."
