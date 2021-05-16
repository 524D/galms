// Copyright 2021 Rob Marissen
// SPDX-License-Identifier: MIT

package mzxml

import (
	"encoding/xml"
	"errors"
	"os"
)

type MzXML struct {
	mzXMLfile  *os.File
	decoder    *xml.Decoder
	content    mzXMLContent
	index2id   []int64
	id2Index   map[int64]int64
	index2Scan []*scan
}

// Peak contains the actual ms peak info
type Peak struct {
	Mz     float64
	Intens float64
}

type mzXMLContent struct {
	XMLName xml.Name `xml:"mzXML"`
	Run     msRun    `xml:"msRun"`
}

type msRun struct {
	ScanCount      int64            `xml:"scanCount,attr"`
	StartTime      string           `xml:"startTime,attr,omitEmpty"`
	EndTime        string           `xml:"endTime,attr,omitEmpty"`
	ParentFile     []parentFile     `xml:"parentFile,omitEmpty"`
	MsInstrument   []msInstrument   `xml:"msInstrument,omitEmpty"`
	DataProcessing []dataProcessing `xml:"dataProcessing,omitEmpty"`
	Specs          []scan           `xml:"scan,omitEmpty"`
}

type msInstrument struct {
	MsInstrumentID int                `xml:"msInstrumentID,attr"`
	MsManufacturer msManufacturer     `xml:"msManufacturer,omitEmpty"`
	MsModel        msModel            `xml:"msModel,omitEmpty"`
	MsIonisation   msIonisation       `xml:"msIonisation,omitEmpty"`
	MsMassAnalyzer msMassAnalyzer     `xml:"msMassAnalyzer,omitEmpty"`
	MsDetector     msDetector         `xml:"msDetector,omitEmpty"`
	Software       instrumentSoftware `xml:"software,omitEmpty"`
}

type msManufacturer struct {
	Category string `xml:"category,attr,omitEmpty"`
	Value    string `xml:"value,attr,omitEmpty"`
}

type msModel struct {
	Category string `xml:"category,attr,omitEmpty"`
	Value    string `xml:"value,attr,omitEmpty"`
}

type msIonisation struct {
	Category string `xml:"category,attr,omitEmpty"`
	Value    string `xml:"value,attr,omitEmpty"`
}

type msMassAnalyzer struct {
	Category string `xml:"category,attr,omitEmpty"`
	Value    string `xml:"value,attr,omitEmpty"`
}

type msDetector struct {
	Category string `xml:"category,attr,omitEmpty"`
	Value    string `xml:"value,attr,omitEmpty"`
}

type instrumentSoftware struct {
	Type    string `xml:"type,attr,omitEmpty"`
	Name    string `xml:"name,attr,omitEmpty"`
	Version string `xml:"version,attr,omitEmpty"`
}

type dataProcessing struct {
	// FIXME: 'omitempty' does not work?
	Centroided          int                 `xml:"centroided,omitEmpty,attr"`
	Software            software            `xml:"software,omitEmpty"`
	Comment             string              `xml:"comment,omitEmpty"`
	ProcessingOperation processingOperation `xml:"processingOperation,omitEmpty"`
	// ProcessingOperation     string `xml:"processingOperation,omitEmpty"`
	// ProcessingOperationName string `xml:"processingOperation>name,attr,omitEmpty"`
}

type software struct {
	Type    string `xml:"type,attr,omitEmpty"`
	Name    string `xml:"name,attr,omitEmpty"`
	Version string `xml:"version,attr,omitEmpty"`
}

type processingOperation struct {
	Name string `xml:"name,attr,omitEmpty"`
}

type parentFile struct {
	FileName string `xml:"fileName,attr,omitEmpty"`
	FileType string `xml:"fileType,attr,omitEmpty"`
	FileSha1 string `xml:"fileSha1,attr,omitEmpty"`
}

type scan struct {
	ScanNum           int64   `xml:"num,attr"`
	RetentionTime     string  `xml:"retentionTime,attr,omitEmpty"`
	Polarity          string  `xml:"polarity,attr,omitEmpty"`
	MsLevel           int     `xml:"msLevel,attr"`
	PeaksCount        int64   `xml:"peaksCount,attr"`
	LowMz             float64 `xml:"lowMz,attr,omitEmpty"`
	HighMz            float64 `xml:"highMz,attr,omitEmpty"`
	BasePeakMz        float64 `xml:"basePeakMz,attr,omitEmpty"`
	BasePeakIntensity float64 `xml:"basePeakIntensity,attr,omitEmpty"`
	TotIonCurrent     float64 `xml:"totIonCurrent,attr,omitEmpty"`
	Peaks             peaks   `xml:"peaks,omitEmpty"`
	FragScans         []scan  `xml:"scan,omitEmpty"`
}

// <scan num="9" retentionTime="PT5.16998400S" polarity="+" msLevel="2" peaksCount="104" lowMz="121.17511749" highMz="669.60003662" basePeakMz="355.42813110" basePeakIntensity="16858.51367188" totIonCurrent="199931.18437386">
// <precursorMz precursorIntensity="25302.23828125">353.31530762</precursorMz>
// <peaks precision="32" byteOrder="network" pairOrder="m/z-int">

type precursorMz struct {
	PrecursorIntensity float64 `xml:"precursorIntensity,attr,omitEmpty"`
	MzStr              string  `xml:",chardata"`
}

type peaks struct {
	Precision       int64  `xml:"precision,attr,omitEmpty"`
	ByteOrder       string `xml:"byteOrder,attr,omitEmpty"`
	PairOrder       string `xml:"pairOrder,attr,omitEmpty"`
	CompressionType string `xml:"compressionType,attr,omitEmpty"`
	Base64Str       string `xml:",chardata"`
}

var (
	ErrInvalidScanID    = errors.New("mzxml: invalid scan id")
	ErrInvalidScanIndex = errors.New("mzxml: invalid scan index")
	ErrInvalidFormat    = errors.New("mzxml: invalid data format")
)
