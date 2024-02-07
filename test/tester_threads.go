package test

import "context"

func (test Tester) TestThreads() {
	test.runTest("threads", func(ctx context.Context) error {

		return nil
	})
}
