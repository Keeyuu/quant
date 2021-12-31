package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"gopkg.in/yaml.v2"
)

type App struct {
	Secret    string    `yaml:"secret"`
	Name      string    `yaml:"name"`
	Http      Http      `yaml:"http"`
	GRpc      GRpc      `yaml:"grpc"`
	Log       Log       `yaml:"log"`
	DB        DB        `yaml:"db"`
	Cache     Cache     `yaml:"cache"`
	Redis     Redis     `yaml:"redis"`
	Scheduler Scheduler `yaml:"scheduler"`
}

type OrderStatsConfig struct {
	AlarmThreshold int64  `yaml:"alarmThreshold"`
	S3Bucket       string `yaml:"s3Bucket"`
	S3Path         string `yaml:"s3Path"`
	FileUrl        string `yaml:"fileUrl"`
}

type Dingtalk struct {
	NotifyUrl string `yaml:"notifyUrl"`
	SecretKey string `yaml:"secretKey"`
}

type MerchantComment struct {
	Bonus      map[string]int64   `yaml:"bonus"`
	BonusRange map[string][]int64 `yaml:"bonusRange"`
}

type Es struct {
	Url        string `yaml:"url"`
	StoreIndex string `yaml:"store_index"`
	Distance   int64  `yaml:"distance"`
}

type Login struct {
	ExpireInterval int64 `yaml:"expireInterval"`
}

type BusinessConfig struct {
	FirestoreConfig FirestoreConfig `yaml:"firestoreConfig"`
}

type FirestoreConfig struct {
	C map[string]string `yaml:"c"`
	B map[string]string `yaml:"b"`
}

type DataSync struct {
	Enable   bool   `yaml:"enable"`
	Project  string `yaml:"project"`
	EndPoint string `yaml:"endPoint"`
	Region   string `yaml:"region"`
}

type SqsConfig struct {
	Enable    bool   `yaml:"enable"`
	QueueUrl  string `yaml:"queueUrl"`
	Region    string `yaml:"region"`
	SendDelay int64  `yaml:"sendDelay"`
	GroupId   string `yaml:"groupId"`
}

type Http struct {
	Port        string     `yaml:"port"`
	Clients     HttpClient `yaml:"clients"`
	WaitTimeout int        `yaml:"waitTimeout"`
	Api         Api        `yaml:"api"`
}

type Api struct {
	UserBalance UserBalance `yaml:"userBalance"`
}

type HttpClient struct {
	FileCenter string `yaml:"fileCenter"`
}

type GRpc struct {
	Port      string     `yaml:"port"`
	Keepalive int        `yaml:"keepalive"`
	Clients   GRpcClient `yaml:"clients"`
}

type GRpcClient struct {
	PMS         string `yaml:"pms"`
	Order       string `yaml:"order"`
	User        string `yaml:"user"`
	RebateOrder string `yaml:"rebateOrder"`
	Balance     string `yaml:"balance"`
	Activity    string `yaml:"activity"`
	DishesTable string `yaml:"dishesTable"`
}

type Log struct {
	Level         int  `yaml:"level"`
	DisableCaller bool `yaml:"disableCaller"`
}

type Scheduler struct {
	DemoJobCron string `yaml:"demoJobCron"`
}

type DB struct {
	Name                             string        `yaml:"name"`
	URI                              string        `yaml:"uri"`
	ConnectionTimeout                time.Duration `yaml:"connectTimeout"`
	UserTable                        TableConfig   `yaml:"userTable"`
	ShopTable                        TableConfig   `yaml:"shopTable"`
	SupplierTable                    TableConfig   `yaml:"supplierTable"`
	CategoryTable                    TableConfig   `yaml:"categoryTable"`
	HistoryTable                     TableConfig   `yaml:"historyTable"`
	CommentTable                     TableConfig   `yaml:"commentTable"`
	OfflineOrderTable                TableConfig   `yaml:"offlineOrderTable"`
	OfflineOrderOperationRecordTable TableConfig   `yaml:"offlineOrderOperationRecordTable"`
	ShopDistrictTable                TableConfig   `yaml:"shopDistrictTable"`
	ApplyRecordTable                 TableConfig   `yaml:"applyRecordTable"`
	ShopTagTable                     TableConfig   `yaml:"shopTagTable"`
	SettlementTable                  TableConfig   `yaml:"settlementTable"`
	UserOrderStatsTable              TableConfig   `yaml:"userOrderStatsTable"`
	DeviceOrderStatsTable            TableConfig   `yaml:"deviceOrderStatsTable"`
	DeliverySettingTable             TableConfig   `yaml:"deliverySettingTable"`
}

type Redis struct {
	Addr              string        `yaml:"addr"`
	Password          string        `yaml:"password"`
	ConnectionTimeout time.Duration `yaml:"connectTimeout"`
}

type AppMessage struct {
	Api    map[string]string `yaml:"api"`
	PicUrl map[string]string `yaml:"picUrl"`
	Link   string            `yaml:"link"`
	Secret string            `yaml:"secret"`
}

const (
	EmptyString       = ""
	EnvAppConfigURL   = "APP_CONFIG_URL"
	EnvAppConfigPath  = "APP_CONFIG_PATH"
	DefaultConfigPath = "config.yaml"
)

var appConfigMutex sync.Mutex
var appConfig App

func Get() App {
	return appConfig
}

func init() {
	appConfigMutex.Lock()
	defer appConfigMutex.Unlock()
	// load readConfig file
	configPath := os.Getenv(EnvAppConfigURL)
	if configPath == EmptyString {
		configPath = os.Getenv(EnvAppConfigPath)
	}
	if configPath == EmptyString {
		configPath = DefaultConfigPath
	}

	//parse readConfig file
	var err error
	var fileBytes []byte
	if fileBytes, err = readConfig(configPath); err != nil {
		fmt.Printf("Load readConfig File Error: %v\n", err)
		os.Exit(3)
		return
	}

	//unmarshal readConfig file
	if err = unmarshal(fileBytes, &appConfig, false); err != nil {
		fmt.Printf("unmarshal readConfig File Error: %v\n", err)
		os.Exit(3)
		return
	}

	// set default value when value is nil
	setDefaultValue(&appConfig)
}

// internal
func readConfig(location string) (bytes []byte, err error) {
	if strings.HasPrefix(location, "http") {
		return remoteConfig(location)
	}
	return localConfig(location)
}

func localConfig(filePath string) (bytes []byte, err error) {
	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return fileBytes, nil
}

func remoteConfig(url string) (bytes []byte, err error) {
	if url == "" {
		return nil, errors.New("can't get readConfig url")
	}
	result, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return nil, err
	}
	return b, err
}

func unmarshal(in []byte, out interface{}, isStrict bool) (err error) {
	if in == nil {
		err = error(errors.New("can't unmarshal empty byte"))
		return err
	}
	if isStrict == true {
		err = yaml.UnmarshalStrict(in, out)
		if err != nil {
			return err
		}
	} else {
		err = yaml.Unmarshal(in, out)
		if err != nil {
			return err
		}
	}
	return nil
}

func setDefaultValue(config *App) {
	if config.Cache.ExpireInterval <= 0 {
		fmt.Println("load cache:expireInterval fail, use default value")
		config.Cache.ExpireInterval = 5
	}
	if config.Http.WaitTimeout == 0 {
		fmt.Println("load http:waitTimeout fail, use default value")
		config.Http.WaitTimeout = 10
	}
	if config.Obs.S3.Region == "" {
		fmt.Println("load s3 region fail, use default value")
		config.Obs.S3.Region = "ap-southeast-1"
	}
}
