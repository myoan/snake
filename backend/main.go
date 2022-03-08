package main

import (
	"context"
	"net/http"

	"agones.dev/agones/pkg/client/clientset/versioned"
	"agones.dev/agones/pkg/util/runtime"
	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

type GameServerSchema struct {
	IP   string `json:"ip"`
	Port int    `json:"port"`
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

	for _, item := range result.Items {
		ports := item.Status.Ports
		if len(ports) < 1 {
			continue
		}

		port := int(item.Status.Ports[0].Port)
		s := GameServerSchema{
			IP:   item.Status.Address,
			Port: port,
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
