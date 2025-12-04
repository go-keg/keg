package workwechat

import (
	"context"
	"errors"
	"fmt"
	nethttp "net/http"
	"time"

	"github.com/go-keg/keg/contrib/cache"
	"github.com/go-keg/keg/contrib/http"
)

type Client struct {
	corpID     string
	corpSecret string
	agentID    int
	http       *http.Client
}

func NewClient(corpID string, corpSecret string, agentID int) *Client {
	return &Client{
		corpID:     corpID,
		corpSecret: corpSecret,
		agentID:    agentID,
		http:       http.NewClient(),
	}
}

type Response struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}
type AccessTokenResp struct {
	Response
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

func (r Client) AccessToken(ctx context.Context) (string, error) {
	return cache.LocalRemember("workwechat:access_token", time.Hour, func() (string, error) {
		response, err := r.http.Get(ctx, "https://qyapi.weixin.qq.com/cgi-bin/gettoken",
			http.SetQueryParams(map[string]any{
				"corpid":     r.corpID,
				"corpsecret": r.corpSecret,
			}),
		)
		if err != nil {
			return "", err
		}
		var resp AccessTokenResp
		err = response.Unmarshal(&resp)
		if err != nil {
			return "", err
		}
		if resp.Errcode == 0 {
			return resp.AccessToken, nil
		}
		return "", errors.New(response.Content())
	})
}

type SendMessageParams struct {
	// Message
	MsgType    string           `json:"msgtype"`
	AgentID    int              `json:"agentid"`
	Text       *MessageContent  `json:"text,omitempty"`
	Markdown   *MessageContent  `json:"markdown,omitempty"`
	MarkdownV2 *MessageContent  `json:"markdown_v2,omitempty"`
	Textcard   *MessageTextCard `json:"textcard,omitempty"`

	ToUser                 string `json:"touser,omitempty"`
	ToParty                string `json:"toparty,omitempty"`
	ToTag                  string `json:"totag,omitempty"`
	Safe                   int    `json:"safe,omitempty"`
	EnableIdTrans          int    `json:"enable_id_trans,omitempty"`
	EnableDuplicateCheck   int    `json:"enable_duplicate_check,omitempty"`
	DuplicateCheckInterval int    `json:"duplicate_check_interval,omitempty"`
}

type SendMessageResp struct {
	Response
	MsgID string `json:"msgid"`
}

func (r Client) SendMessage(ctx context.Context, params SendMessageParams) error {
	if params.MsgType == "" {
		if params.Text != nil {
			params.MsgType = "text"
		}
		if params.Markdown != nil {
			params.MsgType = "markdown"
		}
		if params.MarkdownV2 != nil {
			params.MsgType = "markdown_v2"
		}
		if params.Textcard != nil {
			params.MsgType = "textcard"
		}
	}
	accessToken, err := r.AccessToken(ctx)
	if err != nil {
		return err
	}
	params.AgentID = r.agentID
	response, err := r.http.Post(ctx, "https://qyapi.weixin.qq.com/cgi-bin/message/send",
		http.SetBody(params),
		http.SetQueryParams(map[string]any{
			"access_token": accessToken,
		}),
	)
	if err != nil {
		return err
	}
	if response.StatusCode() != nethttp.StatusOK {
		return fmt.Errorf("statusCode: %d", response.StatusCode())
	}
	var resp SendMessageResp
	err = response.Unmarshal(&resp)
	if err != nil {
		return err
	}
	if resp.Errcode != 0 {
		return errors.New(response.Content())
	}
	return nil
}

type GetUserIDResp struct {
	Response
	UserID string `json:"userid"`
}

func (r Client) GetUserID(ctx context.Context, mobile string) (string, error) {
	accessToken, err := r.AccessToken(ctx)
	if err != nil {
		return "", err
	}
	response, err := r.http.Post(ctx, "https://qyapi.weixin.qq.com/cgi-bin/user/getuserid", http.SetQueryParams(map[string]any{
		"access_token": accessToken,
	}), http.SetBody(map[string]string{
		"mobile": mobile,
	}))
	if err != nil {
		return "", err
	}
	if response.StatusCode() != nethttp.StatusOK {
		return "", fmt.Errorf("statusCode: %d", response.StatusCode())
	}
	var resp GetUserIDResp
	err = response.Unmarshal(&resp)
	if err != nil {
		return "", err
	}
	if resp.Errcode != 0 {
		return "", errors.New(response.Content())
	}
	return resp.UserID, nil
}

func (r Client) GetUserInfo(ctx context.Context, code string) (string, error) {
	accessToken, err := r.AccessToken(ctx)
	if err != nil {
		return "", err
	}
	response, err := r.http.Get(ctx, "https://qyapi.weixin.qq.com/cgi-bin/auth/getuserinfo", http.SetQueryParams(map[string]any{
		"access_token": accessToken,
		"code":         code,
	}))
	if err != nil {
		return "", err
	}
	if response.StatusCode() != nethttp.StatusOK {
		return "", fmt.Errorf("statusCode: %d", response.StatusCode())
	}
	var resp GetUserIDResp
	err = response.Unmarshal(&resp)
	if err != nil {
		return "", err
	}
	if resp.Errcode != 0 {
		return "", errors.New(response.Content())
	}
	return resp.UserID, nil
}
