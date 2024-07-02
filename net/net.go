package net

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// GetPublicIP Retrieve public IP address.
func GetPublicIP() (string, error) {
	c := http.Client{}
	c.Timeout = time.Second * 10
	rsp, err := c.Get("http://47.112.241.125:3001/ip")
	if err != nil {
		return "", errors.New("Failed to retrieve external IP. Please check the network.")
	}
	defer rsp.Body.Close()
	body, _ := ioutil.ReadAll(rsp.Body)

	return string(body), nil
}

// DownloadFile Download file from ssh.Session
func DownloadFile(session *ssh.Session, remotePath, localPath string) error {
	srcFile, err := session.StdoutPipe()
	if err != nil {
		return err
	}

	dstFile, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer func(dstFile *os.File) {
		_ = dstFile.Close()
	}(dstFile)

	go func() {
		_, _ = io.Copy(dstFile, srcFile)
	}()

	err = session.Run("cat " + remotePath)
	return err
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
