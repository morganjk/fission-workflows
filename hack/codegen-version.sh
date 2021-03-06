#!/usr/bin/env bash

# Handle arguments
while [[ $# -gt 0 ]]; do
  case $1 in
    -d| --date)
      GIT_DATE=$2
      shift
      ;;
    -c| --commit)
      GIT_COMMIT=$2
      shift
      ;;
    -v| --version)
      VERSION=$2
      shift
      ;;
    -o| --output)
      OUTPUT=$2
      shift
      ;;
    -h| --help)
      echo "usage: codegen-version.sh [-c commit] [-d RFC3339] [-v version]"
      exit 0;
  esac
  shift
done

if [ -z ${VERSION:-} ]; then
  VERSION=$(git rev-parse HEAD)
fi

if [ -z ${GIT_DATE:-} ] ; then
  # or use `date -R`
  GIT_DATE=$(git show -s --format="%cD" HEAD)
fi

if [ -z ${GIT_COMMIT:-} ] ; then
  GIT_COMMIT=$(git rev-parse HEAD)
fi

read -r -d '' VERSION << EOM
// Code generated by hack/codegen-version.sh. DO NOT EDIT.
package version

import "time"

const (
	dateFormat string = time.RFC1123Z

	// Git commit (e.g. 1b4716ab84903b2e477135a3dc5afdb07f685cb7)
	GitCommit string = "${GIT_COMMIT}"

	// Version contains a (potentially) human-readable version
	// For example 1.1.0 or 1b4716ab84903b2e477135a3dc5afdb07f685cb7
	Version string = "${VERSION}"

	// gitDate is a date in RFC1123Z format
	gitDate string = "${GIT_DATE}"
)

var (
	GitDate time.Time
)

func init() {

	d, err := time.Parse(dateFormat, gitDate)
	if err != nil {
		panic(err)
	}
	GitDate = d
}


EOM

if [ -z ${OUTPUT:-} ] ; then
    echo "${VERSION}"
else
    echo "${VERSION}" > ${OUTPUT}
fi