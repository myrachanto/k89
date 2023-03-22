package main

import (
	"fmt"
	"log"
	"strings"
	"time"
)

const (
	name             = "K89: "
	address          = "127.0.0.10"
	address1         = "127.0.0.1"
	defaultNameSpace = "defaultNamespace"
)

type (
	SystemControllerInterface interface {
		StopSystem()
		SystemStatus()
		AddNode(string) error
		NodeStatus(string) (bool, error)
		DeleteNode(string)
		SchedulePod(string) (bool, error)
		DeletePod(string) error
		PodStatus(string) (bool, error)
	}
	NodeInterface interface {
		StartBackground() (bool, error)
	}
)
type (
	SystemController struct {
		Name    string
		Nodes   map[string]*Node
		Address string
		Status  bool
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
		Ports       []string
		StartTime   time.Time
		CreatedTime time.Time
	}
)

func New() *SystemController {
	log.SetPrefix(name)
	log.Println("Started.......")
	return &SystemController{
		Name:    name,
		Address: address,
		Nodes:   make(map[string]*Node),
		Status:  true,
	}
}
func (sc *SystemController) StopSystem() {
	sc.Status = false
	log.Println("Stoped!")
}
func (sc *SystemController) SystemStatus() {
	if sc.Status {
		log.Println("is running!...")
	} else {
		log.Println("has being Stoped!")
	}
}
func (sc *SystemController) AddNode(name string) error {
	schedulable := IsSchedulable(name)
	node := sc.NewNode(name, schedulable)
	err := node.StartBackground()
	if err != nil {
		return err
	} else {
		sc.Nodes[node.Name] = node
		return nil
	}
}
func (sc *SystemController) NewNode(name string, schedulable bool) *Node {
	if name == "" {
		log.Panic("Please use a valid Name for a Node!")
	}
	address, err := sc.CreateAddress(name)
	if err != nil {
		log.Panic("something Went wrong with creating an Address")
	}
	return &Node{
		Name:        name,
		Namespace:   defaultNameSpace,
		Pods:        make(map[string]*Pod),
		Address:     address,
		Schedulable: schedulable,
		CreatedTime: time.Now(),
	}
}
func (sc *SystemController) NodeStatus(name string) {
	node, ok := sc.Nodes[name]
	if ok && node.Status {
		log.Printf("%s has been running since %v", name, node.StartTime)
		return
	}
	log.Printf("Node %s has failed working!", name)
}
func (sc *SystemController) DeleteNode(name string) {
	_, ok := sc.Nodes[name]
	if !ok {
		log.Println("Node not found!")
		return
	}
	delete(sc.Nodes, name)
}
func (sc *SystemController) SchedulePod(namespace, name, image string, ports []string) (bool, error) {
	var (
		res bool = false
	)
	pod := NewPod(name, image, ports)
	for _, v := range sc.Nodes {
		if v.Namespace == namespace && v.Schedulable {
			v.Pods[name] = pod
			res = true
			break
		}
	}
	return res, fmt.Errorf("could not schedule the pod")
}
func (sc *SystemController) DeletePod(namespace, name string) error {
	for _, v := range sc.Nodes {
		if v.Namespace == namespace {
			_, ok := v.Pods[name]
			if ok {
				delete(v.Pods, name)
				break
			}
		}
	}
	return fmt.Errorf("could not find the pod")
}
func (sc *SystemController) PodStatus(namespace, name string) {
	for _, v := range sc.Nodes {
		if v.Namespace == namespace {
			pod, ok := v.Pods[name]
			if ok && pod.Status {
				log.Printf("%s has been running since %v", name, pod.StartTime)
				break
			}
		}
	}
}
func NewPod(name, image string, port []string) *Pod {
	if name == "" {
		log.Panic("Please use a valid Name for a Node!")
	}
	return &Pod{
		Name:        name,
		Image:       image,
		Ports:       port,
		CreatedTime: time.Now(),
	}
}
func (n *Node) StartBackground() error {
	n.StartTime = time.Now()
	n.Status = true
	log.Printf("%s background services started \n", n.Name)
	return nil
}
func (sc *SystemController) CreateAddress(name string) (string, error) {
	if len(sc.Nodes) == 0 {
		return address1, nil
	} else {
		l := len(sc.Nodes)
		return fmt.Sprintf("127.0.0.%d", l), nil
	}
}
func IsSchedulable(name string) bool {
	res := false
	lastString := strings.Split(name, "_")
	lastStringSchedulable := lastString[len(lastString)-1]
	if lastStringSchedulable == "master" {
		res = true
		return res
	}
	return res
}

func main() {

}
