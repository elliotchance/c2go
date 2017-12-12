echo "----------------------"

export SOURCE_ARCHIVE_REMOTE="https://codeload.github.com/gist/74e0d4313e8068507693d97dc8813fa1/zip/6aa6a6c1820b79986004506becd5f8a00c17b4ea"
export SOURCE="knightstour"

# Variable for location of temp files
TEMP_FOLDER="/tmp/$SOURCE"
mkdir -p $TEMP_FOLDER

# Download/unpack if required.
if [ ! -e "$TEMP_FOLDER/$SOURCE.c" ]; then
    curl -o "$TEMP_FOLDER/$SOURCE.zip" $SOURCE_ARCHIVE_REMOTE
    unzip $TEMP_FOLDER/$SOURCE.zip -d $TEMP_FOLDER
	cp "$TEMP_FOLDER/74e0d4313e8068507693d97dc8813fa1-6aa6a6c1820b79986004506becd5f8a00c17b4ea/$SOURCE.c" "$TEMP_FOLDER/$SOURCE.c"
fi

OUTFILE="$1"

if [[ -z $OUTFILE ]] ; then
    OUTFILE="$TEMP_FOLDER/out.txt"
fi

# Clean generated files. This should not be required, but it's polite.
rm -f $TEMP_FOLDER/*.go

# Transpile the C files.
# If transpiling write to stderr, then it will be append into OUTFILE
echo "Transpiling $SOURCE ..."
./c2go transpile -o=$TEMP_FOLDER/$SOURCE.go $TEMP_FOLDER/$SOURCE.c >> $OUTFILE 2>&1

# Show amount "Warning" in Go code
SQLITE_WARNINGS=`cat $TEMP_FOLDER/$SOURCE.go | grep "// Warning" | wc -l`
echo "In files ($SOURCE.go) summary : $SQLITE_WARNINGS warnings."
