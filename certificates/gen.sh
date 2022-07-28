#!/usr/bin/env bash
#
# generate certificates for mTLS
# certificates include the root ca, server and client
#
# certificates are valid for 30 days

set -euo pipefail

function generate_root_ca() {
  echo Generate Root CA

  mkdir -p "ca"

  openssl req -newkey rsa:4096 \
    -new -nodes -x509 \
    -days 30 \
    -out ca/ca.crt \
    -keyout ca/ca.key \
    -subj "/C=AU/ST=New South Wales/O=calvinbui/OU=root/CN=localhost"
}

function generate_certificate() {
  local type=$1

  mkdir -p "$type"

  echo Generate "${type}" cert
  openssl genrsa -out "${type}/${type}".key 2048

  echo Generate "${type}" key and CSR
  openssl req -new \
    -key "${type}/${type}".key \
    -days 30 \
    -out "${type}/${type}".csr \
    -subj "/C=AU/ST=New South Wales/O=calvinbui/OU=${type}/CN=localhost"

  echo Sign CSR with Root CA
  openssl x509 -req \
    -in "${type}/${type}".csr \
    -CA ca/ca.crt \
    -CAkey ca/ca.key \
    -extfile <(printf "subjectAltName=DNS:localhost") \
    -days 30 \
    -sha256 -CAcreateserial \
    -out "${type}/${type}".crt
}

generate_root_ca
generate_certificate server
generate_certificate client
