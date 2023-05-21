package scenario

import (
	"errors"
	"fmt"
	"ohs30359/vegeta-cli/pkg/config"
	"ohs30359/vegeta-cli/pkg/dir"
	"os"
	"strings"
)

func Output(path string) error {
	conf, e := config.Load(fmt.Sprintf("%s/scenario.yaml", path))
	if e != nil {
		return e
	}

	builder := strings.Builder{}

	for _, scenario := range conf.Scenario {
		var values []string

		// valueが指定されている場合は、指定されたパス配下のファイル一覧を取得する
		if scenario.Value != nil {
			values, e = dir.Scan(fmt.Sprintf("%s/%s", path, *scenario.Value))

			if e != nil {
				return fmt.Errorf("the file does not exist in the %s directory or the target directory does not exist", *scenario.Value)
			}
		}

		switch scenario.Method {
		case "GET":
			scenarioTxt, e := buildGetScenario(scenario, values)
			if e != nil {
				return errors.New("GET scenario create fail")
			}
			builder.WriteString(scenarioTxt)
		}
	}

	file, e := os.Create("scenario.txt")
	file.Write([]byte(builder.String()))
	if e != nil {
		return e
	}
	defer file.Close()

	return nil
}

func buildGetScenario(scenario config.Scenario, values []string) (string, error) {
	builder := strings.Builder{}

	// パラメータ一覧がある場合は ファイルを読み込んでクエリパラメータをURLに付与する
	if len(values) != 0 {
		content, e := os.ReadFile(values[0])
		if e != nil {
			return "", fmt.Errorf(" %s does not exist", values[0])
		}

		queryParams := strings.Split(string(content), "\n")

		// TODO: ここで割合計算をする
		for _, param := range queryParams {
			if len(param) == 0 {
				continue
			}
			builder.WriteString(fmt.Sprintf("GET %s?%s \n", scenario.Url, strings.ReplaceAll(param, "?", "")))
		}
	} else {
		// パラメータがなければ単純に URLをセット
		builder.WriteString(fmt.Sprintf("GET %s \n", scenario.Url))
	}

	return builder.String(), nil
}
