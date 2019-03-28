package blockchain

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	common2 "github.com/elastos/Elastos.ELA.Elephant.Node/common"
	"github.com/elastos/Elastos.ELA.Elephant.Node/ela/core/types"
	"github.com/elastos/Elastos.ELA/common"
	"github.com/elastos/Elastos.ELA/common/log"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"
)

var i *int = common2.NewInt(0)

func (c ChainStoreExtend) begin() {
	c.NewBatch()
}

func (c ChainStoreExtend) commit() {
	c.BatchCommit()
}

func (c ChainStoreExtend) rollback() {
	c.rollback()
}

// key: DataEntryPrefix + height + address
// value: serialized history
func (c ChainStoreExtend) persistTransactionHistory(txhs []types.TransactionHistory) error {
	c.begin()
	for _, txh := range txhs {
		err := c.doPersistTransactionHistory(txh)
		if err != nil {
			c.rollback()
			log.Fatal("Error persist transaction history")
			os.Exit(-1)
		}
	}
	c.commit()
	return nil
}

func (c ChainStoreExtend) doPersistTransactionHistory(history types.TransactionHistory) error {
	key := new(bytes.Buffer)
	key.WriteByte(byte(DataTxHistoryPrefix))
	err := common.WriteVarString(key, history.Address)
	if err != nil {
		return err
	}
	err = common.WriteUint64(key, history.Height)
	if err != nil {
		return err
	}

	value := new(bytes.Buffer)
	history.Serialize(value)
	c.BatchPut(key.Bytes(), value.Bytes())
	return nil
}

func (c ChainStoreExtend) initCmc() {
	c.AddFunc("@every "+common2.Conf.Cmc.Inteval, c.renewCmcPrice)
	c.Start()
}

func (c *ChainStoreExtend) renewCmcPrice() {
	c.mu.Lock()
	defer c.mu.Unlock()
	cmcResponseUSD, err := fetchPrice(*i, "USD")
	if err != nil {
		log.Warnf("Error in cmc price %s", err.Error())
		return
	}
	cmcResponseCNY, err := fetchPrice(*i, "CNY")
	if err != nil {
		log.Warnf("Error in cmc price %s", err.Error())
		return
	}
	cmcResponseBTC, err := fetchPrice(*i, "BTC")
	if err != nil {
		log.Warnf("Error in cmc price %s", err.Error())
		return
	}
	cmcResponseBGX, err := fetchBGXPrice()
	if err != nil {
		log.Warnf("Error in bgx price %s", err.Error())
		return
	}
	c.begin()
	err = c.saveToDb(cmcResponseUSD, cmcResponseCNY, cmcResponseBTC, cmcResponseBGX)
	if err != nil {
		log.Warnf("Error in cmc price %s", err.Error())
		c.rollback()
		return
	}
	c.commit()
	if *i == len(common2.Conf.Cmc.ApiKey)-1 {
		i = common2.NewInt(0)
	} else {
		i = common2.NewInt(*i + 1)
	}
}

