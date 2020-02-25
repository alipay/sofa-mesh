package dubbo

import (
	"fmt"
	"testing"
	"time"

	"strings"

	"github.com/pkg/errors"
	"istio.io/istio/pkg/log"
)

type versionRoutingRule struct {
	key     string
	version string
}

func TestVersionRouting(t *testing.T) {
	var rules = []versionRoutingRule{
		{
			key:     versionV1Rule,
			version: "v1",
		},
		{
			key:     versionV2Rule,
			version: "v2",
		},
		{
			key:     versionV1Rule,
			version: "v1",
		},
	}

	inspect(applyRules([]string{destRule}), "failed to apply rules", "", t)
	defer func() {
		inspect(deleteRules([]string{destRule}), "failed to delete rules", "", t)
	}()

	testVersionRoutingRules(t, rules)
}

func testVersionRoutingRules(t *testing.T, rules []versionRoutingRule) {
	for _, rule := range rules {
		testVersionRoutingRule(t, rule)
	}
}

func testVersionRoutingRule(t *testing.T, rule versionRoutingRule) {
	inspect(applyRules([]string{rule.key}), "failed to apply rules", "", t)
	defer func() {
		inspect(deleteRules([]string{rule.key}), "failed to delete rules", "", t)
	}()

	standby := 0
	totalShot := 10

	for i := 0; i < testRetryTimes; i++ {
		time.Sleep(time.Duration(standby) * time.Second)
		standby += 5

		for c := 0; c < totalShot; c++ {
			url := getConsumerTargetUrl(queryName)
			log.Infof("%d time fetch url: %s \n", (c+1)+(i*testRetryTimes), url)
			actualResponseContent, err := fetchConsumerResult(url)

			if err != nil {
				log.Errorf("error when fetch %s. Error %s", url, err)
				continue
			}

			log.Infof("success when fetch %s. response content: %s \n", url, actualResponseContent)

			if msg, err := verifyVersion(actualResponseContent, rule.version); err != nil {
				log.Errora(err)
				t.Fatal(err)
			} else {
				log.Info(msg)
			}
		}
	}
}

func verifyVersion(actualResponseContent string, expectedVersion string) (string, error) {
	if strings.Index(actualResponseContent, expectedVersion) == -1 {
		return "", errors.New(fmt.Sprintf("response content is not correct, expect version %s, but %s", expectedVersion, actualResponseContent))
	} else {
		return fmt.Sprintf("response content is correct: %s", actualResponseContent), nil
	}
}
