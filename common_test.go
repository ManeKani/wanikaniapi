package wanikaniapi_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/sixels/wanikaniapi"
	"github.com/sixels/wanikaniapi/wktesting"
	assert "github.com/stretchr/testify/require"
)

func TestClientContext(t *testing.T) {
	client := wktesting.LocalClient()

	ctx := context.Background()

	_, err := client.SubjectList(&wanikaniapi.SubjectListParams{
		Params: wanikaniapi.Params{
			Context: &ctx,
		},
	})

	// This doesn't really text anything, we may want to test this with with an
	// actual API call.
	assert.NoError(t, err)
}

func TestClientError(t *testing.T) {
	client := wktesting.LocalClient()

	client.RecordedResponses = []*wanikaniapi.RecordedResponse{
		{StatusCode: http.StatusTooManyRequests, Body: []byte(`{
			"code": 429,
			"error": "You are rate limited"
		}`)},
	}

	_, err := client.SubjectList(&wanikaniapi.SubjectListParams{})

	assert.Equal(t, &wanikaniapi.APIError{
		StatusCode: http.StatusTooManyRequests,
		Message:    "You are rate limited",
	}, err)
}

// See also `samples/if_modified_since/main.go` to test this for real.
func TestClientIfModifiedSince(t *testing.T) {
	client := wktesting.LocalClient()

	client.RecordedResponses = []*wanikaniapi.RecordedResponse{
		{StatusCode: http.StatusNotModified, Body: []byte(`{}`)},
	}

	obj, err := client.SubjectList(&wanikaniapi.SubjectListParams{
		Params: wanikaniapi.Params{
			IfModifiedSince: wanikaniapi.Time(time.Now()),
		},
	})
	assert.NoError(t, err)
	assert.True(t, obj.NotModified)
}

// See also `samples/if_none_match/main.go` to test this for real.
func TestClientIfNoneMatch(t *testing.T) {
	client := wktesting.LocalClient()

	client.RecordedResponses = []*wanikaniapi.RecordedResponse{
		{StatusCode: http.StatusNotModified, Body: []byte(`{}`)},
	}

	obj, err := client.SubjectList(&wanikaniapi.SubjectListParams{
		Params: wanikaniapi.Params{
			IfNoneMatch: wanikaniapi.String("an-etag"),
		},
	})
	assert.NoError(t, err)
	assert.True(t, obj.NotModified)
}

func TestClientRetry(t *testing.T) {
	client := wktesting.LocalClient()
	client.MaxRetries = 2
	client.NoRetrySleep = true

	client.RecordedResponses = []*wanikaniapi.RecordedResponse{
		{StatusCode: http.StatusTooManyRequests, Body: []byte(`{
			"code": 429,
			"error": "You are rate limited"
		}`)},
		{StatusCode: http.StatusTooManyRequests, Body: []byte(`{
			"code": 429,
			"error": "You are rate limited"
		}`)},
		{StatusCode: http.StatusOK, Body: []byte(`{}`)},
	}

	subjects, err := client.SubjectList(&wanikaniapi.SubjectListParams{})

	assert.NoError(t, err)
	assert.Equal(t, 0, len(subjects.Data))
}

func TestClientNoRetry(t *testing.T) {
	client := wktesting.LocalClient()

	client.RecordedResponses = []*wanikaniapi.RecordedResponse{
		{StatusCode: http.StatusTooManyRequests, Body: []byte(`{
			"code": 429,
			"error": "You are rate limited"
		}`)},
	}

	_, err := client.SubjectList(&wanikaniapi.SubjectListParams{})

	assert.Equal(t, &wanikaniapi.APIError{
		StatusCode: http.StatusTooManyRequests,
		Message:    "You are rate limited",
	}, err)
}

