package tui

import (
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/textfuel/lazyjira/v2/pkg/git"
	"github.com/textfuel/lazyjira/v2/pkg/jira"
)

// SetStartIssue opens the given issue once the app starts.
func (a *App) SetStartIssue(key string) {
	a.startIssue = key
}

// resolveGotoKey turns "PROJ-123", "proj-123", a /browse/ URL, or a bare
// number (prefixed with the active project) into an issue key.
func resolveGotoKey(input, projectKey string) string {
	input = strings.TrimSpace(input)
	if _, err := strconv.Atoi(input); err == nil && projectKey != "" {
		return projectKey + "-" + input
	}
	return git.ExtractIssueKey(input)
}

// gotoIssue selects the issue if it's in the current list, otherwise loads
// it into the JQL tab so tab refreshes keep it. Focuses Info to show its fields.
func (a *App) gotoIssue(input string) tea.Cmd {
	key := resolveGotoKey(input, a.projectKey)
	if key == "" {
		a.statusPanel.SetError("not an issue key: " + strings.TrimSpace(input))
		return nil
	}
	if a.issuesList.SelectByKey(key) {
		a.side = sideLeft
		a.leftFocus = focusInfo
		a.updateFocusState()
		return tea.Batch(a.previewSelectedIssue(), fetchIssueDetail(a.client, key))
	}
	return fetchGotoIssue(a.client, key)
}

func (a *App) handleGotoIssueLoaded(msg gotoIssueLoadedMsg) (tea.Model, tea.Cmd) {
	a.statusPanel.SetError("")
	a.issueCache[msg.issue.Key] = msg.issue
	a.issuesList.AddJQLTab("key = " + msg.issue.Key)
	a.issuesList.SetIssues([]jira.Issue{*msg.issue})
	a.issuesList.SelectByKey(msg.issue.Key)
	a.side = sideLeft
	a.leftFocus = focusInfo
	a.updateFocusState()
	return a, a.previewSelectedIssue()
}
