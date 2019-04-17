package service

import (
	"github.com/elastos/Elastos.ELA.Elephant.Node/id/blockchain"
	"github.com/elastos/Elastos.ELA.SideChain.ID/service"
	"github.com/elastos/Elastos.ELA/common"
	"github.com/elastos/Elastos.ELA/errors"
	"github.com/elastos/Elastos.ELA/utils/http"
)

type HttpServiceEx struct {
	*service.HttpService
	store *blockchain.IDChainStoreEx
}

func NewHttpServiceEx(service *service.HttpService, idChainStoreEx *blockchain.IDChainStoreEx) *HttpServiceEx {
	return &HttpServiceEx{
		service,
		idChainStoreEx,
	}
}

func (s *HttpServiceEx) GetDidProperty(param http.Params) (interface{}, error) {
	did, ok := param.String("did")
	if !ok {
		return nil, http.NewError(int(errors.InvalidParams), "Did is not a valid parameter")
	}
	_, err := common.Uint168FromAddress(did)
	if err != nil {
		return nil, http.NewError(int(errors.InvalidParams), "Did is not a valid parameter")
	}
	return s.store.GetDidProperty(did), nil
}

func (s *HttpServiceEx) GetDidPropertyByKey(param http.Params) (interface{}, error) {
	did, ok := param.String("did")
	if !ok {
		return nil, http.NewError(int(errors.InvalidParams), "Did is not a valid parameter")
	}
	_, err := common.Uint168FromAddress(did)
	if err != nil {
		return nil, http.NewError(int(errors.InvalidParams), "Did is not a valid parameter")
	}
	key, ok := param.String("key")
	if !ok {
		return nil, http.NewError(int(errors.InvalidParams), "key is not a valid parameter")
	}
	return s.store.GetDidPropertyByKey(did, key), nil
}
