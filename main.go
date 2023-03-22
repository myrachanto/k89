package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/fatih/color"
)

const (
	name             = "K89: "
	address          = "127.0.0.0"
	address1         = "127.0.1.0"
	address2         = "127.0.2.0"
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

func New() *SystemController {
	log.SetPrefix(name)
	log.Println("Started.......")
	return &SystemController{
		Name:            name,
		Address:         address,
		Nodes:           make(map[string]*Node),
		NodesAssingment: make(map[string]int),
		Status:          true,
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
	schedulable := isSchedulable(name)
	node := sc.newNode(name, schedulable)
	node.Status = true
	err := node.StartBackground()
	if err != nil {
		return err
	} else {
		sc.Nodes[node.Name] = node
		sc.NodeNumber++
		if schedulable {
			sc.NodesAssingment[node.Name] = 0
		}
		return nil
	}
}
func (sc *SystemController) newNode(name string, schedulable bool) *Node {
	if name == "" {
		log.Panic("Please use a valid Name for a Node!")
	}
	address, err := sc.CreateNodeAddress()
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
func (sc *SystemController) SchedulePod(name, image string, ports []string) (bool, error) {
	var (
		res bool = false
	)
	pod := newPod(name, image, ports)
	pod.Status = true
	podaddr, err := sc.CreatePodAddress()
	if err != nil {
		return res, err
	}
	pod.Address = podaddr
	bestNodeCadidates := sc.bestNodeCadidate()
	fmt.Printf("Best candidate for this schedule is %s \n", bestNodeCadidates)
	for _, v := range sc.Nodes {
		if v.Namespace == defaultNameSpace && v.Schedulable && v.Name == bestNodeCadidates {
			v.Pods[name] = pod
			res = true
			sc.PodNumber++
			if isSchedulable(v.Name) {
				sc.NodesAssingment[v.Name]++
			}
			break
		}
	}
	return res, fmt.Errorf("could not schedule the pod")
}
func (sc *SystemController) bestNodeCadidate() string {
	// color.Red("The list of acceptable schedulables", sc.NodesAssingment)
	var bestPotentialCadidate string
	var res []int
	for _, v := range sc.NodesAssingment {
		res = append(res, v)
	}
	sorted(res)
	// color.Red("The list of acceptable schedulables", res)
	bestPotentialCadidateint := res[0]
	for k, v := range sc.NodesAssingment {
		if v == bestPotentialCadidateint {
			bestPotentialCadidate = k
			break
		}
	}
	return bestPotentialCadidate
}
func sorted(a []int) {
	for i := 0; i < len(a); i++ {
		for j := i + 1; j < len(a); j++ {
			if a[i] > a[j] {
				a[i], a[j] = a[j], a[i]
			}
		}
	}
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
func newPod(name, image string, port []string) *Pod {
	if name == "" {
		log.Panic("Please use a valid Name for a Node!")
	}
	return &Pod{
		Name:        name,
		Image:       image,
		Ports:       port,
		CreatedTime: time.Now(),
		StartTime:   time.Now(),
	}
}
func (n *Node) StartBackground() error {
	n.StartTime = time.Now()
	n.Status = true
	log.Printf("%s background services started \n", n.Name)
	return nil
}
func (sc *SystemController) CreateNodeAddress() (string, error) {
	// fmt.Println(">>>>>>>>>>>>>>>", sc.NodeNumber)
	if sc.NodeNumber == 0 {
		return address1, nil
	} else {
		return fmt.Sprintf("127.0.1.%d", sc.NodeNumber), nil
	}
}
func (sc *SystemController) CreatePodAddress() (string, error) {
	if sc.PodNumber == 0 {
		return address2, nil
	} else {
		return fmt.Sprintf("127.0.2.%d", sc.PodNumber), nil
	}
}
func (sc *SystemController) ShowSystemControllerInfo() {
	fmt.Println("")
	color.Blue(" %s %s Cluster %s \n", strings.Repeat(" ", 40), sc.Name, strings.Repeat(" ", 40))
	color.Red("%s \n", strings.Repeat("_", 100))
	color.Green("%15.10s %30s %30s %20s  \n", "Name", "Address", "Namespace", "Start Time")
	for k, v := range sc.Nodes {
		fmt.Printf("%15.10s %30s %30s %20s \n", k, v.Address, v.Namespace, v.StartTime.Format("2006-01-02"))
		color.Blue("%s %s \n", strings.Repeat(" ", 20), strings.Repeat("#", 80))
		color.Green("%s %s %s Node %s \n", strings.Repeat(" ", 20), strings.Repeat(" ", 30), v.Name, strings.Repeat(" ", 30))
		color.Red("%s %s \n", strings.Repeat(" ", 20), strings.Repeat("_", 80))
		fmt.Printf(" %s %15.10s %20s %10s %10s %20s \n", strings.Repeat(" ", 20), "Name", "Image", "Status", "Address", "Start Time")
		for s, g := range v.Pods {
			fmt.Printf("%s %15.10s  %20s %10v %10s %20s \n", strings.Repeat(" ", 20), s, g.Image, g.Status, g.Address, g.StartTime.Format("2006-01-02"))
		}
		fmt.Printf("%s \n", strings.Repeat("+", 100))
	}
}
func isSchedulable(name string) bool {
	res := true
	lastString := strings.Split(name, "_")
	lastStringSchedulable := lastString[len(lastString)-1]
	if lastStringSchedulable == "m11" {
		res = false
		return res
	}
	return res
}

func main() {
	systemcontrol := New()
	systemcontrol.AddNode("master_m11")
	systemcontrol.AddNode("worker1")
	systemcontrol.AddNode("worker2")
	ports := []string{"4000"}
	systemcontrol.SchedulePod("goapp", "myrachanto/goapp", ports)
	ports1 := []string{"6379"}
	systemcontrol.SchedulePod("redis", "redis", ports1)

	systemcontrol.ShowSystemControllerInfo()

}
