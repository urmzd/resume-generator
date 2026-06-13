package services

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSyncLatestDir(t *testing.T) {
	baseDir := t.TempDir()
	runDir := filepath.Join(baseDir, "2026-06-11_10-00")
	if err := os.MkdirAll(filepath.Join(runDir, "modern-html", ArtifactsDirName), 0755); err != nil {
		t.Fatal(err)
	}

	writeOutput := func(relPath, content string) string {
		path := filepath.Join(runDir, relPath)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
		return path
	}

	results := []GenerationResult{
		{Template: "modern-html", OutputFormat: OutputFormatHTML, Artifact: true,
			OutputPath: writeOutput(filepath.Join("modern-html", ArtifactsDirName, "Jane_Doe.html"), "<html>")},
		// Uniqueness suffix from a same-minute rerun should be dropped in latest.
		{Template: "modern-html", OutputFormat: OutputFormatPDF,
			OutputPath: writeOutput(filepath.Join("modern-html", "Jane_Doe_2.pdf"), "%PDF")},
	}

	// Seed a stale file from a previous run that must be cleared.
	latestDir := filepath.Join(baseDir, LatestDirName)
	if err := os.MkdirAll(filepath.Join(latestDir, "old-template"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(latestDir, "old-template", "Jane_Doe.pdf"), []byte("stale"), 0644); err != nil {
		t.Fatal(err)
	}

	if err := SyncLatestDir(runDir, "Jane_Doe", results); err != nil {
		t.Fatalf("SyncLatestDir failed: %v", err)
	}

	wantPaths := []string{
		filepath.Join(latestDir, "modern-html", ArtifactsDirName, "Jane_Doe.html"),
		filepath.Join(latestDir, "modern-html", "Jane_Doe.pdf"),
	}
	for i, want := range wantPaths {
		if results[i].LatestPath != want {
			t.Errorf("results[%d].LatestPath = %s, want %s", i, results[i].LatestPath, want)
		}
		if _, err := os.Stat(want); err != nil {
			t.Errorf("expected latest copy at %s: %v", want, err)
		}
	}
	if _, err := os.Stat(filepath.Join(latestDir, "old-template")); !os.IsNotExist(err) {
		t.Error("stale entries from previous run were not cleared")
	}

	for i, r := range results {
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
	if err := os.MkdirAll(filepath.Join(runDir, "modern-html"), 0755); err != nil {
		t.Fatal(err)
	}
	src := filepath.Join(runDir, "modern-html", "out.pdf")
	if err := os.WriteFile(src, []byte("%PDF"), 0644); err != nil {
		t.Fatal(err)
	}

	results := []GenerationResult{{Template: "modern-html", OutputFormat: OutputFormatPDF, OutputPath: src}}
	if err := SyncLatestDir(runDir, "", results); err != nil {
		t.Fatalf("SyncLatestDir failed: %v", err)
	}

	want := filepath.Join(baseDir, LatestDirName, "modern-html", "Resume.pdf")
	if results[0].LatestPath != want {
		t.Errorf("LatestPath = %s, want %s", results[0].LatestPath, want)
	}
	if _, err := os.Stat(want); err != nil {
		t.Errorf("expected latest copy at %s: %v", want, err)
	}
}

func TestOutputPathFor(t *testing.T) {
	runDir := t.TempDir()

	deliverable, err := OutputPathFor(runDir, "modern-html", "Jane_Doe", ".pdf", false)
	if err != nil {
		t.Fatal(err)
	}
	if want := filepath.Join(runDir, "modern-html", "Jane_Doe.pdf"); deliverable != want {
		t.Errorf("deliverable path = %s, want %s", deliverable, want)
	}

	artifact, err := OutputPathFor(runDir, "modern-html", "Jane_Doe", ".html", true)
	if err != nil {
		t.Fatal(err)
	}
	if want := filepath.Join(runDir, "modern-html", ArtifactsDirName, "Jane_Doe.html"); artifact != want {
		t.Errorf("artifact path = %s, want %s", artifact, want)
	}

	// Both parent directories must exist after the call.
	for _, p := range []string{deliverable, artifact} {
		if _, err := os.Stat(filepath.Dir(p)); err != nil {
			t.Errorf("directory for %s not created: %v", p, err)
		}
	}

	// A second call with the file present must add a uniqueness suffix.
	if err := os.WriteFile(deliverable, []byte("%PDF"), 0644); err != nil {
		t.Fatal(err)
	}
	second, err := OutputPathFor(runDir, "modern-html", "Jane_Doe", ".pdf", false)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(second, "Jane_Doe_2.pdf") {
		t.Errorf("rerun path = %s, want a Jane_Doe_2.pdf suffix", second)
	}
}
