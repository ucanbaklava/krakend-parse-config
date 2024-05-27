package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

type SecurityCors struct {
	AllowOrigins     []string `json:"allow_origins,omitempty"`
	AllowMethods     []string `json:"allow_methods,omitempty"`
	ExposeHeaders    []string `json:"expose_headers,omitempty"`
	AllowHeaders     []string `json:"allow_headers,omitempty"`
	MaxAge           string   `json:"max_age,omitempty"`
	AllowCredentials bool     `json:"allow_credentials,omitempty"`
	Debug            bool     `json:"debug,omitempty"`
}
type JWT struct {
	Alg                string     `json:"alg"`
	JwkURL             string     `json:"jwk_url"`
	DisableJwkSecurity bool       `json:"disable_jwk_security"`
	Cache              bool       `json:"cache"`
	RolesKey           string     `json:"roles_key"`
	Roles              []string   `json:"roles"`
	PropagateClaims    [][]string `json:"propagate_claims"`
}
type Endpoint struct {
	Endpoint          string   `json:"endpoint"`
	InputQueryStrings []string `json:"input_query_strings,omitempty"`
	Method            string   `json:"method"`
	OutputEncoding    string   `json:"output_encoding"`
	ConcurrentCalls   int      `json:"concurrent_calls"`
	InputHeaders      []string `json:"input_headers,omitempty"`
	ExtraConfig       struct {
		AuthValidator      map[string]interface{} `json:"auth/validator,omitempty"`
		QosRatelimitRouter *QosRatelimitRouter    `json:"qos/ratelimit/router,omitempty"`
	} `json:"extra_config"`
	HeadersToPass []string  `json:"headers_to_pass,omitempty"`
	Backend       []Backend `json:"backend"`
}

type QosRatelimitRouter struct {
	Key           string `json:"key,omitempty"`
	ClientMaxRate int    `json:"clientMaxRate,omitempty"`
	MaxRate       int    `json:"maxRate,omitempty"`
	Strategy      string `json:"strategy,omitempty"`
}

type TokenModifierConfig struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type ReqRespModifier struct {
	Name                []string            `json:"name,omitempty"`
	TokenModifierConfig TokenModifierConfig `json:"shipink-token-modifier"`
}

type TokenExtraConfig struct {
	ReqRespModifier ReqRespModifier `json:"plugin/req-resp-modifier,omitempty"`
}

type Backend struct {
	Host        []string          `json:"host"`
	URLPattern  string            `json:"url_pattern"`
	Method      string            `json:"method"`
	Encoding    string            `json:"encoding"`
	ExtraConfig *TokenExtraConfig `json:"extra_config,omitempty"`
}

type Plugin struct {
	Pattern string `json:"pattern"`
	Folder  string `json:"folder"`
}

type TelemetryMetrics struct {
	CollectionTime string `json:"collection_time"`
	ListenAddress  string `json:"listen_address"`
}

type TelemetryLogging struct {
	Level  string `json:"level"`
	Prefix string `json:"prefix"`
	Syslog bool   `json:"syslog"`
	Stdout bool   `json:"stdout"`
}

type TelemetryOpencensus struct {
	SampleRate      int       `json:"sample_rate"`
	ReportingPeriod int       `json:"reporting_period"`
	Exporters       Exporters `json:"exporters"`
}

type Exporters struct {
	Jaeger Jaeger `json:"jaeger"`
}

type Jaeger struct {
	Endpoint    string `json:"endpoint"`
	ServiceName string `json:"service_name"`
}

type ExtraConfig struct {
	TelemetryMetrics    *TelemetryMetrics   `json:"telemetry/metrics,omitempty"`
	TelemetryLogging    TelemetryLogging    `json:"telemetry/logging"`
	TelemetryOpencensus TelemetryOpencensus `json:"telemetry/opencensus"`
	SecurityCors        *SecurityCors       `json:"security/cors,omitempty"`
}

type KrakendConfig struct {
	Version     int         `json:"version"`
	Port        int         `json:"port"`
	Timeout     string      `json:"timeout"`
	CacheTTL    string      `json:"cache_ttl"`
	Plugin      Plugin      `json:"plugin"`
	ExtraConfig ExtraConfig `json:"extra_config"`
	Endpoints   []Endpoint  `json:"endpoints"`
}

type Config struct {
	Roles          []string
	Method         string
	Endpoint       string
	ServiceName    string
	QueryStrings   []string
	RateLimitSpecs []string
}

