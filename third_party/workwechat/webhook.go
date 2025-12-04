package workwechat

import (
	"context"
	"fmt"
	"time"

	"github.com/go-keg/keg/contrib/http"
	"golang.org/x/time/rate"
)

type Webhook struct {
	key     string
	client  *http.Client
	limiter *rate.Limiter
}

const wechatWebHookURL = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?"

type WebhookOption func(webhook *Webhook)

func WithWebhookLimiter(limiter *rate.Limiter) WebhookOption {
	return func(webhook *Webhook) {
		webhook.limiter = limiter
	}
}

func NewWebhook(key string, opts ...WebhookOption) *Webhook {
	webhook := &Webhook{
		key:     key,
		client:  http.NewClient(),
		limiter: rate.NewLimiter(rate.Every(time.Second*10), 1),
	}
	for _, opt := range opts {
		opt(webhook)
	}
	return webhook
}

func (r Webhook) SendText(ctx context.Context, content string) error {
	return r.sendMessage(ctx, Message{
		MsgType: "text",
		Text: &MessageContent{
			Content: content,
		},
	})
}

func (r Webhook) SendMarkdown(ctx context.Context, content string) error {
	return r.sendMessage(ctx, Message{
		MsgType: "markdown",
		Markdown: &MessageContent{
			Content: content,
		},
	})
}

func (r Webhook) Alert(ctx context.Context, content string) error {
	return r.SendMarkdown(ctx, content)
}

func (r Webhook) sendMessage(ctx context.Context, message Message) error {
	_ = r.limiter.Wait(ctx)
	resp, err := r.client.Post(ctx, wechatWebHookURL+"key="+r.key, http.SetBody(message))
	if err != nil {
		return err
	}
	var response struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}
	err = resp.Unmarshal(&response)
	if err != nil {
		return err
	}
	if response.ErrCode != 0 {
		return fmt.Errorf("workwechat webhook: error response: [%s]", resp.Content())
	}
	return nil
}
