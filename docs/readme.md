# Surgeon

Modify your fork of an upstream repository with surgical precision.

See [IncusScripts](https://github.com/bketelsen/IncusScripts) which uses `surgeon`
to modify [community-scripts for Proxmox](https://github.com/community-scripts/ProxmoxVE)
to run on Incus.

Use `surgeon` to keep your fork up-to-date with an upstream, but make modifications during the process.

For example:

- use the `inject` module to add a license or credit to the upstream
- use the `sed` module to change a function invocation to a different one
- use the `replacefile` module to completely replace a single file

Codemods implement a small-ish interface so it's easy to add new ones. PR's accepted if you find something you need.

``` go
type CodeMod interface {
	Apply(source string, target string, match string, args ...string) error
	Validate(source string, target string, match string, args ...string) error
	Description() string
	Usage() string
}
```

---

## 🚀 Project Overview

**surgeon** Modify your fork of an upstream repository with surgical precision.

See [IncusScripts](https://github.com/bketelsen/IncusScripts) which uses `surgeon`
to modify [community-scripts for Proxmox](https://github.com/community-scripts/ProxmoxVE)
to run on Incus.

---

## 🚀 Known Issues

**git push/commit/stage** the git commit/push/stage settings don't work correctly. Don't use them yet.

---

## 🚀 Installation / Usage

See [installation documentation](https://bketelsen.github.io/surgeon/installation/)

---

## Config

``` yaml
upstream: https://some.repository.com/upstream/repo
stage: true
commit: true
push: false
modsdir: mymods
codemods:
- description: Modify URLS
  mod: sed
  match: cmd/*.go
  args:
  - github.com/upstream/repo
  - github.com/myfork/repo
ignorelist:
- prefix: ct/

```

See a [real world example](https://github.com/bketelsen/IncusScripts/blob/main/.surgeon.yaml)

IncusScripts uses surgeon as a GitHub Action. See the [action](https://github.com/bketelsen/IncusScripts/blob/main/.github/workflows/surgeon.yml)

## ❤️ Community and Contributions

We appreciate any contributions to the project—whether it's bug reports, feature requests, documentation improvements, or spreading the word. Your involvement helps keep the project alive and sustainable.

---

## 🤝 Report a Bug or Feature Request

If you encounter any issues or have suggestions for improvement, file a new issue on our [GitHub issues page](https://github.com/bketelsen/surgeon/issues). You can also submit pull requests with solutions or enhancements!

---

## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=bketelsen/surgeon&type=Date)](https://star-history.com/#bketelsen/surgeon&Date)

---

## 📜 License

This project is licensed under the [MIT License](LICENSE).

