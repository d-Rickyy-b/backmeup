# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Added
### Changed
### Fixed
### Docs

## [1.0.0] - 2021-07-05
New major version because of the switch to Go 1.16! The features are not backwards compatible.

### Added
- Support for single files as sources ([620fd62](https://github.com/d-Rickyy-b/backmeup/commit/620fd620d13a3687015d61f1bf7f3d89cbedf3a5))
### Changed
- Improved speed by using WalkDir instead of Walk ([9604540](https://github.com/d-Rickyy-b/backmeup/commit/96045409b099ca77f24cf43c442762aeb87ea62d))
- Updated Go version to 1.16 ([52b9b34](https://github.com/d-Rickyy-b/backmeup/commit/52b9b34767166f910466799e69d8499fbfa8db16))

## [0.1.4] - 2021-02-11
### Added
- Symlink support for tar files ([5d75752](https://github.com/d-Rickyy-b/backmeup/commit/5d757525bbde26429e90a30ea5fba8d721db6f72))
- Following symlinks (aka replacing a symlink to a file with the actual file) ([a98fe65](https://github.com/d-Rickyy-b/backmeup/commit/a98fe65d8188cd8f5abac2d766cffa594c032757))
- Ability to only run backups of certain units via `-u`/`--unit` CLI parameter ([92db794](https://github.com/d-Rickyy-b/backmeup/commit/92db794365448c67379f20ff3e2d6bfb998f1f57))
- Check if archive already exists ([fd88626](https://github.com/d-Rickyy-b/backmeup/commit/fd886263038d6c97cb0f481e9ff0140187d5283e))
- Add `-t`/`--test-path` CLI parameter for checking exclusion for given paths ([1b13e44](https://github.com/d-Rickyy-b/backmeup/commit/1b13e44a38faa0e472ecaea4b8864cfffc2ab147))
- Add `-v`/`--version` CLI parameter to just print the tool's version ([2c51b05](https://github.com/d-Rickyy-b/backmeup/commit/2c51b058723e1eb3e46ba2e8ee0b2260ad39b362))
### Changed
- Move archive code to archiver package ([d8666cb](https://github.com/d-Rickyy-b/backmeup/commit/d8666cb5d3acc25a77f3d84f92c52301687dd6ae))
- Move config code to config package ([0a03807](https://github.com/d-Rickyy-b/backmeup/commit/0a038077a21c88781abf77b85a6a9da7b60df9f6))
### Fixed
- Add compression for zip files ([52733bc](https://github.com/d-Rickyy-b/backmeup/commit/52733bc0dc4e1378e02467c3712ffe05b6cb3fd2))
- Replace Fatalln with Println ([721f6b2](https://github.com/d-Rickyy-b/backmeup/commit/721f6b27d1501b403d94f1273639a7a1a92b8b76))
- Correctly assign 'verbose' and 'debug' variables ([cd11006](https://github.com/d-Rickyy-b/backmeup/commit/cd110062d8f619ead0f63b4a663c3a46aedbd228))
- Only store regular files in tar archives ([d9b26fc](https://github.com/d-Rickyy-b/backmeup/commit/d9b26fc5d0b465bebec05454fecbe4b5b14538b9))

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

[unreleased]: https://github.com/d-Rickyy-b/backmeup/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/d-Rickyy-b/backmeup/compare/v0.1.4...v1.0.0
[0.1.4]: https://github.com/d-Rickyy-b/backmeup/compare/v0.1.3...v0.1.4
[0.1.3]: https://github.com/d-Rickyy-b/backmeup/compare/v0.1.2...v0.1.3
[0.1.2]: https://github.com/d-Rickyy-b/backmeup/compare/v0.1.1...v0.1.2
[0.1.1]: https://github.com/d-Rickyy-b/backmeup/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/d-Rickyy-b/backmeup/tree/v0.1.0
