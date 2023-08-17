#!/bin/bash
set -eu
s3simple() {
  local command="$1"
  local url="$2"
  local file="${3:--}"

  if [ "${url:0:5}" != "s3://" ]; then
    echo "Need an s3 url"
    return 1
  fi
  local path="${url:4}"

  if [ -z "${S3_SERVER-}" ]; then
    echo "Need S3_SERVER to be set"
    return 1
  fi

  if [ -z "${S3_ACCESS_KEY_ID-}" ]; then
    echo "Need S3_ACCESS_KEY_ID to be set"
    return 1
  fi

  if [ -z "${S3_SECRET_ACCESS_KEY-}" ]; then
    echo "Need S3_SECRET_ACCESS_KEY to be set"
    return 1
  fi

  local method md5 args params
  case "$command" in
  get)
    method="GET"
    md5=""
    args="-o $file"
    params=""
    ;;
  pre)
    method="GET"
    md5=""
    params="prefix=$(echo ${path} | sed 's/\(\/[^/]*\)\/\(.*\)$/\2/')"
    path="$(echo ${path} | sed 's/\(\/[^/]*\)\/\(.*\)$/\1/')"
    args="-o $file"
    ;;
  put)
    method="PUT"
    if [ ! -f "$file" ]; then
      echo "file not found"
      exit 1
    fi
    md5="$(openssl md5 -binary $file | openssl base64)"
    args="-T $file -H Content-MD5:$md5"
    params=""
    ;;
  *)
    echo "Unsupported command"
    return 1
  esac

  local date="$(date -u '+%a, %e %b %Y %H:%M:%S +0000')"
  local string_to_sign
  printf -v string_to_sign "%s\n%s\n\n%s\n%s" "$method" "$md5" "$date" "$path"
  local signature=$(echo -n "$string_to_sign" | openssl sha1 -binary -hmac "${S3_SECRET_ACCESS_KEY}" | openssl base64)
  local authorization="AWS ${S3_ACCESS_KEY_ID}:${signature}"
  set -x
  curl $args -f -H Date:"${date}" -H Authorization:"${authorization}" ${S3_SERVER}"${path}"?"${params}"
}

s3simple "$@"

# Anatomy of a command
# ./s3cmd.sh [method] [s3_url] [output_file->leave empty for stdout]
# List all buckets
# ./s3cmd.sh get s3:// ls.xml
# List bucket
# ./s3cmd.sh get s3://maxgds ls.xml
# List bucket with prefix
# ./s3cmd.sh pref s3://maxgds/lumen1/ ls.xml
# Get file
# ./s3cmd.sh get s3://maxgds/lumen1/28-02/cmd_tlm_db.zip cmd_tlm_db.zip
# Put file
# ./s3cmd.sh put s3://maxgds/lumen1/02-03/cmd_tlm_db.zip cmd_tlm_db.zip

# export S3_SERVER=https://flashblade.rocketlab.local
# export S3_ACCESS_KEY_ID=
# export S3_SECRET_ACCESS_KEY=
