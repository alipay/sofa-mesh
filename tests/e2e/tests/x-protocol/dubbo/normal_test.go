package dubbo

import (
	"fmt"
	"strings"
	"testing"

	"time"

	"github.com/pkg/errors"
	"istio.io/istio/pkg/log"
)

func TestNormal(t *testing.T) {
	log.Infof("Begin to start test case")

	standby := 0

	for i := 0; i <= testRetryTimes; i++ {
		if i > testRetryTimes {
			t.Fatalf("has been try %d times, but never success", i)
		}

		standby += 5

		url := getConsumerTargetUrl(queryName)
		log.Infof("%d time fetch url: %s \n", i, url)
		actualResponseContent, err := fetchConsumerResult(url)

		if err != nil {
			log.Errorf("error when fetch %s. Error %s", url, err)
			continue
		}

		log.Infof("success when fetch %s. response content: %s \n", url, actualResponseContent)

		if msg, err := verifyQuery(actualResponseContent, expectedResponseContent); err != nil {
			log.Errora(err)
			t.Fatal(err)
		} else {
			log.Info(msg)
			break
		}

		time.Sleep(time.Duration(standby) * time.Second)
	}
}

func verifyQuery(actualResponseContent string, expectedResponseContent string) (string, error) {
	if strings.Index(actualResponseContent, expectedResponseContent) == -1 {
		return "", errors.New(fmt.Sprintf("response content is not correct, expect include %s, but %s", expectedResponseContent, actualResponseContent))
	} else {
		return fmt.Sprintf("response content is correct: %s", actualResponseContent), nil
	}
}
