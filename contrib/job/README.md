# go job

## Usage

### Install

```shell
go get github.com/go-keg/keg
```

### Example
```go
package main

import (
	"context"
	"errors"
	syslog "log"
	"os"
	"time"

	"github.com/go-keg/keg/contrib/alert"
	"github.com/go-keg/keg/contrib/job"
	"github.com/go-keg/keg/third_party/workwechat"
	"github.com/go-kratos/kratos/v2/log"
	"golang.org/x/time/rate"
)

func main() {
	j := job.NewJob(
		log.DefaultLogger,
		job.NewWorker("test", example),
		job.NewWorker(
			"test-with-limiter",
			example,
			job.WithLimiter(rate.NewLimiter(rate.Every(time.Second), 1)),
		),
		job.NewWorker(
			"test-with-report-error",
			example,
			job.WithLimiter(rate.NewLimiter(rate.Every(time.Second), 1)),
			job.WithAlert(
				alert.NewDeduper(
					alert.SetAlert(
						workwechat.NewWebhook(os.Getenv("WORK_WECHAT_TOKEN")),
					),
				),
			),
		),
		job.NewWorker(
			"test-with-report-panic",
			example2,
			job.WithLimiter(rate.NewLimiter(rate.Every(time.Second), 1)),
			job.WithAlert(
				alert.NewDeduper(
					alert.SetAlert(
						workwechat.NewWebhook(os.Getenv("WORK_WECHAT_TOKEN")),
					),
				),
			),
		),
	)
	err := j.Start(context.Background())
	if err != nil {
		panic(err)
	}
}

func example(ctx context.Context) error {
	syslog.Println("do something example...")
	if time.Now().Second()%10 == 1 {
		// test report error
		return errors.New("test err")
	}
	return nil
}

func example2(ctx context.Context) error {
	syslog.Println("do something example2...")
	if time.Now().Second()%10 == 2 {
		panic("test panic")
	}
	return nil
}

```