package data

import "embed"

var (
	//go:embed web/*
	WebFs embed.FS
)
