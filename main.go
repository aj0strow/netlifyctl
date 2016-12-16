package main

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/aj0strow/netlifyctl/commands"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:    true,
		DisableTimestamp: false,
		TimestampFormat:  time.RFC822Z,
	})
	commands.Execute()
}
