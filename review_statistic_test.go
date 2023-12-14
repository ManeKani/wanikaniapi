package wanikaniapi_test

import (
	"net/http"
	"testing"

	"github.com/sixels/wanikaniapi"
	"github.com/sixels/wanikaniapi/wktesting"
	assert "github.com/stretchr/testify/require"
)

func TestReviewStatisticList(t *testing.T) {
	client := wktesting.LocalClient()

	_, err := client.ReviewStatisticList(&wanikaniapi.ReviewStatisticListParams{
		Hidden: wanikaniapi.Bool(true),
		IDs:    []wanikaniapi.WKID{1, 2, 3},
	})
	assert.NoError(t, err)

	req := client.RecordedRequests[0]
	assert.Equal(t, "", string(req.Body))
	assert.Equal(t, http.MethodGet, req.Method)
	assert.Equal(t, "/v2/review_statistics", req.Path)
	assert.Equal(t, "hidden=true&ids=1,2,3", wktesting.MustQueryUnescape(req.Query))
}

func TestReviewStatisticGet(t *testing.T) {
	client := wktesting.LocalClient()

	_, err := client.ReviewStatisticGet(&wanikaniapi.ReviewStatisticGetParams{ID: wanikaniapi.ID(123)})
	assert.NoError(t, err)

	req := client.RecordedRequests[0]
	assert.Equal(t, "", string(req.Body))
	assert.Equal(t, http.MethodGet, req.Method)
	assert.Equal(t, "/v2/review_statistics/123", req.Path)
	assert.Equal(t, "", req.Query)
}
