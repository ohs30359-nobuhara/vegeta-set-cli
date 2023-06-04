package internal

import (
	"fmt"
	"ohs30359/vegeta-cli/pkg/config"
	"ohs30359/vegeta-cli/pkg/dir"
	"ohs30359/vegeta-cli/pkg/file"
	scenarioBuilder "ohs30359/vegeta-cli/pkg/scenario"
	"os"
	"strconv"
)

func Output(path string) error {
	conf, e := config.Load(fmt.Sprintf("%s/scenario.yaml", path))
	if e != nil {
		return e
	}

	if e := os.MkdirAll("dist", 0755); e != nil {
		return e
	}

	for i, scenario := range conf.Scenario {
		// ./dist/scenario_x 配下にscenarioファイルを作成
		scenarioDirPath := fmt.Sprintf("./dist/scenario_%s", strconv.Itoa(i))
		if e := os.MkdirAll(scenarioDirPath, 0755); e != nil {
			return e
		}

		var valueFilePaths []string

		if scenario.Value != nil {
			values, e := dir.Scan(fmt.Sprintf("%s/%s", path, *scenario.Value))
			if e != nil {
				return e
			}

			switch scenario.Method {
			case "GET":
				// GETの場合はパラメータファイルの複製は不要なため元のパラメータファイルの絶対パスを渡す
				// → Builderの中でファイルを読み込んでURLパラメータが自動付与される
				valueFilePaths = values
			case "POST":
				// POSTの場合はvalues配下のパラメータファイルをscenario dir配下にコピーする
				targetDirPath := scenarioDirPath + "/values"
				if e := os.MkdirAll(targetDirPath, 0755); e != nil {
					return e
				}
				for _, value := range values {
					fileState, e := file.Copy(value, targetDirPath)
					if e != nil {
						return e
					}
					valueFilePaths = append(valueFilePaths, "./values/"+fileState.Name)
				}
			}

		}

		builder := scenarioBuilder.NewBuilder(scenario, conf.Rate, conf.Duration, path)
		targetBuf, e := builder.CreateTargetBuffer(scenario, valueFilePaths)
		if e != nil {
			return e
		}

		// "target" option で指定されるファイルを作成
		if e := file.Write(scenarioDirPath+"/target.txt", []byte(targetBuf)); e != nil {
			return e
		}

		if e := file.Write(scenarioDirPath+"/scenario.sh", []byte(builder.CreateScenarioBuffer("./target.txt"))); e != nil {
			return e
		}
	}

	return nil
}
