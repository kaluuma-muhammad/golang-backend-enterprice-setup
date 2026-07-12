package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/mssola/user_agent"
)

type DeviceInfo struct {
	DeviceID   string
	Platform   string
	Browser    string
	DeviceName string
}

func ParseDeviceInfo(userAgentStr, ip string) DeviceInfo {
	fmt.Println("USER AGENT:", userAgentStr)

	ua := user_agent.New(userAgentStr)

	browserName, browserVersion := ua.Browser()

	platform := ua.Platform()

	deviceName := "Unknown Device"

	if ua.Mobile() {
		deviceName = "Mobile Device"
	} else {
		deviceName = platform + " PC"
	}

	hash := sha256.Sum256([]byte(userAgentStr + ip))
	deviceID := hex.EncodeToString(hash[:])

	return DeviceInfo{
		DeviceID:   deviceID,
		Platform:   platform,
		Browser:    browserName + " " + browserVersion,
		DeviceName: deviceName,
	}
}
