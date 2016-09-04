package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/mgoodness/k8s-bot/k8s"

	"github.com/go-chat-bot/bot/slack"
	"github.com/namsral/flag"
)

var (
	debug      bool
	kubeconfig string
	name       string
	token      string
)

func init() {
	flag.BoolVar(&debug, "debug", false, "Debugging mode")
	flag.StringVar(&name, "bot-name", "go-bot", "Bot name")
	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to kubeconfig file")
	flag.StringVar(&token, "slack-token", "", "Slack API token")

	log.SetFormatter(&log.JSONFormatter{})
}

func main() {
	flag.Parse()
	if debug {
		log.SetLevel(log.DebugLevel)
	}

	contextLogger := log.WithFields(log.Fields{
		"app": name,
	})

	contextLogger.Info("Starting up...")
	k8s.New(kubeconfig, contextLogger)

	slack.Run(token)
}
