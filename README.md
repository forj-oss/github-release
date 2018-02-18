# Warning!!! Development suspended.

> Small DevOps Story: **Speedy gonzales developer**!
>
>
>
> This Repository has been strongly inspired by [aktau/github-release](https://github.com/aktau/github-release)
>
> Several reasons were listed which make me decide to create my own
> `light github-release`
>
> But one extremely important were wrong and removed. `aktau/github-release`
> were considered to fail against private github.
>
> This one was a mistake from the url I gave. :-[
>
> In short, if you set `http` instead of `https`, github return a `301
> moved permanently`. Internally, before creating or updating, or even
> connecting, `github-release` were success until call to the create/update
> API.

> So, it worked at 80%, because the github library used did the move to
> the `https` protocol automatically. But not on last API call. And it fails.
> For example, the connection succeeded, but the create failed. Very
> frustrating.
>
> Of course, it was a good reason to try to fix it, directly, as library
> is already there and easy to use. So I created my own code.
> But unfortunately, I finally lost my time here... atkau/github-release
> is still valid and work as well.
>
> Note that I'm still maintaining it, today. But I may decide to stop it
> definitely.
> the list of reason, against the cost to maintain it could close this
> this repo soon.
>
> I would suggest you to avoid using it as of now.
> If you consider that this code make sense for you, fork it.
> Or even, discuss that with me so that I can change my mind and maintain
> it long term.
>
> But for now, I 'm providing a last update, and I think I will re-use
> `aktau/github-release`.
>
> Thank you


# Introduction

This Repository has been strongly inspired by [aktau/github-release](https://github.com/aktau/github-release)

The core code has been completely rewrote for our simple need.
So, this is not exactly the same tool.

In short, it provides simple Github release management from CLI.
(create/update/delete/check)

The github code is based on https://github.com/google/go-github library
which is well maintained.

Why did I created a new one instead of using aktau/github-release?

Several reasons:
- The project seems not too much active
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

## `release` command - To create or update a release

No need to explain too much here.
Provide the appropriate list of parameters to create or update the
github release.

It does NOT upload artifacts.

## `delete` command - To delete a release

No need to explain too much here.
It will remove it if found.

## `has-release` command - To check if a release exist or not

This new last feature added, help to check if a particular release
exist or not.
It returns 1, if fails, and 0 if succeed.

You can use it in a if command directly.

```bash
if github_release has-release MyRelease && [[ ... ]] ...
then
  # Do something ...
else
  # Do something ...
fi
```

You can get more information with GOTRACE=info

```bash
$ GOTRACE=info github_release has-release latest
INFO: Github API URL used : https://api.github.com/

INFO: Connection successful. Token given by user '<TheTokenOwner>'

INFO: Tag '0.9a' not found! Valid tags are 'latest', '0.9', '0.9-a', '0.1'

```

If github fails, it returns 255 with error messages displayed.

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
