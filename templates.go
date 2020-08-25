// +build !release

package main

import "net/http"

// Assets contains project assets.
var templates http.FileSystem = http.Dir("templates")
