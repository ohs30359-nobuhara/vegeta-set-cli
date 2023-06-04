# vegeta-set-cli

## Description
vegetaの実行に必要な各種設定ファイルを自動生成するためのツールです

## Usage

### シナリオの準備
以下のようなディレクトリを用意してください。

```
scenario
    ├── scenario.yaml -- テストシナリオ
    └── values -- 攻撃に必要なパラメータ管理用ディレクトリ
```

#### scenario.yaml 

```
scenario:
  - url: http://localhost:3000
    method: GET
    ratio: 10 -- 全体負荷のn%を割り当てるか
  - url: http://localhost:3000/item
    method: GET
    ratio: 45
    value: ./values/item
  - url: http://localhost:3000/user
    method: POST
    ratio: 45
    value: ./values/user
tester:
  limit: 100 -- vegetaの性能限界 (req/sec)
rate: 1000  -- req/sec 
duration: 100s -- 攻撃の実行時間
```


#### values dir  
このディレクトリは各scenarioで用いる各種パラメータを管理するために用います。  

```
scenario
    ├── scenario.yaml
    └── values
        ├── item
        │   └── queryParams.txt
        └── user
            ├── body_1.json
            └── body_2.json
```


■ GET (query params)   
クエリパラメータは `.txt` で用意してください。

```
name=本
name=紙&price=100
```


■ POST (body)    
bodyは各パラメータをcontent-typeに合わせたファイルで用意してください

```
{
  "name": "山田"
}
```

### シナリオ作成
以下のコマンドで vegeta用の実行ファイルが`dist`ディレクトリに生成されます
```
go run main.go -s ./scenario
```

上記のサンプルであれば以下のような実行ファイルが生成されます   
※ scenario.sh が複数生成されているのは `rate`が `tester.limit` を超えているため複数台での実行を前提とし結果
```  
dist
├── scenario_0
│   ├── scenario_0.sh
│   └── target.txt
├── scenario_1
│   ├── scenario_0.sh
│   ├── scenario_1.sh
│   ├── scenario_2.sh
│   ├── scenario_3.sh
│   ├── scenario_4.sh
│   └── target.txt
└── scenario_2
    ├── scenario_0.sh
    ├── scenario_1.sh
    ├── scenario_2.sh
    ├── scenario_3.sh
    ├── scenario_4.sh
    ├── target.txt
    └── values
        ├── val1.json
        └── val2.json
```

後は `dist`配下に移動して scenario.sh を実行すればvegetaが実行されます。   

