package main

import "testing"

func TestAddNode(t *testing.T) {
	cluster := New()
	var n Node
	cluster.AddNode("vamos", &n)
	if len(cluster.Nodes) != 1 {
		t.Errorf("expected %d but got %d", 1, len(cluster.Nodes))
	}
	if cluster.NodeNumber != 1 {
		t.Errorf("expected %d but got %d", 1, cluster.NodeNumber)
	}
}
func TestDeleteNode(t *testing.T) {
	cluster := New()
	var n Node
	cluster.AddNode("vamos", &n)
	cluster.DeleteNode("vamos")
	if len(cluster.Nodes) != 0 {
		t.Errorf("expected %d but got %d", 0, len(cluster.Nodes))
	}
	if cluster.NodeNumber != 0 {
		t.Errorf("expected %d but got %d", 0, cluster.NodeNumber)
	}
}
func TestMasterNodeBackGroundProcesesses(t *testing.T) {
	cluster := New()
	var n Node
	cluster.AddNode("vamos_m11", &n)
	cluster.MasterNodeBackGroundProcesesses()
	if len(cluster.Nodes) != 1 {
		t.Errorf("expected %d but got %d", 1, len(cluster.Nodes))
	}
	if cluster.NodeNumber != 1 {
		t.Errorf("expected %d but got %d", 1, cluster.NodeNumber)
	}
	if cluster.PodNumber != 3 {
		t.Errorf("expected %d but got %d", 3, cluster.PodNumber)
	}
}
func TestSchedulePod(t *testing.T) {
	cluster := New()
	var n Node
	cluster.AddNode("vamos", &n)
	var p Pod
	p.Name = "mongo"
	p.Image = "mongo"
	p.Ports = []string{"mongo"}
	resp, err := cluster.SchedulePod(p.Name, p.Image, p.Ports)
	if err != nil {
		t.Errorf("err should be nil %s", err)
	}
	if p.Name != resp.Name {
		t.Errorf("expected %s but got %s", resp.Name, p.Name)
	}
	if cluster.PodNumber != 1 {
		t.Errorf("expected %d but got %d", 1, cluster.PodNumber)
	}
}
func TestDeletePod(t *testing.T) {
	cluster := New()
	var n Node
	cluster.AddNode("vamos", &n)
	var p Pod
	p.Name = "mongo"
	p.Image = "mongo"
	p.Ports = []string{"mongo"}
	cluster.DeletePod(defaultNameSpace, p.Name)
	if cluster.PodNumber != 0 {
		t.Errorf("expected %d but got %d", 0, cluster.PodNumber)
	}
}
