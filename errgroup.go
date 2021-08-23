// 对golang.org/x/sync/errgroup调用的简单封装，让其可以自动超时
// 用法和errgroup一样
package errgroup

import (
	"context"
	"time"

	"golang.org/x/sync/errgroup"
)

type Group struct {
	eg  *errgroup.Group
	ctx context.Context
}

func New() *Group {
	return &Group{eg: &errgroup.Group{}, ctx: context.Background()}
}

func WithContext(ctx context.Context) (*Group, context.Context) {
	g, ctx := errgroup.WithContext(ctx)
	return &Group{eg: g, ctx: ctx}, ctx
}
func (g *Group) Go(f func() error) {
	g.eg.Go(func() error {
		function := func() <-chan error {
			cherr := make(chan error, 1)
			go func() {
				defer func() {
					time.Sleep(10 * time.Millisecond) // 防止chan关闭过快, 导致外部读到的错误都是空的，概率很小。
					close(cherr)
				}()
				cherr <- f()
			}()
			return cherr
		}

		select {
		case <-g.ctx.Done():
			return g.ctx.Err()
		case err := <-function():
			if err != nil {
				return err
			}
		}
		return nil
	})
}
func (g *Group) Wait() error {
	return g.eg.Wait()
}
