package cli

import (
	"testing"
)

func TestHandleConfigCombinedFlags(t *testing.T) {
	config := NewConfig()
	flags := []string{"-nr"}

	err := config.HandleConfig(flags)
	if err != nil {
		t.Errorf("Got error: %s", err)
	}

	if !config.Cli {
		t.Error("-n failed")
	}

	if !config.random {
		t.Error("-r failed")
	}
}

func TestHandleConfigSeparatedFlags(t *testing.T) {
	config := NewConfig()
	flags := []string{"-n", "-r"}

	err := config.HandleConfig(flags)
	if err != nil {
		t.Errorf("Got error: %s", err)
	}

	if !config.Cli {
		t.Error("-n failed")
	}

	if !config.random {
		t.Error("-r failed")
	}
}
