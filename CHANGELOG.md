# Change Log
All notable changes to this project will be documented in this file.
This project adheres to [Semantic Versioning](http://semver.org/).

## [Unreleased]

### Added
- Add `(*pq.Reader).Begin/Done` to reuse a read transaction for multiple reads. PR #4
- Add `Flags` to txfile.Options. PR #5
- Add support to increase a file's maxSize on open. PR #5
- Add support to pre-allocate the meta area. PR #7
- Begin returns an error if transaction is not compatible to file open mode. PR #17
- Introduce Error type to txfile and pq package. PR #17, #18

### Changed
- Refine platform dependent file syncing. PR #10

### Deprecated

### Removed

### Fixed


[Unreleased]: https://github.com/elastic/go-structform/compare/v0.0.1...HEAD
