package ipdb

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
)

const (
	NULL = "N/A"

	STRING_EMPTY    = ""
	STRING_CHINA    = "中国"
	STRING_CHINA_DC = "香港/澳门/台湾"

	FILENAME_COUNTRY  = "utils/codeCountry.json"
	FILENAME_REGION   = "utils/codeRegion.json"
	FILENAME_CITY     = "utils/codeCity.json"
	FILENAME_DATABASE = "utils/ipipfree.ipdb"
)

var (
	ERR_EMPTY_CITYINFO            = errors.New("empty CityInfo!")
	ERR_NOT_FOUND_COUNTRY         = errors.New("not found ContinentCode & IddCode")
	ERR_NOT_FOUND_MAINLAND_REGION = errors.New("not found ChinaAdminCode")
	ERR_NOT_FOUND_MAINLAND_CITY   = errors.New("not found ChinaCityCode")
)

var IPupgrader *Upgrader

func init() {
	IPupgrader, _ = NewUpgrader(FILENAME_DATABASE)
}

// Upgrader upgrades the free version ipdb CityInfo to custom UpgradeCityInfo
type Upgrader struct {
	db          *City
	codeCountry map[string][]string
	codeRegion  map[string]string
	codeCity    map[string]map[string]string
}

func NewUpgrader(dbName string) (u *Upgrader, err error) {
	var b []byte
	u = new(Upgrader)
	// init db
	u.db, err = NewCity(dbName)
	if err != nil {
		return nil, err
	}
	// load country code file
	b, err = ioutil.ReadFile(FILENAME_COUNTRY)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(string(b)), &u.codeCountry)
	if err != nil {
		return nil, err
	}
	// load region code
	b, err = ioutil.ReadFile(FILENAME_REGION)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(string(b)), &u.codeRegion)
	if err != nil {
		return nil, err
	}
	// load city code
	b, err = ioutil.ReadFile(FILENAME_CITY)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(string(b)), &u.codeCity)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// FindCityInfo is a wrapper of original FindCityInfo but returns a custom UpgradeCityInfo
func (u *Upgrader) FindCityInfo(addr string) (*UpgradeCityInfo, error) {
	// default db.Languages is "CN"
	info, err := u.db.FindInfo(addr, u.db.Languages()[0])
	if err != nil {
		return nil, err
	}
	return u.NewUpgradeCityInfo(info)
}

// NewUpgradeCityInfo returns an UpgradeCityInfo by using an upgrader to upgrade given CityInfo
func (u *Upgrader) NewUpgradeCityInfo(info *CityInfo) (uInfo *UpgradeCityInfo, err error) {
	if info == nil {
		return nil, ERR_EMPTY_CITYINFO
	}
	uInfo = &UpgradeCityInfo{
		ContinentCode:  NULL,
		CountryName:    info.CountryName,
		IddCode:        NULL,
		RegionName:     info.RegionName,
		ChinaAdminCode: NULL,
		CityName:       info.CityName,
		ChinaCityCode:  NULL,
	}
	// upgrade ContinentCode & IddCode
	err = uInfo.upgradeCountry(u)
	if err != nil {
		return uInfo, err
	}
	// upgrade ChinaAdminCode & ChinaCityCode
	err = uInfo.upgradeMainland(u)
	if err != nil {
		return uInfo, err
	}
	return uInfo, nil
}

// custom cityInfo
type UpgradeCityInfo struct {
	ContinentCode  string `json:"continent_code"`   // Continent name
	CountryName    string `json:"country_name"`     // Country name
	IddCode        string `json:"idd_code"`         // Prefix of National Telephone Number
	RegionName     string `json:"region_name"`      // Mainland Region name
	ChinaAdminCode string `json:"china_admin_code"` // Provincial administrative unit code
	CityName       string `json:"city_name"`        // Mainland City name
	ChinaCityCode  string `json:"china_city_code"`  // Municipal administrative unit code
}

// String prettify print in console
func (u *UpgradeCityInfo) String() string {
	b, err := json.Marshal(*u)
	if err != nil {
		return fmt.Sprintf("%+v", *u)
	}
	var out bytes.Buffer
	err = json.Indent(&out, b, "", "    ")
	if err != nil {
		return fmt.Sprintf("%+v", *u)
	}
	return out.String()
}

// upgradeCountry uses CountryName to upgrade corresponding ContinentCode & IddCode
func (u *UpgradeCityInfo) upgradeCountry(upgrader *Upgrader) error {
	// map cotinentCode & iddCode
	if countryValue, ok := upgrader.codeCountry[u.CountryName]; ok {
		u.IddCode, u.ContinentCode = countryValue[0], countryValue[1]
		return nil
	}
	return ERR_NOT_FOUND_COUNTRY
}

// upgradeMainland uses RegionName & CityName to upgrade corresponding ChinaAdminCode & ChinaCityCode
func (u *UpgradeCityInfo) upgradeMainland(upgrader *Upgrader) error {
	if strings.Compare(u.CountryName, STRING_CHINA) != 0 {
		u.RegionName, u.CityName = NULL, NULL
		return nil
	}
	// map adminCode & cityCode(the same as stampCode)
	if code, ok := upgrader.codeRegion[u.RegionName]; ok { // Provincial administrative unit
		if strings.Contains(STRING_CHINA_DC, u.RegionName) { // HK/MO/TW
			u.ChinaAdminCode, u.CityName = code, NULL
			return nil
		}
		if strings.Compare(u.CityName, STRING_EMPTY) != 0 { // city name is not empty
			for cityCode, city := range upgrader.codeCity[code] {
				if strings.Contains(city, u.CityName) { // Municipal administrative unit
					u.ChinaAdminCode, u.ChinaCityCode = code, cityCode
					return nil
				}
			}
		}
		u.ChinaAdminCode, u.CityName = code, NULL
		return ERR_NOT_FOUND_MAINLAND_CITY
	}
	u.RegionName, u.CityName = NULL, NULL
	return ERR_NOT_FOUND_MAINLAND_REGION
}
