package runscope

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestCreateBucket(t *testing.T) {
	testPreCheck(t)
	client := clientConfigure()
	bucket, err := client.CreateBucket(&Bucket{Name: "test", Team: &Team{ID: teamID}})

	if err != nil {
		t.Error(err)
	}

	client.DeleteBucket(bucket.Key)
}
func TestDeleteBuckets(t *testing.T) {
	testPreCheck(t)
	client := clientConfigure()
	bucket, err := client.CreateBucket(&Bucket{Name: "test-fred", Team: &Team{ID: teamID}})
	defer client.DeleteBucket(bucket.Key)
	if err != nil {
		t.Error(err)
	}

	bucket2, err := client.CreateBucket(&Bucket{Name: "test-bob", Team: &Team{ID: teamID}})
	defer client.DeleteBucket(bucket2.Key)
	if err != nil {
		t.Error(err)
	}

	client.DeleteBuckets(func(bucket *Bucket) bool { return bucket.Name == "test-bob" })

	fredBucket, err := client.ReadBucket(bucket.Key)
	if fredBucket == nil {
		t.Errorf("Bucket key: %v should not be deleted", bucket.Key)
	}

	bobBucket, err := client.ReadBucket(bucket2.Key)
	if bobBucket != nil {
		t.Errorf("Bucket key: %v should be deleted", bobBucket.Key)
	}

}

func TestListBuckets(t *testing.T) {
	testPreCheck(t)
	client := clientConfigure()
	bucket, err := client.CreateBucket(&Bucket{Name: "test", Team: &Team{ID: teamID}})
	defer client.DeleteBucket(bucket.Key)
	if err != nil {
		t.Error(err)
	}

	bucket2, err := client.CreateBucket(&Bucket{Name: "test2", Team: &Team{ID: teamID}})
	defer client.DeleteBucket(bucket2.Key)
	if err != nil {
		t.Error(err)
	}

	results, err := client.ListBuckets()

	if err != nil {
		t.Error(err)
	}

	if results == nil {
		t.Error("list buckets result cannot be nil")
	}

	if len(results) < 2 {
		t.Errorf("Length of buckets expected more than 1, actual:%v", len(results))
	}
}

func TestReadBucket(t *testing.T) {
	testPreCheck(t)
	client := clientConfigure()

	createdBucket, err := client.CreateBucket(&Bucket{Name: "terraform-client.go-test", Team: &Team{ID: teamID}})
	if err != nil {
		t.Error(err)
	}

	readBucket, err := client.ReadBucket(createdBucket.Key)
	if err != nil {
		t.Error(err)
	}

	if readBucket.Key != createdBucket.Key {
		t.Errorf("Bucket createdBucket expected %s was %s.", createdBucket.Key, readBucket.Key)
	}

	if readBucket.TestsURL != fmt.Sprintf("https://api.runscope.com/buckets/%s/tests", readBucket.Key) {
		t.Errorf("Bucket url expected %s was %s.",
			fmt.Sprintf("https://api.runscope.com/buckets/%s/tests", readBucket.Key), readBucket.TestsURL)
	}

	client.DeleteBucket(createdBucket.Key)
}

func TestBucketReadFromResponse(t *testing.T) {
	responseBody := `
{
  "meta": {
    "status": "success"
  },
  "data": {
    "verify_ssl": true,
    "trigger_url": "https://api.runscope.com/radar/bucket/f2f4dbbb-7bf0-4528-bf51-eb3d06a20423/trigger",
    "name": "Sample Name",
    "key": "z3n32gktzx94",
    "team": {
      "name": "form3",
      "id": "870ed937-bc6e-4d8b-a9a5-d7f9f2412fa3"
    },
    "default": false,
    "auth_token": null,
    "tests_url": "https://api.runscope.com/buckets/z3n32gktzx94/tests",
    "collections_url": "https://api.runscope.com/buckets/z3n32gktzx94/collections",
    "messages_url": "https://api.runscope.com/buckets/z3n32gktzx94/stream"
  },
  "error": null
}
`
	responseMap := make(map[string]interface{})
	if err := json.Unmarshal([]byte(responseBody), &responseMap); err != nil {
		t.Error(err)
	}

	dataMap := responseMap["data"].(map[string]interface{})
	bucket, err := getBucketFromResponse(dataMap)
	if err != nil {
		t.Error(err)
	}

	if bucket.Name != "Sample Name" {
		t.Errorf("Expected name %s, actual %s", "Sample Name", bucket.Name)
	}

	if len(bucket.TestsURL) == 0 {
		t.Error("Missing test url")
	}
}
