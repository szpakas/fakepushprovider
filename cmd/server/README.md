# server

Command runs fake push server.

## Usage

Configuration is following 12-factor app methodology and is using ENV as config source.

### FCM/GCM server

```bash
#!/usr/bin/env bash

export APP_SERVICE="fcm"
export APP_APPS_FILE="data/fcm-apps.json"
export APP_INSTANCES_FILE="data/fcm-inst.json"
export APP_LOG_LEVEL="all"

./server
```

### APNS server

```bash
#!/usr/bin/env bash

export APP_SERVICE="apns"
export APP_APPS_FILE="data/apns-apps.json"
export APP_INSTANCES_FILE="data/apns-inst.json"
export APP_APNS_CERT_FILE="data/self-001-cert.pem"
export APP_APNS_KEY_FILE="data/self-001-key.pem"
export APP_LOG_LEVEL="all"

./server
```

#### certificate generation

Generation of the SSL certificates which are supporting Subject Alternative Names for local use (also IP) without disabling TLS/SSL security in client.

    TODO(szpakas): add support for ASN1 commonName field (2.5.4.3) for passing APNS environment (Apple Development IOS Push Services / Apple Production IOS Push Services / Apple Push Services).

openssl.cnf
```
[req]
distinguished_name = req_distinguished_name
x509_extensions = v3_req
prompt = no

[req_distinguished_name]
# userID OID extension (0.9.2342.19200300.100.1.1) for passing bundleID for topic selection
userId = pl.example.prod
# As defined in 4.1.2.4 and appendix A of RFC 5280
countryName = PL
stateOrProvinceName = Silesia
localityName = Gliwice
organizationName = Nyota

[ v3_req ]
subjectKeyIdentifier = hash
authorityKeyIdentifier = keyid,issuer
basicConstraints = CA:TRUE
subjectAltName = @alt_names

[alt_names]
DNS.1 = localhost
DNS.2 = *.push.local
IP.1 = 192.168.99.100
IP.2 = 127.0.0.1
```

generate.sh
```bash
#!/usr/bin/env bash

# Generation of the certificate with Subject Alternative Names (SAN).
# This allows using IPs.
#
# @see:
#  http://apetec.com/support/GenerateSAN-CSR.htm
#  http://pro-tips-dot-com.tumblr.com/post/65411476159/self-signed-ssl-certificates-with-multiple-hostnames

NO=`printf "%0*d\n" 3 $1`
BASE_CERT_NAME="self-${NO}"

openssl req -x509 -nodes -days 3650 \
    -newkey rsa:2048 \
    -config openssl.cnf \
    -keyout ${BASE_CERT_NAME}-key.pem \
    -out ${BASE_CERT_NAME}-cert.pem

# create pair
cat ${BASE_CERT_NAME}-cert.pem ${BASE_CERT_NAME}-key.pem > ${BASE_CERT_NAME}-pair.pem

# examine CERT
openssl x509 -in ${BASE_CERT_NAME}-cert.pem -noout -text
```

certificate example with UID set to bundleID

    Certificate:
        Data:
            Version: 3 (0x2)
            Serial Number:
                b6:f3:af:b7:3b:b1:b5:7e
            Signature Algorithm: sha1WithRSAEncryption
            Issuer: UID=pl.example.prod, C=PL, ST=Silesia, L=Gliwice, O=Nyota
            Validity
                Not Before: Jun  7 19:44:19 2016 GMT
                Not After : Jun  5 19:44:19 2026 GMT
            Subject: UID=pl.example.prod, C=PL, ST=Silesia, L=Gliwice, O=Nyota
            Subject Public Key Info:
            ...

## Middlewares

Both APNS and FCM handlers allow for use of middlewares.
Currently following middlewares are available:
- logging
- delay
- instrumentation

## Response Headers

There are extra response headers introduced which are outside of the APNS/FCM specifications.

- X-delayed-by-ms: number of milliseconds the response was delayed (set by DelayMiddleware)
- X-requests-total: number of requests served by server till its start regardless of request outcome

## Instrumentation

Instrumentation is based on [Prometheus](https://prometheus.io).

Metrics:
- push_requests_total: 
-- type: CounterVec
-- desc: How many push requests were processed, partitioned by provider, status code and error reason.
