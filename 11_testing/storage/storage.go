// Package storage sends an automated email when users near quota.
package storage

import (
	"fmt"
	"log"
	"net/smtp"
)

var usage = make(map[string]int64)

func bytesInUse(username string) int64 {
	return usage[username]
}

// Email example sender convifguration.
// Never actually put passwords in source code.
const (
	sender   = "notifications@example.com"
	password = "fakepasswd"
	hostname = "stmp.example.com"
)

const template = `Warning: you are using %d bytes of storage,
%d%% of your quota.`

var notifyUser = func(username, msg string) {
	auth := smtp.PlainAuth("", sender, password, hostname)
	err := smtp.SendMail(hostname+":587", auth, sender,
		[]string{username}, []byte(msg))
	if err != nil {
		log.Printf("smtp.SendMail(%s) failed: %s", username, err)
	}
}

// CheckQuota calls notifyUser if less than 10% storage remains.
func CheckQuota(username string) {
	used := bytesInUse(username)
	const quota = 1e9 // 1e9 bytes = 1GB
	percent := 100 * used / quota
	if percent < 90 {
		return // OK
	}
	msg := fmt.Sprintf(template, used, percent)
	notifyUser(username, msg)
}
