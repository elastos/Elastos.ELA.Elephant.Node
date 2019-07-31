package types

import (
	"errors"
	"github.com/elastos/Elastos.ELA/common"
	"io"
)

const (
	CMC_ENDPOINT_URL = "https://pro-api.coinmarketcap.com/v1/cryptocurrency/listings/latest?limit=%d&convert=%s"
	HBG_ENDPOINT_URL = "https://api.huobi.pro/market/history/trade?symbol=elabtc"
)

type Status struct {
	Timestamp     string
	Error_code    int
	Error_message string
	Elapsed       int
	Credit_count  int
}

type Price struct {
	Price              float64
	Volume_24h         float64
	Percent_change_1h  float64
	Percent_change_24h float64
	Percent_change_7d  float64
	Market_cap         float64
	Last_updated       string
}

type Quote struct {
	USD, CNY, BTC Price
}

type Plateform struct {
	Id            int64
	Name          string
	Symbol        string
	Slug          string
	Token_Address string
}

type Data struct {
	Id                 int64
	Name               string
	Symbol             string
	Slug               string
	Circulating_supply float64
	Total_supply       float64
	Max_supply         float64
	Date_added         string
	Num_market_pairs   int64
	Tags               []string
	Platform           Plateform
	Cmc_rank           int
	Last_updated       string
	Quote              Quote
}

type CmcResponse struct {
	Status Status
	Data   []Data
}

type Cmc struct {
	Volume_btc             string `json:"24h_volume_btc"`
	Volume_cny             string `json:"24h_volume_cny"`
	Volume_usd             string `json:"24h_volume_usd"`
	Available_supply       string `json:"available_supply"`
	Id                     string `json:"id"`
	Last_updated           string `json:"last_updated"`
	Market_cap_btc         string `json:"market_cap_btc"`
	Market_cap_cny         string `json:"market_cap_cny"`
	Market_cap_usd         string `json:"market_cap_usd"`
	Max_supply             string `json:"max_supply"`
	Name                   string `json:"name"`
	Num_market_pairs       string `json:"num_market_pairs"`
	Percent_change_1h      string `json:"percent_change_1h"`
	Percent_change_24h     string `json:"percent_change_24h"`
	Percent_change_7d      string `json:"percent_change_7d"`
	Platform_symbol        string `json:"platform_symbol"`
	Platform_token_address string `json:"platform_token_address"`
	Price_btc              string `json:"price_btc"`
	Price_cny              string `json:"price_cny"`
	Price_usd              string `json:"price_usd"`
	Rank                   string `json:"rank"`
	Symbol                 string `json:"symbol"`
	Total_supply           string `json:"total_supply"`
}

type Cmcs struct {
	C []Cmc
}

