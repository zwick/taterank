package main

// Enum for the different modes the application can run in. Choices are http and lambda
type AppMode string

const (
	HTTP   AppMode = "http"
	Lambda AppMode = "lambda"
)
