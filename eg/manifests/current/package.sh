#!/bin/bash
#####
# Packages the given artifacts into a single gziped tarbal and a MANIFEST file
#####

set -e

#############
# Functions #
#############

usage() {
cat <<EOF
$0 [options] PACKAGE VERSION -- ARTIFACT_PATH[S]

Options:
  -n	dry-run script
  -h	prints this message

EOF
}

parse_args() {
	_dry_run=""

	while getopts "hn" opt ; do
		case "$opt" in
			h)
				usage
				exit 0
				;;
			n)
				_dry_run=1
				shift
				;;
		esac
	done

	if [[ $# -lt 2 ]]; then
		usage
		exit 1
	fi

	PACKAGE="$1"

	VERSION="$2"

	shift
	shift

	if [[ "$1" != "--" ]]; then
		echo "Missing Artifacts separator"
		usage
		exit 1
	fi
	shift

	ARTIFACTS="$@"
}

manifest_file() {
# Outputs the manifest.yaml file to STDOUT


# Header
cat <<EOF
name: $PACKAGE
version: $VERSION
EOF

}

sha1sum_file() {
	# Outputs the sha1sums of artifacts to STDOUT
	for a in $ARTIFACTS; do
		find "$a" -type f | sort | xargs sha1sum
	done
}

################
# Start Script #
################

parse_args $@

if [[ -n $_dry_run ]]; then
	echo "### Manifest File"

	manifest_file

	echo ""
	echo "### Artifacts"
	for i in $ARTIFACTS; do
		echo $i
	done

	echo ""
	echo "### Sha1sum"
	sha1sum_file
	exit 0

fi

TAR_FILE="${PACKAGE}_${VERSION}.tgz"

echo "### Cleaning up old Metadata files ###"
rm -f $TAR_FILE
rm -rf .bundle

echo "### Generating Metadata files ###"
mkdir .bundle
manifest_file > ./.bundle/MANIFEST.yaml
sha1sum_file | tee ./.bundle/sha1sums.txt

echo "### Creating Bundle File ###"
tar -czvf "${TAR_FILE}" .bundle $ARTIFACTS
