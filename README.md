# Surgeon

Modify your fork of an upstream repository with surgical precision.

See [IncusScripts](https://github.com/bketelsen/IncusScripts) which uses `surgeon`
to modify [community-scripts for Proxmox](https://github.com/community-scripts/ProxmoxVE)
to run on Incus.

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

Clone and build the repository with Go

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

