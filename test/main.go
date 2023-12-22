package main

import (
	"github.com/crashappsec/go-log"
	level "github.com/crashappsec/go-log/log"
)

func main() {
	level.Info("default logger")

	parent := log.NewLogger().With(level.String("test", "pArEnT"))
	parent.Info("parent logger")
	parent.Info("parent logger", level.String("test", "parent"))

	child := parent.With(level.String("test", "cHiLd"), level.String("hello", "world"))
	child.Info("child log")
	child.Info("child log", level.String("test", "child"))

	parent.Info("after child")
	parent.Info("after child", level.String("test", "parent"))
}