func (c ChainStoreExtend) saveToDb(cmcResponseUSD, cmcResponseCNY, cmcResponseBTC, cmcResponseBGX types.CmcResponse) error {
	data := cmcResponseUSD.Data
	cmcs := make([]types.Cmc, 0)
	for i := 0; i < len(data); i++ {
		cmc := types.Cmc{}
		cmc.Id = strconv.Itoa(int(data[i].Id))
		cmc.Name = data[i].Name
		cmc.Symbol = data[i].Symbol
		cmc.Rank = strconv.Itoa(data[i].Cmc_rank)
		cmc.Price_usd = strconv.FormatFloat(data[i].Quote.USD.Price, 'f', 8, 64)
		cmc.Price_cny = strconv.FormatFloat(cmcResponseCNY.Data[i].Quote.CNY.Price, 'f', 8, 64)
		cmc.Price_btc = strconv.FormatFloat(cmcResponseBTC.Data[i].Quote.BTC.Price, 'f', 8, 64)
		cmc.Volume_usd = strconv.FormatFloat(data[i].Quote.USD.Volume_24h, 'f', 8, 64)
		cmc.Market_cap_usd = strconv.FormatFloat(data[i].Quote.USD.Market_cap, 'f', 8, 64)
		cmc.Available_supply = strconv.FormatFloat(data[i].Circulating_supply, 'f', 8, 64)
		cmc.Total_supply = strconv.FormatFloat(data[i].Total_supply, 'f', 8, 64)
		cmc.Max_supply = strconv.FormatFloat(data[i].Max_supply, 'f', 8, 64)
		cmc.Percent_change_1h = strconv.FormatFloat(data[i].Quote.USD.Percent_change_1h, 'f', 8, 64)
		cmc.Percent_change_24h = strconv.FormatFloat(data[i].Quote.USD.Percent_change_24h, 'f', 8, 64)
		cmc.Percent_change_7d = strconv.FormatFloat(data[i].Quote.USD.Percent_change_7d, 'f', 8, 64)
		cmc.Last_updated = data[i].Quote.USD.Last_updated
		cmc.Volume_btc = strconv.FormatFloat(cmcResponseBTC.Data[i].Quote.BTC.Volume_24h, 'f', 8, 64)
		cmc.Market_cap_btc = strconv.FormatFloat(cmcResponseBTC.Data[i].Quote.BTC.Market_cap, 'f', 8, 64)
		cmc.Volume_cny = strconv.FormatFloat(cmcResponseCNY.Data[i].Quote.CNY.Volume_24h, 'f', 8, 64)
		cmc.Market_cap_cny = strconv.FormatFloat(cmcResponseCNY.Data[i].Quote.CNY.Market_cap, 'f', 8, 64)
		cmc.Platform_symbol = data[i].Platform.Symbol
		cmc.Platform_token_address = data[i].Platform.Token_Address
		cmc.Num_market_pairs = strconv.Itoa(int(data[i].Num_market_pairs))
		cmcs = append(cmcs, cmc)
		// put price that not in the cmc at rank 100
		if i == 99 {
			bgxcmc := types.Cmc{}
			bgxcmc.Id = strconv.Itoa(int(cmcResponseBGX.Data[0].Id))
			bgxcmc.Name = cmcResponseBGX.Data[0].Name
			bgxcmc.Symbol = cmcResponseBGX.Data[0].Symbol
			bgxcmc.Rank = strconv.Itoa(cmcResponseBGX.Data[0].Cmc_rank)
			bgxcmc.Price_usd = strconv.FormatFloat(cmcResponseBGX.Data[0].Quote.USD.Price, 'f', 8, 64)
			bgxcmc.Price_cny = strconv.FormatFloat(cmcResponseBGX.Data[0].Quote.CNY.Price, 'f', 8, 64)
			bgxcmc.Price_btc = strconv.FormatFloat(cmcResponseBGX.Data[0].Quote.BTC.Price, 'f', 8, 64)
			bgxcmc.Volume_usd = strconv.FormatFloat(cmcResponseBGX.Data[0].Quote.USD.Volume_24h, 'f', 8, 64)
			bgxcmc.Market_cap_usd = strconv.FormatFloat(cmcResponseBGX.Data[0].Quote.USD.Market_cap, 'f', 8, 64)
			bgxcmc.Available_supply = strconv.FormatFloat(cmcResponseBGX.Data[0].Circulating_supply, 'f', 8, 64)
			bgxcmc.Total_supply = strconv.FormatFloat(cmcResponseBGX.Data[0].Total_supply, 'f', 8, 64)
			bgxcmc.Max_supply = strconv.FormatFloat(cmcResponseBGX.Data[0].Max_supply, 'f', 8, 64)
			bgxcmc.Percent_change_1h = strconv.FormatFloat(cmcResponseBGX.Data[0].Quote.USD.Percent_change_1h, 'f', 8, 64)
			bgxcmc.Percent_change_24h = strconv.FormatFloat(cmcResponseBGX.Data[0].Quote.USD.Percent_change_24h, 'f', 8, 64)
			bgxcmc.Percent_change_7d = strconv.FormatFloat(cmcResponseBGX.Data[0].Quote.USD.Percent_change_7d, 'f', 8, 64)
			bgxcmc.Last_updated = data[i].Quote.USD.Last_updated
			bgxcmc.Volume_btc = strconv.FormatFloat(cmcResponseBGX.Data[0].Quote.BTC.Volume_24h, 'f', 8, 64)
			bgxcmc.Market_cap_btc = strconv.FormatFloat(cmcResponseBGX.Data[0].Quote.BTC.Market_cap, 'f', 8, 64)
			bgxcmc.Volume_cny = strconv.FormatFloat(cmcResponseBGX.Data[0].Quote.CNY.Volume_24h, 'f', 8, 64)
			bgxcmc.Market_cap_cny = strconv.FormatFloat(cmcResponseBGX.Data[0].Quote.CNY.Market_cap, 'f', 8, 64)
			bgxcmc.Platform_symbol = cmcResponseBGX.Data[0].Platform.Symbol
			bgxcmc.Platform_token_address = cmcResponseBGX.Data[0].Platform.Token_Address
			bgxcmc.Num_market_pairs = strconv.Itoa(int(cmcResponseBGX.Data[0].Num_market_pairs))
			cmcs = append(cmcs, bgxcmc)
		}
	}
	c.persistCmc(types.Cmcs{cmcs})
	return nil
}

