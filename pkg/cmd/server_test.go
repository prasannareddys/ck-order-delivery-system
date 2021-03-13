package cmd

import (
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewServerCommand(t *testing.T) {

	cmd := NewServerCommand()

	err := cmd.Flags().Set("order-file-path", "../../data/test/orders_test.json")
	require.NoError(t, err)

	err = cmd.Flags().Set("ops", "2")
	require.NoError(t, err)

	err = cmd.RunE(&cobra.Command{}, nil)
	if err != nil {
		t.FailNow()
	}

}