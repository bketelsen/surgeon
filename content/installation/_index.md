---
title: 'Installation'
weight: 1
---

You can install surgeon by downloading a release from GitHub or by using our installer script.

Choose your adventure below.

## Direct Download

You can download the binary from the [surgeon releases page](https://github.com/bketelsen/surgeon/releases) on GitHub and add to your `$PATH`.

The surgeon_VERSION_checksums.txt file contains the SHA-256 checksum for each file.

## Installer Script

We also have an [install script](https://github.com/bketelsen/surgeon/blob/main/install.sh) which is very useful in scenarios like CI.

By default, it installs on the `./bin` directory relative to the working directory:

```bash
sh -c "$(curl --location https://bketelsen.github.io/surgeon/install.sh)" -- -d
```

It is possible to override the installation directory with the -b parameter. On Linux, common choices are `~/.local/bin` and `~/bin` to install for the current user or `/usr/local/bin` to install for all users:

```bash
sh -c "$(curl --location https://bketelsen.github.io/surgeon/install.sh)" -- -d -b ~/.local/bin
```

> [!WARNING]
> On macOS and Windows, ~/.local/bin and ~/bin are not added to $PATH by default.

By default, it installs the latest version available. You can also specify a tag ([available in releases](https://github.com/bketelsen/surgeon/releases)) to install a specific version:

```bash
sh -c "$(curl --location https://bketelsen.github.io/surgeon/install.sh)" -- -d v0.2.2
```