# mvisonneau/s5 - Safely Store Super Sensitive Stuff

[![GoDoc](https://godoc.org/github.com/mvisonneau/s5?status.svg)](https://godoc.org/github.com/mvisonneau/s5)
[![Go Report Card](https://goreportcard.com/badge/github.com/mvisonneau/s5)](https://goreportcard.com/report/github.com/mvisonneau/s5)
[![Docker Pulls](https://img.shields.io/docker/pulls/mvisonneau/s5.svg)](https://hub.docker.com/r/mvisonneau/s5/)
[![Build Status](https://cloud.drone.io/api/badges/mvisonneau/s5/status.svg)](https://cloud.drone.io/mvisonneau/s5)
[![Coverage Status](https://coveralls.io/repos/github/mvisonneau/s5/badge.svg?branch=master)](https://coveralls.io/github/mvisonneau/s5?branch=master)

`s5` is a very small binary that allows you to easily cipher/decipher content within your files.

## Encryption backends supported

- [Hashicorp Vault](https://www.vaultproject.io) - [Transit secret engine](https://www.vaultproject.io/docs/secrets/transit/index.html)

## TL:DR

```bash
# Configure Vault
~$ export VAULT_ADDR=https://vault.rocks
~$ export VAULT_TOKEN=f4262de2-4e07-5b85-98ea-7702e2c7cdb9

# Encrypt text
~$ s5 cipher vault very_sensitive_value
{{ s5:sIPFWfAcBvOnOtVcs65QGh+S3af4Wo= }}

# Store it anywhere in your files
~$ cat example.yml
---
var1: {{ s5:EtWnJ8ZyuwzRn8I3jw== }}
var2: {{ s5:8tceTb9yc0CBgEqrpw== }}
{{ s5:Glv1MRAuNOorI3oJA== }}: {{ s5:S4Lfavx2svWlSAD8sWHV }}

# Render!
~$ s5 render vault example.yml
---
var1: foo
var2: bar
secret_key: secret_value

# s5 can also read from stdin
~$ echo "foo" | s5 cipher vault | s5 decipher vault
foo
~$ echo "foo: {{ s5:8tceTb9yc0CBgEqrpw== }}" | s5 render vault
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
# Binary (eg: freebsd/amd64)
~$ wget https://github.com/mvisonneau/strongbox/releases/download/${S5_VERSION}/strongbox_${S5_VERSION}_freebsd_arm64.tar.gz
~$ tar zxvf strongbox_${S5_VERSION}_freebsd_amd64.tar.gz -C /usr/local/bin

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

~$ s5 render --in-place example.yml

~$ cat example.yml
foo: bar
```

### Render in a new file

```bash
~$ cat example.yml
foo: {{ s5:8tceTb9yc0CBgEqrpw== }}

~$ s5 render example.yml --output example-dec.yml

~$ cat example-dec.yml
foo: bar
```

## Troubleshoot

You can use the `--log-level debug` flag in order to troubleshoot

```bash
~$ cat example.yml
foo: {{ s5:8tceTb9yc0CBgEqrpw== }}

~$ s5 --log-level debug render example.yml
s5 --log-level debug render secrets.yml
DEBU[2018-07-09T15:06:49Z] Configuring Vault
DEBU[2018-07-09T15:06:49Z] Executing function 'render'
DEBU[2018-07-09T15:06:49Z] Opening input file : example.yml
DEBU[2018-07-09T15:06:49Z] Starting deciphering
DEBU[2018-07-09T15:06:49Z] found: s5:8tceTb9yc0CBgEqrpw==
DEBU[2018-07-09T15:06:49Z] Outputing to stdout
foo: bar
DEBU[2018-07-09T15:06:49Z] Executed in 13.1337ms, exiting..
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

## BONUS

If you are using [atom.io](https://atom.io) as you IDE. You can have a look onto a [module I have written that integrates s5 with it](https://github.com/mvisonneau/atom-s5).

## Contribute

Contributions are more than welcome! Feel free to submit a [PR](https://github.com/mvisonneau/s5/pulls).
