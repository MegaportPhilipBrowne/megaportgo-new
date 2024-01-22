package megaport

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/megaport/megaportgo/mega_err"
	"github.com/megaport/megaportgo/shared"
	"github.com/megaport/megaportgo/types"
)

// PortService is an interface for interfacing with the Port endpoints
// of the Megaport API.

type PortService interface {
	BuyPort(ctx context.Context, req *BuyPortRequest) (*types.PortOrderConfirmation, error)
	BuySinglePort(ctx context.Context, req *BuySinglePortRequest) (*types.PortOrderConfirmation, error)
	BuyLAGPort(ctx context.Context, req *BuyLAGPortRequest) (*types.PortOrderConfirmation, error)
	ListPorts(ctx context.Context) ([]*types.Port, error)
	ModifyPort(ctx context.Context, req *ModifyPortRequest) (*ModifyPortResponse, error)
	DeletePort(ctx context.Context, req *DeletePortRequest) (*DeletePortResponse, error)
}

type ParsedProductsResponse struct {
	Message string        `json:"message"`
	Terms   string        `json:"terms"`
	Data    []interface{} `json:"data"`
}

// PortServiceOp handles communication with Port methods of the Megaport API.
type PortServiceOp struct {
	Client           *Client
	ProductServiceOp *ProductServiceOp
}

type BuyPortRequest struct {
	Name       string
	Term       int
	PortSpeed  int
	LocationId int
	Market     string
	IsLag      bool
	LagCount   int
	IsPrivate  bool
}

type BuySinglePortRequest struct {
	Name       string
	Term       int
	PortSpeed  int
	LocationId int
	Market     string
	IsPrivate  bool
}

type BuyLAGPortRequest struct {
	Name       string
	Term       int
	PortSpeed  int
	LocationId int
	Market     string
	LagCount   int
	IsPrivate  bool
}

type ModifyPortRequest struct {
	PortID                string
	Name                  string
	MarketplaceVisibility bool
	CostCentre            string
}

type ModifyPortResponse struct {
	IsUpdated bool
}

type DeletePortRequest struct {
	PortID    string
	DeleteNow bool
}

type DeletePortResponse struct{}

func NewPortServiceOp(c *Client) *PortServiceOp {
	return &PortServiceOp{
		Client: c,
	}
}

func (svc *PortServiceOp) BuyPort(ctx context.Context, req *BuyPortRequest) (*types.PortOrderConfirmation, error) {
	var buyOrder []types.PortOrder
	if req.Term != 1 && req.Term != 12 && req.Term != 24 && req.Term != 36 {
		return nil, errors.New(mega_err.ERR_TERM_NOT_VALID)
	}
	if req.IsLag {
		buyOrder = []types.PortOrder{
			{
				Name:                  req.Name,
				Term:                  req.Term,
				ProductType:           "MEGAPORT",
				PortSpeed:             req.PortSpeed,
				LocationID:            req.LocationId,
				CreateDate:            shared.GetCurrentTimestamp(),
				Virtual:               false,
				Market:                req.Market,
				LagPortCount:          req.LagCount,
				MarketplaceVisibility: !req.IsPrivate,
			},
		}
	} else {
		buyOrder = []types.PortOrder{
			{
				Name:                  req.Name,
				Term:                  req.Term,
				ProductType:           "MEGAPORT",
				PortSpeed:             req.PortSpeed,
				LocationID:            req.LocationId,
				CreateDate:            shared.GetCurrentTimestamp(),
				Virtual:               false,
				Market:                req.Market,
				MarketplaceVisibility: !req.IsPrivate,
			},
		}
	}

	requestBody, _ := json.Marshal(buyOrder)
	responseBody, responseError := svc.Client.ProductService.ExecuteOrder(ctx, &requestBody)
	if responseError != nil {
		return nil, responseError
	}
	orderInfo := types.PortOrderResponse{}
	unmarshalErr := json.Unmarshal(*responseBody, &orderInfo)
	if unmarshalErr != nil {
		return nil, unmarshalErr
	}
	return &types.PortOrderConfirmation{
		TechnicalServiceUID: orderInfo.Data[0].TechnicalServiceUID,
	}, nil
}

