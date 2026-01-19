package main

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

func main() {
	dashboard := New(
		mockProfile(),
		mockOrder(),
		WithTimeout(time.Second),
		WithLog(slog.DiscardHandler),
	)

	s, err := dashboard.Aggregate(context.Background(), 1)
	if err != nil {
		slog.Error("Error", "", err)
		return
	}
	fmt.Println(s)
}
