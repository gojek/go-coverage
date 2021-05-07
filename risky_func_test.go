package main

import "testing"

func Test_trimString(t *testing.T) {

	tests := []struct {
		in   string
		want string
	}{
		{"test", "test"},
		{"testtesttesttesttesttest", "...testtesttesttesttest"},
		{"testtesttesttesttesttesttesttesttesttesttestbest", "...testtesttesttestbest"},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			if got := trimString(tt.in); got != tt.want {
				t.Errorf("trimString() = %v, want %v", got, tt.want)
			}
		})
	}
}
