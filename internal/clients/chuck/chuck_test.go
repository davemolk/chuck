package chuck

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestSearch(t *testing.T) {
	const validData = "testdata/chuck.json"
	const emptyData = "testdata/empty.json"
	t.Run("success", func(t *testing.T) {
		ts := httptest.NewTLSServer(http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				http.ServeFile(w, r, validData)
			}))
		defer ts.Close()

		c := NewClient(zap.NewNop())
		c.baseURL = ts.URL
		c.client = ts.Client()

		got, err := c.Search(context.Background(), "foo", 10)
		require.NoError(t, err)
		require.Len(t, got, 4)
		require.Equal(t, "c-3yrrglr0ouxifeo2rzsw", got[0].ExternalID)
	})

	t.Run("success, limit < results", func(t *testing.T) {
		ts := httptest.NewTLSServer(http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				http.ServeFile(w, r, validData)
			}))
		defer ts.Close()

		c := NewClient(zap.NewNop())
		c.baseURL = ts.URL
		c.client = ts.Client()

		got, err := c.Search(context.Background(), "foo", 2)
		require.NoError(t, err)
		require.Len(t, got, 2)
		require.Equal(t, "c-3yrrglr0ouxifeo2rzsw", got[0].ExternalID)
	})

	t.Run("handles empty", func(t *testing.T) {
		ts := httptest.NewTLSServer(http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				http.ServeFile(w, r, emptyData)
			}))
		defer ts.Close()

		c := NewClient(zap.NewNop())
		c.baseURL = ts.URL
		c.client = ts.Client()

		got, err := c.Search(context.Background(), "foo", 2)
		require.NoError(t, err)
		require.Len(t, got, 0)
	})
}
