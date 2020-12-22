# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Added
### Changed
### Fixed
### Docs

## [0.1.3] - 2020-12-22
### Added
- Ability to generate an archive with relative paths ([#13](https://github.com/d-Rickyy-b/backmeup/pull/13))
### Fixed
- Add check for duplicate source paths([#14](https://github.com/d-Rickyy-b/backmeup/pull/14))

## [0.1.2] - 2020-12-10
### Added
- Implement -d/--debug switch to enable debug logging for printing exclude matches ([21a201a](https://github.com/d-Rickyy-b/backmeup/commit/21a201a7fa7013aee2159cd18d4672ada65442b0))
### Fixed
- Don't return error in filepath.Walk() for excluded files ([94fba5c](https://github.com/d-Rickyy-b/backmeup/commit/94fba5cab11d3dc07b2ef613e81455b2c1c215bc))

## [0.1.1] - 2020-12-10
### Fixed
- Use 24h format for backup file names ([#8](https://github.com/d-Rickyy-b/backmeup/pull/8))
- Use proper file globbing for matching excludes ([0b25309](https://github.com/d-Rickyy-b/backmeup/commit/0b2530989232f7082f14e79f1036cb8f7ee6053c))
### Docs
- Document limitations in README ([9f33d9a](https://github.com/d-Rickyy-b/backmeup/commit/9f33d9adaa81c90ddd5b9b166ac61cee46317175), [8bac46a](https://github.com/d-Rickyy-b/backmeup/commit/8bac46ac6272f29e2b8b3555fcbae36619732d5c))
- Several fixes in README ([01951dc](https://github.com/d-Rickyy-b/backmeup/commit/01951dc4273ab968d616d839d7b66fdee6d69371), [0cbf045](https://github.com/d-Rickyy-b/backmeup/commit/0cbf045898889e808d462dfd0452b6a9d2715579))


## [0.1.0] - 2020-11-17
Initial release! First usable version of backmeup is published as v0.1.0 

[unreleased]: https://github.com/d-Rickyy-b/backmeup/compare/v0.1.3...HEAD
[0.1.3]: https://github.com/d-Rickyy-b/backmeup/tree/v0.1.3
[0.1.2]: https://github.com/d-Rickyy-b/backmeup/tree/v0.1.2
[0.1.1]: https://github.com/d-Rickyy-b/backmeup/tree/v0.1.1
[0.1.0]: https://github.com/d-Rickyy-b/backmeup/tree/v0.1.0
