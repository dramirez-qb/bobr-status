package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	clientset *kubernetes.Clientset
	namespace string
	port      string
)

func init() {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		config, err = clientcmd.BuildConfigFromFlags("", os.Getenv("KUBE_CONFIG"))
		if err != nil {
			config, err = clientcmd.BuildConfigFromFlags("", filepath.Join(os.Getenv("HOME"), ".kube", "config"))
			if err != nil {
				panic(err.Error())
			}
		}
	}
	// creates the clientset
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	namespace = os.Getenv("POD_NAMESPACE")
	if namespace == "" {
		namespace = apiv1.NamespaceDefault
	}
	port = os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

}

type Info struct {
	Name              string            `json:"name,omitempty"`
	Replicas          int32             `json:"replicas,omitempty"`
	CreationTimestamp time.Time         `json:"creationtimestamp,omitempty"`
	PodIPs            []apiv1.PodIP     `json:"podips,omitempty"`
	ClusterIP         string            `json:"ip,omitempty"`
	Type              apiv1.ServiceType `json:"type,omitempty"`
}

func NewRouter() *gin.Engine {
	router := gin.Default()
	deploymentsClient := clientset.AppsV1().Deployments(namespace)
	statefulSetsClient := clientset.AppsV1().StatefulSets(namespace)
	podsClient := clientset.CoreV1().Pods(namespace)
	serviceClient := clientset.CoreV1().Services(namespace)
	router.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"POST", "GET", "OPTIONS"},
		AllowHeaders:  []string{"Origin", "Content-Type"},
		ExposeHeaders: []string{"Content-Length"},
	}))

	// Ping handler
	router.Any("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	// Health handler
	router.Any("/healthz", func(c *gin.Context) {
		c.String(200, "healthy")
	})

	// Simple group: v1
	api := router.Group("/v1")
	{
		status := api.Group("/status")
		{
			status.GET("/deployments", func(c *gin.Context) {
				list, err := deploymentsClient.List(context.TODO(), metav1.ListOptions{})
				deployments := make([]Info, 0)
				if err != nil {
					panic(err)
				}
				for _, d := range list.Items {
					deployments = append(deployments, Info{
						Name:              d.Name,
						Replicas:          *d.Spec.Replicas,
						CreationTimestamp: d.ObjectMeta.CreationTimestamp.Time,
					})
					fmt.Printf(" * %s (%d replicas)\n", d.Name, *d.Spec.Replicas)
				}
				c.JSON(http.StatusOK, deployments)
			})
			status.GET("/statefulsets", func(c *gin.Context) {
				list, err := statefulSetsClient.List(context.TODO(), metav1.ListOptions{})
				statefulsets := make([]Info, 0)
				if err != nil {
					panic(err)
				}
				for _, d := range list.Items {
					statefulsets = append(statefulsets, Info{
						Name:              d.Name,
						Replicas:          *d.Spec.Replicas,
						CreationTimestamp: d.ObjectMeta.CreationTimestamp.Time,
					})
					fmt.Printf(" * %s (%d replicas)\n", d.Name, *d.Spec.Replicas)
				}
				c.JSON(http.StatusOK, statefulsets)
			})
			status.GET("/pods", func(c *gin.Context) {
				list, err := podsClient.List(context.TODO(), metav1.ListOptions{})
				pods := make([]Info, 0)
				if err != nil {
					panic(err)
				}
				for _, p := range list.Items {
					pods = append(pods, Info{
						Name:              p.Name,
						CreationTimestamp: p.ObjectMeta.CreationTimestamp.Time,
						PodIPs:            p.Status.PodIPs,
					})
					fmt.Printf(" * %s\n", p.Name)
				}
				c.JSON(http.StatusOK, pods)
			})
			status.GET("/services", func(c *gin.Context) {
				list, err := serviceClient.List(context.TODO(), metav1.ListOptions{})
				services := make([]Info, 0)
				if err != nil {
					panic(err)
				}
				for _, s := range list.Items {
					services = append(services, Info{
						Name:              s.Name,
						CreationTimestamp: s.ObjectMeta.CreationTimestamp.Time,
						Type:              s.Spec.Type,
						ClusterIP:         s.Spec.ClusterIP,
					})
					fmt.Printf(" * %s\n", s.Name)
				}
				c.JSON(http.StatusOK, services)
			})
		}
	}
	return router
}

func startServer() {
	// get and start router
	router := NewRouter()
	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}

func main() {
	startServer()
}
