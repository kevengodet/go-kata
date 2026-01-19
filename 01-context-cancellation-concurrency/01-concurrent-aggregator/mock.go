package main

import (
	"context"
	"errors"
	"time"
)

type MockService struct {
	response string
	timeout  time.Duration // Service call will time out if timeout > 0 (default: 100ms)
	error    string        // If provided, the service will encounter an error
}

type OrderService = MockService
type ProfileService = MockService

type MockOption func(*MockService)

func MockTimeout(timeout time.Duration) MockOption {
	return func(m *MockService) {
		m.timeout = timeout
	}
}

func MockError(error string) MockOption {
	return func(m *MockService) {
		m.error = error
	}
}

func MockResponse(response string) MockOption {
	return func(m *MockService) {
		m.response = response
	}
}

func (m MockService) Fetch(ctx context.Context, id int) (string, error) {
	err := make(chan error)
	res := make(chan string)

	go func() {
		if m.error != "" {
			err <- errors.New(m.error)
		}

		if m.timeout > 0 {
			time.Sleep(m.timeout)
		} else {
			time.Sleep(50 * time.Millisecond)
		}

		res <- m.response
	}()

	select {
	case msg := <-err:
		return "", msg
	case msg := <-res:
		return msg, nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

func mockOrder(opts ...MockOption) *OrderService {
	o := &OrderService{
		response: "Order: 5",
	}

	for _, opt := range opts {
		opt(o)
	}

	return o
}

func mockProfile(opts ...MockOption) *ProfileService {
	o := &ProfileService{
		response: "User: Alice",
	}

	for _, opt := range opts {
		opt(o)
	}

	return o
}
