package net

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"os"
)

// DownloadFile Download file from ssh.Session
func DownloadFile(session *ssh.Session, remotePath, localPath string) error {
	dstFile, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("failed to create local file: %w", err)
	}
	defer dstFile.Close()

	stdout, err := session.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout pipe: %w", err)
	}
	stderr, err := session.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to get stderr pipe: %w", err)
	}

	go io.Copy(os.Stderr, stderr)

	if err := session.Start("cat " + remotePath); err != nil {
		return fmt.Errorf("failed to start remote command: %w", err)
	}

	if _, err := io.Copy(dstFile, stdout); err != nil {
		return fmt.Errorf("failed to copy data: %w", err)
	}

	if err := session.Wait(); err != nil {
		return fmt.Errorf("failed to wait for remote command: %w", err)
	}

	return nil
}

// UploadFile Upload file to ssh.Session
func UploadFile(session *ssh.Session, localPath string, remotePath string) error {
	srcFile, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("failed to open local file: %w", err)
	}
	defer srcFile.Close()

	dstFile, err := session.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdin pipe: %w", err)
	}

	go func() {
		defer dstFile.Close()
		io.Copy(dstFile, srcFile)
	}()

	err = session.Run(fmt.Sprintf("cat > %s", remotePath))
	if err != nil {
		return fmt.Errorf("failed to run remote command: %w", err)
	}
	return nil
}
