# TelShell

Tiny Telnet shell server in Go

![alt text](./docs/preview.png)

## Download

Grab latest release from [here](https://github.com/x1unix/telshell/releases/latest)

## Usage

```bash
./telshell -listen=:5000
```

**Optional arguments:**

* `shell` - Login shell
* `-s` - Set login shell parameter (e.g. `-s=-r` for restricted shell)
* `-buffer` - TCP connection read buffer size
* `-auth` - Require authorization (experimental)
