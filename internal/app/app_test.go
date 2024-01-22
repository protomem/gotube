package application_test

import (
	"testing"

	application "github.com/protomem/gotube/internal/app"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
)

func TestValidateApp(t *testing.T) {
	err := fx.ValidateApp(application.Create())
	require.NoError(t, err)
}
