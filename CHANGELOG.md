<a name="v0.0.1-pre"></a>
# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

## [v0.0.1-pre] - 2025-12-20

### Added
- Initial implementation of `terraform-provider-idgen`
- [Proquint](https://arxiv.org/html/0901.4016) ID generation with configurable entropy and grouping
- [NanoID](https://github.com/ai/nanoid) generation with custom alphabets (alphanumeric, numeric, readable)
- Random word selection from a very basic built-in or custom wordlists
- Templated ID composition combining multiple ID types and basic string functions
  - Template functions: `upper`, `lower`, `replace`, `prepend`, `append`, `substr`
- Deterministic seeding for reproducible IDs
- Support for direct IPv4-to-Proquint encoding per specification


### Documentation
- README and docs with usage examples
- Example configurations for various use cases
- Contribution guidelines
- MIT License

### Notes
- ⚠️ Pre-release version - API may change
- Not suitable for cryptographic purposes

[Unreleased]: https://github.com/iilei/terraform-provider-idgen/compare/v0.0.1-pre.1...HEAD
[0.0.1-pre.1]:  https://github.com/iilei/terraform-provider-idgen/releases/tag/v0.0.1-pre.1
