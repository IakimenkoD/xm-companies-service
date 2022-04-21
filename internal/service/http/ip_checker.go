package http

import (
	"context"
	"github.com/IakimenkoD/xm-companies-service/internal/config"
	ierr "github.com/IakimenkoD/xm-companies-service/internal/errors"
	"github.com/IakimenkoD/xm-companies-service/internal/service"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
)

const (
	undefined = "Undefined"
	localhost = "127.0.0.1"
)

func NewIpChecker(conf *config.Config, log *zap.Logger) service.IpChecker {
	return &IpApi{
		client: &http.Client{
			Timeout: conf.API.ReadTimeout,
		},
		Url: conf.IpApi.Address,
		log: log,
	}
}

type IpApi struct {
	client *http.Client
	log    *zap.Logger

	Url string
}

// GetUserLocation gets location from IpApi service by ip.
func (i *IpApi) GetUserLocation(ctx context.Context, ip string) (location string, err error) {
	//debug
	//if ip == localhost {
	//	return "CY", err
	//}

	reqUrl := i.Url + ip + "/country/"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqUrl, nil)
	req.Header.Set("User-Agent", "ipapi.co/#go-v1.5")

	i.log.Debug("location request",
		zap.String("ip", ip),
		zap.String("request", reqUrl),
	)

	resp, err := i.client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	res := string(body)

	if resp.StatusCode != http.StatusOK {
		i.log.Error("ip_api service responded non 200 status code",
			zap.Int("status_code", resp.StatusCode),
			zap.String("body", res))
		return "", errors.New("unexpected response from external service")
	}

	if res == undefined {
		return "", ierr.UnknownLocation
	}

	return res, nil
}
