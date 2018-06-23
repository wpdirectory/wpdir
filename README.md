<p align="center"><img src="web/public/static/assets/horizontal.svg" alt="wpdir" height="200px"></p>


**Fast searching of the WordPress Plugin & Theme Repositories.**

[![License](https://img.shields.io/badge/license-MIT-red.svg)](https://github.com/wpdirectory/wpdir/blob/master/LICENSE) [![Go Report Card](https://goreportcard.com/badge/github.com/wpdirectory/wpdir)](https://goreportcard.com/report/github.com/wpdirectory/wpdir) [![Build Status](https://travis-ci.org/wpdirectory/wpdir.svg?branch=master)](https://travis-ci.org/wpdirectory/wpdir)

WPDirectory provides a web interface for viewing and searching the WordPress.org [Plugin](https://plugins.svn.wordpress.org/) and [Theme](https://themes.svn.wordpress.org/) Repositories. It is spawned from the [WPDS](https://github.com/PeterBooker/wpds) project and aims to avoid many people having to download and search the data locally.

WPDirectory is still in early development but now has an alpha version intermittently live at [https://wpdirectory.net/](https://wpdirectory.net/). I am currently working towards improving the user-experience and then I will arrange production hosting.

## FAQs

### What problem does WPDirectory solve?

While working with [WPDS](https://github.com/PeterBooker/wpds) it seemed wasteful for so many people to be using Slurpers to download all theme/plugin files and perform individual searches. I realise that a web based tool providing the same service might be better faster and allow for easier collaboration, and so WPDirectory was borne.

### What does WPDirectory actually do?

It creates and maintains an up-to-date copy of all current Plugin and Theme data from the official WordPress Directories. It then performs Trigram indexing on the files, allowing the frontend to perform lightning fast regular expression based searches (inspired by [etsy/hound](https://github.com/etsy/hound)).
