package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/fselich/dossier/internal/openspec"
	"github.com/fselich/dossier/internal/settings"
	"github.com/fselich/dossier/internal/ui"
)

var version string

type cliOptions struct {
	themeName   string
	themeSet    bool
	configPath  string
	showVersion bool
	showHelp    bool
	readOnly    bool
}

func parseOptions(args []string, output io.Writer) (cliOptions, []string, error) {
	var opts cliOptions
	var shortHelp bool

	flags := flag.NewFlagSet("dossier", flag.ContinueOnError)
	flags.SetOutput(output)
	flags.StringVar(&opts.configPath, "config", "", "Configuration file (default: XDG config path)")
	flags.StringVar(&opts.themeName, "theme", "none", "Base theme (dark, none, light, dracula, gruvbox-light-soft)")
	flags.BoolVar(&opts.showVersion, "version", false, "Print version and exit")
	flags.BoolVar(&opts.showHelp, "help", false, "Print help and exit")
	flags.BoolVar(&shortHelp, "h", false, "Print help and exit")
	flags.BoolVar(&opts.readOnly, "read-only", false, "Disable task toggles, editor launch, and Git stage/unstage")
	flags.Usage = func() { writeUsage(output) }

	if err := flags.Parse(args); err != nil {
		return cliOptions{}, nil, err
	}
	opts.showHelp = opts.showHelp || shortHelp
	flags.Visit(func(f *flag.Flag) {
		if f.Name == "theme" {
			opts.themeSet = true
		}
	})
	return opts, flags.Args(), nil
}

func writeUsage(output io.Writer) {
	_, _ = fmt.Fprint(output, `Usage: dossier [options] [path]

A keyboard-driven TUI for navigating OpenSpec project artifacts.

Options:
  -h, --help          Print help and exit
      --config <path> Configuration file (default: XDG config path)
      --read-only     Disable task toggles, editor launch, and Git stage/unstage
      --theme <name>  Base theme: dark, none, light, dracula, or gruvbox-light-soft (default: none)
      --version       Print version and exit
`)
}

func main() {
	var (
		project *openspec.Project
		err     error
		model   ui.Model
	)

	opts, args, err := parseOptions(os.Args[1:], os.Stderr)
	if err != nil {
		os.Exit(2)
	}

	if opts.showVersion {
		fmt.Println("dossier", version)
		os.Exit(0)
	}

	if opts.showHelp {
		writeUsage(os.Stderr)
		os.Exit(0)
	}

	appConfig, _, err := settings.Load(opts.configPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	if opts.themeSet {
		appConfig.Theme.Base = opts.themeName
	}
	theme, err := ui.BuildTheme(appConfig.Theme)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error: loading theme:", err)
		os.Exit(1)
	}

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error: cannot determine working directory:", err)
		os.Exit(1)
	}

	cfg, err := openspec.LoadConfigFrom(cwd)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error: loading config:", err)
		os.Exit(1)
	}

	loader := openspec.NewLoader(openspec.OSFS{})

	var pathArg string
	if len(args) > 0 {
		pathArg = args[0]
	}
	if pathArg != "" {
		project, err = openspec.LoadFromPath(pathArg)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
		model = ui.NewSinglePath(project, cfg, pathArg, loader, theme, appConfig.Keys, opts.readOnly)
	} else {
		project, err = openspec.LoadFrom(cwd)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		model = ui.New(project, cfg, cwd, loader, theme, appConfig.Keys, opts.readOnly)
	}

	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
