package errgroup

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func now() string {
	return time.Now().Format("2006-01-2 15:04:05")
}

func Test_GroupTimeout(t *testing.T) {
	t.Log(now(), "begin")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	g, _ := WithContext(ctx)
	g.Go(func() error {
		t.Log(now(), "------>1 begin")

		time.Sleep(1 * time.Second)

		t.Log(now(), "------>1 end")

		return fmt.Errorf("1 err")
		// return nil
	})
	g.Go(func() error {
		t.Log(now(), "------>2 begin")

		time.Sleep(2 * time.Second)

		t.Log(now(), "------>2 end")
		return nil
	})

	err := g.Wait()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(now(), "end")
}
