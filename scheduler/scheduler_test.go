package scheduler_test

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/platforma-dev/platforma/application"
	"github.com/platforma-dev/platforma/scheduler"
)

func TestSuccessRun(t *testing.T) {
	t.Parallel()

	buf := bytes.Buffer{}
	s := scheduler.New(1*time.Second, application.RunnerFunc(func(ctx context.Context) error {
		buf.WriteString("1")
		return nil
	}))

	go s.Run(context.TODO())

	time.Sleep(3500 * time.Millisecond)

	if buf.String() != "111" {
		t.Errorf("wrong buffer content. expected %v, got %v", "111", buf.String())
	}
}

func TestErrorRun(t *testing.T) {
	t.Parallel()

	buf := bytes.Buffer{}
	s := scheduler.New(1*time.Second, application.RunnerFunc(func(ctx context.Context) error {
		buf.WriteString("1")
		return errors.New("some error")
	}))

	go s.Run(context.TODO())

	time.Sleep(3500 * time.Millisecond)

	if buf.String() != "111" {
		t.Errorf("wrong buffer content. expected %v, got %v", "111", buf.String())
	}
}

func TestContextDecline(t *testing.T) {
	t.Parallel()

	buf := bytes.Buffer{}
	s := scheduler.New(1*time.Second, application.RunnerFunc(func(ctx context.Context) error {
		buf.WriteString("1")
		return nil
	}))

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(3*time.Second + 10*time.Millisecond)
		cancel()
	}()

	err := s.Run(ctx)

	if buf.String() != "111" {
		t.Errorf("wrong buffer content. expected %v, got %v", "111", buf.String())
	}

	if err == nil {
		t.Error("expected error, got nil")
	}
}
