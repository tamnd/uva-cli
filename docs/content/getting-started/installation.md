---
title: "Installation"
description: "Install uva from a release, with go install, or from source."
weight: 20
---

## Prebuilt binaries

Every [release](https://github.com/tamnd/uva-cli/releases) carries archives for Linux, macOS,
and Windows on amd64 and arm64, plus deb, rpm, and apk packages for Linux.
Download, unpack, put `uva` on your `PATH`, done. The `checksums.txt`
on each release is signed with keyless [cosign](https://docs.sigstore.dev/) if
you want to verify before running.

## With Go

```bash
go install github.com/tamnd/uva-cli/cmd/uva@latest
```

That puts `uva` in `$(go env GOPATH)/bin`, which is `~/go/bin` unless
you moved it. Make sure that directory is on your `PATH`.

## From source

```bash
git clone https://github.com/tamnd/uva-cli
cd uva-cli
make build        # produces ./bin/uva
./bin/uva version
```

## Container image

```bash
docker run --rm ghcr.io/tamnd/uva:latest --help
```

## Checking the install

```bash
uva version
```

prints the version and exits.
