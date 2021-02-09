package mzxml

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"encoding/binary"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"regexp"
	"strconv"

	"golang.org/x/net/html/charset"
)

func (f *MzXML) Read(reader io.Reader) error {
	f.decoder = xml.NewDecoder(reader)
	f.decoder.CharsetReader = charset.NewReaderLabel
	err := f.decoder.Decode(&f.content)
	if err != nil {
		return nil
	}
	err = f.traverseScan()
	return err
}

// ReadScan reads a single scan
// n is the sequence number of the scan in the mzXML file,
// This is not the same as the scan number that is specified
// in the mzMXL file! To read a scan using the mzXML number,
// use ReadScan(f, ScanIndex(f, scanNum))
func (f *MzXML) ReadScan(scanIndex int64) ([]Peak, error) {
	if scanIndex < 0 || scanIndex >= f.content.Run.ScanCount {
		return nil, ErrInvalidScanIndex
	}
	scan := f.index2Scan[scanIndex]
	cnt := scan.PeaksCount
	p := make([]Peak, cnt)
	data, err := base64.StdEncoding.DecodeString(scan.Peaks.Base64Str)
	if err != nil {
		fmt.Println("error:", err)
		return nil, err
	}
	// Only "network" byteorder is allowed according to the schema
	// if (scan.Peaks.ByteOrder=="network"){
	// }
	// Only "m/z-int" pairorder is allowed according to the schema
	// if scan.Peaks.pairOrder == "m/z-int" {
	// }
	if scan.Peaks.CompressionType == "zlib" {
		b := bytes.NewReader(data)
		z, err := zlib.NewReader(b)
		if err != nil {
			return nil, err
		}
		defer z.Close()
		p, err := ioutil.ReadAll(z)
		if err != nil {
			return nil, err
		}
		data = p
	}
	if scan.Peaks.Precision == 64 {
		for i := int64(0); i < cnt; i++ {
			bits := binary.BigEndian.Uint64(data[i*16:])
			float := math.Float64frombits(bits)
			p[i].Mz = float64(float)
			bits = binary.BigEndian.Uint64(data[i*16+8:])
			float = math.Float64frombits(bits)
			p[i].Intens = float64(float)
		}
	} else {
		for i := int64(0); i < cnt; i++ {
			bits := binary.BigEndian.Uint32(data[i*8:])
			float := math.Float32frombits(bits)
			p[i].Mz = float64(float)
			bits = binary.BigEndian.Uint32(data[i*8+4:])
			float = math.Float32frombits(bits)
			p[i].Intens = float64(float)

		}
	}
	return p, nil
}

// NumSpecs returns the number of spectra
func (f *MzXML) NumSpecs() int64 {
	return int64(len(f.index2Scan))
}

// Regular expression to match "duration"
var reDuration = regexp.MustCompile(`^PT(.*)S`)

// RetentionTime returns the retention time of a spectrum
// If no retention time is present, return -1
func (f *MzXML) RetentionTime(scanIndex int64) (float64, error) {
	if scanIndex < 0 || scanIndex >= f.NumSpecs() {
		return 0.0, ErrInvalidScanIndex
	}
	rtStr := f.content.Run.Specs[scanIndex].RetentionTime

	// Retention time is specified in "duration" format: https://www.w3schools.com/xml/schema_dtypes_date.asp
	// For now, we only accept PT<float>S format
	m := reDuration.FindStringSubmatch(rtStr)
	if len(m) <= 1 {
		return 0.0, ErrInvalidFormat
	}
	rt, err := strconv.ParseFloat(m[1], 64)
	if err != nil {
		return 0.0, err
	}
	return rt, nil
}

// Centroid returns true is the spectrum contains centroid peaks
// In mzXML, centroided applied to the whole file, not individual spectra.
// For compatibility with mzML, we keep the same interface
func (f *MzXML) Centroid(scanIndex int64) (bool, error) {
	if scanIndex < 0 || scanIndex >= f.NumSpecs() {
		return false, ErrInvalidScanIndex
	}
	if len(f.content.Run.DataProcessing) > 0 {
		return (f.content.Run.DataProcessing[0].Centroided != 0), nil
	}
	return false, nil
}

// MSLevel returns the MS Level of a spectrum
// If no MS level is present, return 0
func (f *MzXML) MSLevel(scanIndex int64) (int, error) {
	if scanIndex < 0 || scanIndex >= f.NumSpecs() {
		return 0, ErrInvalidScanIndex
	}
	return f.content.Run.Specs[scanIndex].MsLevel, nil
}

// traverseScan traverses all (recursive)scans and fills the
// arrays f.index2id and f.id2Index to make scans accessible
func (f *MzXML) traverseScan() error {
	f.index2id = make([]int64, f.content.Run.ScanCount, f.content.Run.ScanCount)
	f.id2Index = make(map[int64]int64, f.content.Run.ScanCount)
	f.index2Scan = make([]*scan, f.content.Run.ScanCount, f.content.Run.ScanCount)
	scanIndex := int64(0)
	err := error(nil)

	for i := range f.content.Run.Specs {
		scanIndex, err = f.addSpecToIndex(scanIndex, &f.content.Run.Specs[i])
		if err != nil {
			return err
		}
	}
	return err
}

func (f *MzXML) addSpecToIndex(scanIndex int64,
	scan *scan) (int64, error) { //x
	err := error(nil)
	f.index2id[scanIndex] = scan.ScanNum
	f.id2Index[scan.ScanNum] = scanIndex
	f.index2Scan[scanIndex] = scan
	scanIndex++
	if scan.FragScans != nil {
		for i := range scan.FragScans {
			scanIndex, err = f.addSpecToIndex(scanIndex, &scan.FragScans[i])
		}
	}
	return scanIndex, err
}

// ScanIndex converts a scan identifier (the number used in the mzXML file)
// into an index that is used to access the scans
func (f *MzXML) ScanIndex(scanID int64) (int64, error) {
	if index, ok := f.id2Index[scanID]; ok {
		return index, nil
	}
	return 0, ErrInvalidScanID
}

// ScanID converts a scan index (used to access the scan data) into a scan id
// (used in the mzxml file)
func (f *MzXML) ScanID(scanIndex int64) (int64, error) {
	if scanIndex > 0 && scanIndex < f.content.Run.ScanCount {
		return f.index2id[scanIndex], nil
	}
	return 0, ErrInvalidScanIndex
}
