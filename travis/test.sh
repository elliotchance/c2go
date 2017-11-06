#!/bin/bash

set -e

OUTFILE=/tmp/out.txt

function cleanup {
    EXIT_STATUS=$?

    if [ $EXIT_STATUS != 0 ]; then
        [ ! -f $OUTFILE ] || cat $OUTFILE
    fi

    exit $EXIT_STATUS
}
trap cleanup EXIT

echo "" > coverage.txt

# The code below was copied from:
# https://github.com/golang/go/issues/6909#issuecomment-232878416
#
# As in @rodrigocorsi2 comment above (using full path to grep due to 'grep -n'
# alias).
export PKGS=$(go list ./... | grep -v c2go/build | grep -v /vendor/)

# Make comma-separated.
export PKGS_DELIM=$(echo "$PKGS" | paste -sd "," -)

# Run tests and append all output to out.txt. It's important we have "-v" so
# that all the test names are printed. It's also important that the covermode be
# set to "count" so that the coverage profiles can be merged correctly together
# with gocovmerge.
#
# Exit code 123 will be returned if any of the tests fail.
rm -f $OUTFILE
go list -f 'go test -v -tags integration -race -covermode atomic -coverprofile {{.Name}}.coverprofile -coverpkg $PKGS_DELIM {{.ImportPath}}' $PKGS | xargs -I{} bash -c "{} >> $OUTFILE"

# Merge coverage profiles.
COVERAGE_FILES=`ls -1 *.coverprofile 2>/dev/null | wc -l`
if [ $COVERAGE_FILES != 0 ]; then
	# check program `gocovmerge` is exist
	if which gocovmerge >/dev/null 2>&1; then
		gocovmerge `ls *.coverprofile` > coverage.txt
		rm *.coverprofile
	fi
fi

# Print stats
echo "Unit tests: " $(grep "=== RUN" $OUTFILE | wc -l)
echo "Integration tests: " $(grep "# Total tests" $OUTFILE | cut -c21-)

# These steps are from the README to verify it can be installed and run as
# documented.
go build

export C2GO_DIR=$GOPATH/src/github.com/elliotchance/c2go
export C2GO=$C2GO_DIR/c2go

echo "Run: c2go transpile prime.c"
$C2GO transpile -o=/tmp/prime.go $C2GO_DIR/examples/prime.c
echo "47" | go run /tmp/prime.go
if [ $($C2GO -v | wc -l) -ne 1 ]; then exit 1; fi
if [ $(cat /tmp/prime.go | wc -l) -eq 0 ]; then exit 1; fi
if [ $($C2GO ast $C2GO_DIR/examples/prime.c | wc -l) -eq 0 ]; then exit 1; fi

# This will have to be updated every so often to the latest version. You can
# find the latest version here: https://sqlite.org/download.html
export SQLITE3_FILE=sqlite-amalgamation-3190300

# Variable for location of temp sqlite files
SQLITE_TEMP_FOLDER="/tmp/SQLITE"
mkdir -p $SQLITE_TEMP_FOLDER

# Download/unpack SQLite if required.
if [ ! -e $SQLITE_TEMP_FOLDER/$SQLITE3_FILE.zip ]; then
    curl https://sqlite.org/2017/$SQLITE3_FILE.zip > $SQLITE_TEMP_FOLDER/$SQLITE3_FILE.zip
    unzip $SQLITE_TEMP_FOLDER/$SQLITE3_FILE.zip -d $SQLITE_TEMP_FOLDER
fi

# Clean generated files. This should not be required, but it's polite.
rm -f $SQLITE_TEMP_FOLDER/sqlite3.go $SQLITE_TEMP_FOLDER/shell.go

# Transpile the SQLite3 files.
# If transpiling write to stderr, then it will be append into OUTFILE
echo "Transpiling shell.c..."
./c2go transpile -o=$SQLITE_TEMP_FOLDER/shell.go   $SQLITE_TEMP_FOLDER/$SQLITE3_FILE/shell.c   >> $OUTFILE 2>&1
echo "Transpiling sqlite3.c..."
./c2go transpile -o=$SQLITE_TEMP_FOLDER/sqlite3.go $SQLITE_TEMP_FOLDER/$SQLITE3_FILE/sqlite3.c >> $OUTFILE 2>&1

# Show amount "Warning" in sqlite Go codes
SQLITE_WARNING_SQLITE3=`cat $SQLITE_TEMP_FOLDER/sqlite3.go | grep "// Warning" | wc -l`
echo "In file sqlite3.go : $SQLITE_WARNING_SQLITE3 warnings."

SQLITE_WARNING_SHELL=`cat $SQLITE_TEMP_FOLDER/shell.go | grep "// Warning" | wc -l`
echo "In file shell.go   : $SQLITE_WARNING_SHELL warnings."
