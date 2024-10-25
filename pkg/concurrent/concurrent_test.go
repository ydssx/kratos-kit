package concurrent

import (
	"context"
	"errors"
	"log"
	"os"
	"testing"
	"time"

	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)

	os.Exit(m.Run())
}

func TestGroup_Run(t *testing.T) {
	// Test case 1: Limiting concurrent execution to 2
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	g := NewGroup(ctx, WithSemaphore(3), WithFastFail(true))
	err1 := errors.New("Error 1")
	err2 := errors.New("Error 2")
	err3 := errors.New("Error 3")
	log.Print("start")
	gErr := g.Run(
		func() error {
			time.Sleep(3 * time.Second)
			log.Print(err1)
			return err1
		},
		func() error {
			time.Sleep(3 * time.Second)
			log.Print(err2)
			return err2
		},
		func() error {
			time.Sleep(5 * time.Second)
			log.Print(err3)
			return err3
		},
	)
	log.Print("result: ", gErr)
	if gErr == nil {
		t.Errorf("Expected gErr to be %v, got %v", err1, gErr)
	}
	time.Sleep(3 * time.Second)
}
