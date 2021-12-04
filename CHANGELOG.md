# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.3] - 2021-12-04
### Added
- 
### Changed
- fix permission for $HOME/.ssh
### Removed
-
## [0.1.2] - 2021-12-04
### Added
- 
### Changed
- fix permission for /usr/local/etc/gossh
### Removed
-
## [0.1.1] - 2021-12-01
### Added
- 
### Changed
- Change the behaviour of ini parser to only treat '#' and ';' as comment indicators when a space is present before them. ('xxx;yyy' is parsed as it is. 'xxx ;yyy' is parsed to 'xxx'). This is useful when a password contains one of those characters.
### Removed
-

## [0.1.0] - 2021-11-21
### Added
- Support for local port forwarding
- CHANGELOG.md
### Changed
- Start following [SemVer](https://semver.org) properly.
- Updated README.md 
### Removed
-