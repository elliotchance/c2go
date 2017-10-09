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
echo "PKGS : $PKGS"

# Make comma-separated.
export PKGS_DELIM=$(echo "$PKGS" | paste -sd "," -)
echo "PKGS_DELIM : $PKGS_DELIM"

# Run tests and append all output to out.txt. It's important we have "-v" so
# that all the test names are printed. It's also important that the covermode be
# set to "count" so that the coverage profiles can be merged correctly together
# with gocovmerge.
#
# Exit code 123 will be returned if any of the tests fail.
echo "Run: go test"
rm -f $OUTFILE

while read -r line
do
	echo "Starting : $line"
	xargs -I{} bash -c "{} >> $OUTFILE"
done < "go list -f 'go test -v -tags integration -race -covermode atomic -coverprofile {{.Name}}.coverprofile -coverpkg $PKGS_DELIM {{.ImportPath}}' $PKGS" 
#GOLIST=/tmp/golist.txt
#go list -f 'go test -v -tags integration -race -covermode atomic -coverprofile {{.Name}}.coverprofile -coverpkg $PKGS_DELIM {{.ImportPath}}' $PKGS > $GOLIST
#echo "Show go list fully:"
#cat $GOLIST
#while read -r line
#do 
#	echo "Starting : $line"
#	xargs -I{} bash -c "{} >> $OUTFILE" "$line" 
#done < "$GOLIST"
#rm $GOLIST

# Merge coverage profiles.
echo "Run: cover profile"
COVERLIST=/tmp/coverlist.txt
ls -la *.coverprofile > $COVERLIST
echo "Show coverprofile files:"
cat $COVERLIST
COVERAGE_FILES=`ls -1 *.coverprofile 2>/dev/null | wc -l`
echo "Cover files : $COVERAGE_FILES"
if [ $COVERAGE_FILES != 0 ]; then
    gocovmerge `ls *.coverprofile` > coverage.txt
	echo "Show summary coverage doc:"
	cat coverage.txt
    rm *.coverprofile
fi
rm $COVERLIST

# Print stats
echo "Run: print stats"
echo "Unit tests: " $(grep "=== RUN" $OUTFILE | wc -l)
echo "Integration tests: " $(grep "# Total tests" $OUTFILE | cut -c21-)

# Remove the outfile so it is not printed when an error happens beyond this
# point.
rm $OUTFILE

# These steps are from the README to verify it can be installed and run as
# documented.
echo "Run: go build"
go build

export C2GO_DIR=$GOPATH/src/github.com/elliotchance/c2go
export C2GO=$C2GO_DIR/c2go

echo "Run: transpile prime.c"
$C2GO transpile -V $C2GO_DIR/examples/prime.c
echo "47" | go run prime.go
if [ $($C2GO -v | wc -l) -ne 1 ]; then exit 1; fi
if [ $(cat prime.go | wc -l) -eq 0 ]; then exit 1; fi
if [ $($C2GO ast $C2GO_DIR/examples/prime.c | wc -l) -eq 0 ]; then exit 1; fi

# This will have to be updated every so often to the latest version. You can
# find the latest version here: https://sqlite.org/download.html
export SQLITE3_FILE=sqlite-amalgamation-3190300

# Download Sqlite3 amalgamated source.
curl https://sqlite.org/2017/$SQLITE3_FILE.zip > /tmp/$SQLITE3_FILE.zip
unzip /tmp/$SQLITE3_FILE.zip -d /tmp

# Transpile the SQLite3 files.
# Add flag `-keep-unused` because linter have too many warning and
# removing unused elements is not provided.
./c2go transpile -V -keep-unused /tmp/sqlite-amalgamation-3190300/shell.c
./c2go transpile -V -keep-unused /tmp/sqlite-amalgamation-3190300/sqlite3.c
