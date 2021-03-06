# mp

[![Go Report Card](https://goreportcard.com/badge/github.com/martinplaner/mp)](https://goreportcard.com/report/github.com/martinplaner/mp)

Highly sophisticated [initialism](https://en.wikipedia.org/wiki/Acronym#Nomenclature) to hyphenated [compound words](https://en.wikipedia.org/wiki/Compound_(linguistics)) generator.

## Example

*MP* -> *Melanchton-Paralogismus*

## Disclaimer

All this effort for a stupid inside joke. Don't worry if you don't get it, it's not that funny anyway...

## Usage

```
Usage of mp:
  -debug
        Enable verbose debug mode
  -file string
        Path to word list (one word per line) (default "words.txt")
  -listen string
        TCP address for the server to listen on, in the form 'host:port' (default ":8080")
```

The `PORT` environment variable can also be used to set the default listening port.

### Docker

A Docker image is available through Docker Hub: https://hub.docker.com/r/martinplaner/mp 

```
$ docker run -it --rm -v $PWD/words.txt:/data/words.txt -p 8080:8080 martinplaner/mp:latest
```

⚠ Don't use the old image from GitHub Packages Docker Registry (aka docker.pkg.github.com), since it is deprecated and GitHub Packages Docker Registry will sunset early next year. Unfortunately it is not possible to delete public images from GitHub Packages, so just be aware.

## Licenses

### Source code

Copyright 2020 Martin Planer. All rights reserved. Use of this source code is governed by a BSD-style license that can be found in the LICENSE file.

### Word list

Korpusbasierte Wortgrundformenliste DEREWO, v-ww-bll-320000g-2012-12-31-1.0, mit Benutzerdokumentation, http://www.ids-mannheim.de/derewo, © Institut für Deutsche Sprache, Programmbereich Korpuslinguistik, Mannheim, Deutschland, 2013.

### Twemoji (Favicon)

Copyright 2019 Twitter, Inc and other contributors. 

Code licensed under the MIT License: <http://opensource.org/licenses/MIT>

Graphics licensed under CC-BY 4.0: <https://creativecommons.org/licenses/by/4.0/>
