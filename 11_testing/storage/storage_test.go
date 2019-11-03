package storage

import (
	"strings"
	"testing"
)

func TestCheckQuotaNotifiesUser(t *testing.T) {
	// Save and restore original func notifyUser after tests
	original := notifyUser
	defer func() { notifyUser = original }()

	// Override original notifyUser with "fake"/simple implementation.
	var notifiedUser, notifiedMsg string
	notifyUser = func(user, msg string) {
		notifiedUser, notifiedMsg = user, msg
	}

	const user = "jane@example.org"
	usage[user] = 980000000 // simulate 980MB-used condition

	CheckQuota(user)

	if notifiedUser == "" && notifiedMsg == "" {
		t.Fatalf("notifyUser not called") // immedaiate termination
	}
	if notifiedUser != user {
		t.Errorf("wrong user (%s) notified, want %s",
			notifiedUser, user)
	}
	const wantSubstring = "98% of your quota"
	if !strings.Contains(notifiedMsg, wantSubstring) {
		t.Errorf("unexpected notification message <<%s>>, "+
			"want substring %q", notifiedMsg, wantSubstring)
	}
}
