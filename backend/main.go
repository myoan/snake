package main

import (
	"context"
	"net/http"

	v1 "agones.dev/agones/pkg/apis/agones/v1"
	"agones.dev/agones/pkg/client/clientset/versioned"
	"agones.dev/agones/pkg/util/runtime"
	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

type GameServerSchema struct {
	IP    string `json:"ip"`
	State string `json:"state"`
	Port  int    `json:"port"`
}

func RoomHandler(c *gin.Context) {
	config, err := rest.InClusterConfig()
	logger := runtime.NewLoggerWithSource("main")
	if err != nil {
		logger.WithError(err).Fatal("Could not create in cluster config")
	}

	agonesClient, err := versioned.NewForConfig(config)
	if err != nil {
		logger.WithError(err).Fatal("Could not create the agones api clientset")
	}

	result, err := agonesClient.AgonesV1().GameServers("default").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	c.Writer.Header().Set("Access-Control-Allow-Origin", "http://snake-game.myoan.dev")
	for _, item := range result.Items {
		if item.Status.State != v1.GameServerStateReady {
			continue
		}

		ports := item.Status.Ports
		if len(ports) < 1 {
			continue
		}

		port := int(item.Status.Ports[0].Port)
		s := GameServerSchema{
			IP:    item.Status.Address,
			State: string(item.Status.State),
			Port:  port,
		}
		c.JSON(http.StatusOK, s)
		return
	}

	c.JSON(http.StatusNotFound, gin.H{"msg": "Ready state server not found."})
}

func main() {
	r := gin.Default()
	r.GET("/room", RoomHandler)
	r.Run()
}
