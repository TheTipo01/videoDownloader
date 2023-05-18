package main

type config struct {
	Token    string   `fig:"token" validate:"required"`
	LogLevel string   `fig:"loglevel" validate:"required"`
	URLs     []string `fig:"urls" validate:"required"`
}
