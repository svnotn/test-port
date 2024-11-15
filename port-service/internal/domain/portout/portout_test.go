package portout

import (
	"testing"

	"github.com/svnotn/test-port/port-service/internal/domain"
)

func TestNewPort(t *testing.T) {
	port := New(1)

	if port.State() != domain.Closed {
		t.Errorf("expected port state to be Closed, got %v", port.State())
	}
}

func TestOpenPort(t *testing.T) {
	port := New(1)
	err := port.Open()

	if err != nil {
		t.Fatalf("expected no error when opening the port, got %v", err)
	}

	if port.State() != domain.Opened {
		t.Errorf("expected port state to be Opened, got %v", port.State())
	}

	// Test opening an already opened port
	err = port.Open()
	if err == nil {
		t.Errorf("expected error when opening an already opened port, got nil")
	}
}

func TestClosePort(t *testing.T) {
	port := New(1)

	if err := port.Open(); err != nil {
		t.Fatalf("expected no error when opening the port, got %v", err)
	}
	err := port.Close()

	if err != nil {
		t.Fatalf("expected no error when closing the port, got %v", err)
	}

	if port.State() != domain.Closed {
		t.Errorf("expected port state to be Closed, got %v", port.State())
	}

	// Test closing an already closed port
	err = port.Close()
	if err == nil {
		t.Errorf("expected error when closing an already closed port, got nil")
	}
}

func TestPortWrite(t *testing.T) {
	port := New(1)

	// Пытаемся записать в закрытый порт
	if err := port.Write(42); err == nil {
		t.Fatalf("expected error when writing to a closed port, got none")
	}

	// Открываем порт
	if err := port.Open(); err != nil {
		t.Fatalf("failed to open port: %v", err)
	}

	// Теперь можем писать в порт
	if err := port.Write(42); err != nil {
		t.Fatalf("failed to write to port: %v", err)
	}
}

func TestRunConcurrency(t *testing.T) {
	port := New(1)
	err := port.Open()
	if err != nil {
		t.Fatalf("expected no error when opening the port, got %v", err)
	}

	// Read multiple times in a loop
	for i := 0; i < 100; i++ {
		_, err := port.Read()
		if err != nil {
			t.Fatalf("expected no error when reading from the port, got %v", err)
		}
	}

	if err := port.Close(); err != nil {
		t.Fatalf("expected no error when closing the port, got %v", err)
	}
}
