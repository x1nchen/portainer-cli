package dockerhub

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDockerhubClient(t *testing.T) {
	if testing.Short() {
		t.Skip("skip test")
	}

	c, err := NewClient("https://dockerhub.com", "test", "test")
	require.NoError(t, err)

	t.Run("Auth", func(t *testing.T) {
		ctx := context.Background()
		err = c.Auth(ctx)
		require.NoError(t, err)
	})

	t.Run("FindImageTagList", func(t *testing.T) {
		ctx := context.Background()
		tags, err := c.FindImageTagList(ctx, "copytrade/go-report-grpc-srv")
		require.NoError(t, err)
		for _, tag := range tags {
			t.Logf("%+v", tag)
		}
	})
}
