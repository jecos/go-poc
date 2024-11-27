package models

type Occurrence struct {
	SeqId                 int     `json:"seq_id,omitempty"`
	LocusId               string  `json:"locus_id,omitempty"`
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
