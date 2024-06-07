package domain

// Core represents data from the Core API.
type IpInfo struct {
	IP             string `json:"ip"`
	Hostname       string `json:"hostname,omitempty" csv:"hostname" yaml:"hostname,omitempty"`
	Bogon          bool   `json:"bogon,omitempty" csv:"bogon" yaml:"bogon,omitempty"`
	Anycast        bool   `json:"anycast,omitempty" csv:"anycast" yaml:"anycast,omitempty"`
	City           string `json:"city,omitempty" csv:"city" yaml:"city,omitempty"`
	Region         string `json:"region,omitempty" csv:"region" yaml:"region,omitempty"`
	Country        string `json:"country,omitempty" csv:"country" yaml:"country,omitempty"`
	CountryFlagURL string `json:"country_flag_url,omitempty" csv:"country_flag_url" yaml:"countryFlagURL,omitempty"`
	IsEU           bool   `json:"isEU,omitempty" csv:"isEU" yaml:"isEU,omitempty"`
	Location       string `json:"loc,omitempty" csv:"loc" yaml:"location,omitempty"`
	Org            string `json:"org,omitempty" csv:"org" yaml:"org,omitempty"`
	Postal         string `json:"postal,omitempty" csv:"postal" yaml:"postal,omitempty"`
	Timezone       string
}
