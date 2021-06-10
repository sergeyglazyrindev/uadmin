package config

import (
	"github.com/go-openapi/loads"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

// DBSettings !
type DBSettings struct {
	Type     string `json:"type"` // sqlite, mysql
	Name     string `json:"name"` // File/DB name
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
}

type UadminConfigOptions struct {
	Theme string `yaml:"theme"`
	SiteName string `yaml:"site_name"`
	ReportingLevel int `yaml:"reporting_level"`
	ReportTimeStamp bool `yaml:"report_timestamp"`
	DebugDB bool `yaml:"debug_db"`
	PageLength int `yaml:"page_length"`
	MaxImageHeight int `yaml:"max_image_height"`
	MaxImageWidth int `yaml:"max_image_width"`
	MaxUploadFileSize int64 `yaml:"max_upload_file_size"`
	EmailFrom string `yaml:"email_from"`
	EmailUsername string `yaml:"email_username"`
	EmailPassword string `yaml:"email_password"`
	EmailSmtpServer string `yaml:"email_smtp_server"`
	EmailSmtpServerPort int `yaml:"email_smtp_server_port"`
	RootURL string `yaml:"root_url"`
	OTPAlgorithm string `yaml:"otp_algorithm"`
	OTPDigits int `yaml:"otp_digits"`
	OTPPeriod uint `yaml:"otp_period"`
	OTPSkew uint `yaml:"otp_skew"`
	PublicMedia bool `yaml:"public_media"`
	LogDelete bool `yaml:"log_delete"`
	LogAdd bool `yaml:"log_add"`
	LogEdit bool `yaml:"log_edit"`
	LogRead bool `yaml:"log_read"`
	CacheTranslation bool `yaml:"cache_translation"`
	AllowedIPs string `yaml:"allowed_ips"`
	BlockedIPs string `yaml:"blocked_ips"`
	RestrictSessionIP bool `yaml:"restrict_session_ip"`
	RetainMediaVersions bool `yaml:"retain_media_versions"`
	RateLimit uint `yaml:"rate_limit"`
	RateLimitBurst uint `yaml:"rate_limit_burst"`
	APILogRead bool `yaml:"api_log_read"`
	APILogDelete bool `yaml:"api_log_delete"`
	APILogAdd bool `yaml:"api_log_add"`
	APILogEdit bool `yaml:"api_log_edit"`
	LogHTTPRequests bool `yaml:"log_http_requests"`
	HTTPLogFormat string `yaml:"http_log_format"`
	LogTrail bool `yaml:"log_trail"`
	TrailLoggingLevel int `yaml:"trail_logging_level"`
	SystemMetrics bool `yaml:"system_metrics"`
	UserMetrics bool `yaml:"user_metrics"`
	PasswordAttempts int `yaml:"password_attempts"`
	PasswordTimeout int `yaml:"password_timeout"`
	AllowedHosts string `yaml:"allowed_hosts"`
	Logo string `yaml:"logo"`
	FavIcon string `yaml:"fav_icon"`
}

type UadminDbOptions struct {
	Default *DBSettings
}

type UadminAuthOptions struct {
	JWT_SECRET_TOKEN string `yaml:"jwt_secret_token"`
	MinUsernameLength int `yaml:"min_username_length"`
	MaxUsernameLength int `yaml:"max_username_length"`
	MinPasswordLength int `yaml:"min_password_length"`
	SaltLength int `yaml:"salt_length"`
}

type UadminAdminOptions struct {
	ListenPort int `yaml:"listen_port"`
	SSL        struct {
		ListenPort int `yaml:"listen_port"`
	} `yaml:"ssl"`
}

type UadminApiOptions struct {
	ListenPort int `yaml:"listen_port"`
	SSL        struct {
		ListenPort int `yaml:"listen_port"`
	} `yaml:"ssl"`
}

type UadminSwaggerOptions struct {
	ListenPort int `yaml:"listen_port"`
	SSL        struct {
		ListenPort int `yaml:"listen_port"`
	} `yaml:"ssl"`
	PathToSpec string `yaml:"path_to_spec"`
	ApiEditorListenPort int `yaml:"api_editor_listen_port"`
}

type UadminConfigurableConfig struct {
	Uadmin *UadminConfigOptions `yaml:"uadmin"`
	Test string `yaml:"test"`
	Db *UadminDbOptions `yaml:"db"`
	Auth *UadminAuthOptions `yaml:"auth"`
	Admin *UadminAdminOptions `yaml:"admin"`
	Api *UadminApiOptions `yaml:"api"`
	Swagger *UadminSwaggerOptions `yaml:"swagger"`
}

// Info from config file
type UadminConfig struct {
	ApiSpec *loads.Document
	D *UadminConfigurableConfig
}

func (ucc *UadminConfigurableConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type rawStuff UadminConfigurableConfig
	raw := rawStuff{
		Auth: &UadminAuthOptions{SaltLength: 16},
		Uadmin: &UadminConfigOptions{
			Theme: "default",
			SiteName: "uadmin",
			ReportingLevel: 0,
			ReportTimeStamp: false,
			DebugDB: false,
			PageLength: 100,
			MaxImageHeight: 600,
			MaxImageWidth: 800,
			MaxUploadFileSize: int64(25 * 1024 * 1024),
			RootURL: "/",
			OTPAlgorithm: "sha1",
			OTPDigits: 6,
			OTPPeriod: uint(30),
			OTPSkew: uint(5),
			PublicMedia: false,
			LogDelete: true,
			LogAdd: true,
			LogEdit: true,
			LogRead: false,
			CacheTranslation: false,
			AllowedIPs: "*",
			BlockedIPs: "",
			RestrictSessionIP: false,
			RetainMediaVersions: true,
			RateLimit: uint(3),
			RateLimitBurst: uint(3),
			APILogRead: false,
			APILogEdit: true,
			APILogAdd: true,
			APILogDelete: true,
			LogHTTPRequests: true,
			HTTPLogFormat: "%a %>s %B %U %D",
			LogTrail: false,
			TrailLoggingLevel: 2,
			SystemMetrics: false,
			UserMetrics: false,
			PasswordAttempts: 5,
			PasswordTimeout: 15,
			AllowedHosts: "0.0.0.0,127.0.0.1,localhost,::1",
			Logo: "/static-inbuilt/uadmin/logo.png",
			FavIcon: "/static-inbuilt/uadmin/favicon.ico",
		},
	}
	// Put your defaults here
	if err := unmarshal(&raw); err != nil {
		return err
	}

	*ucc = UadminConfigurableConfig(raw)
	return nil

}

var CurrentConfig *UadminConfig

// Reads info from config file
func NewConfig(file string) *UadminConfig {
	_, err := os.Stat(file)
	if err != nil {
		log.Fatal("Config file is missing: ", file)
	}
	content, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	c := UadminConfig{}
	err = yaml.Unmarshal([]byte(content), &c.D)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	CurrentConfig = &c
	return &c
}

// Reads info from config file
func NewSwaggerSpec(file string) *loads.Document {
	_, err := os.Stat(file)
	if err != nil {
		log.Fatal("Config file is missing: ", file)
	}
	doc, err := loads.Spec(file)
	if err != nil {
		panic(err)
	}
	return doc
}
