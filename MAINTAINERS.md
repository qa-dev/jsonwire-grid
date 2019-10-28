# Maintainers Guide
A guide for maintainers.

## How to release
1. Get [GoReleaser](https://goreleaser.com).
1. Create and push a tag.
1. Run `export GITHUB_TOKEN=<your_token>`.
1. Run `goreleaser --rm-dist`.
