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
UNIT_TESTS=$(grep "=== RUN" $OUTFILE | wc -l | tr -d '[:space:]')
INT_TESTS=$(grep "# Total tests" $OUTFILE | cut -c21- | tr -d '[:space:]')

echo "Unit tests: ${UNIT_TESTS}"
echo "Integration tests: ${INT_TESTS}"

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

echo "----------------------"
# This will have to be updated every so often to the latest version. You can
# find the latest version here: https://sqlite.org/download.html
export SQLITE3_FILE=sqlite-amalgamation-3240000

# Variable for location of temp sqlite files
SQLITE_TEMP_FOLDER="/tmp/SQLITE"
mkdir -p $SQLITE_TEMP_FOLDER

# Download/unpack SQLite if required.
if [ ! -e $SQLITE_TEMP_FOLDER/$SQLITE3_FILE.zip ]; then
    curl https://sqlite.org/2018/$SQLITE3_FILE.zip > $SQLITE_TEMP_FOLDER/$SQLITE3_FILE.zip
    unzip $SQLITE_TEMP_FOLDER/$SQLITE3_FILE.zip -d $SQLITE_TEMP_FOLDER
fi

# Clean generated files. This should not be required, but it's polite.
rm -f $SQLITE_TEMP_FOLDER/sqlite3.go $SQLITE_TEMP_FOLDER/shell.go

# Transpile the SQLite3 files.
# If transpiling write to stderr, then it will be append into OUTFILE
# shell.c
echo "Transpiling shell.c..."
./c2go transpile -o=$SQLITE_TEMP_FOLDER/shell.go   $SQLITE_TEMP_FOLDER/$SQLITE3_FILE/shell.c   >> $OUTFILE 2>&1

# sqlite3.c
echo "Transpiling sqlite3.c..."
./c2go transpile -o=$SQLITE_TEMP_FOLDER/sqlite3.go $SQLITE_TEMP_FOLDER/$SQLITE3_FILE/sqlite3.c >> $OUTFILE 2>&1

# Show amount "Warning" in sqlite Go codes
SQLITE_WARNINGS=`cat $SQLITE_TEMP_FOLDER/sqlite3.go $SQLITE_TEMP_FOLDER/shell.go | grep "// Warning" | wc -l`
echo "In files (sqlite3.go and shell.go) summary : $SQLITE_WARNINGS warnings."

# Update Github PR statuses. These two statuses will always pass but will show
# information about the number of tests run and how many warnings are generated
# in the SQLite3 transpile.
if [ "$TRAVIS_OS_NAME" == "osx" ]; then
    curl -H "Authorization: token ${GITHUB_API_TOKEN}" -H "Content-Type: application/json" https://api.github.com/repos/elliotchance/c2go/statuses/${TRAVIS_COMMIT} -d "{\"state\": \"success\",\"target_url\": \"https://travis-ci.org/elliotchance/c2go/builds/${TRAVIS_JOB_ID}\", \"description\": \"$(($UNIT_TESTS + $INT_TESTS)) tests passed (${UNIT_TESTS} unit + ${INT_TESTS} integration)\", \"context\": \"c2go/tests\"}"

    curl -H "Authorization: token ${GITHUB_API_TOKEN}" -H "Content-Type: application/json" https://api.github.com/repos/elliotchance/c2go/statuses/${TRAVIS_COMMIT} -d "{\"state\": \"success\",\"target_url\": \"https://travis-ci.org/elliotchance/c2go/builds/${TRAVIS_JOB_ID}\", \"description\": \"$(($SQLITE_WARNINGS)) warnings\", \"context\": \"c2go/sqlite3\"}" 
fi

# SQLITE
c2go transpile -o="$SQLITE_TEMP_FOLDER/sqlite.go" -clang-flag="-DSQLITE_THREADSAFE=0" -clang-flag="-DSQLITE_OMIT_LOAD_EXTENSION" $SQLITE_TEMP_FOLDER/$SQLITE3_FILE/shell.c $SQLITE_TEMP_FOLDER/$SQLITE3_FILE/sqlite3.c

# Show amount "Warning":
SQLITE_WARNINGS=`cat $SQLITE_TEMP_FOLDER/sqlite.go | grep "// Warning" | wc -l`
echo "After transpiling shell.c and sqlite3.c together, have summary: $SQLITE_WARNINGS warnings."

# Show amount error from `go build`:
SQLITE_WARNINGS_GO=`go build -gcflags="-e" $SQLITE_TEMP_FOLDER/sqlite.go 2>&1 | wc -l`
echo "In file sqlite.go summary : $SQLITE_WARNINGS_GO warnings in go build."

# Amount warning from gometalinter
echo "Calculation warnings by gometalinter"
GOMETALINTER_WARNINGS=`$GOPATH/bin/gometalinter $SQLITE_TEMP_FOLDER/sqlite.go 2>&1 | wc -l`
echo "Amount found warnings by gometalinter at 30 second : $GOMETALINTER_WARNINGS warnings."
