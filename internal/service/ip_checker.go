package service

import "context"

//go:generate minimock -i IpChecker -g -o ip_checker_mock.go

// IpChecker interface provides methods to interacts with IpChecker service.
type IpChecker interface {
	GetUserLocation(ctx context.Context, ip string) (location string, err error)
}
