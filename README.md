# acdcli [![Build Status](https://travis-ci.org/sgeb/acdcli.svg?branch=master)](https://travis-ci.org/sgeb/acdcli)

Command Line Interface tool for Amazon Cloud Drive.

## Installation

Binary releases are made available upon reaching a milestone. No milestone has
been reached yet, therefore no binaries. Fear not, acdcli can still be installed
from source:

```bash
# clone this repo
% git clone https://github.com/sgeb/acdcli.git

# get the right versions of dependencies (through gpm)
# optional, is part of the two next commands
% make updatedeps

# compile binaries for your current platform
% ACD_API_CLIENTID="your clientID" ACD_API_SECRET="your secret" make dev

# compile binaries for all supported platforms
% ACD_API_CLIENTID="your clientID" ACD_API_SECRET="your secret" make bin
```

You will need to register an API client for the Amazon Cloud Drive to get a
clientID and secret. You can do so
[here](https://developer.amazon.com/public/apis/experience/cloud-drive/content/getting-started).

Note that invoking the line as written above will store your clientID and secret
in your shell's history. It's safer to use a commandline-accessible password
store (such as [pass](http://www.passwordstore.org/)) and invoke the line as
follows:

```bash
% ACD_API_CLIENTID=$(pass ApiKeys/acdcli/clientid) ACD_API_SECRET=$(pass ApiKeys/acdcli/secret) make bin
```

Make sure to setup `pass` and adjust the path to the password file.

## Features

Still work in progress. Focusing on read-only operations at first. See next
section for planned features.

Help:

```
% acdcli help
usage: acdcli [--version] [--help] <command> [<args>]

Available commands are:
    auth       Authorizes access to your Amazon Cloud Drive account
    get        Download files
    info       Display a node's metadata
    ls         List files and folder in the root folder of the drive
    storage    Prints information on storage usage and quota
    version    Prints the acdcli version
```

Storage information, the forth and fifth columns are billable storage (displays
zeroes in this example taken on an Umlimited Everything plan):

```
% acdcli storage
Quota (last calculated 1 second ago)
Size: 100TiB, Available: 100TiB, Used: 0%

Usage (last calculated now)
   Photos   31GiB  10,820      0B       0
    Video   56GiB     121      0B       0
      Doc  9.3MiB      40      0B       0
    Other   37GiB     361      0B       0
    Total  124GiB  11,342      0B       0
```

Listing items at the top-level folder:

```
% acdcli ls
      - Archives/
      - Backups/
      - Documents/
      - Pictures/
      - Shared/
      - Videos/
 700KiB example.jpg
 1.2MiB sample.doc
```

Listing items by specifying the folder (initial `/` is optional):

```
% acdcli ls /Documents/Samples
      - Subfolder/
 2.7KiB Lorem Ipsum.txt
  26KiB sample.jpg
  20MiB sample.mp4
    15B sample.txt
```

Listing a specific file (initial `/` is optional):

```
% acdcli ls /Documents/Samples/sample.mp4
  20MiB sample.mp4
```

Download a file to current directory:

```
% acdcli get /Documents/Samples/sample.jpg
```

Download a file to `/tmp/some/folder` (output directory must exist):

```
% acdcli get /Documents/Samples/sample.jpg /tmp/some/folder
```

Download a file to a local name (output directory must exist):

```
% acdcli get /Documents/Samples/sample.jpg /tmp/some/folder/downloaded.jpg
```

Display any node's metadata, works for files and folders. In case of a file,
Amazon Cloud Drive might extract additional information.

```json
% acdcli info /Documents/Samples/sample.mp4
{
    "id": "6wDU-sqpRXSWOq1zYYPM0A",
    "kind": "FILE",
    "version": 12,
    "labels": [],
    "contentProperties": {
        "contentDate": "1970-01-01T00:00:00.000Z",
        "extension": "mp4",
        "size": 21069678,
        "video": {
            "rotate": 0,
            "bitrate": 1436830,
            "audioCodec": "aac",
            "creationDate": "1970-01-01T00:00:00.000Z",
            "videoFrameRate": 25,
            "encoder": "Lavf53.24.2",
            "audioBitrate": 383869,
            "audioSampleRate": 48000,
            "audioChannelLayout": "5.1",
            "duration": 117.312,
            "videoBitrate": 1048652,
            "audioChannels": 6,
            "width": 1280,
            "height": 720,
            "videoCodec": "h264"
        },
        "contentType": "video/mp4",
        "version": 1,
        "md5": "442a2dc932b8a4a26ca12b73e796507b"
    },
    "createdDate": "2015-05-17T01:54:05.291Z",
    "createdBy": "CloudDriveWeb",
    "restricted": false,
    "modifiedDate": "2015-05-17T01:55:16.621Z",
    "name": "sample.mp4",
    "isShared": false,
    "properties": {
        "CloudDrive": {
            "Processing": "VIDEO_PROCESSED"
        }
    },
    "parents": [
        "L_9wVa74QWyYmqJITeYmhw"
    ],
    "status": "AVAILABLE"
}
```

Version:

```
% acdcli version
acdcli v0.1-dev (6c9065149252edcb9ff212296bc4b514e0bb12bb)
using go-acd v0.1
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
