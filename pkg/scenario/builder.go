package scenario

import (
	"errors"
	"fmt"
	"ohs30359/vegeta-cli/pkg/config"
	"os"
	"strconv"
	"strings"
	"time"
)

type Builder struct {
	scenario config.Scenario
	max      int
	duration time.Duration
	root     string
}

func NewBuilder(scenario config.Scenario, max int, duration time.Duration, root string) Builder {
	return Builder{
		scenario: scenario,
		max:      max,
		duration: duration,
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

func (own *Builder) CreateScenarioBuffer(targetFile string) string {
	return fmt.Sprintf("vegeta attack -targets=%s -rate=%s/s -duration %ss", targetFile, strconv.Itoa(own.max), own.duration)
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
