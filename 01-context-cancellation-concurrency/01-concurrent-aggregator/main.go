package main

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	users := New(
		&ProfileService{},
		&OrderService{},
		WithTimeout(time.Second),
		WithLog(slog.DiscardHandler),
	)

	s, err := users.Aggregate(context.Background(), 1)
	if err != nil {
		slog.Error("Error", "Yo!", err)
		return
	}
	fmt.Println(s)
}

type UserAggregator struct {
	profiles *ProfileService
	orders   *OrderService
	timeout  time.Duration
	log      slog.Handler
}

func (u UserAggregator) Aggregate(ctx context.Context, id int) (string, error) {
	var profile, order string
	var err error

	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	group, gCtx := errgroup.WithContext(ctx)

	group.Go(func() error {
		order, err = u.orders.Fetch(gCtx, id)
		return err
	})

	group.Go(func() error {
		profile, err = u.profiles.Fetch(gCtx, id)
		return err
	})

	if err := group.Wait(); err != nil {
		return "", fmt.Errorf("failed to aggregate user: %w", err)
	}

	return profile + " | " + order, nil
}

type Option func(*UserAggregator)

func WithTimeout(timeout time.Duration) Option {
	return func(u *UserAggregator) {
		u.timeout = timeout
	}
}

func WithLog(log slog.Handler) Option {
	return func(u *UserAggregator) {
		u.log = log
	}
}

func New(
	profiles *ProfileService,
	orders *OrderService,
	opts ...Option,
) *UserAggregator {
	u := &UserAggregator{
		profiles: profiles,
		orders:   orders,
		timeout:  time.Second,
	}

	// Apply options
	for _, opt := range opts {
		opt(u)
	}

	if u.log == nil {
		u.log = slog.DiscardHandler
	}

	return u
}

type ProfileService struct{}

func (p ProfileService) Fetch(ctx context.Context, id int) (string, error) {
	err := make(chan error)
	res := make(chan string)

	go func() {
		//err <- errors.New("Boo!")
		time.Sleep(500 * time.Millisecond)
		res <- "Name: Alice"
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

type OrderService struct{}

func (p OrderService) Fetch(ctx context.Context, id int) (string, error) {
	err := make(chan error)
	res := make(chan string)

	go func() {
		//time.Sleep(500 * time.Millisecond)
		time.Sleep(10 * time.Second)
		res <- "Order: 5"
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
