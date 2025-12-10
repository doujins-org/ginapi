// Package response provides Stripe-style response helpers for Gin APIs.
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Object sends a single object response.
// The object should have an "object" field identifying its type.
func Object(c *gin.Context, obj any) {
	c.JSON(http.StatusOK, obj)
}

// Created sends a 201 Created response with the created object.
func Created(c *gin.Context, obj any) {
	c.JSON(http.StatusCreated, obj)
}

// NoContent sends a 204 No Content response.
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// Deleted sends a Stripe-style deletion confirmation.
func Deleted(c *gin.Context, objectType string, id string) {
	c.JSON(http.StatusOK, DeletedObject{
		Object:  objectType,
		ID:      id,
		Deleted: true,
	})
}

// DeletedObject represents a Stripe-style deletion response.
type DeletedObject struct {
	Object  string `json:"object"`
	ID      string `json:"id"`
	Deleted bool   `json:"deleted"`
}

// Message represents a simple message response.
type Message struct {
	Object  string `json:"object"`
	Message string `json:"message"`
}

// Success sends a 200 OK response with a success message.
func Success(c *gin.Context, message string) {
	c.JSON(http.StatusOK, Message{
		Object:  "message",
		Message: message,
	})
}
