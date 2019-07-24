# mvisonneau/s5 - Safely Store Super Sensitive Stuff

[![GoDoc](https://godoc.org/github.com/mvisonneau/s5?status.svg)](https://godoc.org/github.com/mvisonneau/s5)
[![Go Report Card](https://goreportcard.com/badge/github.com/mvisonneau/s5)](https://goreportcard.com/report/github.com/mvisonneau/s5)
[![Docker Pulls](https://img.shields.io/docker/pulls/mvisonneau/s5.svg)](https://hub.docker.com/r/mvisonneau/s5/)
[![Build Status](https://cloud.drone.io/api/badges/mvisonneau/s5/status.svg)](https://cloud.drone.io/mvisonneau/s5)
[![Coverage Status](https://coveralls.io/repos/github/mvisonneau/s5/badge.svg?branch=master)](https://coveralls.io/github/mvisonneau/s5?branch=master)

`s5` is a small binary (~5MB) that works on most platforms (`AIX`, `Linux`, `Mac OS X`, `Windows`, `Freebsd` and `Solaris`!) and allows you to easily cipher/decipher content within your files.

## Encryption backends supported

- [AES](https://en.wikipedia.org/wiki/Advanced_Encryption_Standard) - [GCM](https://en.wikipedia.org/wiki/Galois/Counter_Mode) (using hexadecimal keys >= 128b) ([example usage](examples/aes-gcm.md))
- [AWS](https://aws.amazon.com) - [KMS](https://aws.amazon.com/kms/) ([example usage](examples/aws-kms.md))
- [GCP](https://cloud.google.com) - [KMS](https://cloud.google.com/kms/) ([example usage](examples/gcp-kms.md))
- [Hashicorp Vault](https://www.vaultproject.io) - [Transit secret engine](https://www.vaultproject.io/docs/secrets/transit/index.html) ([example usage](examples/vault.md))
- [PGP](https://www.openpgp.org/) ([example usage](examples/pgp.md))

## TL:DR

Example using AES-GCM as the encryption backend

```bash
# Generate and use a 128B random hexadecimal key
~$ export S5_AES_KEY=$(openssl rand -hex 16)

# Encrypt text
~$ s5 cipher aes very_sensitive_value
{{ s5:sIPFWfAcBvOnOtVcs65QGh+S3af4Wo= }}

# Store it anywhere in your files
~$ cat example.yml
---
var1: {{ s5:EtWnJ8ZyuwzRn8I3jw== }}
var2: {{ s5:8tceTb9yc0CBgEqrpw== }}
{{ s5:Glv1MRAuNOorI3oJA== }}: {{ s5:S4Lfavx2svWlSAD8sWHV }}

# Render!
~$ s5 render aes example.yml
---
var1: foo
var2: bar
secret_key: secret_value

# s5 can also read from stdin
~$ echo "foo" | s5 cipher aes | s5 decipher aes
foo
~$ echo "foo: {{ s5:8tceTb9yc0CBgEqrpw== }}" | s5 render aes
foo: bar
```

## Usage

```bash
~$ s5
NAME:
   s5 - cipher/decipher text within a file

USAGE:
   s5 [global options] command [command options] [arguments...]

COMMANDS:
     cipher    return an encrypted s5 pattern that can be included in any file
     decipher  return an unencrypted s5 value from a given pattern
     render    render a file that (may) contain s5 encrypted patterns
     help, h   Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --log-level level    log level (debug,info,warn,fatal,panic) (default: "info") [$S5_LOG_LEVEL]
   --log-format format  log format (json,text) (default: "text") [$S5_LOG_FORMAT]
   --help, -h           show help
   --version, -v        print the version
```

## Install

Have a look onto the [latest release page](https://github.com/mvisonneau/s5/releases/latest) and pick your flavor.

### Go

```bash
~$ go get -u github.com/mvisonneau/s5
```

### Homebrew

```bash
~$ brew install mvisonneau/tap/s5
```

### Docker

```bash
~$ docker run -it --rm mvisonneau/s5
```

### Scoop

```bash
~$ scoop bucket add https://github.com/mvisonneau/scoops
~$ scoop install s5
```

### Binaries, DEB and RPM packages

For the following ones, you need to know which version you want to install, to fetch the latest available :

```bash
~$ export S5_VERSION=$(curl -s "https://api.github.com/repos/mvisonneau/s5/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
```

```bash
# Binary (eg: aix/ppc64)
~$ wget https://github.com/mvisonneau/strongbox/releases/download/${S5_VERSION}/strongbox_${S5_VERSION}_aix_ppc64.tar.gz
~$ tar zxvf strongbox_${S5_VERSION}_aix_ppc64.tar.gz -C /usr/local/bin

# DEB package (eg: linux/386)
~$ wget https://github.com/mvisonneau/strongbox/releases/download/${S5_VERSION}/strongbox_${S5_VERSION}_linux_386.deb
~$ dpkg -i strongbox_${S5_VERSION}_linux_386.deb

# RPM package (eg: linux/arm64)
~$ wget https://github.com/mvisonneau/strongbox/releases/download/${S5_VERSION}/strongbox_${S5_VERSION}_linux_arm64.rpm
~$ rpm -ivh strongbox_${S5_VERSION}_linux_arm64.rpm
```

## Examples

### Render in-place

```bash
~$ cat example.yml
foo: {{ s5:8tceTb9yc0CBgEqrpw== }}

~$ s5 render aes --in-place example.yml

~$ cat example.yml
foo: bar
```

### Render in a new file

```bash
~$ cat example.yml
foo: {{ s5:8tceTb9yc0CBgEqrpw== }}

~$ s5 render aes example.yml --output example-dec.yml

~$ cat example-dec.yml
foo: bar
```

## Troubleshoot

You can use the `--log-level debug` flag in order to troubleshoot

```bash
~$ cat example.yml
foo: {{ s5:8tceTb9yc0CBgEqrpw== }}

~$ s5 --log-level=debug render aes example.yml
DEBU[2019-07-24T18:16:44+02:00] Opening input file : example.yml
DEBU[2019-07-24T18:16:44+02:00] Starting deciphering
DEBU[2019-07-24T18:16:44+02:00] found: 8tceTb9yc0CBgEqrpw
DEBU[2019-07-24T18:16:44+02:00] Deciphering '8tceTb9yc0CBgEqrpw' using AES
DEBU[2019-07-24T18:16:44+02:00] Outputing to stdout
foo: bar
DEBU[2019-07-24T18:16:44+02:00] Executed in 1.803305ms, exiting..
```

## Develop / Test

If you use docker, you can easily get started using :

```bash
~$ make dev-env
# You should then be able to use go commands to work onto the project, eg:
~docker$ make fmt
~docker$ s5
```

This command will spin up a `Vault` container and build another one with everything required in terms of **golang** dependencies in order to get started.

## Build / Release

If you want to build and/or release your own version of `s5`, you need the following prerequisites :

- [git](https://git-scm.com/)
- [golang](https://golang.org/)
- [make](https://www.gnu.org/software/make/)
- [goreleaser](https://goreleaser.com/)

```bash
~$ git clone git@github.com:mvisonneau/strongbox.git && cd strongbox

# Build the binaries locally
~$ make build

# Build the binaries and release them (you will need a GITHUB_TOKEN and to reconfigure .goreleaser.yml)
~$ make release
```

## BONUS

If you are using [atom.io](https://atom.io) as you IDE. You can have a look onto a [module I have written that integrates s5 with it](https://github.com/mvisonneau/atom-s5).

## Contribute

Contributions are more than welcome! Feel free to submit a [PR](https://github.com/mvisonneau/s5/pulls).
