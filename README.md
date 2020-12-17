# 🔐 s5 - Safely Store Super Sensitive Stuff

[![PkgGoDev](https://pkg.go.dev/badge/github.com/mvisonneau/s5)](https://pkg.go.dev/mod/github.com/mvisonneau/s5)
[![Go Report Card](https://goreportcard.com/badge/github.com/mvisonneau/s5)](https://goreportcard.com/report/github.com/mvisonneau/s5)
[![Docker Pulls](https://img.shields.io/docker/pulls/mvisonneau/s5.svg)](https://hub.docker.com/r/mvisonneau/s5/)
[![Build Status](https://github.com/mvisonneau/s5/workflows/test/badge.svg?branch=main)](https://github.com/mvisonneau/s5/actions)
[![Coverage Status](https://coveralls.io/repos/github/mvisonneau/s5/badge.svg?branch=master)](https://coveralls.io/github/mvisonneau/s5?branch=master)

`s5` allows you to easily cipher/decipher content within your files. It is a cli tool that works on most platforms (`Linux`, `Mac OS X`, `Windows`, `Freebsd` and `Solaris`!)

## Encryption backends supported

- [AES](https://en.wikipedia.org/wiki/Advanced_Encryption_Standard) - [GCM](https://en.wikipedia.org/wiki/Galois/Counter_Mode) (using hexadecimal keys >= 128b) ([example usage](examples/aes-gcm.md))
- [AWS](https://aws.amazon.com) - [KMS](https://aws.amazon.com/kms/) ([example usage](examples/aws-kms.md))
- [GCP](https://cloud.google.com) - [KMS](https://cloud.google.com/kms/) ([example usage](examples/gcp-kms.md))
- [Hashicorp Vault](https://www.vaultproject.io) - [Transit secret engine](https://www.vaultproject.io/docs/secrets/transit/index.html) ([example usage](examples/vault.md))
- [PGP](https://www.openpgp.org/) ([example usage](examples/pgp.md))

## TL:DR

Example using AES-GCM as the encryption backend

[![asciicast](https://asciinema.org/a/gmKNYVb49Vzp3SFpeqiavvVe5.svg)](https://asciinema.org/a/gmKNYVb49Vzp3SFpeqiavvVe5)

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

Checksums are signed with the [following GPG key](https://keybase.io/mvisonneau/pgp_keys.asc): `C09C A9F7 1C5C 988E 65E3  E5FC ADEA 38ED C46F 25BE`

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
~$ docker run -it --rm docker.io/mvisonneau/s5
or
~$ docker run -it --rm ghcr.io/mvisonneau/s5
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
~$ wget https://github.com/mvisonneau/strongbox/releases/download/${S5_VERSION}/strongbox_${S5_VERSION}_freebsd_amd64.tar.gz
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

## IDEs

We have modules/extensions for both [atom.io](https://atom.io) and [vscode](https://code.visualstudio.com/)

- [mvisonneau/vscode-s5](https://github.com/mvisonneau/vscode-s5)
- [mvisonneau/atom-s5](https://github.com/mvisonneau/atom-s5)

## Contribute

Contributions are more than welcome! Feel free to submit a [PR](https://github.com/mvisonneau/s5/pulls).
