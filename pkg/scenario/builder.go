package scenario

import (
	"errors"
	"fmt"
	"math"
	"ohs30359/vegeta-cli/pkg/config"
	"os"
	"strconv"
	"strings"
)

type Builder struct {
	scenario config.Scenario
	conf     config.Config
	root     string
}

func NewBuilder(scenario config.Scenario, conf config.Config, root string) Builder {
	return Builder{
		scenario: scenario,
		conf:     conf,
		root:     root,
	}
}

func (own *Builder) CreateTargetBuffer(scenario config.Scenario, values []string) (string, error) {
	switch scenario.Method {
	case "GET":
		return own.createGetBuffer(scenario, values)
	case "POST":
		return own.createPostBuffer(scenario, values)
	}

	return "", errors.New("scenario Method must be GET or POST")
}

func (own *Builder) CreateScenarioBuffer(scenario config.Scenario, targetFile string) []string {
	// testerの性能上限を超えている場合は シナリオを分割する
	rps := int(float64(own.conf.Rate) * (float64(scenario.Ratio) / 100))
	if rps > own.conf.Tester.Limit {
		var results []string
		rounded := int(math.Ceil(float64(rps) / float64(own.conf.Tester.Limit)))
		for i := 0; i < rounded; i++ {
			results = append(results, fmt.Sprintf("vegeta attack -targets=%s -rate=%s/s -duration %s", targetFile, strconv.Itoa(rps/rounded), own.conf.Duration))
		}

		return results
	}

	return []string{fmt.Sprintf("vegeta attack -targets=%s -rate=%s/s -duration %ss", targetFile, strconv.Itoa(rps), own.conf.Duration)}
}

func (own *Builder) createGetBuffer(scenario config.Scenario, values []string) (string, error) {
	builder := strings.Builder{}

	// パラメータ一覧がある場合は ファイルを読み込んでクエリパラメータをURLに付与する
	if len(values) != 0 {
		content, e := os.ReadFile(values[0])
		if e != nil {
			return "", fmt.Errorf(" %s does not exist", values[0])
		}

		queryParams := strings.Split(string(content), "\n")

		for _, param := range queryParams {
			if len(param) == 0 {
				continue
			}
			builder.WriteString(fmt.Sprintf("GET %s?%s \n\n", scenario.Url, strings.ReplaceAll(param, "?", "")))
		}
	} else {
		// パラメータがなければ単純に URLをセット
		builder.WriteString(fmt.Sprintf("GET %s \n\n", scenario.Url))
	}

	return builder.String(), nil

}

func (own *Builder) createPostBuffer(scenario config.Scenario, values []string) (string, error) {
	builder := strings.Builder{}

	for _, val := range values {
		builder.WriteString(fmt.Sprintf("POST %s \n@%s \n\n", scenario.Url, val))
	}

	return builder.String(), nil
}
