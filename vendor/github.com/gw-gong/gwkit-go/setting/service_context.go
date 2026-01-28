package setting

import "context"

var serviceContext = context.Background()

func ResetServiceContext(ctx context.Context) {
	serviceContext = ctx
}

func GetServiceContext() context.Context {
	return serviceContext
}
