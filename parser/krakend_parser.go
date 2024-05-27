package parser

import (
	"encoding/json"
	"fmt"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

func Parse(filePath string) {
	fs := token.NewFileSet()
	f, err := parser.ParseFile(fs, filePath, nil, parser.ParseComments)
	if err != nil {
		fmt.Println(err)
		return
	}
	krakendRegex, _ := regexp.Compile("(([^/ :]+):*?):")
	getFieldRegex, _ := regexp.Compile(":(.*):")
	roleFindRegex, _ := regexp.Compile("Role:(.*)")
	queryStringFindRegex, _ := regexp.Compile("QueryStrings:(.*)")

	methodFindRegex, _ := regexp.Compile("Method:(.*)")
	endpointFindRegex, _ := regexp.Compile("Endpoint:(.*)")
	serviceFindRegex, _ := regexp.Compile("ServiceName:(.*)")
	rateLimitFindRegex, _ := regexp.Compile("RateLimit:(.*)")
	// for _, element := range f.Comments[0].List {
	// 	if krakendRegex.FindStringSubmatch(element.Text)[1] == "krakend" {
	// 		field := getFieldRegex.FindStringSubmatch(element.Text)[1]
	// 		switch field {
	// 		case "Role":
	// 			r := roleFindRegex.FindStringSubmatch(element.Text)[1]
	// 			roles := strings.Split(r, ",")
	// 			c.Roles = roles
	// 		case "Method":
	// 			c.Method = methodFindRegex.FindStringSubmatch(element.Text)[1]
	// 			fmt.Println("Method is ", methodFindRegex.FindStringSubmatch(element.Text)[1])
	// 		case "Endpoint":
	// 			c.Endpoint = endpointFindRegex.FindStringSubmatch(element.Text)[1]
	// 			fmt.Println("Endpoint is ", endpointFindRegex.FindStringSubmatch(element.Text)[1])
	// 		}
	// 	}
	// }

	var krakendConf []Endpoint
	for _, method := range f.Comments {
		c := Config{}

		for _, element := range method.List {
			if krakendRegex.FindStringSubmatch(element.Text)[1] == "krakend" {
				field := getFieldRegex.FindStringSubmatch(element.Text)[1]
				switch field {
				case "Role":
					r := roleFindRegex.FindStringSubmatch(element.Text)[1]
					roles := strings.Split(r, ",")
					c.Roles = roles
				case "QueryStrings":
					r := queryStringFindRegex.FindStringSubmatch(element.Text)[1]
					qStrings := strings.Split(r, ",")
					c.QueryStrings = qStrings
				case "Method":
					c.Method = strings.TrimSpace(methodFindRegex.FindStringSubmatch(element.Text)[1])
					fmt.Println("Method is ", methodFindRegex.FindStringSubmatch(element.Text)[1])
				case "Endpoint":
					c.Endpoint = strings.TrimSpace(endpointFindRegex.FindStringSubmatch(element.Text)[1])
					fmt.Println("Endpoint is ", endpointFindRegex.FindStringSubmatch(element.Text)[1])
				case "ServiceName":
					c.ServiceName = strings.TrimSpace(serviceFindRegex.FindStringSubmatch(element.Text)[1])
					fmt.Println("Endpoint is ", serviceFindRegex.FindStringSubmatch(element.Text)[1])
				case "RateLimit":
					c.RateLimitSpecs = strings.Split(strings.TrimSpace(rateLimitFindRegex.FindStringSubmatch(element.Text)[1]), ",")

				}
			}

		}
		endpoint := NewConfig(c)
		krakendConf = append(krakendConf, endpoint)
	}

	if _, err := os.Stat("krakend.json"); err == nil {
		fmt.Printf("File exists\n")
		file, _ := ioutil.ReadFile("krakend.json")

		var data []Endpoint

		_ = json.Unmarshal([]byte(file), &data)

		for _, endpoint := range krakendConf {
			data = append(data, endpoint)
		}
		f, _ := json.MarshalIndent(data, "", " ")
		_ = ioutil.WriteFile("krakend.json", f, 0644)

	} else {
		file, _ := json.MarshalIndent(krakendConf, "", " ")
		_ = ioutil.WriteFile("krakend.json", file, 0644)
	}

	var test KrakendConfig

	test.Endpoints = krakendConf
}

func Concat(clientID, clientSecret string) {
	file, _ := ioutil.ReadFile("krakend.json")
	var data []Endpoint

	_ = json.Unmarshal([]byte(file), &data)

	data = appendTokenEndpoint(data, clientID, clientSecret)
	krakendConfig := DefaultKrakenConfig()
	krakendConfig.Endpoints = data

	f, _ := json.MarshalIndent(krakendConfig, "", " ")
	_ = ioutil.WriteFile(fmt.Sprintf("%s.krakend.json", os.Getenv("ENVIRONMENT")), f, 0644)

}

func appendTokenEndpoint(e []Endpoint, clientID, clientSecret string) []Endpoint {
	tokenEndpoint := Endpoint{
		Method:          "POST",
		Endpoint:        "/token",
		OutputEncoding:  "no-op",
		ConcurrentCalls: 1,
		InputHeaders:    []string{"Authorization", "X-User", "X-Client", "X-Role", "Content-Type", "X-Language"},
		Backend: []Backend{{
			Host:       []string{fmt.Sprintf("%s-keycloak-http.auth.svc.cluster.local", os.Getenv("ENVIRONMENT"))},
			URLPattern: "/auth/realms/master/protocol/openid-connect/token",
			Method:     "POST",
			Encoding:   "no-op",
			ExtraConfig: &TokenExtraConfig{
				ReqRespModifier: ReqRespModifier{
					Name: []string{"shipink-token-modifier"},
					TokenModifierConfig: TokenModifierConfig{
						ClientID:     clientID,
						ClientSecret: clientSecret,
					},
				},
			}},
		}}

	e = append(e, tokenEndpoint)
	return e

}
