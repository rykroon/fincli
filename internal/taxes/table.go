package taxes

type FilingStatus string

const (
	Single          FilingStatus = "single"
	MarriedJoint    FilingStatus = "married_joint"
	MarriedSeparate FilingStatus = "married_separate"
	HeadOfHouse     FilingStatus = "head_of_household"
)

type TaxTable map[int]map[FilingStatus]FilingConfig

func (t TaxTable) GetConfig(year int, status FilingStatus) (*FilingConfig, bool) {
	taxYear, ok := t[year]
	if !ok {
		return nil, false
	}
	config, ok := taxYear[status]
	if !ok {
		return nil, false
	}
	return &config, true
}
