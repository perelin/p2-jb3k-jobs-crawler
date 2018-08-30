package responses

import (
	"encoding/json"
)

type MonsterJobAdListEntry struct {
	JobID                 int    `json:"JobID"`
	Title                 string `json:"Title"`
	TitleLink             string `json:"TitleLink"`
	IsBolded              bool   `json:"IsBolded"`
	DatePostedText        string `json:"DatePostedText"`
	DatePosted            string `json:"DatePosted"`
	LocationText          string `json:"LocationText"`
	LocationLink          string `json:"LocationLink"`
	JobViewURL            string `json:"JobViewUrl"`
	ImpressionTracking    string `json:"ImpressionTracking"`
	HasLocationAddress    bool   `json:"HasLocationAddress"`
	IsSavedJob            bool   `json:"IsSavedJob"`
	IsAppliedJob          bool   `json:"IsAppliedJob"`
	IsNewJob              bool   `json:"IsNewJob"`
	HasAdapt              bool   `json:"HasAdapt"`
	HasProDiversity       bool   `json:"HasProDiversity"`
	HasSpecialCommitments bool   `json:"HasSpecialCommitments"`
	Company               struct {
		Name              string `json:"Name"`
		CompanyLink       string `json:"CompanyLink"`
		HasCompanyAddress bool   `json:"HasCompanyAddress"`
	} `json:"Company"`
	Text                       string      `json:"Text"`
	LocationClickJsFunction    string      `json:"LocationClickJsFunction"`
	CompanyClickJsFunction     string      `json:"CompanyClickJsFunction"`
	JobTitleClickJsFunction    string      `json:"JobTitleClickJsFunction"`
	JobDescription             string      `json:"JobDescription"`
	ApplyMethod                int         `json:"ApplyMethod"`
	ApplyType                  string      `json:"ApplyType"`
	IsAggregated               string      `json:"IsAggregated"`
	CityText                   string      `json:"CityText"`
	StateText                  string      `json:"StateText"`
	JobDescriptionMeta         string      `json:"JobDescriptionMeta"`
	EmploymentTypeMeta         string      `json:"EmploymentTypeMeta"`
	IndustryTypeMeta           string      `json:"IndustryTypeMeta"`
	JobViewURLMeta             string      `json:"JobViewUrlMeta"`
	IsFastApply                bool        `json:"IsFastApply"`
	Target                     interface{} `json:"Target"`
	IsSecondaryJob             bool        `json:"IsSecondaryJob"`
	JobIDCloud                 int         `json:"JobIdCloud"`
	MusangKingID               string      `json:"MusangKingId"`
	IsSecondarySearchResultJob bool        `json:"IsSecondarySearchResultJob"`
	InlineAdIndex              int         `json:"InlineAdIndex"`
	ShowCompanyAsLink          bool        `json:"ShowCompanyAsLink"`
	ShowLocationAsLink         bool        `json:"ShowLocationAsLink"`
	HideCompanyLogo            bool        `json:"HideCompanyLogo"`
	ShowMultilocHover          bool        `json:"ShowMultilocHover"`
	MultilocHoverTitle         interface{} `json:"MultilocHoverTitle"`
	MultilocHover              interface{} `json:"MultilocHover"`
}

type JobPosting struct {
	Type        string   `json:"@type"`
	Context     string   `json:"@context"`
	Title       string   `json:"title"`
	DatePosted  string   `json:"datePosted"`
	Description string   `json:"description"`
	Industry    []string `json:"industry"`
	JobLocation struct {
		Type string `json:"@type"`
		Geo  struct {
			Type      string `json:"@type"`
			Latitude  string `json:"latitude"`
			Longitude string `json:"longitude"`
		} `json:"geo"`
		Address struct {
			Type            string `json:"@type"`
			AddressLocality string `json:"addressLocality"`
			AddressRegion   string `json:"addressRegion"`
			PostalCode      string `json:"postalCode"`
			AddressCountry  string `json:"addressCountry"`
		} `json:"address"`
	} `json:"jobLocation"`
	URL                string `json:"url"`
	HiringOrganization struct {
		Type string `json:"@type"`
		Name string `json:"name"`
		Logo string `json:"logo"`
	} `json:"hiringOrganization"`
	ValidThrough   string `json:"validThrough"`
	SalaryCurrency string `json:"salaryCurrency"`
	Identifier     struct {
		Type  string `json:"@type"`
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"identifier"`
}

type IndustryList struct {
	Industry []string
}

func IndustryToJSON(industries IndustryList) string {
	industryJSON, _ := json.Marshal(industries)
	//spew.Dump(string(industryJSON))
	return string(industryJSON)
}
