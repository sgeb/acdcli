# acdcli [![Build Status](https://travis-ci.org/sgeb/acdcli.svg?branch=master)](https://travis-ci.org/sgeb/acdcli)

Command Line Interface tool for Amazon Cloud Drive.

## Features

Still work in progress. Focusing on read-only operations at first. See next
section for planned features.

Help:

```
% acdcli help
usage: acdcli [--version] [--help] <command> [<args>]

Available commands are:
    auth       Authorizes access to your Amazon Cloud Drive account
    ls         List files and folder in the root folder of the drive
    storage    Prints information on storage usage and quota
    version    Prints the acdcli version
```

Version:

```
% acdcli version
acdcli v0.1-dev (b35d3ba166af52def9356a5e05f812e56be5ef81)
using go-acd v0.1
```

Storage information, the forth and fifth columns are billable storage (displays
zeroes in this example taken on an Umlimited Everything plan):

```
% acdcli storage
Quota (last calculated 3 minutes ago)
Size: 100TiB, Available: 100TiB, Used: 0%

Usage (last calculated 3 minutes ago)
   Photos    31GiB    10,820        0B         0
    Video    56GiB       121        0B         0
      Doc   9.3MiB        39        0B         0
    Other    37GiB       361        0B         0
    Total   124GiB    11,341        0B         0
```

Listing items at the top-level folder (navigation to follow soon):

```
% acdcli ls
Archives/
Backups/
Documents/
Pictures/
Shared/
Videos/
example.jpg
sample.txt
```

## Planned features

Refer to the [milestones](https://github.com/sgeb/acdcli/milestones) and
[issues](https://github.com/sgeb/acdcli/issues) for more information on planned
features.

The following is a rough list of milestones:

* v0.1: read-only
* v0.2: read-write
* v1.0: caching and improvements
* v1.1: encrypted folders
* v1.2: FUSE filesystem
* v1.3: multi-account

## Contributing

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request

## Credit where credit is due

acdcli uses code, ideas and inspiration from the following projects:

* [go-github](https://github.com/google/go-github)
* [Terraform](https://www.terraform.io/)
* [google-api-go-client](https://github.com/google/google-api-go-client)

Thanks to the original authors for making their code available. Without their
contributions to open source, acdcli would not have turned out the way it has.

## License

Copyright (c) 2015 Serge Gebhardt. All rights reserved.

Use of this source code is governed by the ISC license that can be found in the
LICENSE file.
