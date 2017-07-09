#!/bin/bash

set -e
echo "" > coverage.txt

# The code below was copied from:
# https://github.com/golang/go/issues/6909#issuecomment-232878416
#
# As in @rodrigocorsi2 comment above (using full path to grep due to
# 'grep -n' alias).
export PKGS=$(go list ./... | grep -v c2go/build | grep -v /vendor/)

# Make comma-separated.
export PKGS_DELIM=$(echo "$PKGS" | paste -sd "," -)

# Run tests and append all output to out.txt. It's important we have "-v"
# so that all the test names are printed. It's also important that the
# covermode be set to "count" so that the coverage profiles can be merged
# correctly together with gocovmerge.
#
# Exit code 123 will be returned if any of the tests fail.
rm -f /tmp/out.txt
go list -f 'go test -v -tags=integration -race -covermode count -coverprofile {{.Name}}.coverprofile -coverpkg $PKGS_DELIM {{.ImportPath}}' $PKGS |
  xargs -I{} bash -c '{}'

# Merge coverage profiles.
COVERAGE_FILES=`ls -1 *.coverprofile 2>/dev/null | wc -l`
if [ $COVERAGE_FILES != 0 ]; then
  gocovmerge `ls *.coverprofile` > coverage.txt
  rm *.coverprofile
fi

# Print stats
echo "Unit tests: " $(grep "=== RUN" /tmp/out.txt | wc -l)
echo "Integration tests: " $(grep "# Total tests" /tmp/out.txt | cut -c21-)

# These steps are from the README to verify it can be installed and run as
# documented.
- cd /tmp
- export C2GO=$GOPATH/src/github.com/elliotchance/c2go
- c2go transpile $C2GO/examples/prime.c
- echo "47" | go run prime.go
- if [ $(c2go -v | wc -l) -ne 1 ]; then exit 1; fi
- if [ $(cat prime.go | wc -l) -eq 0 ]; then exit 1; fi
- if [ $(c2go ast $C2GO/examples/prime.c | wc -l) -eq 0 ]; then exit 1; fi
