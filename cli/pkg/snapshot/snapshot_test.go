package snapshot

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestCreateAndRestorePreservesSnapshotUntrackedFiles(t *testing.T) {
	workspace := newGitWorkspace(t)

	writeFile(t, workspace, "tracked.txt", "tracked at commit\n")
	runGit(t, workspace, "add", "tracked.txt")
	runGit(t, workspace, "-c", "user.name=StuntDouble", "-c", "user.email=stuntdouble@example.test", "commit", "-m", "initial")

	writeFile(t, workspace, "tracked.txt", "tracked at snapshot\n")
	writeFile(t, workspace, "scratch.txt", "untracked at snapshot\n")
	writeFile(t, workspace, "nested/keep.txt", "nested untracked at snapshot\n")

	if err := Create(workspace); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	writeFile(t, workspace, "tracked.txt", "tracked after agent\n")
	writeFile(t, workspace, "scratch.txt", "untracked after agent\n")
	writeFile(t, workspace, "nested/keep.txt", "nested after agent\n")
	writeFile(t, workspace, "new.txt", "created after snapshot\n")
	writeFile(t, workspace, "nested/new.txt", "nested created after snapshot\n")

	if err := Restore(workspace); err != nil {
		t.Fatalf("Restore() error = %v", err)
	}

	assertFile(t, workspace, "tracked.txt", "tracked at snapshot\n")
	assertFile(t, workspace, "scratch.txt", "untracked at snapshot\n")
	assertFile(t, workspace, "nested/keep.txt", "nested untracked at snapshot\n")
	assertMissing(t, workspace, "new.txt")
	assertMissing(t, workspace, "nested/new.txt")
	assertMissing(t, workspace, ".stuntdouble/stunt.index")
}

func TestCreateNoopsOutsideGitRepository(t *testing.T) {
	workspace := t.TempDir()
	if err := Create(workspace); err != nil {
		t.Fatalf("Create() outside git repo error = %v", err)
	}
	assertMissing(t, workspace, ".stuntdouble/latest_snapshot.txt")
}

func TestRestoreWithoutSnapshotReturnsError(t *testing.T) {
	workspace := newGitWorkspace(t)
	if err := Restore(workspace); err == nil || !strings.Contains(err.Error(), "no StuntDouble snapshot found") {
		t.Fatalf("Restore() error = %v, want no snapshot error", err)
	}
}

func newGitWorkspace(t *testing.T) string {
	t.Helper()
	workspace := t.TempDir()
	runGit(t, workspace, "init")
	return workspace
}

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s failed: %v\n%s", strings.Join(args, " "), err, out)
	}
}

func writeFile(t *testing.T, workspace, rel, content string) {
	t.Helper()
	path := filepath.Join(workspace, rel)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("MkdirAll(%q): %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("WriteFile(%q): %v", path, err)
	}
}

func assertFile(t *testing.T, workspace, rel, want string) {
	t.Helper()
	path := filepath.Join(workspace, rel)
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile(%q): %v", path, err)
	}
	if string(got) != want {
		t.Fatalf("%s = %q, want %q", rel, got, want)
	}
}

func assertMissing(t *testing.T, workspace, rel string) {
	t.Helper()
	path := filepath.Join(workspace, rel)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("%s exists or stat failed with non-ENOENT error: %v", rel, err)
	}
}