func (cmcs *Cmcs) Serialize(w io.Writer) error {
	err := common.WriteVarUint(w, uint64(len(cmcs.C)))
	if err != nil {
		return errors.New("[TransactionHistory], Address serialize failed.")
	}
	for _, cmc := range cmcs.C {
		err := common.WriteVarString(w, cmc.Volume_btc)
		if err != nil {
			return err
		}
		err = common.WriteVarString(w, cmc.Volume_cny)
		if err != nil {
			return err
		}
		err = common.WriteVarString(w, cmc.Volume_usd)
		if err != nil {
			return err
		}
		err = common.WriteVarString(w, cmc.Available_supply)
		if err != nil {
			return err
		}
		err = common.WriteVarString(w, cmc.Id)
		if err != nil {
			return err
		}
		err = common.WriteVarString(w, cmc.Last_updated)
		if err != nil {
			return err
		}
		err = common.WriteVarString(w, cmc.Market_cap_btc)
		if err != nil {
			return err
		}
		err = common.WriteVarString(w, cmc.Market_cap_cny)
		if err != nil {
			return err
		}
		err = common.WriteVarString(w, cmc.Market_cap_usd)
		if err != nil {
			return err
		}
		err = common.WriteVarString(w, cmc.Max_supply)
		if err != nil {
			return err
		}
		err = common.WriteVarString(w, cmc.Name)
		if err != nil {
			return err
		}
		err = common.WriteVarString(w, cmc.Num_market_pairs)
		if err != nil {
			return err
		}
		err = common.WriteVarString(w, cmc.Percent_change_1h)
		if err != nil {
			return err
		}
		err = common.WriteVarString(w, cmc.Percent_change_24h)
		if err != nil {
			return err
		}
		err = common.WriteVarString(w, cmc.Percent_change_7d)
		if err != nil {
			return err
		}
		err = common.WriteVarString(w, cmc.Platform_symbol)
		if err != nil {
			return err
		}
		err = common.WriteVarString(w, cmc.Platform_token_address)
		if err != nil {
			return err
		}
		err = common.WriteVarString(w, cmc.Price_btc)
		if err != nil {
			return err
		}
		err = common.WriteVarString(w, cmc.Price_cny)
		if err != nil {
			return err
		}
		err = common.WriteVarString(w, cmc.Price_usd)
		if err != nil {
			return err
		}
		err = common.WriteVarString(w, cmc.Rank)
		if err != nil {
			return err
		}
		err = common.WriteVarString(w, cmc.Symbol)
		if err != nil {
			return err
		}
		err = common.WriteVarString(w, cmc.Total_supply)
		if err != nil {
			return err
		}
	}
	return nil
}

func (cmcs *Cmcs) Deserialize(r io.Reader) error {
	n, err := common.ReadVarUint(r, 0)
	if err != nil {
		return err
	}
	for i := 0; i < int(n); i++ {
		c := Cmc{}
		c.Volume_btc, err = common.ReadVarString(r)
		if err != nil {
			return err
		}
		c.Volume_cny, err = common.ReadVarString(r)
		if err != nil {
			return err
		}
		c.Volume_usd, err = common.ReadVarString(r)
		if err != nil {
			return err
		}
		c.Available_supply, err = common.ReadVarString(r)
		if err != nil {
			return err
		}
		c.Id, err = common.ReadVarString(r)
		if err != nil {
			return err
		}
		c.Last_updated, err = common.ReadVarString(r)
		if err != nil {
			return err
		}
		c.Market_cap_btc, err = common.ReadVarString(r)
		if err != nil {
			return err
		}
		c.Market_cap_cny, err = common.ReadVarString(r)
		if err != nil {
			return err
		}
		c.Market_cap_usd, err = common.ReadVarString(r)
		if err != nil {
			return err
		}
		c.Max_supply, err = common.ReadVarString(r)
		if err != nil {
			return err
		}
		c.Name, err = common.ReadVarString(r)
		if err != nil {
			return err
		}
		c.Num_market_pairs, err = common.ReadVarString(r)
		if err != nil {
			return err
		}
		c.Percent_change_1h, err = common.ReadVarString(r)
		if err != nil {
			return err
		}
		c.Percent_change_24h, err = common.ReadVarString(r)
		if err != nil {
			return err
		}
		c.Percent_change_7d, err = common.ReadVarString(r)
		if err != nil {
			return err
		}
		c.Platform_symbol, err = common.ReadVarString(r)
		if err != nil {
			return err
		}
		c.Platform_token_address, err = common.ReadVarString(r)
		if err != nil {
			return err
		}
		c.Price_btc, err = common.ReadVarString(r)
		if err != nil {
			return err
		}
		c.Price_cny, err = common.ReadVarString(r)
		if err != nil {
			return err
		}
		c.Price_usd, err = common.ReadVarString(r)
		if err != nil {
			return err
		}
		c.Rank, err = common.ReadVarString(r)
		if err != nil {
			return err
		}
		c.Symbol, err = common.ReadVarString(r)
		if err != nil {
			return err
		}
		c.Total_supply, err = common.ReadVarString(r)
		if err != nil {
			return err
		}
		cmcs.C = append(cmcs.C, c)
	}
	return nil
}