func TestPageFullyLocal(t *testing.T) {
	client := wktesting.LocalClient()

	client.RecordedResponses = []*wanikaniapi.RecordedResponse{
		{StatusCode: http.StatusOK, Body: []byte(`{
			"pages": {
				"per_page": 1000,
				"next_url": "https://api.wanikani.com/v2/subjects?page_after_id=125",
				"previous_url": null
			},
			"data": [
				{"id": 123, "object": "kanji"},
				{"id": 124, "object": "kanji"},
				{"id": 125, "object": "kanji"}
			]
		}`)},
		{StatusCode: http.StatusOK, Body: []byte(`{
			"pages": {
				"per_page": 1000,
				"next_url": "https://api.wanikani.com/v2/subjects?page_after_id=128",
				"previous_url": null
			},
			"data": [
				{"id": 126, "object": "kanji"},
				{"id": 127, "object": "kanji"},
				{"id": 128, "object": "kanji"}
			]
		}`)},
		{StatusCode: http.StatusOK, Body: []byte(`{
			"pages": {
				"per_page": 1000,
				"next_url": "https://api.wanikani.com/v2/subjects?page_after_id=129",
				"previous_url": null
			},
			"data": [
				{"id": 129, "object": "kanji"},
				{"id": 130, "object": "kanji"},
				{"id": 131, "object": "kanji"}
			]
		}`)},
	}

	var subjects []*wanikaniapi.Subject
	err := client.PageFully(func(id *wanikaniapi.WKID) (*wanikaniapi.PageObject, error) {
		page, err := client.SubjectList(&wanikaniapi.SubjectListParams{
			ListParams: wanikaniapi.ListParams{
				PageAfterID: id,
			},
		})
		if err != nil {
			return nil, err
		}

		subjects = append(subjects, page.Data...)
		return &page.PageObject, nil
	})
	assert.NoError(t, err)

	assert.Equal(t,
		[]*wanikaniapi.Subject{
			// page 1
			{Object: wanikaniapi.Object{ID: 123, ObjectType: wanikaniapi.ObjectTypeKanji}},
			{Object: wanikaniapi.Object{ID: 124, ObjectType: wanikaniapi.ObjectTypeKanji}},
			{Object: wanikaniapi.Object{ID: 125, ObjectType: wanikaniapi.ObjectTypeKanji}},

			// page 2
			{Object: wanikaniapi.Object{ID: 126, ObjectType: wanikaniapi.ObjectTypeKanji}},
			{Object: wanikaniapi.Object{ID: 127, ObjectType: wanikaniapi.ObjectTypeKanji}},
			{Object: wanikaniapi.Object{ID: 128, ObjectType: wanikaniapi.ObjectTypeKanji}},

			// page 3
			{Object: wanikaniapi.Object{ID: 129, ObjectType: wanikaniapi.ObjectTypeKanji}},
			{Object: wanikaniapi.Object{ID: 130, ObjectType: wanikaniapi.ObjectTypeKanji}},
			{Object: wanikaniapi.Object{ID: 131, ObjectType: wanikaniapi.ObjectTypeKanji}},
		},
		subjects,
	)
}

func TestPageFullyLive(t *testing.T) {
	client := wktesting.LiveClient()
	if client == nil {
		return
	}

	var i int
	var subjects []*wanikaniapi.Subject
	err := client.PageFully(func(id *wanikaniapi.WKID) (*wanikaniapi.PageObject, error) {
		// Quit after convincing ourselves that we can page by going through a
		// few pages. This saves time and API calls.
		i++
		if i > 2 {
			return nil, nil
		}

		page, err := client.SubjectList(&wanikaniapi.SubjectListParams{
			ListParams: wanikaniapi.ListParams{
				PageAfterID: id,
			},
		})
		if err != nil {
			return nil, err
		}

		subjects = append(subjects, page.Data...)
		return &page.PageObject, nil
	})
	assert.NoError(t, err)

	t.Logf("num subjects paged before quitting: %v", len(subjects))
}

func TestWKTimeMarshalJSON(t *testing.T) {
	goT := time.Now()
	wkT := wanikaniapi.WKTime(goT)
	marshaled, err := wkT.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf(`"%s"`, goT.Format(time.RFC3339)), string(marshaled))
}
