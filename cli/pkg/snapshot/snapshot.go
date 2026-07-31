package snapshot

import (
	"bytes"
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

	index, err := os.CreateTemp("", "stuntdouble-index-*")
	if err != nil {
		return err
	}
	indexFile := index.Name()
	index.Close()
	defer os.Remove(indexFile)

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
	stuntDir := filepath.Join(workspace, ".stuntdouble")
	os.MkdirAll(stuntDir, 0755)
	snapshotFile := filepath.Join(stuntDir, "latest_snapshot.txt")
	if err := os.WriteFile(snapshotFile, []byte(treeHash), 0644); err != nil {
		return err
	}

	fmt.Printf("📸 [StuntDouble] Zero-copy workspace snapshot captured (Tree: %s)\n", treeHash[:8])
	return nil
}

// SaveCheckpoint creates a named workspace snapshot saved in .stuntdouble/checkpoints/<name>.txt
func SaveCheckpoint(workspace, name string) error {
	if name == "" {
		return fmt.Errorf("checkpoint name cannot be empty")
	}

	treeHash, err := captureTreeHash(workspace)
	if err != nil {
		return fmt.Errorf("failed to capture workspace tree: %w", err)
	}

	checkpointDir := filepath.Join(workspace, ".stuntdouble", "checkpoints")
	if err := os.MkdirAll(checkpointDir, 0755); err != nil {
		return err
	}

	checkpointFile := filepath.Join(checkpointDir, name+".txt")
	if err := os.WriteFile(checkpointFile, []byte(treeHash), 0644); err != nil {
		return err
	}

	fmt.Printf("🔖 [StuntDouble] Checkpoint %q saved (Tree: %s)\n", name, treeHash[:8])
	return nil
}

// RestoreCheckpoint rewinds workspace to a previously saved named checkpoint
func RestoreCheckpoint(workspace, name string) error {
	checkpointFile := filepath.Join(workspace, ".stuntdouble", "checkpoints", name+".txt")
	data, err := os.ReadFile(checkpointFile)
	if err != nil {
		return fmt.Errorf("checkpoint %q not found: %w", name, err)
	}

	treeHash := strings.TrimSpace(string(data))
	if err := restoreTreeHash(workspace, treeHash); err != nil {
		return fmt.Errorf("failed to restore checkpoint %q: %w", name, err)
	}

	fmt.Printf("⏪ [StuntDouble] Workspace rewound to checkpoint %q (Tree: %s)\n", name, treeHash[:8])
	return nil
}

// ListCheckpoints returns all saved checkpoint names
func ListCheckpoints(workspace string) ([]string, error) {
	checkpointDir := filepath.Join(workspace, ".stuntdouble", "checkpoints")
	entries, err := os.ReadDir(checkpointDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}

	var names []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".txt") {
			names = append(names, strings.TrimSuffix(entry.Name(), ".txt"))
		}
	}
	return names, nil
}

func captureTreeHash(workspace string) (string, error) {
	if _, err := os.Stat(filepath.Join(workspace, ".git")); os.IsNotExist(err) {
		return "", fmt.Errorf("not a git repository")
	}

	index, err := os.CreateTemp("", "stuntdouble-index-*")
	if err != nil {
		return "", err
	}
	indexFile := index.Name()
	index.Close()
	defer os.Remove(indexFile)

	cmd1 := exec.Command("git", "read-tree", "HEAD")
	cmd1.Env = append(os.Environ(), "GIT_INDEX_FILE="+indexFile)
	cmd1.Dir = workspace
	if err := cmd1.Run(); err != nil {
		return "", err
	}

	cmd2 := exec.Command("git", "add", "-A")
	cmd2.Env = append(os.Environ(), "GIT_INDEX_FILE="+indexFile)
	cmd2.Dir = workspace
	if err := cmd2.Run(); err != nil {
		return "", err
	}

	cmd3 := exec.Command("git", "write-tree")
	cmd3.Env = append(os.Environ(), "GIT_INDEX_FILE="+indexFile)
	cmd3.Dir = workspace
	out, err := cmd3.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}

func restoreTreeHash(workspace, treeHash string) error {
	cmd1 := exec.Command("git", "restore", "--source", treeHash, "--worktree", ".")
	cmd1.Dir = workspace
	if err := cmd1.Run(); err != nil {
		return fmt.Errorf("failed to restore snapshot: %w", err)
	}

	return removeUntrackedAbsentFromSnapshot(workspace, treeHash)
}

// Restore rewinds the workspace to the latest captured snapshot
func Restore(workspace string) error {
	snapshotFile := filepath.Join(workspace, ".stuntdouble", "latest_snapshot.txt")
	data, err := os.ReadFile(snapshotFile)
	if err != nil {
		return fmt.Errorf("no StuntDouble snapshot found: %w", err)
	}

	treeHash := strings.TrimSpace(string(data))
	return restoreTreeHash(workspace, treeHash)
}

func removeUntrackedAbsentFromSnapshot(workspace, treeHash string) error {
	snapshotFiles, err := filesInTree(workspace, treeHash)
	if err != nil {
		return err
	}

	untracked, err := untrackedFiles(workspace)
	if err != nil {
		return err
	}

	for _, rel := range untracked {
		if _, ok := snapshotFiles[rel]; ok {
			continue
		}

		cleanRel := filepath.Clean(rel)
		if cleanRel == "." || filepath.IsAbs(cleanRel) || strings.HasPrefix(cleanRel, ".."+string(os.PathSeparator)) || cleanRel == ".." {
			return fmt.Errorf("refusing to remove unsafe path %q", rel)
		}
		if err := os.RemoveAll(filepath.Join(workspace, cleanRel)); err != nil {
			return err
		}
	}
	return nil
}

func filesInTree(workspace, treeHash string) (map[string]struct{}, error) {
	cmd := exec.Command("git", "ls-tree", "-r", "-z", "--name-only", treeHash)
	cmd.Dir = workspace
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	files := make(map[string]struct{})
	for _, part := range bytes.Split(out, []byte{0}) {
		if len(part) == 0 {
			continue
		}
		files[string(part)] = struct{}{}
	}
	return files, nil
}

func untrackedFiles(workspace string) ([]string, error) {
	cmd := exec.Command("git", "ls-files", "--others", "--exclude-standard", "-z")
	cmd.Dir = workspace
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var files []string
	for _, part := range bytes.Split(out, []byte{0}) {
		if len(part) == 0 {
			continue
		}
		files = append(files, string(part))
	}
	return files, nil
}
