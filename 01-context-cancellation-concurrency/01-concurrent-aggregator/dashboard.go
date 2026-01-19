package main

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"golang.org/x/sync/errgroup"
)

type Dashboard struct {
	profiles *ProfileService
	orders   *OrderService
	timeout  time.Duration
	log      slog.Handler
}

func (u Dashboard) Aggregate(ctx context.Context, id int) (string, error) {
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

type Option func(*Dashboard)

func WithTimeout(timeout time.Duration) Option {
	return func(u *Dashboard) {
		u.timeout = timeout
	}
}

func WithLog(log slog.Handler) Option {
	return func(u *Dashboard) {
		u.log = log
	}
}

func New(
	profiles *ProfileService,
	orders *OrderService,
	opts ...Option,
) *Dashboard {
	u := &Dashboard{
		profiles: profiles,
		orders:   orders,
		timeout:  time.Second,
	}

	for _, opt := range opts {
		opt(u)
	}

	if u.log == nil {
		u.log = slog.DiscardHandler
	}

	return u
}
