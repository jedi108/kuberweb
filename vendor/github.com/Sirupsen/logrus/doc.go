/*
Package logger is a structured logger for Go, completely API compatible with the standard library logger.


The simplest way to use logger is simply the package-level exported logger:

  package main

  import (
    log "github.com/sirupsen/logger"
  )

  func main() {
    log.WithFields(log.Fields{
      "animal": "walrus",
      "number": 1,
      "size":   10,
    }).Info("A walrus appears")
  }

Output:
  time="2015-09-07T08:48:33Z" level=info msg="A walrus appears" animal=walrus number=1 size=10

For a full guide visit https://github.com/sirupsen/logger
*/
package logger
