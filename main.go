package main

import (
	"github.com/xbt573/project-example/cmd"
	"log/slog"
)

func main() {
	if err := cmd.Execute(); err != nil {
		slog.Error("Failed to start cmd!", slog.String("err", err.Error()))
	}
}
