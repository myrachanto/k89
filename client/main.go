package main

import (
	"flag"
	"fmt"
	"log"
	"net/rpc"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

var (
	Client  *rpc.Client
	cluster SystemController
)

type CommandLine struct{}

type (
	SystemController struct {
		Name                     string
		Nodes                    map[string]*Node
		Address                  string
		Status                   bool
		PodNumber                int
		NodeNumber               int
		NodesAssingment          map[string]int
		NextViableNodeToSchedule string
	}
	Node struct {
		Name        string
		Namespace   string
		Pods        map[string]*Pod
		Address     string
		Status      bool
		Schedulable bool
		StartTime   time.Time
		CreatedTime time.Time
	}
	Pod struct {
		Name        string
		Image       string
		Status      bool
		Address     string
		Ports       []string
		StartTime   time.Time
		CreatedTime time.Time
	}
)

func New() *CommandLine {
	log.SetPrefix("Client: ")
	return &CommandLine{}
}
func (cli *CommandLine) PrintUsage() {
	fmt.Println("Usage:")
	fmt.Println(" showcluster -Show Cluster - show information about the cluster")
	fmt.Println(" createNode -nodename NodeName creates a node")
	fmt.Println(" createPod -podname PodName creates a pod")
	fmt.Println(" nodelists - Prints the Nodes in the cluster")
	fmt.Println(" podlists - Prints the Pods in the cluster")
	fmt.Println(" deleteNode -nodename NodeName delete a node")
	fmt.Println(" deletePod -podname PodName delete a pod")
}
func (cli *CommandLine) ValidateArgs() {
	if len(os.Args) < 2 {
		cli.PrintUsage()
		log.Fatal("not enough arguments")
		// os.Exit(1)
		// runtime.Goexit()
	}
}
func (cli *CommandLine) ShowCluster() {
	// cluster := GetCluster()
	Client.Call("SystemController.ShowSystemControllerInfo", "", &cluster)
	fmt.Println("")
	// fmt.Println("debuging cluster step 2")
	color.Green(" %s %s Cluster %s \n", strings.Repeat(" ", 40), cluster.Name, strings.Repeat(" ", 40))
	color.Red("%s \n", strings.Repeat("_", 100))
	color.Green("%15.10s %30s %30s %20s  \n", "Name", "Address", "Namespace", "Start Time")
	for k, v := range cluster.Nodes {
		color.Green("%15.10s %30s %30s %20s \n", k, v.Address, v.Namespace, v.StartTime.Format("2006-01-02"))
		color.Blue("%s %s \n", strings.Repeat(" ", 20), strings.Repeat("#", 80))
		color.Green("%s %s %s Node %s \n", strings.Repeat(" ", 20), strings.Repeat(" ", 30), v.Name, strings.Repeat(" ", 30))
		color.Red("%s %s \n", strings.Repeat(" ", 20), strings.Repeat("_", 80))
		fmt.Printf(" %s %15.10s %20s %10s %15s %20s \n", strings.Repeat(" ", 10), "Name", "Image", "Status", "Address", "Start Time")
		for s, g := range v.Pods {
			color.Cyan("%s %15.10s|  %20s| %10v| %15s| %20s \n", strings.Repeat(" ", 10), s, g.Image, g.Status, g.Address, g.StartTime.Format("2006-01-02"))
		}
		fmt.Printf("%s \n", strings.Repeat("+", 100))
	}
}
func (cli *CommandLine) CreateNode(name string) {
	var node Node
	Client.Call("SystemController.AddNode", name, node)
	log.Println("node created successifuly")
}
func (cli *CommandLine) CreatePod(name, image, port string) {
	var pod Pod
	pod.Name = name
	pod.Image = image
	pod.Ports = append(pod.Ports, port)
	fmt.Println("-----------------------", pod)
	Client.Call("SystemController.AddPod", pod, &pod)
	log.Println("Pod created successifuly")
}
func (cli *CommandLine) DeleteNode(name string) {
	log.Println("node deleted successifuly")
}
func (cli *CommandLine) DeletePod(name, image, port string) {
	log.Println("Pod deleted successifuly")
}
func (cli *CommandLine) NodeList() {
	// cluster := GetCluster()
	Client.Call("SystemController.ShowSystemControllerInfo", "", &cluster)
	fmt.Println("")
	color.Green(" %s %s Cluster %s \n", strings.Repeat(" ", 40), cluster.Name, strings.Repeat(" ", 40))
	color.Red("%s \n", strings.Repeat("_", 100))
	color.Green("%15.10s %30s %30s %20s  \n", "Name", "Address", "Namespace", "Start Time")
	for k, v := range cluster.Nodes {
		color.Green("%15.10s %30s %30s %20s \n", k, v.Address, v.Namespace, v.StartTime.Format("2006-01-02"))
		color.Blue("%s %s \n", strings.Repeat(" ", 20), strings.Repeat("#", 80))
		color.Green("%s %s %s Node %s \n", strings.Repeat(" ", 20), strings.Repeat(" ", 30), v.Name, strings.Repeat(" ", 30))
		color.Red("%s %s \n", strings.Repeat(" ", 20), strings.Repeat("_", 80))
		// fmt.Printf(" %s %15.10s %20s %10s %15s %20s \n", strings.Repeat(" ", 10), "Name", "Image", "Status", "Address", "Start Time")
		// for s, g := range v.Pods {
		// 	color.Cyan("%s %15.10s|  %20s| %10v| %15s| %20s \n", strings.Repeat(" ", 10), s, g.Image, g.Status, g.Address, g.StartTime.Format("2006-01-02"))
		// }
		fmt.Printf("%s \n", strings.Repeat("+", 100))
	}
}

func (cli *CommandLine) PodList() {
	// cluster := GetCluster()
	Client.Call("SystemController.ShowSystemControllerInfo", "", &cluster)
	fmt.Println("")
	color.Green(" %s %s Cluster %s \n", strings.Repeat(" ", 40), cluster.Name, strings.Repeat(" ", 40))
	color.Red("%s \n", strings.Repeat("_", 100))
	// color.Green("%15.10s %30s %30s %20s  \n", "Name", "Address", "Namespace", "Start Time")
	for _, v := range cluster.Nodes {
		// color.Green("%15.10s %30s %30s %20s \n", k, v.Address, v.Namespace, v.StartTime.Format("2006-01-02"))
		// color.Blue("%s %s \n", strings.Repeat(" ", 20), strings.Repeat("#", 80))
		// color.Green("%s %s %s Node %s \n", strings.Repeat(" ", 20), strings.Repeat(" ", 30), v.Name, strings.Repeat(" ", 30))
		color.Red("%s %s \n", strings.Repeat(" ", 20), strings.Repeat("_", 80))
		fmt.Printf(" %s %15.10s %20s %10s %15s %20s \n", strings.Repeat(" ", 10), "Name", "Image", "Status", "Address", "Start Time")
		for s, g := range v.Pods {
			color.Cyan("%s %15.10s|  %20s| %10v| %15s| %20s \n", strings.Repeat(" ", 10), s, g.Image, g.Status, g.Address, g.StartTime.Format("2006-01-02"))
		}
		fmt.Printf("%s \n", strings.Repeat("+", 100))
	}
}
func (cli *CommandLine) Run() {
	cli.ValidateArgs()

	showClusterCmd := flag.NewFlagSet("showcluster", flag.ExitOnError)
	nodeListsCmd := flag.NewFlagSet("nodelists", flag.ExitOnError)
	podListsCmd := flag.NewFlagSet("podlists", flag.ExitOnError)
	CreateNodeCmd := flag.NewFlagSet("createNode", flag.ExitOnError)
	CreatePodCmd := flag.NewFlagSet("createPod", flag.ExitOnError)

	createNodeName := CreateNodeCmd.String("name", "", "The name of the node")
	PodName := CreatePodCmd.String("name", "", "The name of the Pod")
	PodImage := CreatePodCmd.String("image", "", "The image of the Pod")
	PodPort := CreatePodCmd.String("port", "", "The port of the Pod")

	switch os.Args[1] {
	case "showcluster":
		// fmt.Println("debuging cluster step 1")
		err := showClusterCmd.Parse(os.Args[2:])
		Handle(err)
	case "nodelists":
		err := nodeListsCmd.Parse(os.Args[2:])
		Handle(err)
	case "podlists":
		err := podListsCmd.Parse(os.Args[2:])
		Handle(err)
	case "createNode":
		err := CreateNodeCmd.Parse(os.Args[2:])
		Handle(err)
	case "createPod":
		err := CreatePodCmd.Parse(os.Args[2:])
		Handle(err)
	default:
		cli.ShowCluster()
		log.Fatal("not enough arguments - default case")
		// runtime.Goexit()
	}

	if showClusterCmd.Parsed() {
		cli.ShowCluster()
	}
	if nodeListsCmd.Parsed() {
		cli.NodeList()
	}
	if podListsCmd.Parsed() {
		cli.PodList()
	}
	if CreateNodeCmd.Parsed() {
		if *createNodeName == "" {
			CreateNodeCmd.Usage()
			log.Fatal("not enough arguments create node")
			// runtime.Goexit()
		}
		cli.CreateNode(*createNodeName)
	}

	if CreatePodCmd.Parsed() {
		if *PodName == "" || *PodImage == "" || *PodPort == "" {
			CreatePodCmd.Usage()
			log.Fatal("not enough arguments create pod")
			// runtime.Goexit()
		}

		cli.CreatePod(*PodName, *PodImage, *PodPort)
	}

}

func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}
func main() {
	client, err := rpc.DialHTTP("tcp", "localhost:2300")
	if err != nil {
		log.Fatal("something went wrong with Dialing: ", err)
	}
	Client = client
	cli := New()
	cli.Run()
}
