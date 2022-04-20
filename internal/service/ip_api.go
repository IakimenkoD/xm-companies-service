package service

import "context"

//go:generate minimock -i IpApi -g -o ip_api_mock.go

// IpChecker interface provides methods to interacts with IpChecker service.
type IpChecker interface {
	GetUserLocation(ctx context.Context, ip string) (location string, err error)
}
