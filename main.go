package main

import (
    "github.com/msoulier/go-taskwarrior"
	"slices"
	"fmt"
	"cmp"
)

var (
	TASKRC = taskwarrior.PathExpandTilda("~/.taskrc")
	desclen int = 40
)

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
	// Sort by urgency
	slices.SortFunc(tw.Tasks, func(a, b taskwarrior.Task) int {
		return cmp.Compare[float32](a.Urgency, b.Urgency)
	})
	slices.Reverse(tw.Tasks)
	//tw.PrintTasks()
	sformat := "| %40s | %10s | %4s | %7s |"
	fformat := "| %40s | %10s | %4s | %7.2f |"
	fmt.Printf(sformat + "\n", "Description", "Status", "Pri", "Urg")
	fmt.Printf(sformat + "\n", "----------------------------------------", "----------", "----", "-------")
	for _, task := range tw.Tasks {
		if task.Status != "completed" && task.Status != "deleted" {
			description := task.Description
			if len(description) > desclen {
				description = description[:desclen]
			}
			fmt.Printf(fformat + "\n", description, task.Status, task.Priority, task.Urgency)
		}
	}
}
