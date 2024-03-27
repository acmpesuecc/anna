---
title: Week-1 Progress
date: 2024-03-25
type: post
tags:
  - progress
---

# Week-1

# Anirudh’s Progress

- Built the markdown to HTML renderer using the [Goldmark](https://github.com/yuin/goldmark) library
- Set up a web server to preview the rendered content locally
- Implemented a front matter YAML parser to retrieve page metadata
- Designed a reusable system of partials to build page layouts
- Implemented a recursive renderer for the content/ and static/ directories

## What's next?

- Restructure the project
- Refactor and improve the live reload functionality
- Implement the post tagging system

---

# Adhesh’s Progress

- Migrated to [Cobra](https://cobra.dev) for better **CLI integration**.
- Implemented important flags for serving a local **HTTP server** (---serve), and to explicitly mention the **port** to use (---addr).
- Implemented **Real-Time Directory/File change** watcher using [fsnotify](https://pkg.go.dev/github.com/fsnotify/fsnotify), and Hot-Reload system for reserving updated files instantaneously.
- Cleaned and sanitised some parts of the codebase.
- Included early optimizations to some functions of the codebase using goroutines and sync operations.

## What's next?

- Making **main.go** the entry point for the code.
- Fix issues with the watcher, and clean up the goroutine issues.
- Try studying about **parallelizing code functions**, and implementing it.
- Add a **developer** mode/flag for profiling the performance of the application (---dev)
- Look into [Huge theme](https://themes.gohugo.io) compatibility.
- Try integrating **JavaScript** in templates.

---

# Hegde’s Progress

- Switched to automatic filename parsing
- Implementing Draft Posts
- Complete CSS styling
- Cleaned up unnecessary post rendering

## Whats Next:

- JS injection as plugins into pages and individual posts
- Draft post rendering
- Chronological Feed for posts

# Nathan’s Progress

- Setup CI using GitHub actions which builds and deploys the SSG to [gh-pages](https://ssg-test-org.github.io)
- Setup `robots.txt` to (currently set to prevent indexing until we reach a v1.0 release)
- Implemented `sitemap.xml` to tell search engines how our site is structured
- Fixed `baseURL` it now, uses absolute paths from root `/` when loading stylesheets

## What's next?

- SEO optimization (verify semantic HTML generation)
- OGP tags in page headers

---
