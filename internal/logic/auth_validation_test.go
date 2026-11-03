package logic

import "testing"

func TestValidUsername(t *testing.T) {
	tests := []struct {
		name     string
		username string
		want     bool
	}{
		{name: "minimum length", username: "abc", want: true},
		{name: "maximum length", username: "user_123456789012345678901234567", want: true},
		{name: "too short", username: "ab", want: false},
		{name: "invalid character", username: "user-name", want: false},
		{name: "non ASCII", username: "用户abc", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validUsername(tt.username); got != tt.want {
				t.Fatalf("validUsername(%q) = %t, want %t", tt.username, got, tt.want)
			}
		})
	}
}

func TestValidPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		want     bool
	}{
		{name: "minimum length", password: "12345678", want: true},
		{name: "maximum length", password: string(make([]byte, 72)), want: true},
		{name: "too short", password: "1234567", want: false},
		{name: "too long", password: string(make([]byte, 73)), want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validPassword(tt.password); got != tt.want {
				t.Fatalf("validPassword(%q) = %t, want %t", tt.password, got, tt.want)
			}
		})
	}
}
