package k8s

import (
	"fmt"

	"k8s.io/client-go/1.4/kubernetes"
	"k8s.io/client-go/1.4/pkg/api"
	"k8s.io/client-go/1.4/tools/clientcmd"

	log "github.com/Sirupsen/logrus"

	"github.com/go-chat-bot/bot"
)

var (
	clientSet *kubernetes.Clientset
	logger    *log.Entry
)

func init() {
	bot.RegisterCommand(
		"k8s",
		"Issues commands to a Kubernetes API server.",
		"",
		k8s)
}

// New configures the Kubernetes client.
func New(k string, l *log.Entry) {
	logger = l

	config, err := clientcmd.BuildConfigFromFlags("", k)
	if err != nil {
		logger.WithFields(log.Fields{
			"source": "k8s_client",
		}).Fatalln(err)
	}

	logger.WithFields(log.Fields{
		"source": "k8s_client",
	}).Debugln("Loaded configuration")

	if c, err := kubernetes.NewForConfig(config); err != nil {
		logger.WithFields(log.Fields{
			"source": "k8s_client",
		}).Fatalln(err)
	} else {
		clientSet = c
		logger.WithFields(log.Fields{
			"source": "k8s_client",
		}).Debugln("Created clientset")
	}
}

func logCommand(command bot.Cmd) {
	logger.WithFields(log.Fields{
		"channel":   command.Channel,
		"command":   command.Command,
		"raw_args":  command.RawArgs,
		"source":    "slack_client",
		"user_name": command.User.RealName,
		"user_nick": command.User.Nick,
	}).Debugln("Received command")
}

func handleGet(args []string) (msg string) {
	msg = ""

	switch args[1] {
	case "po", "pod", "pods":
		if len(args) <= 2 {
			pods, err := clientSet.Core().Pods("").List(api.ListOptions{})
			if err != nil {
				logger.WithFields(log.Fields{
					"source": "k8s_client",
				}).Errorln(err)
				msg = "Error retrieving pods"
			}

			for _, pod := range pods.Items {
				podName := pod.ObjectMeta.GetName()
				podStatus := pod.Status.Phase
				containersReady := len(pod.Status.ContainerStatuses)
				containersTotal := len(pod.Spec.Containers)
				msg += fmt.Sprintf("%s\t%s\t%d/%d\n", podName, podStatus,
					containersReady, containersTotal)
				logger.WithFields(log.Fields{
					"pod_name":         podName,
					"pod_status":       podStatus,
					"containers_ready": containersReady,
					"containers_total": containersTotal,
				}).Debugln("Retrieved pod metadata")
			}
		}
	}

	return
}

func k8s(command *bot.Cmd) (msg string, err error) {
	msg = ""

	logCommand(*command)
	switch command.Args[0] {
	case "get":
		msg = handleGet(command.Args)
	}

	return
}
