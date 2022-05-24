package simulator_go_sdk

import (
	"fmt"
	"github.com/deng00/req"
)

type hApi struct {
	api    *net
	url    string
	header map[string]string
}

func newHApi(url string, header map[string]string) *hApi {
	api := newNet(url, header, nil)
	return &hApi{api, url, header}
}

func (h *hApi) Reset(count uint64) error {
	h.url = fmt.Sprintf("%s/v1/reset", h.url)
	h.api.Params = req.Param{
		"blockNumber": count,
	}

	_, err := h.api.Request(PostTy)
	return err
}
