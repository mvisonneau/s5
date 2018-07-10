# mvisonneau/s5 - Safely Store Super Sensitive Stuff

[![GoDoc](https://godoc.org/github.com/mvisonneau/s5?status.svg)](https://godoc.org/github.com/mvisonneau/s5/app)
[![Go Report Card](https://goreportcard.com/badge/github.com/mvisonneau/s5)](https://goreportcard.com/report/github.com/mvisonneau/s5)
[![Docker Pulls](https://img.shields.io/docker/pulls/mvisonneau/s5.svg)](https://hub.docker.com/r/mvisonneau/s5/)
[![Build Status](https://travis-ci.org/mvisonneau/s5.svg?branch=master)](https://travis-ci.org/mvisonneau/s5)
[![Coverage Status](https://coveralls.io/repos/github/mvisonneau/s5/badge.svg?branch=master)](https://coveralls.io/github/mvisonneau/s5?branch=master)

`s5` is a very small binary that allows you to easily cipher/decipher content within your files. For the moment it only supports [Vault transit secret engine](https://www.vaultproject.io/docs/secrets/transit/index.html) (Hashicorp) but it could be ported to additional ones as well.

## TL:DR

```bash
# Configure Vault
~$ export VAULT_ADDR=https://vault.rocks
~$ export VAULT_TOKEN=f4262de2-4e07-5b85-98ea-7702e2c7cdb9

# Encrypt text
~$ s5 cipher very_sensitive_value
{{ s5:sIPFWfAcBvOnOtVcs65QGh+S3af4Wo= }}

# Store it anywhere in your files
~$ cat example.yml
---
var1: {{ s5:EtWnJ8ZyuwzRn8I3jw== }}
var2: {{ s5:8tceTb9yc0CBgEqrpw== }}
{{ s5:Glv1MRAuNOorI3oJA== }}: {{ s5:S4Lfavx2svWlSAD8sWHV }}

# Render!
~$ s5 render example.yml
---
var1: foo
var2: bar
secret_key: secret_value
```

## Usage

```bash
~$ s5
NAME:
   s5 - cipher/decipher text within a file from a (Hashicorp) Vault transit key

USAGE:
   s5 [global options] command [command options] [arguments...]

VERSION:
   <devel>

COMMANDS:
     cipher    return an encrypted s5 pattern that can be included in any file
     decipher  return an unencrypted s5 value from a given pattern
     render    render a file that (may) contain s5 encrypted patterns
     help, h   Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --vault-addr address         vault address [$VAULT_ADDR]
   --vault-token token          vault token [$VAULT_TOKEN]
   --transit-key name, -k name  name of the transit key used by s5 to cipher/decipher data (default: "default") [$S5_TRANSIT_KEY]
   --log-level level            log level (debug,info,warn,fatal,panic) (default: "info") [$S5_LOG_LEVEL]
   --log-format format          log format (json,text) (default: "text") [$S5_LOG_FORMAT]
   --help, -h                   show help
   --version, -v                print the version
```

## Install

You can have a look at the [release page](https://github.com/mvisonneau/s5/releases) of the project, we currently build it for **Linux**, **Darwin** and **Windows** platforms.

```
~$ wget https://github.com/mvisonneau/s5/releases/download/0.1.1/s5_linux_amd64 -O /usr/local/bin/s5; chmod +x /usr/local/bin/s5
```

You can also use the [docker version](https://hub.docker.com/r/mvisonneau/s5):

```
~$ docker run -it --rm mvisonneau/s5
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

## Contribute

Contributions are more than welcome! Feel free to submit a [PR](https://github.com/mvisonneau/s5/pulls).
