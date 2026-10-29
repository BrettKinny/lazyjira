package tui

import (
	"testing"

	"github.com/textfuel/lazyjira/v2/pkg/internal/testkit"
	"github.com/textfuel/lazyjira/v2/pkg/jira"
	"github.com/textfuel/lazyjira/v2/pkg/jira/jiratest"
)

func TestResolveGotoKey(t *testing.T) {
	t.Parallel()
	cases := []struct {
		input, project, want string
	}{
		{"PLAT-1", testProject, "PLAT-1"},
		{" plat-1 ", testProject, "PLAT-1"},
		{"42", testProject, "PLAT-42"},
		{"42", "", ""},
		{"https://example.atlassian.net/browse/PLAT-7?focusedCommentId=1", testProject, "PLAT-7"},
		{"nonsense", testProject, ""},
		{"", testProject, ""},
	}
	for _, c := range cases {
		testkit.AssertEqual(t, c.input, resolveGotoKey(c.input, c.project), c.want)
	}
}

func TestGotoIssue(t *testing.T) {
	t.Parallel()

	t.Run("invalid input sets error and no command", func(t *testing.T) {
		t.Parallel()
		app := newAppWithFake(t, &jiratest.FakeClient{T: t})

		if cmd := app.gotoIssue("nonsense"); cmd != nil {
			t.Error("expected no command for invalid key")
		}
		if app.statusPanel.ErrorMessage() == "" {
			t.Error("expected error for invalid key")
		}
	})

	t.Run("bare start number waits for project", func(t *testing.T) {
		t.Parallel()
		app := newAppWithFake(t, &jiratest.FakeClient{T: t})
		app.demoMode = true
		app.SetStartIssue("42")

		_ = app.Init()
		testkit.AssertEqual(t, "startIssue before projects", app.startIssue, "42")

		_, _ = app.handleProjectsLoaded(projectsLoadedMsg{projects: []jira.Project{{Key: testProject}}})
		testkit.AssertEqual(t, "startIssue after projects", app.startIssue, "")
		if app.statusPanel.ErrorMessage() != "" {
			t.Errorf("unexpected error: %s", app.statusPanel.ErrorMessage())
		}
	})

	t.Run("loaded issue opens in JQL tab", func(t *testing.T) {
		t.Parallel()
		app := newAppWithFake(t, &jiratest.FakeClient{T: t})

		_, _ = app.handleGotoIssueLoaded(gotoIssueLoadedMsg{issue: &jira.Issue{Key: testKey}})

		if !app.issuesList.IsJQLTab() {
			t.Error("expected JQL tab to be active")
		}
		testkit.AssertEqual(t, "jql", app.issuesList.JQLQuery(), "key = "+testKey)
		if sel := app.issuesList.SelectedIssue(); sel == nil || sel.Key != testKey {
			t.Errorf("selected issue = %v, want %s", sel, testKey)
		}
		testkit.AssertEqual(t, "previewKey", app.previewKey, testKey)
		testkit.AssertEqual(t, "leftFocus", app.leftFocus, focusInfo)
	})
}
