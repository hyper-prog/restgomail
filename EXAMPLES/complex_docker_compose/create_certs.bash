#!/bin/bash

# Generate certificates for server and client

SERVERCERTNAME="restgomail"
CLIENTCERTNAME="rgmclient"

openssl genrsa -out ${SERVERCERTNAME}.key 2048
openssl ecparam -genkey -name secp384r1 -out ${SERVERCERTNAME}.key
openssl req -new -x509 -sha256 -key ${SERVERCERTNAME}.key -out ${SERVERCERTNAME}.crt -days 365

openssl genrsa -out ${CLIENTCERTNAME}.key 2048
openssl ecparam -genkey -name secp384r1 -out ${CLIENTCERTNAME}.key
openssl req -new -x509 -sha256 -key ${CLIENTCERTNAME}.key -out ${CLIENTCERTNAME}.crt -days 365

mv ${SERVERCERTNAME}.key ./data/
mv ${SERVERCERTNAME}.crt ./data/

cp ${CLIENTCERTNAME}.crt ./data/${CLIENTCERTNAME}.crt