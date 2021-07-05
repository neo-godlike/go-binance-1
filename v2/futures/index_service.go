package futures

import (
	"context"
)

type GetIndexInfoService struct {
	c      *Client
	symbol string
}

type IndexInfo struct {
	Symbol        string       `json:"symbol"`
	Time          int64        `json:"time"`
	Component     string       `json:"component"`
	BaseAssetList []*BaseAsset `json:"baseAssetList"`
}

type BaseAsset struct {
	BaseAsset          string `json:"baseAsset"`
	QuoteAsset         string `json:"quoteAsset"`
	WeightInQuantity   string `json:"weightInQuantity"`
	WeightInPercentage string `json:"weightInPercentage"`
}

// Symbol set symbol
func (s *GetIndexInfoService) Symbol(symbol string) *GetIndexInfoService {
	s.symbol = symbol
	return s
}

// Do send request
func (s *GetIndexInfoService) Do(ctx context.Context, opts ...RequestOption) (res []*IndexInfo, err error) {
	r := &request{
		method:   "GET",
		endpoint: "/fapi/v1/indexInfo",
		secType:  secTypeSigned,
	}
	r.setParam("symbol", s.symbol)
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}

	j, err := newJSON(data)
	if err != nil {
		return nil, err
	}
	num := len(j.MustArray())
	res = make([]*IndexInfo, num)

	for i := 0; i < num; i++ {
		item := j.GetIndex(i)
		res[i] = &IndexInfo{
			Symbol:    item.Get("symbol").MustString(),
			Time:      item.Get("time").MustInt64(),
			Component: item.Get("component").MustString(),
		}
		baseAssetListLen := len(item.Get("baseAssetList").MustArray())
		res[i].BaseAssetList = make([]*BaseAsset, baseAssetListLen)
		for k := 0; k < baseAssetListLen; k++ {
			asset := item.Get("baseAssetList").GetIndex(k)
			res[i].BaseAssetList[k] = &BaseAsset{
				BaseAsset:          asset.Get("baseAsset").MustString(),
				QuoteAsset:         asset.Get("quoteAsset").MustString(),
				WeightInQuantity:   asset.Get("weightInQuantity").MustString(),
				WeightInPercentage: asset.Get("weightInPercentage").MustString(),
			}
		}
	}

	return res, nil
}
