// +build darwin freebsd openbsd netbsd dragonfly
// +build !appengine,!gopherjs

package logger

import "golang.org/x/sys/unix"

const ioctlReadTermios = unix.TIOCGETA

type Termios unix.Termios
