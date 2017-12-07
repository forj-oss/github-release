# Introduction

This Repository has been strongly inspired by [aktau/github-release](https://github.com/aktau/github-release)

The core code has been completed rewrote for our simple need.
So, this is not exactly the same tool.

In short, it provides simple Github release management from CLI.
(create/update/delete)

The github code is based on https://github.com/google/go-github library
which is well maintained.

Why did I created a new one instead of using aktau/github-release?

Several reasons:
- The project seems not too much active
- creating a release on private github is broken
- The github api management is not based on go-github. So, it may not be
    up to date and/or tested.

# How to use it?

Download the x64 binary from github and use it:

```bash
$ mkdir ~/bin
$ wget -O ~/bin/github-release https://github.com/forj-oss/github-release/releases/download/1.0.0/github-release
$ chmod +x ~/bin/github-release
$ ~/bin/github-release --help
```

NOTE: github-release do not create git tags!

You can set Environment variables instead of setting flags to the cli:

- GITHUB_API: Must be as `https://<server>/api/v3/` If not set,
    default is to use publish github url.
- GITHUB_TOKEN: Github token with release rights given
- GITHUB_USER: Github Organization or User name
- GITHUB_REPO: Organization or User repository to attach releases.

# How to build it?

__Requirements__:
- docker 1.8 or higher
- bash
- git

The project use docker to build and reproduce binary code as it was
designed.
It is not mandatory but recommended but should work well.
If you do not want to use it, ignore following *Build env* sections

To simplify access to this docker environment and setup, a build environment
has been written.

**Please note!** The build env install a `build-env` alias in your .bashrc

## Using build env

### Build env Initial setup

The initial setup will try to determine
- if you run docker directly or through sudo. With sudo, add --sudo
- where is your go GOPATH setting. [GO path] is required if you do not
    have any GOPATH variable setup in bash.

```bash
$ # Assuming go path to be ~/go
$ mkdir -p ~/go/src
$ cd ~/go/src
$ git clone https://github.com/forj-oss/github-release
$ cd github-release
$ source build-env.sh [--sudo] [GO path]
```

### Load and build

The build env must be loaded before running go from docker:

```bash
$ build-env # This is the alias added to your .bashrc
$ glide i
$ go build
```

## Without build env

**Requirements:**
- go 1.7.4
- glide

### install

Warning! Your distribution may not have the right version of go or even
to not having glide at all.

Following instructions are given as is and were not tested.
If there is a need to fix this doc, contribute! Thank you in advance.

On Fedora:
```bash
$ sudo dnf install go glide -y
```

On CentOS/RHEL:
```bash
$ sudo yum install go glide -y
```

On Ubuntu:
```bash
$ sudo apt-get install go glide -y
```

### build

```bash
$ mkdir -p ~/go/src
$ cd ~/go/src
$ git clone https://github.com/forj-oss/github-release
$ cd github-release
$ glide i
$ go build
```

# How to contribute?

To be simple: Create a Pull Request!

To debug, you can set a shell variable `GOTRACE` to `true`/`debug`/`<debug level number>`

Ex:
```bash
GOTRACE=1 ./github-release release ...
```

if you set GOTRACE to `true` or `debug`, the debug level will be set to 0

If you want to see more detailled debug output, set the level number
in GOTRACE.

At level 0:
- Show usual program activity trace.

At Level 1:
- github-release shows REST API calls against github.

Forj team
