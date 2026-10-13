package dtrack

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestTeamService_GetAPIKeys(t *testing.T) {
	target := uuid.MustParse("0d9c5a4e-0000-4000-8000-000000000060")

	// 120 teams, paged the way v5 pages them. The target team is on page 2.
	registerTeams := func(t *testing.T) {
		httpmock.RegisterResponder(http.MethodGet, "http://localhost/api/v1/team",
			func(req *http.Request) (*http.Response, error) {
				pageNumber, _ := strconv.Atoi(req.URL.Query().Get("pageNumber"))
				pageSize, _ := strconv.Atoi(req.URL.Query().Get("pageSize"))

				body := "["
				for i := (pageNumber - 1) * pageSize; i < pageNumber*pageSize && i < 120; i++ {
					if i > (pageNumber-1)*pageSize {
						body += ","
					}
					id := fmt.Sprintf("0d9c5a4e-0000-4000-8000-%012d", i)
					keys := "[]"
					if id == target.String() {
						keys = `[{"publicId":"odt_abc","maskedKey":"odt_abc****"}]`
					}
					body += fmt.Sprintf(`{"uuid":%q,"name":"team-%d","apiKeys":%s}`, id, i, keys)
				}
				body += "]"

				res := httpmock.NewStringResponse(http.StatusOK, body)
				res.Header.Set("X-Total-Count", "120")
				return res, nil
			})
	}

	t.Run("finds a team past the first page", func(t *testing.T) {
		client := newTestClient(t, "5.2.0")
		registerTeams(t)

		keys, err := client.Team.GetAPIKeys(context.TODO(), target)
		require.NoError(t, err)
		require.Len(t, keys, 1)
		require.Equal(t, "odt_abc", keys[0].PublicID)
	})

	t.Run("unknown team is a 404", func(t *testing.T) {
		client := newTestClient(t, "5.2.0")
		registerTeams(t)

		_, err := client.Team.GetAPIKeys(context.TODO(), uuid.New())
		var apiErr *APIError
		require.True(t, errors.As(err, &apiErr))
		require.Equal(t, http.StatusNotFound, apiErr.StatusCode)
	})
}
