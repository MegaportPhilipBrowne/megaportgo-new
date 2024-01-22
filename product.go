package megaport

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/megaport/megaportgo/mega_err"
	"github.com/megaport/megaportgo/types"
)

type ProductService interface {
	ExecuteOrder(ctx context.Context, requestBody *[]byte) (*[]byte, error)
	ModifyProduct(ctx context.Context, req *ModifyProductRequest) (*ModifyProductResponse, error)
	DeleteProduct(ctx context.Context, req *DeleteProductRequest) (*DeleteProductResponse, error)
}

type ModifyProductRequest struct {
	ProductID             string
	ProductType           string
	Name                  string
	CostCentre            string
	MarketplaceVisibility bool
}

type ModifyProductResponse struct {
	IsUpdated bool
}

type DeleteProductRequest struct {
	ProductID string
	DeleteNow bool
}

type DeleteProductResponse struct{}

// ProductServiceOp handles communication with Product methods of the Megaport API.
type ProductServiceOp struct {
	Client *Client
}

func NewProductServiceOp(c *Client) *ProductServiceOp {
	return &ProductServiceOp{
		Client: c,
	}
}

func (svc *ProductServiceOp) ExecuteOrder(ctx context.Context, requestBody *[]byte) (*[]byte, error) {
	path := "/v3/networkdesign/buy"

	url := svc.Client.BaseURL.JoinPath(path).String()

	req, err := svc.Client.NewRequest(ctx, http.MethodPost, url, requestBody)
	if err != nil {
		return nil, err
	}

	response, resErr := svc.Client.Do(ctx, req, nil)
	if err != nil {
		return nil, resErr
	}

	if response != nil {
		svc.Client.Logger.Debug("Executing product order", "url", url, "status_code", response.StatusCode)
		defer response.Body.Close()
	}

	isError, parsedError := svc.Client.IsErrorResponse(response, &resErr, 200)

	if isError {
		return nil, parsedError
	}

	body, fileErr := io.ReadAll(response.Body)
	if fileErr != nil {
		return nil, fileErr
	}

	return &body, nil
}

// ModifyProduct modifies a product. The available fields to modify are Name, Cost Centre, and Marketplace Visibility.
func (svc *ProductServiceOp) ModifyProduct(ctx context.Context, req *ModifyProductRequest) (*ModifyProductResponse, error) {

	if req.ProductType == types.PRODUCT_MEGAPORT || req.ProductType == types.PRODUCT_MCR {
		update := types.ProductUpdate{
			Name:                 req.Name,
			CostCentre:           req.CostCentre,
			MarketplaceVisbility: req.MarketplaceVisibility,
		}
		path := fmt.Sprintf("/v2/product/%s/%s", req.ProductType, req.ProductID)
		url := svc.Client.BaseURL.JoinPath(path).String()

		body, marshalErr := json.Marshal(update)

		if marshalErr != nil {
			return nil, marshalErr
		}

		req, err := svc.Client.NewRequest(ctx, http.MethodPut, url, []byte(body))

		if err != nil {
			return nil, err
		}

		updateResponse, err := svc.Client.Do(ctx, req, nil)

		isResErr, compiledResErr := svc.Client.IsErrorResponse(updateResponse, &err, 200)

		if isResErr {
			return nil, compiledResErr
		} else {
			return &ModifyProductResponse{IsUpdated: true}, nil
		}
	} else {
		return nil, errors.New(mega_err.ERR_WRONG_PRODUCT_MODIFY)
	}
}

// DeleteProduct is responsible for either scheduling a product for deletion "CANCEL" or deleting a product immediately
// "CANCEL_NOW".
func (svc *ProductServiceOp) DeleteProduct(ctx context.Context, req *DeleteProductRequest) (*DeleteProductResponse, error) {
	var action string

	if req.DeleteNow {
		action = "CANCEL_NOW"
	} else {
		action = "CANCEL"
	}

	path := "/v3/product/" + req.ProductID + "/action/" + action
	url := svc.Client.BaseURL.JoinPath(path).String()

	clientReq, err := svc.Client.NewRequest(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}

	deleteResp, err := svc.Client.Do(ctx, clientReq, nil)
	defer deleteResp.Body.Close() // nolint

	isError, errorMessage := svc.Client.IsErrorResponse(deleteResp, &err, 200)
	if isError {
		return nil, errorMessage
	} else {
		return &DeleteProductResponse{}, nil
	}
}