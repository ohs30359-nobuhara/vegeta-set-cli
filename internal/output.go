package internal

import (
	"fmt"
	"ohs30359/vegeta-cli/pkg/config"
	scenarioBuilder "ohs30359/vegeta-cli/pkg/scenario"
	"os"
	"strconv"
)

func Output(path string) error {
	conf, e := config.Load(fmt.Sprintf("%s/scenario.yaml", path))
	if e != nil {
		return e
	}

	for i, scenario := range conf.Scenario {
		builder := scenarioBuilder.NewBuilder(scenario, conf.Rate, conf.Duration, path)

		targetBuf, e := builder.CreateTargetBuffer(scenario)
		if e != nil {
			return e
		}

		// "target"で指定されるファイルを作成
		targetFileName := fmt.Sprintf("scenario_%s.txt", strconv.Itoa(i))
		targetFile, e := os.Create(targetFileName)
		targetFile.Write([]byte(targetBuf))
		if e != nil {
			return e
		}

		if e := targetFile.Close(); e != nil {
			return e
		}
	}

	return nil
}
