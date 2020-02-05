# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [0ver](https://0ver.org).

## [Unreleased]

### Changed

- Fixed silly bugs on `AWS` and `Vault` ciphers introduced in latest release
- Use more coherent flags for `PGP` public and private key paths

## [0.1.7] - 2020-02-05

### Changed

- New input/output syntax : `{{s5:*}}` - retro compatible with the old one (`{{ s5:* }}`)
- Refactored cipher engine implementation to be able to use it as a library
- Bumped dependencies

### Removed

- `aix/ppc64` binaries due to a limitation in the github.com/hashicorp/go-sockaddr library

## [0.1.6] - 2019-09-20

### Added

- Release `aix/ppc64` binaries
- Release `linux/ppc64` binaries
- Release `solaris/amd64` binaries

### Changed

- Upgraded dependencies, go to 1.13 / goreleaser to 0.118.0

## [0.1.5] - 2019-07-24

### Added

- By default, we will now trim input for `cipher` function from whitespaces and return carriages

### Changed

- Updated dependencies

## [0.1.4] - 2019-07-19

## Added

- Added support for **PGP public/private keypair encryption**
- Added support for **AES (GCM) key encryption**
- Added support for **AWS KMS key encryption**
- Added support for **GCP KMS key encryption**
- Release `homebrew` packages
- Release `scoop` packages
- Release `freebsd` packages
- Release `DEB` packages
- Release `RPM` packages
- Parameterized the cipher engine in order to not support only Vault

## [0.1.3] - 2019-07-03

### Added

- Use Vault `0.11.3` in the dev container
- Updated the `golang.org/x/lint/golint` path in the setup make rule
- Upgraded to Go 1.12
- Use Gomodules
- Enhanced Makefile
- Moved CI to drone.io
- Refactored codebase for better modularity and portability

## [0.1.2] - 2018-08-22

### Added

- We can now pass text through `stdin` as well as an argument for all functions (`cipher`/`decipher`/`render`)
- Updated dependencies
- Used `busybox` image instead of **scratch** for the container image
- Use Vault `0.10.4` in the dev container

## [0.1.1] - 2018-07-10

### Added

- Replaced the `vault:v1` prefix on stored secrets with `s5`

### Changed

- Get the linux binaries working on `alpine` by disabling CGO
- Fixed the Dockerfile
- Fixed the `make dev-env` function

## [0.1.0] - 2018-07-09

### Added

- Working state of the app
- cipher function
- decipher function
- render function
- got some tests in place
- Makefile
- LICENSE
- README

[Unreleased]: https://github.com/mvisonneau/s5/compare/0.1.7...HEAD
[0.1.6]: https://github.com/mvisonneau/s5/tree/0.1.7
[0.1.6]: https://github.com/mvisonneau/s5/tree/0.1.6
[0.1.5]: https://github.com/mvisonneau/s5/tree/0.1.5
[0.1.4]: https://github.com/mvisonneau/s5/tree/0.1.4
[0.1.3]: https://github.com/mvisonneau/s5/tree/0.1.3
[0.1.2]: https://github.com/mvisonneau/s5/tree/0.1.2
[0.1.1]: https://github.com/mvisonneau/s5/tree/0.1.1
[0.1.0]: https://github.com/mvisonneau/s5/tree/0.1.0
