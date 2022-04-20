package service

import "context"

//go:generate minimock -i IpApi -g -o ip_api_mock.go

// IpApi interface provides methods to interacts with IpApi service.
type IpApi interface {
	GetUserLocation(ctx context.Context, ip string) (location string, err error)
}
