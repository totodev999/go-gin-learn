package infra

import (
	"flea-market/utils"
	"fmt"

	"github.com/go-resty/resty/v2"
)

func NewBaseAPIClient() *resty.Client {
	client := resty.New()

	// Beforeリクエスト
	client.OnBeforeRequest(func(c *resty.Client, req *resty.Request) error {
		endpoint := req.URL
		method := req.Method
		utils.Logger(
			utils.ExternalAPIRequestStart, nil,
			fmt.Sprintf("Method:%s Path:%v", method, endpoint),
		)
		return nil
	})

	// Afterレスポンス
	client.OnAfterResponse(func(c *resty.Client, resp *resty.Response) error {
		endpoint := resp.Request.Method + " " + resp.Request.URL
		status := resp.StatusCode()

		if resp.IsError() {
			utils.Logger(
				utils.ExternalAPIRequestEnd,
				nil,
				fmt.Sprintf("HTTPステータスエラー Request:%v StatusCode:%d ErrorResponse:%s", endpoint, status, resp.String()),
			)
		} else {
			utils.Logger(
				utils.ExternalAPIRequestEnd,
				nil,
				fmt.Sprintf("成功 Request:%v", endpoint),
			)
		}
		return nil
	})

	return client
}
