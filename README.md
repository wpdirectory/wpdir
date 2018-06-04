# WPdirectory

**Fast searching of the WordPress Plugin & Theme Repositories.**

[![License](https://img.shields.io/badge/license-MIT-red.svg)](https://github.com/wpdirectory/wpdir/blob/master/LICENSE) [![Go Report Card](https://goreportcard.com/badge/github.com/wpdirectory/wpdir)](https://goreportcard.com/report/github.com/wpdirectory/wpdir)

WPDirectory provides a web interface for viewing and searching the WordPress.org [Plugin](https://plugins.svn.wordpress.org/) and [Theme](https://themes.svn.wordpress.org/) Repositories. It is spawned from the [WPDS](https://github.com/PeterBooker/wpds) project and aims to avoid many people having to download and search the data locally.

WPDirectory is still in early development and the website is not yet live. I hope to have a functioning prototype by the end of June. Once it is live it will be available at [wpdirectory.net](https://wpdirectory.net).

## FAQs

### What problem does WPDirectory solve?

While working with [WPDS](https://github.com/PeterBooker/wpds) it seemed wasteful for so many people to be using Slurpers to download all theme/plugin files and perform individual searches. I realise that a web based tool providing the same service might be better faster and allow for easier collaboration, and so WPDirectory was borne.

### What does WPDirectory actually do?

It creates and maintains an up-to-date copy of all current Plugin and Theme data from the official WordPress Directories. It then performs Trigram indexing on the files, allowing the frontend to perform lightning fast regular expression based searches (inspired by [etsy/hound](https://github.com/etsy/hound)).