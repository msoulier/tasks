package main

import (
    "github.com/msoulier/go-taskwarrior"
	"slices"
	"fmt"
	"cmp"
    "strings"
    "flag"
    "os"
)

var (
	TASKRC = taskwarrior.PathExpandTilda("~/.taskrc")
	desclen int = 40
    debug bool = false
    markdown bool = false
    depgraph bool = false
)

func init() {
    flag.BoolVar(&debug, "debug", false, "Debug logging")
    flag.BoolVar(&markdown, "markdown", false, "Output a markdown table (default)")
    flag.BoolVar(&depgraph, "depgraph", false, "Output a dependency graph in graphviz dot format")
    flag.Parse()
    if !markdown && !depgraph {
        markdown = true
    }
}

/*
type Task struct {
	Description string  `json:"description"`
	Project     string  `json:"project,omitempty"`
	Status      string  `json:"status,omitempty"`
	Uuid        string  `json:"uuid,omitempty"`
	Urgency     float32 `json:"urgency,omitempty"`
	Priority    string  `json:"priority,omitempty"`
	Due         string  `json:"due,omitempty"`
	End         string  `json:"end,omitempty"`
	Entry       string  `json:"entry,omitempty"`
	Modified    string  `json:"modified,omitempty"`
}
*/

func main() {
	tw, err := taskwarrior.NewTaskWarrior(TASKRC)
	if err != nil {
		panic(err)
	}
	if err = tw.FetchAllTasks(); err != nil {
		panic(err)
	}
    tasks_by_uuid := make(map[string]taskwarrior.Task)
    for _, task := range tw.Tasks {
        tasks_by_uuid[task.Uuid] = task
    }
	// Sort by urgency
	slices.SortFunc(tw.Tasks, func(a, b taskwarrior.Task) int {
		return cmp.Compare[float32](a.Urgency, b.Urgency)
	})
	slices.Reverse(tw.Tasks)
	//tw.PrintTasks()
    if markdown {
        sformat := "| %40s | %10s | %4s | %7s | %20s |"
        fformat := "| %40s | %10s | %4s | %7.2f | %20s |"
        fmt.Printf(sformat + "\n", "Description", "Status", "Pri", "Urg", "Deps")
        fmt.Printf(sformat + "\n", "----------------------------------------", "----------", "----", "-------", "--------------------")
        for _, task := range tw.Tasks {
            if task.Status != "completed" && task.Status != "deleted" {
                description := task.Description
                if len(description) > desclen {
                    description = description[:desclen]
                }
                depends := ""
                if len(task.Depends) > 0 {
                    depends = strings.Join(task.Depends, ",")
                }
                fmt.Printf(fformat + "\n", description, task.Status, task.Priority, task.Urgency, depends)
            }
        }
    } else if depgraph {
        fmt.Println("digraph Tasks {")
        objects := make([]string, 0)
        edges := make([]string, 0)
        for _, task := range tw.Tasks {
            if task.Status != "completed" && task.Status != "deleted" {
                if len(task.Depends) > 0 {
                    objects = append(objects, fmt.Sprintf("\"%s\" [shape=circle,label=\"%s\"]", task.Uuid, task.Description))
                    for _, uuid := range task.Depends {
                        if _, present := tasks_by_uuid[uuid]; present {
                            deptask := tasks_by_uuid[uuid]
                            objects = append(objects, fmt.Sprintf("\"%s\" [shape=circle,label=\"%s\"]", deptask.Uuid, deptask.Description))
                            edges = append(edges, fmt.Sprintf("\"%s\" -> \"%s\"", task.Uuid, uuid))
                        } else {
                            fmt.Fprintf(os.Stderr, "Can't find uuid %s\n", uuid)
                            panic("can't find uuid")
                        }
                    }
                }
            }
        }
        for _, obj := range objects {
            fmt.Printf("%s\n", obj)
        }
        for _, edge := range edges {
            fmt.Printf("%s\n", edge)
        }
        fmt.Println("}")
    }
}
