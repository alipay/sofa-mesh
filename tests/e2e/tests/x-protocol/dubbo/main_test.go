package dubbo

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"time"

	"github.com/hashicorp/go-multierror"
	"istio.io/istio/pkg/log"
	"istio.io/istio/tests/e2e/framework"
	"istio.io/istio/tests/util"
)

const (
	testDataDir             = "tests/e2e/tests/x-protocol/dubbo/testdata"
	yamlExtension           = "yaml"
	deploymentDir           = "platform/kube"
	routeRulesDir           = "networking"
	consumerYaml            = "dubbo-consumer"
	providerYaml            = "dubbo-provider"
	busyboxYaml             = "busybox"
	consumerHTTPPort        = "8080"
	busyboxName             = "busybox"
	consumerName            = "dubbo-consumer"
	providerName            = "dubbo-provider"
	destRule                = "destination-rule-all"
	versionV1Rule           = "virtual-service-provider-v1"
	versionV2Rule           = "virtual-service-provider-v2"
	weightTwentyRule        = "virtual-service-provider-20-80"
	weightFiftyRule         = "virtual-service-provider-50-50"
	queryName               = "test"
	expectedResponseContent = `Hello, test (from Spring Boot dubbo e2e test)`
)

var (
	tc             *testConfig
	testRetryTimes = 5
)

type testConfig struct {
	*framework.CommonConfig
}

func (t *testConfig) Setup() error {
	if !util.CheckPodsRunning(tc.Kube.Namespace, tc.Kube.KubeConfig) {
		return fmt.Errorf("can't get all pods running")
	}

	time.Sleep(time.Duration(30) * time.Second)

	return nil
}

func (t *testConfig) Teardown() error {
	return nil
}

func TestMain(m *testing.M) {
	flag.Parse()
	check(framework.InitLogging(), "cannot setup logging")
	check(setTestConfig(), "could not create TestConfig")
	tc.Cleanup.RegisterCleanable(tc)
	os.Exit(tc.RunTest(m))
}

func setTestConfig() error {
	cc, err := framework.NewCommonConfig("dubbo_test")
	if err != nil {
		return err
	}
	tc = new(testConfig)
	tc.CommonConfig = cc

	apps := getApps()

	for i := range apps {
		tc.Kube.AppManager.AddApp(&apps[i])
	}

	return nil
}

func getApps() []framework.App {
	return []framework.App{
		{
			AppYaml:    getDeploymentPath(consumerYaml),
			KubeInject: true,
		},
		{
			AppYaml:    getDeploymentPath(providerYaml),
			KubeInject: true,
		},
		{
			AppYaml:    getDeploymentPath(busyboxYaml),
			KubeInject: false,
		},
	}
}

func getConsumerTargetUrl(queryName string) string {
	return fmt.Sprintf("http://%s.%s.svc.cluster.local:%s/sayHello?name=%s", consumerName, tc.Kube.Namespace, consumerHTTPPort, queryName)
}

func fetchConsumerResult(url string) (string, error) {
	namespace := tc.Kube.Namespace
	kubeConfig := tc.Kube.KubeConfig

	podName, err := util.GetPodName(namespace, "app="+busyboxName, kubeConfig)
	if err != nil {
		return "", err
	}

	resp, err := util.PodExec(namespace, podName, "app", "curl --silent "+url, false, kubeConfig)

	if err != nil {
		return "", err
	}

	return resp, nil
}

func getDeploymentPath(deployment string) string {
	return util.GetResourcePath(filepath.Join(testDataDir, deploymentDir, deployment+"."+yamlExtension))
}

func getRulePath(ruleKey string) string {
	return util.GetResourcePath(filepath.Join(testDataDir, routeRulesDir, ruleKey+"."+yamlExtension))
}

func check(err error, msg string) {
	if err != nil {
		log.Errorf("%s. Error %s", msg, err)
		os.Exit(-1)
	}
}

func inspect(err error, fMsg, sMsg string, t *testing.T) {
	if err != nil {
		log.Errorf("%s. Error %s", fMsg, err)
		t.Error(err)
	} else if sMsg != "" {
		log.Info(sMsg)
	}
}

func deleteRules(ruleKeys []string) error {
	var err error
	for _, ruleKey := range ruleKeys {
		rule := getRulePath(ruleKey)
		if e := util.KubeDelete(tc.Kube.Namespace, rule, tc.Kube.KubeConfig); e != nil {
			err = multierror.Append(err, e)
		}
	}
	log.Info("Waiting for rule to be cleaned up...")
	time.Sleep(time.Duration(30) * time.Second)
	return err
}

func applyRules(ruleKeys []string) error {
	for _, ruleKey := range ruleKeys {
		rule := getRulePath(ruleKey)
		if err := util.KubeApply(tc.Kube.Namespace, rule, tc.Kube.KubeConfig); err != nil {
			return err
		}
	}
	log.Info("Waiting for rules to propagate...")
	time.Sleep(time.Duration(30) * time.Second)
	return nil
}
