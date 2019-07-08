# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [0ver](https://0ver.org).

## [Unreleased]

## FEATURES
- Added support for PGP public/private keypair

## ENHANCEMENTS
- Parameterized the cipher engine in order to not support only Vault

## [0.1.3] - 2019-07-03

### ENHANCEMENTS
- Use Vault `0.11.3` in the dev container
- Updated the `golang.org/x/lint/golint` path in the setup make rule
- Upgraded to Go 1.12
- Use Gomodules
- Enhanced Makefile
- Moved CI to drone.io
- Refactored codebase for better modularity and portability

## [0.1.2] - 2018-08-22

### FEATURES
- We can now pass text through `stdin` as well as an argument for all functions (`cipher`/`decipher`/`render`)
### ENHANCEMENTS
- Updated dependencies
- Used `busybox` image instead of **scratch** for the container image
- Use Vault `0.10.4` in the dev container

## [0.1.1] - 2018-07-10
### FEATURES
- Replaced the `vault:v1` prefix on stored secrets with `s5`

### BUGFIXES
- Get the linux binaries working on `alpine` by disabling CGO
- Fixed the Dockerfile
- Fixed the `make dev-env` function

## [0.1.0] - 2018-07-09
### FEATURES
- Working state of the app
- cipher function
- decipher function
- render function
- got some tests in place
- Makefile
- LICENSE
- README

[Unreleased]: https://github.com/mvisonneau/s5/compare/0.1.3...HEAD
[0.1.3]: https://github.com/mvisonneau/s5/tree/0.1.3
[0.1.2]: https://github.com/mvisonneau/s5/tree/0.1.2
[0.1.1]: https://github.com/mvisonneau/s5/tree/0.1.1
[0.1.0]: https://github.com/mvisonneau/s5/tree/0.1.0
