// Package ui
package ui

import (
	"fmt"
	"os"

	"lazytrees/internal/git"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type AppModel struct {
	projectName string
	worktrees   []git.Worktree

	width  int
	height int
	cursor int
}

func Bootstrap() AppModel {
	projectName, err := git.GetProjectName()
	app := AppModel{projectName: projectName}
	if err != nil {
		fmt.Println("FATAL: lazytrees must be called inside a repository")
		os.Exit(1)
	}

	trees, err := git.GetWorktrees()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	app.worktrees = trees

	return app

	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println("---")
	// for i := range worktrees {
	// 	tree := worktrees[i]
	// 	fmt.Println(tree.Alias, tree.Path)
	// }
	// fmt.Println("---")
	// worktree, createErr := git.CreateWorktree()
	// if createErr != nil {
	// 	fmt.Println("Error: creating", worktree.Alias, "\n", createErr)
	// 	return
	// }
	//
	// fmt.Println("Created new worktree:", worktree.Alias)
	//
	// branchErr := worktree.CheckoutBranch("test-branch")
	//
	// if branchErr != nil {
	// 	fmt.Println(branchErr)
	// }
}

func (m AppModel) Init() tea.Cmd {
	return nil
}

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "j", "down":
			if m.cursor >= 0 && m.cursor < len(m.worktrees)-1 {
				m.cursor++
			}

		case "k", "up":
			if m.cursor > 0 && m.cursor <= len(m.worktrees)-1 {
				m.cursor--
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	}
	return m, nil
}

func (m AppModel) View() tea.View {
	leftWidth := m.width / 3
	rightWidth := m.width - leftWidth

	panelHeight := max(1, m.height-3)

	left := PanelStyle.
		Width(max(1, leftWidth-4)).
		Height(panelHeight).
		Render(m.worktreeView())

	right := PanelStyle.
		Width(max(1, rightWidth-4)).
		Height(panelHeight).
		Render(m.detailsView())

	body := lipgloss.JoinHorizontal(
		lipgloss.Top,
		left,
		right,
	)

	view := tea.NewView(body)
	view.AltScreen = true

	return view
}

func (m AppModel) worktreeView() string {
	title := lipgloss.NewStyle().
		Bold(true).
		Render("Worktrees")

	content := title + "\n\n"

	for i, worktree := range m.worktrees {
		cursor := "  "

		if i == m.cursor {
			cursor = "> "
		}

		content += cursor + worktree.Alias + "\n"
	}

	return content
}

func (m AppModel) detailsView() string {
	title := lipgloss.NewStyle().
		Bold(true).
		Render("Details")

	if len(m.worktrees) == 0 {
		return title + "\n\nNo worktrees"
	}

	selected := m.worktrees[m.cursor]
	branch := selected.Branch
	if selected.Detached {
		branch = "DETACHED"
	}
	return fmt.Sprintf(
		"%s\n\nAlias: %s\nBranch: %s",
		title,
		selected.Alias,
		branch,
	)
}
