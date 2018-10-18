# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]

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

[Unreleased]: https://github.com/mvisonneau/s5/compare/0.1.1...HEAD
[0.1.1]: https://github.com/mvisonneau/s5/tree/0.1.1
[0.1.0]: https://github.com/mvisonneau/s5/tree/0.1.0
