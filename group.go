package service

import (
	"context"
	"os"

	"golang.org/x/sync/errgroup"
)

type runFn func(ctx context.Context) error

type Group struct {
	rootCtx    context.Context
	rootCancel func()
	startedCh  chan struct{}
	setupCh    chan struct{}

	setups    []runFn
	processes []runFn
}

func NewSignals(sig ...os.Signal) *Group {
	return NewCtx(signalCtx(context.Background(), sig...))
}

func NewCtx(ctx context.Context) *Group {
	ctx, cancel := context.WithCancel(ctx)
	return &Group{
		rootCtx:    ctx,
		rootCancel: cancel,

		startedCh: make(chan struct{}),
		setupCh:   make(chan struct{}),
	}
}

func (g *Group) Setup(fn func(ctx context.Context) error) {
	g.setups = append(g.setups, fn)
}

func (g *Group) Register(fn func(ctx context.Context) error) {
	g.processes = append(g.processes, fn)
}

func (g *Group) Start() error {
	if g.started() {
		return nil
	}

	close(g.startedCh)

	for i := range g.setups {
		fn := g.setups[i]

		if err := fn(g.rootCtx); err != nil {
			return err
		}
	}

	close(g.setupCh)

	errGrp, ctx := errgroup.WithContext(g.rootCtx)

	for i := range g.processes {
		fn := g.processes[i]

		errGrp.Go(func() error {
			return fn(ctx)
		})
	}

	return errGrp.Wait()
}

func (g *Group) Close() error {
	g.rootCancel()
	return nil
}

func (g *Group) started() bool {
	select {
	case <-g.startedCh:
		return true
	default:
		return false
	}
}

func (g *Group) setup() bool {
	select {
	case <-g.setupCh:
		return true
	default:
		return false
	}
}
