#!/bin/bash

function die() {
  echo $*
  exit 1
}

# Read all local packages and save to a tmp file
tmpfile=$(mktemp)
govendor list +local | awk '{print $2}' > $tmpfile

# Initialize profile.cov
echo "mode: count" > ./profile.cov

# Initialize error tracking
ERROR=""

# Test each package and append coverage profile info to profile.cov
for pkg in `cat $tmpfile`
do
    echo $pkg
    govendor test -v -covermode=count -coverprofile=profile_tmp.cov $pkg || ERROR="Error testing $pkg"
    if [ -f profile_tmp.cov ];
    then
      tail -n +2 profile_tmp.cov >> ./profile.cov || die "Unable to append coverage for $pkg"
      rm profile_tmp.cov
    fi
done

if [ ! -z "$ERROR" ]
then
    die "Encountered error, last error was: $ERROR"
fi

echo "Generated coverage profile ./profile.cov"