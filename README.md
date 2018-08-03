<p align="center"><img src="web/public/assets/horizontal.svg" alt="wpdir" height="120px"></p>
<p align="center"><strong>Lightning fast searching of the WordPress Plugin &#38; Theme Repositories.</strong></p>

[![License](https://img.shields.io/badge/license-MIT-red.svg)](https://github.com/wpdirectory/wpdir/blob/master/LICENSE) [![Go Report Card](https://goreportcard.com/badge/github.com/wpdirectory/wpdir)](https://goreportcard.com/report/github.com/wpdirectory/wpdir) [![Build Status](https://travis-ci.org/wpdirectory/wpdir.svg?branch=master)](https://travis-ci.org/wpdirectory/wpdir)

WPDirectory provides a web interface for regex searching the WordPress.org [Plugin](https://plugins.svn.wordpress.org/) and [Theme](https://themes.svn.wordpress.org/) Repositories. You can read more about how and why it came to be [here](https://www.peterbooker.com/wpdirectory-reveal/).

The project is just transitioning from prototype to production status and is live at [https://wpdirectory.net/](https://wpdirectory.net/). My focus is now on UX and making it easier to host before I begin looking for sponsors to host it.

WPdirectory can do in seconds what was previously taking 15-60 minutes. It aims to significant reduce decision-making delays across Core WordPress teams and to empower uses across the community which were prohibitively difficult before.

## FAQs

### What problem does WPDirectory solve?

While working with [WPDS](https://github.com/PeterBooker/wpds) it seemed wasteful for so many people to be using Slurpers to download all theme/plugin files and perform slow local searches. I realised that a web based tool providing the same service might be faster and allow for easier collaboration, and so WPDirectory as an idea was born.

### What does WPDirectory actually do?

It creates and maintains an up-to-date copy of all current Plugin and Theme files from the official WordPress Directories. It then performs Trigram indexing on the files, allowing the frontend to perform lightning fast regular expression based searches (inspired by [etsy/hound](https://github.com/etsy/hound)).

### Why should I use WPdir?

Most searches can be completed in under a minute and many can be done in less than 10 seconds. You can then share search results by URL (even with searches set to private if they are sensitive). Anyone viewing the searches can sort results to identify the most important matches and use the inbuilt file viewer to review the context of matches.
