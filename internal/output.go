package internal

import (
	"errors"
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

	// ratioがトータルで100%になっているかを判定
	ratioSum := 0
	for _, scenario := range conf.Scenario {
		ratioSum += scenario.Ratio
	}
	if ratioSum != 100 {
		return errors.New("\"ratio\" must sum up to 100")
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

		builder := scenarioBuilder.NewBuilder(scenario, conf, path)
		targetBuf, e := builder.CreateTargetBuffer(scenario, valueFilePaths)
		if e != nil {
			return e
		}

		// "target" option で指定されるファイルを作成
		if e := file.Write(scenarioDirPath+"/target.txt", []byte(targetBuf)); e != nil {
			return e
		}

		// 実行シェルを作成
		files := builder.CreateScenarioBuffer(scenario, "./target.txt")
		for j, buf := range files {
			if e := file.Write(fmt.Sprintf("%s/scenario_%s.sh", scenarioDirPath, strconv.Itoa(j)), []byte(buf)); e != nil {
				return e
			}
		}
	}

	return nil
}
