package models

import "testing"

func TestUser_CheckPassword(t *testing.T) {
	u := User{
		Id:       1,
		Name:     "john doe",
		password: "john1234!?",
	}
	pw, err := u.GetPasswordHash()
	if err != nil {
		t.Fatalf("Password generation failed when it should have worked! (Password: %s)", pw)
	}
	if err := u.CheckPassword(pw); err != nil {
		t.Fatalf("The password check failed while it should have passed (Password: %s)", pw)
	}
	u.SetPassword("john1234!")
	if err := u.CheckPassword(pw); err == nil {
		t.Fatalf("The password check succeeded while it should have failed (Password: %s)", pw)
	}
}
