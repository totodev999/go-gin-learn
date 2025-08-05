package repositories

import (
	"context"
	"errors"
	"flea-market/utils"
	"fmt"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
	"golang.org/x/sync/errgroup"
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

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Address  struct {
		Street  string `json:"street"`
		Suite   string `json:"suite"`
		City    string `json:"city"`
		Zipcode string `json:"zipcode"`
		Geo     struct {
			Lat string `json:"lat"`
			Lng string `json:"lng"`
		} `json:"geo"`
	} `json:"address"`
	Phone   string `json:"phone"`
	Website string `json:"website"`
	Company struct {
		Name        string `json:"name"`
		CatchPhrase string `json:"catchPhrase"`
		Bs          string `json:"bs"`
	} `json:"company"`
}

type UserAndPosts struct {
	User  User   `json:"user"`
	Posts []Post `json:"posts"`
}

var baseURL string

// is it better to use "resty.Client" BaseUrl?
func setURL() {
	baseURL = os.Getenv("BASE_URL")
	if baseURL == "" {
		utils.Logger(utils.ExternalAPIConnectionError, "", "", "", "env BASE_URL is not set")
		panic("env BASE_URL is not set")
	}
}

func NewAPICallRepository(apiClient *resty.Client) *APICallRepository {
	setURL()
	return &APICallRepository{apiClient: apiClient}
}

func (r *APICallRepository) GetAllPosts(ctx context.Context) (*[]Post, error) {
	var result []Post
	endpoint := baseURL + "/posts"

	apiReqCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
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

func (r *APICallRepository) GetUserAndPosts(ctx context.Context, userId uint) (*UserAndPosts, error) {

	var (
		user  User
		posts []Post
	)
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		reqCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
		endpoint := fmt.Sprintf("%s/users/%d", baseURL, userId)
		defer cancel()
		resp, err := r.apiClient.R().
			SetContext(reqCtx).
			SetResult(&user).
			Get(endpoint)

		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				return errors.New("API呼び出しがタイムアウトしました:" + err.Error())
			}
			if errors.Is(err, context.Canceled) {
				return errors.New("API呼び出しがキャンセルされました" + err.Error())
			}
			return utils.NewExternalAPIConnectionError(fmt.Sprintf("Method:GET Path:%v", endpoint), err)
		}

		if resp.IsError() {
			return fmt.Errorf("ユーザー取得APIエラー: %v", resp.Status())
		}
		return nil
	})

	g.Go(func() error {
		reqCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
		endpoint := fmt.Sprintf("%s/posts?userId=%d", baseURL, userId)
		defer cancel()
		resp, err := r.apiClient.R().
			SetContext(reqCtx).
			SetResult(&posts).
			Get(endpoint)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				return errors.New("API呼び出しがタイムアウトしました:" + err.Error())
			}
			if errors.Is(err, context.Canceled) {
				return errors.New("API呼び出しがキャンセルされました" + err.Error())
			}
			return utils.NewExternalAPIConnectionError(fmt.Sprintf("Method:GET Path:%v", endpoint), err)
		}
		if resp.IsError() {
			return fmt.Errorf("ポスト取得APIエラー: %v", resp.Status())
		}
		return nil
	})

	// if one of request fails, then an another request will be canceled and return an error.
	if err := g.Wait(); err != nil {
		return nil, err
	}

	userAndPosts := &UserAndPosts{
		User:  user,
		Posts: posts,
	}
	return userAndPosts, nil
}
