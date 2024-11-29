package models

type Occurrence struct {
	SeqId                 int     `json:"seq_id,omitempty"`
	LocusId               int64   `json:"locus_id,omitempty"`
	Quality               int     `json:"quality,omitempty"`
	Filter                string  `json:"filter,omitempty"`
	Zygosity              string  `json:"zygosity,omitempty"`
	Pf                    float64 `json:"pf,omitempty"`
	Af                    float64 `json:"af,omitempty"`
	GnomadV3Af            float64 `json:"gnomad_v3_af,omitempty"`
	Hgvsg                 string  `json:"hgvsg,omitempty"`
	OmimInheritanceCode   string  `json:"omim_inheritance_code,omitempty"`
	AdRatio               float64 `json:"ad_ratio,omitempty"`
	VariantClass          string  `json:"variant_class,omitempty"`
	VepImpact             string  `json:"vep_impact,omitempty"`
	Symbol                string  `json:"symbol,omitempty"`
	ClinvarInterpretation string  `json:"clinvar_interpretation,omitempty"`
	ManeSelect            bool    `json:"mane_select,omitempty"`
	Canonical             bool    `json:"canonical,omitempty"`
}

var OccurrenceTable = Table{
	Name:  "occurrences",
	Alias: "o",
}
var VariantTable = Table{
	Name:  "variants",
	Alias: "v",
}

var FilterField = Field{
	Name:          "filter",
	CanBeSelected: true,
	CanBeFiltered: true,
	Table:         OccurrenceTable,
}
var SeqIdField = Field{
	Name:          "seq_id",
	CanBeSelected: true,
	CanBeFiltered: true,
	Table:         OccurrenceTable,
}
var LocusIdField = Field{
	Name:          "locus_id",
	CanBeSelected: true,
	CanBeFiltered: true,
	Table:         OccurrenceTable,
}
var ZygosityField = Field{
	Name:          "zygosity",
	CanBeSelected: true,
	CanBeFiltered: true,
	Table:         OccurrenceTable,
}
var AdRatioField = Field{
	Name:          "ad_ratio",
	CanBeSelected: true,
	CanBeFiltered: true,
	Table:         OccurrenceTable,
}
var PfField = Field{
	Name:          "pf",
	CanBeSelected: true,
	CanBeFiltered: true,
	Table:         VariantTable,
}
var AfField = Field{
	Name:          "af",
	CanBeSelected: true,
	CanBeFiltered: true,
	Table:         VariantTable,
}
var VariantClassField = Field{
	Name:          "variant_class",
	CanBeSelected: true,
	CanBeFiltered: true,
	Table:         VariantTable,
}
var HgvsgField = Field{
	Name:          "hgvsg",
	CanBeSelected: true,
	CanBeFiltered: true,
	Table:         VariantTable,
}
var OccurrencesFields = []Field{
	SeqIdField,
	LocusIdField,
	FilterField,
	ZygosityField,
	AdRatioField,
	PfField,
	AfField,
	VariantClassField,
	HgvsgField,
}
