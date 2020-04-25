package main

type ipLocation struct {
	IP          string `json:"ip"`
	CountryCode string `json:"countryCode"`
}

var data = map[string]ipLocation{
	"49.198.100.160": {
		IP:          "49.198.100.160",
		CountryCode: "AU",
	},
	"130.209.253.33": {
		IP:          "130.209.253.33",
		CountryCode: "GB",
	},
	"20.92.1.80": {
		IP:          "20.92.1.80",
		CountryCode: "US",
	},
	"165.183.1.112": {
		IP:          "165.183.1.112",
		CountryCode: "CL",
	},
	"116.125.192.145": {
		IP:          "116.125.192.145",
		CountryCode: "KR",
	},
	"10.192.168.10": {
		IP:          "10.192.168.10",
		CountryCode: "OTHER",
	},
}
