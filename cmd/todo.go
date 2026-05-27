package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/acemaster-gh/gocode/internal/ui"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var todoCmd = &cobra.Command{
	Use:   "todo",
	Short: "Manage per-project todos",
	RunE:  runTodo,
}

func runTodo(cmd *cobra.Command, args []string) error {
	ui.Banner(Version)

	reg := getRegistry()

	name, err := pickActiveProject("Manage todos for which project?")
	if err != nil || name == "" {
		return err
	}

	proj, ok := reg.Get(name)
	if !ok {
		return fmt.Errorf("project '%s' not found", name)
	}

	for {
		// Render current list
		fmt.Println()
		fmt.Printf("  %s\n", ui.Bold.Render("Todos: "+name))
		ui.Divider()

		if len(proj.Todos) == 0 {
			fmt.Println(ui.Dim.Render("  No todos yet."))
		} else {
			done, total := 0, len(proj.Todos)
			for i, t := range proj.Todos {
				mark := ui.Dim.Render("[ ]")
				task := t.Task
				if t.Done {
					mark = ui.Green.Render("[✓]")
					task = ui.Dim.Render(task)
					done++
				}
				fmt.Printf("  %s %s %s\n",
					ui.Dim.Render(fmt.Sprintf("%2d.", i+1)),
					mark, task)
			}
			fmt.Printf("\n  %s\n",
				ui.Dim.Render(fmt.Sprintf("%d/%d complete", done, total)))
		}
		ui.Divider()

		var action string
		opts := []huh.Option[string]{
			huh.NewOption("+ Add todo", "add"),
		}
		if len(proj.Todos) > 0 {
			opts = append(opts,
				huh.NewOption("✓ Toggle done/undone", "toggle"),
				huh.NewOption("✕ Done — exit", "exit"),
			)
		} else {
			opts = append(opts, huh.NewOption("✕ Exit", "exit"))
		}

		err = huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Action").
					Options(opts...).
					Value(&action),
			),
		).WithTheme(huh.ThemeCharm()).Run()

		if err != nil || action == "exit" {
			return nil
		}

		switch action {
		case "add":
			var task string
			addErr := huh.NewForm(
				huh.NewGroup(
					huh.NewInput().
						Title("New todo").
						Placeholder("implement form validation").
						Value(&task).
						Validate(func(s string) error {
							if strings.TrimSpace(s) == "" {
								return errors.New("task cannot be empty")
							}
							return nil
						}),
				),
			).WithTheme(huh.ThemeCharm()).Run()

			if addErr == nil && task != "" {
				reg.AddTodo(name, strings.TrimSpace(task))
				ui.OK("Todo added.")
			}

		case "toggle":
			if len(proj.Todos) == 0 {
				continue
			}
			opts := make([]huh.Option[int], len(proj.Todos))
			for i, t := range proj.Todos {
				label := t.Task
				if t.Done {
					label = "✓ " + t.Task
				}
				opts[i] = huh.NewOption(label, i)
			}
			var idx int
			toggleErr := huh.NewForm(
				huh.NewGroup(
					huh.NewSelect[int]().
						Title("Toggle which todo?").
						Options(opts...).
						Value(&idx),
				),
			).WithTheme(huh.ThemeCharm()).Run()

			if toggleErr == nil {
				reg.ToggleTodo(name, idx)
			}
		}

		// Reload project for updated todos
		proj, _ = reg.Get(name)
	}
}
