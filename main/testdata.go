package main

type IPLocation struct {
	IP          string `json:"ip"`
	CountryCode string `json:"countryCode"`
}

var data = map[string]IPLocation{
	"49.198.100.160": IPLocation{
		IP:          "49.198.100.160",
		CountryCode: "AU",
	},
	"130.209.253.33": IPLocation{
		IP:          "130.209.253.33",
		CountryCode: "GB",
	},
	"20.92.1.80": IPLocation{
		IP:          "20.92.1.80",
		CountryCode: "US",
	},
	"165.183.1.112": IPLocation{
		IP:          "165.183.1.112",
		CountryCode: "CL",
	},
	"116.125.192.145": IPLocation{
		IP:          "116.125.192.145",
		CountryCode: "KR",
	},
	"10.192.168.10": IPLocation{
		IP:          "10.192.168.10",
		CountryCode: "OTHER",
	},
}
