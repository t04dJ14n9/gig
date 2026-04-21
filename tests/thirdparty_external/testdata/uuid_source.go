package main

import "github.com/google/uuid"

func UUIDNewString() string {
	return uuid.NewString()
}

func UUIDRoundTrip(s string) string {
	u, err := uuid.Parse(s)
	if err != nil {
		return "ERR"
	}
	return u.String()
}

func UUIDURNPrefix() string {
	u := uuid.New()
	return u.URN()[:9]
}
