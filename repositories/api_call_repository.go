package repositories

import (
	"context"
	"errors"
	"flea-market/utils"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

type APICallRepository struct {
	apiClient *resty.Client
}

type Method string

const (
	GET    Method = "GET"
	POST   Method = "POST"
	PUT    Method = "PUT"
	DELETE Method = "DELETE"
)

// If you define like uint and value is not set, zero value will be set.
// It's complicated, so use *uint instead.
type Post struct {
	Id     *uint   `json:"id"`
	Title  *string `json:"title"`
	Body   *string `json:"body"`
	UserId *int    `json:"userId"`
	Dummy  *string `json:"dummy"`
}

const baseURL = "https://jsonplaceholder.typicode.com/"

func NewAPICallRepository(apiClient *resty.Client) *APICallRepository {
	return &APICallRepository{apiClient: apiClient}
}

func (r *APICallRepository) GetAllPosts(ctx context.Context) (*[]Post, error) {
	var result []Post
	endpoint := baseURL + "/posts"

	apiReqCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := r.apiClient.R().
		SetContext(apiReqCtx).
		SetResult(&result).
		Get(baseURL + "/posts")

	// Request itself fails like being unable to connect the server
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, errors.New("API呼び出しがタイムアウトしました:" + err.Error())
		}
		if errors.Is(err, context.Canceled) {
			return nil, errors.New("API呼び出しがキャンセルされました" + err.Error())
		}
		return nil, utils.NewExternalAPIConnectionError(fmt.Sprintf("Method:GET Path:%v", endpoint), err)
	}

	// When status code is greater than 399, handle and error.
	// Logging is done by middleware
	if res.IsError() {
		return nil, utils.NewExternalAPIReturnsError("",
			fmt.Errorf("external API error: endpoint=%s status=%d body=%s", endpoint, res.StatusCode(), res.String()))
	}

	return &result, nil
}
