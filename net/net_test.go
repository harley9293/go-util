package net

import (
	"testing"
)

func TestGetPublicIP(t *testing.T) {
	ip, err := GetPublicIP()
	if err != nil {
		t.Fatalf(err.Error())
	}

	t.Logf("ip: " + ip)
}
