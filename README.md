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
  -default string
        Default fallback query term, if not provided. (default "MP")
  -file string
        Path to word list (one word per line, with optional word type) (default "words.txt")
  -listen string
        TCP address for the server to listen on, in the form 'host:port' (default ":8080")
  -mode value
        Word generation mode, one of a|adjective|c|compound
  -once string
        Run generation once with the given query and print result, then quit
```

The `PORT` environment variable can also be used to set the default listening port.

### Docker

A Docker image is available through Docker Hub: https://hub.docker.com/r/martinplaner/mp 

```
$ docker run -it --rm -v $PWD/words.txt:/data/words.txt -p 8080:8080 martinplaner/mp:latest
```

## Licenses

### Source code

Copyright 2020 Martin Planer. All rights reserved. Use of this source code is governed by a BSD-style license that can be found in the LICENSE file.

### Word list (words_de.txt)

Korpusbasierte Wortgrundformenliste DEREWO, v-ww-bll-320000g-2012-12-31-1.0, mit Benutzerdokumentation, http://www.ids-mannheim.de/derewo, © Institut für Deutsche Sprache, Programmbereich Korpuslinguistik, Mannheim, Deutschland, 2013.

### Word list (word_en.txt)

Part 1: From https://www.desiquintans.com/nounlist, under the public domain.

Part 2: Extracted from Dwarf Fortress file `raw/objects/language_words.txt`, also under the public domain, according to the accompanying `readme.txt`.

### Twemoji (Favicon)

Copyright 2019 Twitter, Inc and other contributors. 

Code licensed under the MIT License: <http://opensource.org/licenses/MIT>

Graphics licensed under CC-BY 4.0: <https://creativecommons.org/licenses/by/4.0/>
