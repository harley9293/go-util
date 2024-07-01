package net

import (
	"errors"
	"io/ioutil"
	"net/http"
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
