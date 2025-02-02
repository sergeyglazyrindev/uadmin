package migrations

import (
	"fmt"
	settingmodel "github.com/sergeyglazyrindev/uadmin/blueprint/settings/models"
	"github.com/sergeyglazyrindev/uadmin/core"
	"strings"
)

type insertall1623263908 struct {
}

func (m insertall1623263908) GetName() string {
	return "settings.1623263908"
}

func (m insertall1623263908) GetID() int64 {
	return 1623263908
}

func (m insertall1623263908) Up(uadminDatabase *core.UadminDatabase) error {
	// Check if the uAdmin category is not there and add it
	db := uadminDatabase.Db
	var uadminSettingcategory settingmodel.SettingCategory
	db.Model(&settingmodel.SettingCategory{}).Where(&settingmodel.SettingCategory{Name: "uAdmin"}).First(&uadminSettingcategory)
	if uadminSettingcategory.ID == 0 {
		uadminSettingcategory = settingmodel.SettingCategory{Name: "uAdmin"}
		db.Create(&uadminSettingcategory)
	}
	t := settingmodel.DataType(0)

	settings := []settingmodel.Setting{
		{
			Name:         "Theme",
			Value:        core.CurrentConfig.D.Uadmin.Theme,
			DefaultValue: "default",
			DataType:     t.String(),
			Help:         "is the name of the theme used in uAdmin",
		},
		{
			Name:         "Site Name",
			Value:        core.CurrentConfig.D.Uadmin.SiteName,
			DefaultValue: "uAdmin",
			DataType:     t.String(),
			Help:         "is the name of the website that shows on title and dashboard",
		},
		{
			Name:         "Reporting Level",
			Value:        fmt.Sprint(core.CurrentConfig.D.Uadmin.ReportingLevel),
			DefaultValue: "0",
			DataType:     t.Integer(),
			Help:         "Reporting level. DEBUG=0, WORKING=1, INFO=2, OK=3, WARNING=4, ERROR=5",
		},
		{
			Name:         "Report Time Stamp",
			Value:        fmt.Sprint(core.CurrentConfig.D.Uadmin.ReportTimeStamp),
			DefaultValue: "0",
			DataType:     t.Boolean(),
			Help:         "set this to true to have a time stamp in your logs",
		},
		{
			Name: "Debug DB",
			Value: func(v bool) string {
				n := 0
				if v {
					n = 1
				}
				return fmt.Sprint(n)
			}(core.CurrentConfig.D.Uadmin.DebugDB),
			DefaultValue: "0",
			DataType:     t.Boolean(),
			Help:         "prints all SQL statements going to DB",
		},
		{
			Name:         "Page Length",
			Value:        fmt.Sprint(core.CurrentConfig.D.Uadmin.PageLength),
			DefaultValue: "100",
			DataType:     t.Integer(),
			Help:         "is the list view max number of records",
		},
		{
			Name:         "Max Image Height",
			Value:        fmt.Sprint(core.CurrentConfig.D.Uadmin.MaxImageHeight),
			DefaultValue: "600",
			DataType:     t.Integer(),
			Help:         "sets the maximum height of an Image",
		},
		{
			Name:         "Max Image Width",
			Value:        fmt.Sprint(core.CurrentConfig.D.Uadmin.MaxImageWidth),
			DefaultValue: "800",
			DataType:     t.Integer(),
			Help:         "sets the maximum width of an image",
		},
		{
			Name:         "Max Upload File Size",
			Value:        fmt.Sprint(core.CurrentConfig.D.Uadmin.MaxUploadFileSize),
			DefaultValue: "26214400",
			DataType:     t.Integer(),
			Help:         "is the maximum upload file size in bytes. 1MB = 1024 * 1024",
		},
		{
			Name:         "Root URL",
			Value:        core.CurrentConfig.D.Uadmin.RootAdminURL,
			DefaultValue: core.CurrentConfig.D.Uadmin.RootAdminURL,
			DataType:     t.String(),
			Help:         "is where the listener is mapped to",
		},
		{
			Name:         "OTP Algorithm",
			Value:        core.CurrentConfig.D.Uadmin.OTPAlgorithm,
			DefaultValue: "sha1",
			DataType:     t.String(),
			Help:         "is the hashing algorithm of OTP. Other options are sha256 and sha512",
		},
		{
			Name:         "OTP Digits",
			Value:        fmt.Sprint(core.CurrentConfig.D.Uadmin.OTPDigits),
			DefaultValue: "6",
			DataType:     t.Integer(),
			Help:         "is the number of digits for the OTP",
		},
		{
			Name:         "OTP Period",
			Value:        fmt.Sprint(core.CurrentConfig.D.Uadmin.OTPPeriod),
			DefaultValue: "30",
			DataType:     t.Integer(),
			Help:         "the number of seconds for the OTP to change",
		},
		{
			Name:         "OTP Skew",
			Value:        fmt.Sprint(core.CurrentConfig.D.Uadmin.OTPSkew),
			DefaultValue: "5",
			DataType:     t.Integer(),
			Help:         "is the number of minutes to search around the OTP",
		},
		{
			Name: "Public Media",
			Value: func(v bool) string {
				n := 0
				if v {
					n = 1
				}
				return fmt.Sprint(n)
			}(core.CurrentConfig.D.Uadmin.PublicMedia),
			DefaultValue: "0",
			DataType:     t.Boolean(),
			Help:         "allows public access to media handler without authentication",
		},
		{
			Name: "Restrict Session IP",
			Value: func(v bool) string {
				n := 0
				if v {
					n = 1
				}
				return fmt.Sprint(n)
			}(core.CurrentConfig.D.Uadmin.RestrictSessionIP),
			DefaultValue: "0",
			DataType:     t.Boolean(),
			Help:         "is to block access of a user if their IP changes from their original IP during login",
		},
		{
			Name: "Retain Media Versions",
			Value: func(v bool) string {
				n := 0
				if v {
					n = 1
				}
				return fmt.Sprint(n)
			}(core.CurrentConfig.D.Uadmin.RetainMediaVersions),
			DefaultValue: "1",
			DataType:     t.Boolean(),
			Help:         "is to allow the system to keep files uploaded even after they are changed. This allows the system to \"Roll Back\" to an older version of the file",
		},
		{
			Name:         "Rate Limit",
			Value:        fmt.Sprint(core.CurrentConfig.D.Uadmin.RateLimit),
			DefaultValue: "3",
			DataType:     t.Integer(),
			Help:         "is the maximum number of requests/second for any unique IP",
		},
		{
			Name:         "Rate Limit Burst",
			Value:        fmt.Sprint(core.CurrentConfig.D.Uadmin.RateLimitBurst),
			DefaultValue: "3",
			DataType:     t.Integer(),
			Help:         "is the maximum number of requests for an idle user",
		},
		{
			Name: "Log HTTP Requests",
			Value: func(v bool) string {
				n := 0
				if v {
					n = 1
				}
				return fmt.Sprint(n)
			}(core.CurrentConfig.D.Uadmin.LogHTTPRequests),
			DefaultValue: "1",
			DataType:     t.Boolean(),
			Help:         "Logs http requests to syslog",
		},
		{
			Name:         "HTTP Log Format",
			Value:        core.CurrentConfig.D.Uadmin.HTTPLogFormat,
			DefaultValue: "",
			DataType:     t.String(),
			Help: `Is the format used to log HTTP access
									%a: Client IP address
									%{remote}p: Client port
									%A: Server hostname/IP
									%{local}p: Server port
									%U: Path
									%c: All coockies
									%{NAME}c: Cookie named 'NAME'
									%{GET}f: GET request parameters
									%{POST}f: POST request parameters
									%B: Response length
									%>s: Response code
									%D: Time taken in microseconds
									%T: Time taken in seconds
									%I: Request length`,
		},
		{
			Name: "Log Trail",
			Value: func(v bool) string {
				n := 0
				if v {
					n = 1
				}
				return fmt.Sprint(n)
			}(core.CurrentConfig.D.Uadmin.LogTrail),
			DefaultValue: "0",
			DataType:     t.Boolean(),
			Help:         "Stores Trail logs to syslog",
		},
		{
			Name:         "Trail Logging Level",
			Value:        fmt.Sprint(core.CurrentConfig.D.Uadmin.TrailLoggingLevel),
			DefaultValue: "2",
			DataType:     t.Integer(),
			Help:         "Is the minimum level to be logged into syslog.",
		},
		{
			Name: "System Metrics",
			Value: func(v bool) string {
				n := 0
				if v {
					n = 1
				}
				return fmt.Sprint(n)
			}(core.CurrentConfig.D.Uadmin.SystemMetrics),
			DefaultValue: "0",
			DataType:     t.Boolean(),
			Help:         "Enables uAdmin system metrics to be recorded",
		},
		{
			Name: "User Metrics",
			Value: func(v bool) string {
				n := 0
				if v {
					n = 1
				}
				return fmt.Sprint(n)
			}(core.CurrentConfig.D.Uadmin.UserMetrics),
			DefaultValue: "0",
			DataType:     t.Boolean(),
			Help:         "Enables the user metrics to be recorded",
		},
		{
			Name:         "Password Attempts",
			Value:        fmt.Sprint(core.CurrentConfig.D.Uadmin.PasswordAttempts),
			DefaultValue: "5",
			DataType:     t.Integer(),
			Help:         "The maximum number of invalid password attempts before the IP address is blocked for some time from usig the system",
		},
		{
			Name:         "Password Timeout",
			Value:        fmt.Sprint(core.CurrentConfig.D.Uadmin.PasswordTimeout),
			DefaultValue: "5",
			DataType:     t.Integer(),
			Help:         "The maximum number of invalid password attempts before the IP address is blocked for some time from usig the system",
		},
		{
			Name:         "Logo",
			Value:        core.CurrentConfig.D.Uadmin.Logo,
			DefaultValue: "/static-inbuilt/uadmin/logo.png",
			DataType:     t.Image(),
			Help:         "the main logo that shows on uAdmin UI",
		},
		{
			Name:         "Fav Icon",
			Value:        core.CurrentConfig.D.Uadmin.FavIcon,
			DefaultValue: "/static-inbuilt/uadmin/favicon.ico",
			DataType:     t.File(),
			Help:         "the fav icon that shows on uAdmin UI",
		},
	}

	// Prepare uAdmin Settings
	for i := range settings {
		settings[i].CategoryID = uadminSettingcategory.ID
		settings[i].Code = "uAdmin." + strings.Replace(settings[i].Name, " ", "", -1)
	}
	// Check if the settings exist in the DB
	var s settingmodel.Setting
	sList := []settingmodel.Setting{}
	db.Model(&settingmodel.Setting{}).Where(&settingmodel.Setting{CategoryID: uadminSettingcategory.ID}).Find(&sList)
	tx := db
	for _, setting := range settings {
		s = settingmodel.Setting{}
		for c := range sList {
			if sList[c].Code == setting.Code {
				s = sList[c]
			}
		}
		if s.ID == 0 {
			tx.Create(&setting)
		} else {
			if s.DefaultValue != setting.DefaultValue || s.Help != setting.Help {
				if s.Help != setting.Help {
					s.Help = setting.Help
				}
				if s.Value == s.DefaultValue {
					s.Value = setting.DefaultValue
				}
				s.DefaultValue = setting.DefaultValue
				tx.Save(s)
				//s.Save()
			}
		}
	}
	return nil
}

func (m insertall1623263908) Down(uadminDatabase *core.UadminDatabase) error {
	db := uadminDatabase.Db
	db.Unscoped().Where("1 = 1").Delete(&settingmodel.Setting{})
	return nil
}

func (m insertall1623263908) Deps() []string {
	return []string{"settings.1623082592"}
}
