package services

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSyncLatestDir(t *testing.T) {
	baseDir := t.TempDir()
	runDir := filepath.Join(baseDir, "2026-06-11_10-00")
	if err := os.MkdirAll(runDir, 0755); err != nil {
		t.Fatal(err)
	}

	writeOutput := func(name, content string) string {
		path := filepath.Join(runDir, name)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
		return path
	}

	results := []GenerationResult{
		{Template: "modern-html", OutputFormat: OutputFormatHTML, OutputPath: writeOutput("Jane_Doe.modern-html.html", "<html>")},
		// Uniqueness suffix from a same-minute rerun should be dropped in latest.
		{Template: "modern-html", OutputFormat: OutputFormatPDF, OutputPath: writeOutput("Jane_Doe.modern-html_2.pdf", "%PDF")},
	}

	// Seed a stale file from a previous run that must be cleared.
	latestDir := filepath.Join(baseDir, LatestDirName)
	if err := os.MkdirAll(latestDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(latestDir, "Jane_Doe.old-template.pdf"), []byte("stale"), 0644); err != nil {
		t.Fatal(err)
	}

	if err := SyncLatestDir(runDir, "Jane_Doe", results); err != nil {
		t.Fatalf("SyncLatestDir failed: %v", err)
	}

	entries, err := os.ReadDir(latestDir)
	if err != nil {
		t.Fatal(err)
	}
	got := make(map[string]bool)
	for _, e := range entries {
		got[e.Name()] = true
	}

	want := []string{"Jane_Doe.modern-html.html", "Jane_Doe.modern-html.pdf"}
	if len(got) != len(want) {
		t.Errorf("latest dir has %d entries, want %d: %v", len(got), len(want), got)
	}
	for _, name := range want {
		if !got[name] {
			t.Errorf("latest dir missing %s", name)
		}
	}
	if got["Jane_Doe.old-template.pdf"] {
		t.Error("stale file from previous run was not cleared")
	}

	for i, r := range results {
		if r.LatestPath == "" {
			t.Errorf("results[%d].LatestPath not set", i)
			continue
		}
		data, err := os.ReadFile(r.LatestPath)
		if err != nil {
			t.Errorf("results[%d] latest copy unreadable: %v", i, err)
			continue
		}
		orig, _ := os.ReadFile(r.OutputPath)
		if string(data) != string(orig) {
			t.Errorf("results[%d] latest copy content differs from original", i)
		}
	}
}

func TestSyncLatestDirEmptyBase(t *testing.T) {
	baseDir := t.TempDir()
	runDir := filepath.Join(baseDir, "2026-06-11_10-00")
	if err := os.MkdirAll(runDir, 0755); err != nil {
		t.Fatal(err)
	}
	src := filepath.Join(runDir, "out.modern-html.pdf")
	if err := os.WriteFile(src, []byte("%PDF"), 0644); err != nil {
		t.Fatal(err)
	}

	results := []GenerationResult{{Template: "modern-html", OutputFormat: OutputFormatPDF, OutputPath: src}}
	if err := SyncLatestDir(runDir, "", results); err != nil {
		t.Fatalf("SyncLatestDir failed: %v", err)
	}

	want := filepath.Join(baseDir, LatestDirName, "Resume.modern-html.pdf")
	if results[0].LatestPath != want {
		t.Errorf("LatestPath = %s, want %s", results[0].LatestPath, want)
	}
	if _, err := os.Stat(want); err != nil {
		t.Errorf("expected latest copy at %s: %v", want, err)
	}
}
