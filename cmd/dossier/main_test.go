package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestParseOptionsReadOnly(t *testing.T) {
	var output bytes.Buffer
	opts, args, err := parseOptions([]string{"--read-only", "/tmp/change"}, &output)
	if err != nil {
		t.Fatalf("parseOptions: %v", err)
	}
	if !opts.readOnly {
		t.Error("expected --read-only to enable read-only mode")
	}
	if len(args) != 1 || args[0] != "/tmp/change" {
		t.Fatalf("expected path argument to be preserved, got %v", args)
	}
}

func TestParseOptionsHelpAliases(t *testing.T) {
	for _, arg := range []string{"-h", "--help"} {
		t.Run(arg, func(t *testing.T) {
			var output bytes.Buffer
			opts, _, err := parseOptions([]string{arg}, &output)
			if err != nil {
				t.Fatalf("parseOptions: %v", err)
			}
			if !opts.showHelp {
				t.Errorf("expected %s to request help", arg)
			}
		})
	}
}

func TestUsageUsesConventionalFlagSpelling(t *testing.T) {
	var output bytes.Buffer
	writeUsage(&output)
	usage := output.String()

	for _, want := range []string{"-h, --help", "--read-only", "--theme <name>", "--version"} {
		if !strings.Contains(usage, want) {
			t.Errorf("expected usage to contain %q, got:\n%s", want, usage)
		}
	}
	if strings.Contains(usage, "\n  -help") {
		t.Errorf("usage should not advertise single-dash long help:\n%s", usage)
	}
}