func (c ChainStoreExtend) persistCmc(cmc types.Cmcs) error {
	println("Persist CMC")
	key := new(bytes.Buffer)
	key.WriteByte(byte(DataCmcPrefix))
	common.WriteVarString(key,"CMC")
	value := new(bytes.Buffer)
	cmc.Serialize(value)
	c.BatchPut(key.Bytes(), value.Bytes())
	return nil
}

func fetchBGXPrice() (types.CmcResponse, error) {
	url := fmt.Sprintf(types.BGX_ENDPOINT_URL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return types.CmcResponse{}, err
	}
	req.Header["Accept-Language"] = []string{"*"}
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return types.CmcResponse{}, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return types.CmcResponse{}, err
	}
	id := 5000
	name := "BIT GAME EXCHANGE"
	symbol := "BGX"
	slug := "BIG GAME"
	circulating_supply := 2000000000
	total_supply := 5000000000
	max_supply := 0
	date_added := "2019-02-27T13:53:00.000Z"
	num_market_pairs := 1
	platform_symbol := "ETH"
	platform_token_address := "0xbf3f09e4eba5f7805e5fac0ee09fd6ee8eebe4cb"
	bgxRespMap := new(map[string]interface{})
	err = json.Unmarshal(body, &bgxRespMap)
	if err != nil {
		return types.CmcResponse{}, err
	}
	data, ok := (*bgxRespMap)["data"].(map[string]interface{})
	if !ok {
		return types.CmcResponse{}, errors.New("BGX Price Error")
	}
	rate, ok := data["rate"].(map[string]interface{})
	if !ok {
		return types.CmcResponse{}, errors.New("BGX Price Error")
	}
	enUS, ok := rate["en_US"].(map[string]interface{})
	if !ok {
		return types.CmcResponse{}, errors.New("BGX Price Error")
	}
	zhCN, ok := rate["zh_CN"].(map[string]interface{})
	if !ok {
		return types.CmcResponse{}, errors.New("BGX Price Error")
	}
	bgxUs, ok1 := enUS["BGX"].(float64)
	btcUs, ok2 := enUS["BTC"].(float64)
	bgxCn, ok3 := zhCN["BGX"].(float64)
	if !(ok1 && ok2 && ok3) {
		return types.CmcResponse{}, errors.New("BGX Price Error")
	}
	bgxbtc := math.Round(bgxUs/btcUs*100000000) / 100000000
	now := time.Now().Format("2006-01-02T15:04:05.000Z")
	return types.CmcResponse{
		Status: types.Status{
			Timestamp:     now,
			Error_code:    0,
			Error_message: "",
			Elapsed:       0,
			Credit_count:  0,
		},
		Data: []types.Data{
			types.Data{
				Id:                 int64(id),
				Name:               name,
				Symbol:             symbol,
				Slug:               slug,
				Circulating_supply: float64(circulating_supply),
				Total_supply:       float64(total_supply),
				Max_supply:         float64(max_supply),
				Date_added:         date_added,
				Num_market_pairs:   int64(num_market_pairs),
				Tags:               nil,
				Platform: types.Plateform{
					Symbol:        platform_symbol,
					Token_Address: platform_token_address,
				},
				Cmc_rank:     0,
				Last_updated: date_added,
				Quote: types.Quote{
					CNY: types.Price{
						Price:              bgxCn,
						Volume_24h:         float64(0),
						Percent_change_1h:  float64(0),
						Percent_change_24h: float64(0),
						Percent_change_7d:  float64(0),
						Market_cap:         float64(0),
						Last_updated:       now,
					},
					USD: types.Price{
						Price:              bgxUs,
						Volume_24h:         float64(0),
						Percent_change_1h:  float64(0),
						Percent_change_24h: float64(0),
						Percent_change_7d:  float64(0),
						Market_cap:         float64(0),
						Last_updated:       now,
					},
					BTC: types.Price{
						Price:              bgxbtc,
						Volume_24h:         float64(0),
						Percent_change_1h:  float64(0),
						Percent_change_24h: float64(0),
						Percent_change_7d:  float64(0),
						Market_cap:         float64(0),
						Last_updated:       now,
					},
				},
			},
		},
	}, nil
}

func fetchPrice(i int, curr string) (types.CmcResponse, error) {
	url := fmt.Sprintf(types.CMC_ENDPOINT_URL, 2000, curr)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return types.CmcResponse{}, err
	}
	req.Header = map[string][]string{
		"X-CMC_PRO_API_KEY": []string{common2.Conf.Cmc.ApiKey[i]},
	}
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return types.CmcResponse{}, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return types.CmcResponse{}, err
	}
	cmcResp := types.CmcResponse{}
	err = json.Unmarshal(body, &cmcResp)
	if err != nil {
		return types.CmcResponse{}, err
	}
	return cmcResp, nil
}
