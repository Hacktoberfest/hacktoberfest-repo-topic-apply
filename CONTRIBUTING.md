## Releasing

This repository contains a GitHub Action configured to automatically build and
publish assets for release when a tag is pushed that matches the pattern v*
(ie. v0.1.0). This is done using [Goreleaser](https://goreleaser.com/).

To build a local test a release, run:

    goreleaser release --snapshot --rm-dist