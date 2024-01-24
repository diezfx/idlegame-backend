package integration

import (
	"context"
	"testing"

	"github.com/diezfx/idlegame-backend/internal/api"
	"github.com/diezfx/idlegame-backend/pkg/httpclient"
)

func TestGathering(t *testing.T) {
	ctx := context.Background()
	const user = 1
	tests := []struct {
		name      string
		monsterID int
		jobDefID  string

		expectedResponseCode int
	}{
		{
			name:      "test",
			monsterID: 1,
			jobDefID:  "spruce",

			expectedResponseCode: 201,
		},

		{
			name:      "test2",
			monsterID: 2,
			jobDefID:  "stone",

			expectedResponseCode: 201,
		},

		{
			name:      "monster not found",
			monsterID: -1,
			jobDefID:  "wood",

			expectedResponseCode: 400,
		},

		{
			name:      "job not found",
			monsterID: 1,
			jobDefID:  "notFound",

			expectedResponseCode: 400,
		},
	}

	client := httpclient.New(httpclient.Config{
		Host: "http://localhost:8080/api/v1.0/jobs/",
	})

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := api.StartGatheringJob{
				UserID:   user,
				Monster:  test.monsterID,
				JobDefID: test.jobDefID,
			}
			response, err := client.Post(ctx, "/gathering", request, nil)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if response.StatusCode != test.expectedResponseCode {
				t.Errorf("expected response code %d, got %d", test.expectedResponseCode, response.StatusCode)
			}
		})
	}

}

// check post to post a new gatheringJob
