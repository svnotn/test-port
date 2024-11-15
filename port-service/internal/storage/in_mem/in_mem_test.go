package in_mem_test

import (
	"testing"

	"github.com/svnotn/test-port/port-service/internal/model"
	"github.com/svnotn/test-port/port-service/internal/storage/in_mem"
)

// TestStorage tests the functionality of the in-memory Storage.
func TestStorage(t *testing.T) {
	storage := in_mem.New(10, 10)

	portIn := model.Port{
		ID:   1,
		Type: model.TypeIN,
	}

	portOut := model.Port{
		ID:   2,
		Type: model.TypeOUT,
	}

	// Test Add
	if err := storage.Add(portIn); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if err := storage.Add(portOut); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Test Add duplicate port
	if err := storage.Add(portIn); err == nil {
		t.Fatalf("Expected error, got none")
	} else if err.Error() != "port type 0 with id 1 already exists" {
		t.Fatalf("Expected duplicate error, got %v", err)
	}
	if err := storage.Add(portOut); err == nil {
		t.Fatalf("Expected error, got none")
	} else if err.Error() != "port type 1 with id 2 already exists" {
		t.Fatalf("Expected duplicate error, got %v", err)
	}

	// Test GetBy
	if _, err := storage.GetBy(portIn); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if _, err := storage.GetBy(portOut); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Test GetBy for non-existent port
	invalidPort := model.Port{
		ID:   3,
		Type: model.TypeIN,
	}
	if _, err := storage.GetBy(invalidPort); err == nil {
		t.Fatalf("Expected error for non-existent port, got none")
	} else if err.Error() != "port type 0 with id 3 not found" {
		t.Fatalf("Expected not found error, got %v", err)
	}

	// Test Remove
	if err := storage.Remove(portIn); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if err := storage.Remove(portOut); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Test Remove for non-existent port
	if err := storage.Remove(portIn); err == nil {
		t.Fatalf("Expected error for non-existent port, got none")
	} else if err.Error() != "port type 0 with id 1 not found" {
		t.Fatalf("Expected not found error, got %v", err)
	}
	if err := storage.Remove(portOut); err == nil {
		t.Fatalf("Expected error for non-existent port, got none")
	} else if err.Error() != "port type 1 with id 2 not found" {
		t.Fatalf("Expected not found error, got %v", err)
	}
}
