// Copyright 2021 Rob Marissen
// SPDX-License-Identifier: MIT

package pepxml

import (
	"encoding/xml"
)

// Types for parsing pepXML

// Content encapsulated the pepXML file content
type Content struct {
	XMLName           xml.Name         `xml:"msms_pipeline_analysis"`
	Date              string           `xml:"date,attr,omitempty"`
	XMLns             string           `xml:"xmlns,attr,omitempty"`
	XMLnsxsi          string           `xml:"xmlns:xsi,attr,omitempty"`
	XSIschemaLocation string           `xml:"xsi:schemaLocation,attr,omitempty"`
	SummaryXML        string           `xml:"summary_xml,attr,omitempty"`
	MsmsRunSummary    []msmsRunSummary `xml:"msms_run_summary"`
}

type msmsRunSummary struct {
	ID             string          `xml:"id,attr,omitempty"`
	BaseName       string          `xml:"base_name,attr,omitempty"`
	MsManufacturer string          `xml:"msManufacturer,attr,omitempty"`
	MsModel        string          `xml:"msModel,attr,omitempty"`
	RawDataType    string          `xml:"raw_data_type,attr,omitempty"`
	RawData        string          `xml:"raw_data,attr,omitempty"`
	SampleEnzyme   sampleEnzyme    `xml:"sample_enzyme"`
	SearchSummary  searchSummary   `xml:"search_summary"`
	SpectrumQuery  []spectrumQuery `xml:"spectrum_query"`
}

type sampleEnzyme struct {
	Name        string      `xml:"name,attr,omitempty"`
	Specificity specificity `xml:"specificity"`
}

type specificity struct {
	Cut   string `xml:"cut,attr,omitempty"`
	NoCut string `xml:"no_cut,attr,omitempty"`
	Sense string `xml:"sense,attr,omitempty"`
}

type searchSummary struct {
	BaseName                  string                    `xml:"base_name,attr,omitempty"`
	SearchEngine              string                    `xml:"search_engine,attr,omitempty"`
	SearchEngineVersion       string                    `xml:"search_engine_version,attr,omitempty"`
	PrecursorMassType         string                    `xml:"precursor_mass_type,attr,omitempty"`
	FragmentMassType          string                    `xml:"fragment_mass_type,attr,omitempty"`
	SearchID                  string                    `xml:"search_id,attr,omitempty"`
	SearchDatabase            searchDatabase            `xml:"search_database"`
	EnzymaticSearchConstraint enzymaticSearchConstraint `xml:"enzymatic_search_constraint"`
	AminoacidModification     []aminoacidModification   `xml:"aminoacid_modification,omitempty"`
	Parameter                 []strParameter            `xml:"parameter,omitempty"`
}

type searchDatabase struct {
	LocalPath string `xml:"local_path,attr,omitempty"`
	Type      string `xml:"type,attr,omitempty"`
}

type enzymaticSearchConstraint struct {
	Enzyme                  string `xml:"enzyme,attr,omitempty"`
	MaxNumInternalCleavages int    `xml:"max_num_internal_cleavages,attr,omitempty"`
	MinNumberTermini        int    `xml:"min_number_termini,attr,omitempty"`
}
type aminoacidModification struct {
	Aminoacid string `xml:"aminoacid,attr,omitempty"`
	Massdiff  string `xml:"massdiff,attr,omitempty"`
	Mass      string `xml:"mass,attr,omitempty"`
	Variable  string `xml:"variable,attr,omitempty"`
	Symbol    string `xml:"symbol,attr,omitempty"`
}

type strParameter struct {
	Name  string `xml:"name,attr,omitempty"`
	Value string `xml:"value,attr,omitempty"`
}

type spectrumQuery struct {
	Spectrum             string       `xml:"spectrum,attr,omitempty"`
	SpectrumNativeID     string       `xml:"spectrumNativeID,attr,omitempty"`
	StartScan            int          `xml:"start_scan,attr"`
	EndScan              int          `xml:"end_scan,attr"`
	PrecursorNeutralMass float64      `xml:"precursor_neutral_mass,attr"`
	AssumedCharge        int          `xml:"assumed_charge,attr"`
	Index                int          `xml:"index,attr"`
	RetentionTime        float64      `xml:"retention_time_sec,attr"`
	SearchResult         searchResult `xml:"search_result"`
}

type searchResult struct {
	SearchHit []searchHit `xml:"search_hit"`
}

type searchHit struct {
	HitRank            int                  `xml:"hit_rank,attr"`
	Peptide            string               `xml:"peptide,attr"`
	PeptidePrevAA      string               `xml:"peptide_prev_aa,attr"`
	PeptideNextAA      string               `xml:"peptide_next_aa,attr"`
	Protein            string               `xml:"protein,attr"`
	NumTotProteins     int                  `xml:"num_tot_proteins,attr"`
	NumMatchedIons     int                  `xml:"num_matched_ions,attr"`
	TotNumIons         int                  `xml:"tot_num_ions,attr"`
	CalcNeutralPepMass float64              `xml:"calc_neutral_pep_mass,attr"`
	Massdiff           float64              `xml:"massdiff,attr"`
	NumTolTerm         int                  `xml:"num_tol_term,attr"`
	NumMissedCleavages int                  `xml:"num_missed_cleavages,attr"`
	NumMatchedPeptides int                  `xml:"num_matched_peptides,attr"`
	AlternativeProtein []alternativeProtein `xml:"alternative_protein"`
	ModificationInfo   []modificationInfo   `xml:"modification_info"`
	SearchScore        []searchScore        `xml:"search_score"`
	AnalysisResult     []analysisResult     `xml:"analysis_result"`
}

type alternativeProtein struct {
	Protein string `xml:"protein,attr"`
}

type modificationInfo struct {
	ModifiedPeptide  string             `xml:"modified_peptide,attr"`
	ModAminoacidMass []modAminoacidMass `xml:"mod_aminoacid_mass"`
}

type modAminoacidMass struct {
	Position int     `xml:"position,attr"`
	Mass     float64 `xml:"mass,attr"`
	Static   float64 `xml:"static,attr"`
}

type searchScore struct {
	Name  string  `xml:"name,attr"`
	Value float64 `xml:"value,attr"`
}

type analysisResult struct {
	PeptideprophetResult peptideprophetResult `xml:"peptideprophet_result"`
}

type peptideprophetResult struct {
	SearchScoreSummary searchScoreSummary `xml:"search_score_summary"`
	Probability        float64            `xml:"probability,attr"`
}

type searchScoreSummary struct {
	Parameter []parameter `xml:"parameter"`
}

type parameter struct {
	Name  string  `xml:"name,attr"`
	Value float64 `xml:"value,attr"`
}
