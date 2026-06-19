package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"

	_ "embed"

	"github.com/charmbracelet/fang"
	"github.com/spf13/cobra"

	"gox/internal/server/config"
	"gox/internal/server/gox"
	"gox/pkg/start"
)

const (
	_ int = iota
	initCode
	fatalCode
)

//go:embed server.key
var key []byte

//go:embed server.crt
var cert []byte

var (
	Username, Password string
)

type exitError struct {
	code int
	err  error
}

func (e *exitError) Error() string {
	if e == nil || e.err == nil {
		return ""
	}

	return e.err.Error()
}

func (e *exitError) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.err
}

func newExitError(code int, err error) error {
	if err == nil {
		return nil
	}

	return &exitError{code: code, err: err}
}

func main() {
	configureLogger()

	if err := fang.Execute(context.Background(), rootCmd(), fang.WithoutVersion()); err != nil {
		var exitErr *exitError
		if errors.As(err, &exitErr) {
			fmt.Fprintln(os.Stderr, exitErr.err)
			os.Exit(exitErr.code)
		}

		fmt.Fprintln(os.Stderr, err)
		os.Exit(fatalCode)
	}
}

func rootCmd() *cobra.Command {
	var path string

	rootCmd := &cobra.Command{
		Use:           "gox",
		Short:         "SOCKS5 and HTTPS proxy server",
		SilenceUsage:  true,
		SilenceErrors: true,
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			return runProxy(path)
		},
	}

	rootCmd.PersistentFlags().StringVar(&path, "config", "config.yaml", "path to config file")
	rootCmd.AddCommand(
		&cobra.Command{
			Use:   "save",
			Short: "Save default config",
			RunE: func(_ *cobra.Command, _ []string) error {
				if err := config.Default(path); err != nil {
					return newExitError(fatalCode, err)
				}

				slog.Info("saved default config", "path", path)
				return nil
			},
		},
		&cobra.Command{
			Use:   "setup",
			Short: "Set autostart via systemd",
			RunE: func(_ *cobra.Command, _ []string) error {
				return runStartAction("set autostart", func(s *start.Start) error {
					return s.Setup()
				})
			},
		},
		&cobra.Command{
			Use:   "remove",
			Short: "Remove autostart via systemd",
			RunE: func(_ *cobra.Command, _ []string) error {
				return runStartAction("remove autostart", func(s *start.Start) error {
					return s.Remove()
				})
			},
		},
	)

	return rootCmd
}

func runProxy(path string) error {
	_config, err := config.New(path, Username, Password)
	if err != nil {
		return newExitError(initCode, fmt.Errorf("fatal config error: %w", err))
	}

	_gox, err := gox.New(_config, key, cert)
	if err != nil {
		return newExitError(fatalCode, fmt.Errorf("create fatal error: %w", err))
	}

	if err := _gox.Listen(); err != nil {
		return newExitError(fatalCode, fmt.Errorf("proxy fatal error: %w", err))
	}

	return nil
}

func runStartAction(action string, run func(*start.Start) error) error {
	_start, err := start.New()
	if err != nil {
		return newExitError(initCode, fmt.Errorf("init autostart error: %w", err))
	}

	if err := run(_start); err != nil {
		return newExitError(fatalCode, fmt.Errorf("%s error: %w", action, err))
	}

	slog.Info("autostart command completed", "action", action)
	return nil
}

func configureLogger() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{})))
}
