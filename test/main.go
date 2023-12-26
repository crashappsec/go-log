package main

import (
	"github.com/crashappsec/go-log/src/log"
)

func main() {
	log.Info("default logger")

	parent := log.NewLogger().With(log.String("test", "pArEnT"))
	parent.Info("parent logger")
	parent.Info("parent logger", log.String("test", "parent"))

	child := parent.With(
		log.String("test", "cHiLd"),
		log.String("hello", "world"),
	)
	child.Info("child log")
	child.Info("child log", log.String("test", "child"))

	parent.Info("after child")
	parent.Info("after child", log.String("test", "parent"))
}
