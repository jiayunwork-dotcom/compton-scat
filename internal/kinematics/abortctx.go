package kinematics

import "context"

func abortedContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	return ctx, cancel
}

func abortKineContext() error {
	ctx, cancel := abortedContext()
	defer cancel()
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}
