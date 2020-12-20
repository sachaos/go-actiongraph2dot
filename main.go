package main

import (
	"encoding/json"
	"fmt"
	"github.com/emicklei/dot"
	"os"
	"time"
)

// NOTE: Copy from https://go.googlesource.com/go/+/refs/tags/go1.16beta1/src/cmd/go/internal/work/action.go#142
type actionJSON struct {
	ID         int
	Mode       string
	Package    string
	Deps       []int     `json:",omitempty"`
	IgnoreFail bool      `json:",omitempty"`
	Args       []string  `json:",omitempty"`
	Link       bool      `json:",omitempty"`
	Objdir     string    `json:",omitempty"`
	Target     string    `json:",omitempty"`
	Priority   int       `json:",omitempty"`
	Failed     bool      `json:",omitempty"`
	Built      string    `json:",omitempty"`
	VetxOnly   bool      `json:",omitempty"`
	NeedVet    bool      `json:",omitempty"`
	NeedBuild  bool      `json:",omitempty"`
	ActionID   string    `json:",omitempty"`
	BuildID    string    `json:",omitempty"`
	TimeReady  time.Time `json:",omitempty"`
	TimeStart  time.Time `json:",omitempty"`
	TimeDone   time.Time `json:",omitempty"`

	Cmd     []string      // `json:",omitempty"`
	CmdReal time.Duration `json:",omitempty"`
	CmdUser time.Duration `json:",omitempty"`
	CmdSys  time.Duration `json:",omitempty"`
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	var actions []*actionJSON
	if err := json.NewDecoder(os.Stdin).Decode(&actions); err != nil {
		return err
	}

	g := dot.NewGraph(dot.Directed)
	nodes := make([]dot.Node, len(actions))
	for i, action := range actions {
		nodes[i] = g.Node(string(action.ID))
		nodes[i].Label(fmt.Sprintf("%s: %s", action.Mode, action.Package))
	}

	for _, action := range actions {
		for _, id := range action.Deps {
			g.Edge(nodes[id], nodes[action.ID])
		}
	}

	fmt.Println(g.String())

	return nil
}
