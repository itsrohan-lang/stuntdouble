package snapshot

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Create takes a virtually instantaneous snapshot of the workspace
// using Git's low-level plumbing commands, ensuring zero pollution to
// the user's branch history or staging area.
func Create(workspace string) error {
	// Check if it's a git repo
	if _, err := os.Stat(filepath.Join(workspace, ".git")); os.IsNotExist(err) {
		fmt.Println("⚠️ [StuntDouble] Not a git repository. Skipping zero-copy snapshot.")
		return nil
	}

	stuntDir := filepath.Join(workspace, ".stuntdouble")
	os.MkdirAll(stuntDir, 0755)

	indexFile := filepath.Join(stuntDir, "stunt.index")
	
	// Create a temporary index loaded with HEAD
	cmd1 := exec.Command("git", "read-tree", "HEAD")
	cmd1.Env = append(os.Environ(), "GIT_INDEX_FILE="+indexFile)
	cmd1.Dir = workspace
	if err := cmd1.Run(); err != nil {
		return err
	}

	// Add all current working directory files to the temporary index
	cmd2 := exec.Command("git", "add", "-A")
	cmd2.Env = append(os.Environ(), "GIT_INDEX_FILE="+indexFile)
	cmd2.Dir = workspace
	if err := cmd2.Run(); err != nil {
		return err
	}

	// Write the index to a tree object in the git database
	cmd3 := exec.Command("git", "write-tree")
	cmd3.Env = append(os.Environ(), "GIT_INDEX_FILE="+indexFile)
	cmd3.Dir = workspace
	out, err := cmd3.Output()
	if err != nil {
		return err
	}

	treeHash := strings.TrimSpace(string(out))
	
	// Save the tree hash
	snapshotFile := filepath.Join(stuntDir, "latest_snapshot.txt")
	if err := os.WriteFile(snapshotFile, []byte(treeHash), 0644); err != nil {
		return err
	}

	fmt.Printf("📸 [StuntDouble] Zero-copy workspace snapshot captured (Tree: %s)\n", treeHash[:8])
	return nil
}

// Restore rewinds the workspace to the latest captured snapshot
func Restore(workspace string) error {
	snapshotFile := filepath.Join(workspace, ".stuntdouble", "latest_snapshot.txt")
	data, err := os.ReadFile(snapshotFile)
	if err != nil {
		return fmt.Errorf("no StuntDouble snapshot found: %w", err)
	}

	treeHash := strings.TrimSpace(string(data))
	
	// Restore tracked files
	cmd1 := exec.Command("git", "restore", "--source", treeHash, "--worktree", ".")
	cmd1.Dir = workspace
	if err := cmd1.Run(); err != nil {
		return fmt.Errorf("failed to restore snapshot: %w", err)
	}

	// Clean untracked files created by the agent
	cmd2 := exec.Command("git", "clean", "-fd")
	cmd2.Dir = workspace
	if err := cmd2.Run(); err != nil {
		return fmt.Errorf("failed to clean untracked files: %w", err)
	}

	return nil
}
