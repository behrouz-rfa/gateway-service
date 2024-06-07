package domain

type GeolocationResponse struct {
	IP           string       `json:"ip"`
	CountryCode  string       `json:"country_code"`
	CountryName  string       `json:"country_name"`
	RegionName   string       `json:"region_name"`
	CityName     string       `json:"city_name"`
	Latitude     float64      `json:"latitude"`
	Longitude    float64      `json:"longitude"`
	ZipCode      string       `json:"zip_code"`
	TimeZone     string       `json:"time_zone"`
	ASN          string       `json:"asn"`
	ISP          string       `json:"isp"`
	Domain       string       `json:"domain"`
	Continent    Continent    `json:"continent"`
	Country      Country      `json:"country"`
	Region       Region       `json:"region"`
	City         City         `json:"city"`
	TimeZoneInfo TimeZoneInfo `json:"time_zone_info"`
	Proxy        Proxy        `json:"proxy"`
}

type Continent struct {
	Name        string      `json:"name"`
	Code        string      `json:"code"`
	Hemisphere  []string    `json:"hemisphere"`
	Translation Translation `json:"translation"`
}

type Country struct {
	Name        string      `json:"name"`
	Alpha3Code  string      `json:"alpha3_code"`
	NumericCode int         `json:"numeric_code"`
	Demonym     string      `json:"demonym"`
	Flag        string      `json:"flag"`
	Capital     string      `json:"capital"`
	TotalArea   int         `json:"total_area"`
	Population  int         `json:"population"`
	Currency    Currency    `json:"currency"`
	Language    Language    `json:"language"`
	TLD         string      `json:"tld"`
	Translation Translation `json:"translation"`
}

type Region struct {
	Name        string      `json:"name"`
	Code        string      `json:"code"`
	Translation Translation `json:"translation"`
}

type City struct {
	Name        string      `json:"name"`
	Translation Translation `json:"translation"`
}

type TimeZoneInfo struct {
	Olson       string `json:"olson"`
	CurrentTime string `json:"current_time"`
	GMTOffset   int    `json:"gmt_offset"`
	IsDST       bool   `json:"is_dst"`
	Sunrise     string `json:"sunrise"`
	Sunset      string `json:"sunset"`
}

type Proxy struct {
	LastSeen                   int    `json:"last_seen"`
	ProxyType                  string `json:"proxy_type"`
	Threat                     string `json:"threat"`
	IsVPN                      bool   `json:"is_vpn"`
	IsTor                      bool   `json:"is_tor"`
	IsDataCenter               bool   `json:"is_data_center"`
	IsPublicProxy              bool   `json:"is_public_proxy"`
	IsWebProxy                 bool   `json:"is_web_proxy"`
	IsWebCrawler               bool   `json:"is_web_crawler"`
	IsResidentialProxy         bool   `json:"is_residential_proxy"`
	IsConsumerPrivacyNetwork   bool   `json:"is_consumer_privacy_network"`
	IsEnterprisePrivateNetwork bool   `json:"is_enterprise_private_network"`
	IsSpammer                  bool   `json:"is_spammer"`
	IsScanner                  bool   `json:"is_scanner"`
	IsBotnet                   bool   `json:"is_botnet"`
}

type Currency struct {
	Code   string `json:"code"`
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}

type Language struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type Translation struct {
	Lang  string `json:"lang"`
	Value string `json:"value"`
}
