package response_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/doujins-org/ginapi/response"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestObject(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	type testObj struct {
		Object string `json:"object"`
		ID     string `json:"id"`
		Name   string `json:"name"`
	}

	obj := testObj{Object: "artist", ID: "art_123", Name: "Test Artist"}
	response.Object(c, obj)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var result testObj
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if result.Object != "artist" {
		t.Errorf("expected object 'artist', got '%s'", result.Object)
	}
	if result.ID != "art_123" {
		t.Errorf("expected id 'art_123', got '%s'", result.ID)
	}
}

func TestCreated(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	obj := map[string]string{"object": "tag", "id": "tag_456"}
	response.Created(c, obj)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}
}

func TestNoContent(t *testing.T) {
	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		response.NoContent(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", w.Code)
	}
}

func TestDeleted(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	response.Deleted(c, "artist", "art_123")

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var result response.DeletedObject
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if result.Object != "artist" {
		t.Errorf("expected object 'artist', got '%s'", result.Object)
	}
	if result.ID != "art_123" {
		t.Errorf("expected id 'art_123', got '%s'", result.ID)
	}
	if !result.Deleted {
		t.Error("expected deleted to be true")
	}
}
