#!/bin/bash
#
# Run this script from <c2go-root-folder>/
# to generate *.expected.c from *.c files
# using the current c2go sources.
#

go build

# Generate code quality Go code
FILES='tests/code_quality/*.c'
for file in $FILES
do
  filename=$(basename "$file")
  ext="${filename#*.}"
  if [ "$ext" = "expected.c" ];  then
	  continue
  fi

  echo "Processing $file file..."
  filename=${file%.*}".expected.c"
  ./c2go transpile -o="$filename" -p="code_quality" $file

  # Normalize transpiled from comments
  sed -i '' -E 's/^\/\/([^\/]*)(.*)tests\/code_quality\/(.*)$/\/\/\1tests\/code_quality\/\3/g' $filename
done
