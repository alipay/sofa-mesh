package dubbo

import (
	"testing"

	"strings"

	"istio.io/istio/pkg/log"
)

type weightRoutingRule struct {
	key  string
	rate float64
}

func TestWeightRouting(t *testing.T) {
	var rules = []weightRoutingRule{
		{
			key:  weightTwentyRule,
			rate: 0.2,
		},
		{
			key:  weightFiftyRule,
			rate: 0.5,
		},
	}

	inspect(applyRules([]string{destRule}), "failed to apply rules", "", t)
	defer func() {
		inspect(deleteRules([]string{destRule}), "failed to delete rules", "", t)
	}()

	testWeightRoutingRules(t, rules)
}

func testWeightRoutingRules(t *testing.T, rules []weightRoutingRule) {
	for _, rule := range rules {
		testWeightRoutingRule(t, rule)
	}
}

func testWeightRoutingRule(t *testing.T, rule weightRoutingRule) {
	inspect(applyRules([]string{rule.key}), "failed to apply rules", "", t)
	defer func() {
		inspect(deleteRules([]string{rule.key}), "failed to delete rules", "", t)
	}()

	tolerance := 0.05
	totalShot := 100

	for i := 0; i < testRetryTimes; i++ {
		cSubnet1, cSubnet2 := 0, 0

		for c := 0; c < totalShot; c++ {
			url := getConsumerTargetUrl(queryName)
			log.Infof("%d time fetch url: %s \n", (c+1)+(i*testRetryTimes), url)
			actualResponseContent, err := fetchConsumerResult(url)

			if err != nil {
				log.Errorf("error when fetch %s. Error %s", url, err)
				continue
			}

			log.Infof("success when fetch %s. response content: %s \n", url, actualResponseContent)

			if isWithRightVersion(actualResponseContent, "v1") {
				cSubnet1 += 1
			} else if isWithRightVersion(actualResponseContent, "v2") {
				cSubnet2 += 1
			} else {
				log.Error("received unexpected version: %s")
			}
		}

		if isWithinPercentage(cSubnet1, totalShot, 1.0-rule.rate, tolerance) &&
			isWithinPercentage(cSubnet2, totalShot, rule.rate, tolerance) {
			log.Infof(
				"Success! Version routing acts as expected for rate %f, "+
					"version v1 hit %d, version v2 hit %d", rule.rate, cSubnet1, cSubnet2)
			break
		}

		if i >= testRetryTimes {
			t.Errorf("Failed version migration test for rate %f, "+
				"version v1 hit %d, version v2 hit %d", rule.rate, cSubnet1, cSubnet2)
		}
	}
}

func isWithRightVersion(actualResponseContent string, expectedVersion string) bool {
	return strings.Index(actualResponseContent, expectedVersion) >= 0
}

func isWithinPercentage(count int, total int, rate float64, tolerance float64) bool {
	minimum := int((rate - tolerance) * float64(total))
	maximum := int((rate + tolerance) * float64(total))
	return count >= minimum && count <= maximum
}
