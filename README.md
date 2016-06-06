# Fake push provider [![Build Status](https://travis-ci.org/szpakas/fakepushprovider.svg?branch=master)](https://travis-ci.org/szpakas/fakepushprovider)

[![Apache 2.0 License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](https://github.com/szpakas/fakepushprovider/blob/master/LICENSE)

`fakepushprovider` is a replacement for the push providers to be used in integration testing of push clients.

Project is in very early stage. Most of the functionality is missing and API breaking changes will happen.

## Overview

There are two major modules: push server and generator. Generator is producing apps and instances (devices) database which can than be shared between fake push server and client under test.

As of now support for FCM(GCM)/HTTP and APNS/2 (HTTP/2) is provided. Support for FCM(GCM)/CCS (XMPP) is planned.

## Commands

- [apps/instances generator](./cmd/generator)
- [fake push server](./cmd/server)
- example of [FCM/HTTP client](./cmd/test-gcmhttp).

## Running tests

    $ go test -v ./apns/...
    $ go test -v ./fcm/...

## Using docker

build locally

    $ docker build -t fakepushprovider .

or pull from docker hub

    $ docker pull szpakas/fakepushprovider

Usage details are in commands packages.

## License

Apache 2.0, see [LICENSE](./LICENSE).
