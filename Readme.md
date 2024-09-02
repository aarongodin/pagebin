# pagebin

Modern and minimalist CMS comparable to Wordpress or Ghost.

## Highlights

- Deployable as a single binary with no dependencies
- Content authoring using a modern block-based page editor
- Unopinionated theme system
- Support for multi-author with fine-grained access control

## Install

Install `pagebin` on your machine through a package manager, precompiled binary, or run as a Docker container.

### homebrew

Available for homebrew users through:

```
brew install pagebin
```

### apt (debian/ubuntu)

Install through aptitude on supported Linux operating systems:

```
todo
```

### Docker

The docker image is provided through Docker hub.

```
todo
```

## Quickstart

You can quickly see pagebin working by running:

```
pagebin serve
```

then open [http://localhost:8080](localhost:8080) to view the default site.

## Overview

Pagebin manages pages through a slightly different approach to many CMS software. The entire site is versioned, so edits to the website do not take effect until the website is published. This allows the authors to make multiple edits, preview the site, and promote all the changes at once when it's ready.

### Pages

Everything in pagebin is a page. Pages have a URL path, a title, and content.

