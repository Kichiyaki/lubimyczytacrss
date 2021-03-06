package lubimyczytac_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Kichiyaki/lubimyczytacrss/internal/lubimyczytac/testdata"

	"github.com/stretchr/testify/assert"

	"github.com/Kichiyaki/lubimyczytacrss/internal/lubimyczytac"
)

func TestClient_GetAuthor(t *testing.T) {
	t.Parallel()

	t.Run("OK", func(t *testing.T) {
		t.Parallel()

		for _, author := range testdata.Authors {
			author := author

			t.Run("authorID="+author.ID, func(t *testing.T) {
				t.Parallel()

				srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path != fmt.Sprintf("/autor/%s/x", author.ID) {
						w.WriteHeader(http.StatusNotFound)
						return
					}
					if r.Method != http.MethodGet {
						w.WriteHeader(http.StatusMethodNotAllowed)
						return
					}

					w.WriteHeader(http.StatusOK)
					_, _ = w.Write(author.HTML)
				}))
				defer srv.Close()

				result, err := lubimyczytac.
					NewClient(srv.Client(), lubimyczytac.WithBaseURL(srv.URL)).
					GetAuthor(context.Background(), author.ID)
				assert.NoError(t, err)
				assert.Equal(t, author.Author, result)
			})
		}
	})

	t.Run("ERR: author not found", func(t *testing.T) {
		t.Parallel()

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		result, err := lubimyczytac.
			NewClient(srv.Client(), lubimyczytac.WithBaseURL(srv.URL)).
			GetAuthor(context.Background(), "123")
		assert.ErrorIs(t, err, lubimyczytac.ErrAuthorNotFound)
		assert.Zero(t, result)
	})
}