func (svc *PortServiceOp) BuySinglePort(ctx context.Context, req *BuySinglePortRequest) (*types.PortOrderConfirmation, error) {
	return svc.BuyPort(ctx, &BuyPortRequest{
		Name:       req.Name,
		Term:       req.Term,
		PortSpeed:  req.PortSpeed,
		LocationId: req.LocationId,
		Market:     req.Market,
		IsLag:      false,
		LagCount:   0,
		IsPrivate:  req.IsPrivate,
	})
}

func (svc *PortServiceOp) BuyLAGPort(ctx context.Context, req *BuyLAGPortRequest) (*types.PortOrderConfirmation, error) {
	return svc.BuyPort(ctx, &BuyPortRequest{
		Name:       req.Name,
		Term:       req.Term,
		PortSpeed:  req.PortSpeed,
		LocationId: req.LocationId,
		Market:     req.Market,
		IsLag:      true,
		LagCount:   req.LagCount,
		IsPrivate:  req.IsPrivate,
	})
}

func (svc *PortServiceOp) ListPorts(ctx context.Context) ([]*types.Port, error) {
	path := "/v2/products"
	url := svc.Client.BaseURL.JoinPath(path).String()
	req, err := svc.Client.NewRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	response, err := svc.Client.Do(ctx, req, nil)

	defer response.Body.Close() // nolint

	isError, parsedError := svc.Client.IsErrorResponse(response, &err, 200)
	if isError {
		return nil, parsedError
	}
	body, fileErr := io.ReadAll(response.Body)

	if fileErr != nil {
		return nil, fileErr
	}

	parsed := ParsedProductsResponse{}

	unmarshalErr := json.Unmarshal(body, &parsed)

	if unmarshalErr != nil {
		return nil, unmarshalErr
	}

	ports := []*types.Port{}

	for _, unmarshaledData := range parsed.Data {
		// The products query response will likely contain non-port objects.  As a result
		// we need to initially Unmarshal as ParsedProductsResponse so that we may iterate
		// over the entries in Data then re-Marshal those entries so that we may Unmarshal
		// them as Port (and `continue` where that doesn't work).  We could write a custom
		// deserializer to avoid this but that is a lot of work for a performance
		// optimization which is likely irrelevant in practice.
		// Unfortunately I know of no better (maintainable) method of making this work.
		remarshaled, err := json.Marshal(unmarshaledData)
		if err != nil {
			svc.Client.Logger.Debug(fmt.Sprintf("Could not remarshal %v as port.", err.Error()))
			continue
		}
		port := types.Port{}
		unmarshalErr = json.Unmarshal(remarshaled, &port)
		if unmarshalErr != nil {
			svc.Client.Logger.Debug(fmt.Sprintf("Could not unmarshal %v as port.", unmarshalErr.Error()))
			continue
		}
		ports = append(ports, &port)
	}
	return ports, nil
}

func (svc *PortServiceOp) ModifyPort(ctx context.Context, req *ModifyPortRequest) (*ModifyPortResponse, error) {
	modifyRes, err := svc.Client.ProductService.ModifyProduct(ctx, &ModifyProductRequest{
		ProductID:             req.PortID,
		ProductType:           types.PRODUCT_MEGAPORT,
		Name:                  req.Name,
		CostCentre:            req.CostCentre,
		MarketplaceVisibility: req.MarketplaceVisibility,
	})
	if err != nil {
		return nil, err
	}
	return &ModifyPortResponse{
		IsUpdated: modifyRes.IsUpdated,
	}, nil
}

func (svc *PortServiceOp) DeletePort(ctx context.Context, req *DeletePortRequest) (*DeletePortResponse, error) {
	_, err := svc.Client.ProductService.DeleteProduct(ctx, &DeleteProductRequest{
		ProductID: req.PortID,
		DeleteNow: req.DeleteNow,
	})
	if err != nil {
		return nil, err
	}
	return &DeletePortResponse{}, nil
}
