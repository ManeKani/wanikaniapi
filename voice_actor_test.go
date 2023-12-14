package wanikaniapi_test

import (
	"net/http"
	"testing"

	"github.com/sixels/wanikaniapi"
	"github.com/sixels/wanikaniapi/wktesting"
	assert "github.com/stretchr/testify/require"
)

func TestVoiceActorList(t *testing.T) {
	client := wktesting.LocalClient()

	_, err := client.VoiceActorList(&wanikaniapi.VoiceActorListParams{
		IDs: []wanikaniapi.WKID{1, 2, 3},
	})
	assert.NoError(t, err)

	req := client.RecordedRequests[0]
	assert.Equal(t, []byte(nil), req.Body)
	assert.Equal(t, http.MethodGet, req.Method)
	assert.Equal(t, "/v2/voice_actors", req.Path)
	assert.Equal(t, "ids=1,2,3", wktesting.MustQueryUnescape(req.Query))
}

func TestVoiceActorGet(t *testing.T) {
	client := wktesting.LocalClient()

	_, err := client.VoiceActorGet(&wanikaniapi.VoiceActorGetParams{ID: wanikaniapi.ID(123)})
	assert.NoError(t, err)

	req := client.RecordedRequests[0]
	assert.Equal(t, []byte(nil), req.Body)
	assert.Equal(t, http.MethodGet, req.Method)
	assert.Equal(t, "/v2/voice_actors/123", req.Path)
	assert.Equal(t, "", req.Query)
}