func NewConfig(conf Config) Endpoint {
	config := Endpoint{}

	// config.Version = 3
	// config.Port = 3000
	// config.Timeout = "3000ms"
	// config.CacheTTL = "300s"
	// config.Plugin.Pattern = ".so"
	// config.Plugin.Folder = "./plugins"
	// config.ExtraConfig.TelemetryMetrics.CollectionTime = "30s"
	// config.ExtraConfig.TelemetryMetrics.ListenAddress = "80"
	// config.ExtraConfig.TelemetryLogging.Level = "DEBUG"
	// config.ExtraConfig.TelemetryLogging.Prefix = "[KRAKEND]"
	// config.ExtraConfig.TelemetryLogging.Syslog = false
	// config.ExtraConfig.TelemetryLogging.Stdout = true
	// config.ExtraConfig.TelemetryOpencensus.SampleRate = 100
	// config.ExtraConfig.TelemetryOpencensus.ReportingPeriod = 1
	// config.ExtraConfig.TelemetryOpencensus.Exporters.Jaeger.Endpoint = "http://dev-jaeger-operator-jaeger-collector.monitoring.svc.cluster.local:14268/api/traces"
	// config.ExtraConfig.TelemetryOpencensus.Exporters.Jaeger.ServiceName = "krakend"

	config.Backend = []Backend{
		{
			Host:       []string{fmt.Sprintf("%s.%s.svc.cluster.local", conf.ServiceName, conf.ServiceName)},
			URLPattern: conf.Endpoint,
			Method:     conf.Method,
			Encoding:   "no-op",
		},
	}

	config.HeadersToPass = []string{"Authorization", "Content-Type", "X-Language"}
	config.Endpoint = conf.Endpoint
	config.Method = conf.Method
	config.OutputEncoding = "no-op"
	config.InputHeaders = []string{"X-User", "X-Client", "X-Role", "X-Country", "Content-Type", "X-Language", "Authorization"}
	config.ConcurrentCalls = 1

	if stringContains(conf.Roles, "Guest") {
		config.ExtraConfig.AuthValidator = map[string]interface{}{}
	} else {
		jwtConfig := JWT{
			Alg:                "RS256",
			JwkURL:             fmt.Sprintf("http://%s-keycloak-http.auth.svc.cluster.local/auth/realms/master/protocol/openid-connect/certs", os.Getenv("ENVIRONMENT")),
			DisableJwkSecurity: true,
			Cache:              true,
			RolesKey:           "roles",
			Roles:              conf.Roles,
			PropagateClaims: [][]string{
				{"sub", "x-user"},
				{"azp", "x-client"},
				{"roles", "x-role"},
				{"country_code", "x-country"},
			},
		}

		var jwt map[string]interface{}
		data, _ := json.Marshal(jwtConfig)
		json.Unmarshal(data, &jwt)
		config.ExtraConfig.AuthValidator = jwt
	}

	if conf.QueryStrings != nil {
		config.InputQueryStrings = conf.QueryStrings
	}

	if conf.RateLimitSpecs != nil {
		clientMaxRate, _ := strconv.Atoi(conf.RateLimitSpecs[1])

		if conf.RateLimitSpecs[0] == "header" {
			rateLimitConf := &QosRatelimitRouter{
				Key:           "X-User",
				Strategy:      conf.RateLimitSpecs[0],
				ClientMaxRate: clientMaxRate,
			}
			config.ExtraConfig.QosRatelimitRouter = rateLimitConf
		} else {
			rateLimitConf := &QosRatelimitRouter{
				Strategy:      conf.RateLimitSpecs[0],
				ClientMaxRate: clientMaxRate,
			}
			config.ExtraConfig.QosRatelimitRouter = rateLimitConf
		}

	}
	return config
}

func stringContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func DefaultKrakenConfig() KrakendConfig {
	return KrakendConfig{
		Version:  3,
		Port:     3000,
		Timeout:  "10000ms",
		CacheTTL: "300s",
		Plugin: Plugin{
			Pattern: ".so",
			Folder:  "./plugins",
		},
		ExtraConfig: ExtraConfig{
			// TelemetryMetrics: TelemetryMetrics{
			// 	CollectionTime: "30s",
			// 	ListenAddress:  ":80",
			// },
			SecurityCors: &SecurityCors{
				AllowOrigins:     []string{"http://localhost:3000"},
				AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
				ExposeHeaders:    []string{"Content-Type", "Content-Length"},
				AllowHeaders:     []string{"Accept-Language"},
				MaxAge:           "12h",
				AllowCredentials: true,
				Debug:            false,
			},
			TelemetryLogging: TelemetryLogging{
				Level:  "DEBUG",
				Prefix: "[KRAKEND]",
				Syslog: false,
				Stdout: true,
			},
			TelemetryOpencensus: TelemetryOpencensus{
				SampleRate:      100,
				ReportingPeriod: 1,
				Exporters: Exporters{
					Jaeger: Jaeger{
						Endpoint:    fmt.Sprintf("http://%s-jaeger-operator-jaeger-collector.monitoring.svc.cluster.local:14268/api/traces", os.Getenv("ENVIRONMENT")),
						ServiceName: "krakend",
					},
				},
			},
		},
	}
}
