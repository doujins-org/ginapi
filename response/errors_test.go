package response_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/doujins-org/ginapi/response"
)

func TestBadRequest(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	response.BadRequest(c, "invalid input")

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}

	var result response.Error
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if result.Object != "error" {
		t.Errorf("expected object 'error', got '%s'", result.Object)
	}
	if result.Error.Type != response.ErrorTypeInvalidRequest {
		t.Errorf("expected type '%s', got '%s'", response.ErrorTypeInvalidRequest, result.Error.Type)
	}
	if result.Error.Message != "invalid input" {
		t.Errorf("expected message 'invalid input', got '%s'", result.Error.Message)
	}
}

func TestBadRequestParam(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	response.BadRequestParam(c, "email", "invalid email format")

	var result response.Error
	json.Unmarshal(w.Body.Bytes(), &result)

	if result.Error.Param != "email" {
		t.Errorf("expected param 'email', got '%s'", result.Error.Param)
	}
}

func TestUnauthorized(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	response.Unauthorized(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}

	var result response.Error
	json.Unmarshal(w.Body.Bytes(), &result)

	if result.Error.Type != response.ErrorTypeAuthentication {
		t.Errorf("expected type '%s', got '%s'", response.ErrorTypeAuthentication, result.Error.Type)
	}
}

func TestForbidden(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	response.Forbidden(c)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected status 403, got %d", w.Code)
	}
}

func TestNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	response.NotFound(c, "artist")

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}

	var result response.Error
	json.Unmarshal(w.Body.Bytes(), &result)

	if result.Error.Type != response.ErrorTypeNotFound {
		t.Errorf("expected type '%s', got '%s'", response.ErrorTypeNotFound, result.Error.Type)
	}
	if result.Error.Message != "artist not found" {
		t.Errorf("expected message 'artist not found', got '%s'", result.Error.Message)
	}
}

func TestConflict(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	response.Conflict(c, "resource already exists")

	if w.Code != http.StatusConflict {
		t.Errorf("expected status 409, got %d", w.Code)
	}
}

func TestTooManyRequests(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	response.TooManyRequests(c, "rate limit exceeded")

	if w.Code != http.StatusTooManyRequests {
		t.Errorf("expected status 429, got %d", w.Code)
	}

	var result response.Error
	json.Unmarshal(w.Body.Bytes(), &result)

	if result.Error.Type != response.ErrorTypeRateLimit {
		t.Errorf("expected type '%s', got '%s'", response.ErrorTypeRateLimit, result.Error.Type)
	}
}

func TestInternalError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	response.InternalError(c, "something went wrong")

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", w.Code)
	}

	var result response.Error
	json.Unmarshal(w.Body.Bytes(), &result)

	if result.Error.Type != response.ErrorTypeAPI {
		t.Errorf("expected type '%s', got '%s'", response.ErrorTypeAPI, result.Error.Type)
	}
}

func TestServiceUnavailable(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	response.ServiceUnavailable(c, "service temporarily unavailable")

	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("expected status 503, got %d", w.Code)
	}
}
