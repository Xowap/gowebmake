# gowebmake

A tool to trigger `make` targets from (GitHub) webhooks.

Endgame: automatically publishing static blogs on `git push`.

## Overview

You need to run `gowebmake` on a server and you need to give it a public URL. It will bind on a TCP port of your
choosing and wait for HTTP. It's up to you to put it in front of the wild web or to put it behind a proxy.

Then you need to configure GitHub to run this webhook on each push.

Once this is done, on each `git push` you make to your repository, `gowebmake` will run the specified `make` target on
an up-to-date version.

All you have to do next is write a `Makefile` that does the deployment for you.

## Usage

`gowebmake` aims at being as simple as possible. It is a simple Go programm with a minimal set of dependencies. And
since it's a Go program, it comes built statically, so you can just copy/paste the binary.

### Prerequisites

A few tools are still needed:

- `make`, since the goal here is to run `make` targets
- `git`, in order to be able to fetch remote code

### Install

You need to install `gowebmake` on a server. You will have to make it live there as a daemon and give it a public URL
(either by running it behind a proxy or by making it face the web).

#### Install binary

If you're using a `x86_64` linux (most servers do), then you can find a pre-compiled binary in the download section of
this repository.

#### Compile from source

You'll need to know how to build Go software. However if you do, it's really easy.

### Configure

You need to create a configuration file on the server that will be doing deployments. The configuration comes as a TOML
file. It looks like this:

```
workdir = "/some/place"

[github."Xowap/TestRepo"]
secret = "tralala"

[github."Xowap/TestRepo".branches.master]
target = "deploy"

[github."Xowap/TestRepo".branches.master.env]
DEPLOY_DIR = "/tmp/yolo"
```

Let's go around the various configuration options you have here.

At root:

- **bind** *(optional, default = `[::1]:8777`)* The TCP address to bind to.
- **workdir** Directory into which all repos will be cloned.
- **github** Map of all GitHub projects *(see next)*
- **gitbin** *(optional, default = `git`)* Path to the `git` binary
- **makebin** *(optional, default = `make`)* Path to the `make` binary

Then, for each GitHub project:

- **protocol** *(optional, default = `auto`)* Protocol to use to clone repos. Possibilities are `ssh`, `http` and
  `auto`. Please note that the standard `git` utility is used, so if you need ssh keys you can just put them wherever
  you usually put them.
- **secret** The secret you provided to GitHub for the webhook.
- **branches** A map of branches you want to deploy. Usually you just want `master`.

And for each branch:

- **target** Name of the `make` target you want to run upon push
- **env** A map of envirnment variables you want to set when running `make`

### Running

It's pretty simple.

```
gowebmake -c /path/to/gowebmake.conf
```

### Configuring webhook

You can configure the webhook in you GitHub project settings. Let's say your server is `example.org` and you've bound
to `[::]:8777`, then the webhook URL will be `http://example.org/github`.

### Making a daemon

Not documented yet. It's pretty easy with systemd units though, have a look at them!


## Example Makefile

Here's how I deploy my Hugo blog:

```
build:
	rm -fr dist
	hugo -d dist

deploy: check-env build
	mkdir -p $(DEPLOY_DIR)
	echo "$(DEPLOY_DIR)"
	rsync -rtv --delete dist/ "$(DEPLOY_DIR)/"

check-env:
ifndef DEPLOY_DIR
	$(error DEPLOY_DIR is undefined)
endif
```

## Contributors and licencing

See CONTRIBUTORS and LICENSE files.
