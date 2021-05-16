package fasta

import (
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type fastaDBMapping struct {
	commonNames []string
	taxonomyID  int64
	URL         string
}

var fastaDBMap = []fastaDBMapping{
	{
		commonNames: []string{`swiss`, `sprot`, `uniprot swissprot`, `sp`},
		taxonomyID:  0,
		URL:         `https://ftp.expasy.org/databases/uniprot/current_release/knowledgebase/complete/uniprot_sprot.fasta.gz`,
	},
	{
		commonNames: []string{`human`, `homo sapiens`},
		taxonomyID:  9606,
		URL:         `https://ftp.expasy.org/databases/uniprot/current_release/knowledgebase/reference_proteomes/Eukaryota/UP000005640/UP000005640_9606.fasta.gz`,
	},
	{
		commonNames: []string{`ecoli`, `e.coli`, `e-coli`, `escherichia coli`, `K12`},
		taxonomyID:  83333,
		URL:         `https://ftp.expasy.org/databases/uniprot/current_release/knowledgebase/reference_proteomes/Bacteria/UP000000625/UP000000625_83333.fasta.gz`,
	},
	{
		commonNames: []string{`yeast`, `s. cerevisiae`, `saccharomyces cerevisiae`},
		taxonomyID:  559292,
		URL:         `https://ftp.expasy.org/databases/uniprot/current_release/knowledgebase/reference_proteomes/Eukaryota/UP000002311/UP000002311_559292.fasta.gz`,
	},
	{
		commonNames: []string{`mouse`, `mus musculus`, `house mouse`},
		taxonomyID:  10090,
		URL:         `https://ftp.expasy.org/databases/uniprot/current_release/knowledgebase/reference_proteomes/Eukaryota/UP000000589/UP000000589_10090.fasta.gz`,
	},
	{
		commonNames: []string{`rat`, `rattus norvegicus`, `norway rat`, `brown rat`},
		taxonomyID:  10116,
		URL:         `https://ftp.expasy.org/databases/uniprot/current_release/knowledgebase/reference_proteomes/Eukaryota/UP000002494/UP000002494_10116.fasta.gz`,
	},
	{
		commonNames: []string{`zebra fish`, `danio rerio`, `zebrafish`},
		taxonomyID:  7955,
		URL:         `https://ftp.expasy.org/databases/uniprot/current_release/knowledgebase/reference_proteomes/Eukaryota/UP000000437/UP000000437_7955.fasta.gz`,
	},
	{
		commonNames: []string{`cow`, `bos taurus`, `dairy cow`, `domestic cow`},
		taxonomyID:  9913,
		URL:         `https://ftp.expasy.org/databases/uniprot/current_release/knowledgebase/reference_proteomes/Eukaryota/UP000009136/UP000009136_9913.fasta.gz`,
	},
	{
		commonNames: []string{`fruit fly`, `drosophila melanogaster`},
		taxonomyID:  7227,
		URL:         `https://ftp.expasy.org/databases/uniprot/current_release/knowledgebase/reference_proteomes/Eukaryota/UP000000803/UP000000803_7227.fasta.gz`,
	},
	{
		commonNames: []string{`Nematode worm`, `Caenorhabditis elegans`},
		taxonomyID:  6239,
		URL:         `https://ftp.expasy.org/databases/uniprot/current_release/knowledgebase/reference_proteomes/Eukaryota/UP000001940/UP000001940_6239.fasta.gz`,
	},
	{
		commonNames: []string{`SARS-CoV-2`, `corona`, `corona virus`, `2019-nCoV`, `covid 19`},
		taxonomyID:  2697049,
		URL:         `https://ftp.expasy.org/databases/uniprot/current_release/knowledgebase/reference_proteomes/Viruses/UP000464024/UP000464024_2697049.fasta.gz`,
	},
	{
		commonNames: []string{`Lambda phage`, `coliphage λ`, `Escherichia virus Lambda`},
		taxonomyID:  10710,
		URL:         `https://ftp.expasy.org/databases/uniprot/current_release/knowledgebase/reference_proteomes/Viruses/UP000001711/UP000001711_10710.fasta.gz`,
	},
	// {
	// 	commonNames: []string{`Phi X 174`, `ΦX174`, `Escherichia virus phiX174`, `Bacteriophage phi-X174`},
	// 	taxonomyID:  10847,
	// 	URL:         ``, // Not listed!
	// },
	{
		commonNames: []string{`SV40`, `simian vacuolating virus 40`, `simian virus 40`, `Macaca mulatta polyomavirus 1`},
		taxonomyID:  1891767,
		URL:         `https://ftp.expasy.org/databases/uniprot/current_release/knowledgebase/reference_proteomes/Viruses/UP000007705/UP000007705_1891767.fasta.gz`,
	},
	{
		commonNames: []string{`Herpes simplex virus 1`, `Herpes simplex 1`, `HSV-1`, `Human herpesvirus 1`},
		taxonomyID:  10298,
		URL:         `https://ftp.expasy.org/databases/uniprot/current_release/knowledgebase/reference_proteomes/Viruses/UP000110586/UP000110586_10298.fasta.gz`,
	},
	{
		commonNames: []string{`Herpes simplex virus 2`, `Herpes simplex 2`, `HSV-2`, `Human herpesvirus 2`, `Human alphaherpesvirus 2`},
		taxonomyID:  10310,
		URL:         `https://ftp.expasy.org/databases/uniprot/current_release/knowledgebase/reference_proteomes/Viruses/UP000270953/UP000270953_10310.fasta.gz`,
	},
	{
		commonNames: []string{`Escherichia virus T4`, `T4 phage`},
		taxonomyID:  10665,
		URL:         `https://ftp.expasy.org/databases/uniprot/current_release/knowledgebase/reference_proteomes/Viruses/UP000009087/UP000009087_10665.fasta.gz`,
	},
	{
		commonNames: []string{`CRAP`},
		taxonomyID:  -1,
		URL:         `http://ftp.thegpm.org/fasta/cRAP/crap.fasta`,
	},
	// {
	// 	commonNames: []string{`Tobacco mosaic virus`, ``, ``},
	// 	taxonomyID:  12242,
	// 	URL:         ``, // Not listed!
	// },
	// {
	// 	commonNames: []string{``, ``, ``},
	// 	taxonomyID:  ,
	// 	URL:         ``,
	// },
}

// UniprotURL converts a string with a taxonomy or common (species) name into a URL where the FASTA file can be downloaded
func UniprotURL(name string) (string, error) {
	tax, err := strconv.ParseInt(name, 10, 64)
	if err == nil {
		for _, DBMap := range fastaDBMap {
			if tax == DBMap.taxonomyID {
				return DBMap.URL, nil
			}
		}
	} else {
		// If string in non-numeric, check if its a common name
		tax = -1
		for _, DBMap := range fastaDBMap {
			for _, n := range DBMap.commonNames {
				if strings.EqualFold(name, n) {
					return DBMap.URL, nil
				}
			}
		}
	}
	return ``, fmt.Errorf("FASTA name unknown: %s", name)
}

// Convert a download URL to the filename that will be used to store the file
func urlToFilename(url string) string {
	fn := url[strings.LastIndex(url, "/")+1:]
	// Files are uncompressed before storing, remove .gz extention
	fn = strings.TrimSuffix(fn, `.gz`)
	return fn
}

func urlToPath(url string, dir string) string {
	pn := filepath.Join(dir, urlToFilename(url))
	return pn
}

// FuzzyNameFile checks a FASTA file exists
// name can be a real filename (full path) or a fuzzy name (Latin species name, taxonomy id or common name)
// dir is the directory to use in case on a fuzzy name
func FuzzyNameFile(name string, dir string) (string, error) {
	// Check if name refers to an existing file
	if _, err := os.Stat(name); !os.IsNotExist(err) {
		return name, nil
	}
	url, err := UniprotURL(name)
	if err != nil {
		return ``, err
	}
	return urlToPath(url, dir), nil
}

// Download gets a file from a given URL, and puts it in the supplied directory
// If the URL ends in ".gz", the file is assumed to be gzip compressed,
// and is uncompressed before writing to disk
func download(url string, dir string) error {
	pn := urlToPath(url, dir)

	// Create a temporary file for download
	fn := urlToFilename(url)
	tmpFile, err := ioutil.TempFile(dir, fn)
	if err != nil {
		return err
	}
	defer tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	r := resp.Body
	// Uncompress if gz file format
	if strings.HasSuffix(url, `.gz`) {
		r, err = gzip.NewReader(resp.Body)
		if err != nil {
			return err
		}
		defer r.Close()
	}
	_, err = io.Copy(tmpFile, r)
	if err != nil {
		return err
	}
	tmpFile.Close()
	os.Rename(tmpFile.Name(), pn)
	return nil
}

// Check if file should be downloaded from the given URL
// If downloading is deemed useless, a non-empty string
// that describes the reason is returned
// Reasons for not downloading are:
// - the file already exists locally
// - there is insufficient disk space to store the file
func checkDownload(url string, dir string) (string, error) {
	pn := urlToPath(url, dir)

	stat, err := os.Stat(pn)
	// If the file doesn't yet exist, no need to check URL timestamp
	if !os.IsNotExist(err) {
		resp, err := http.Head(url)
		if err != nil {
			return ``, err
		}
		if resp.StatusCode >= 400 {
			return ``, errors.New(url + `: ` + resp.Status)
		}
		mod := resp.Header.Get(`Last-Modified`)
		t, err := http.ParseTime(mod)
		// If error, time string can't be parsed and we assume the file must be downloaded
		if err == nil {
			if !t.After(stat.ModTime()) {
				return `local file is newer or equal`, nil
			}
		}
	}
	// TODO: Check if disk space is sufficient
	return ``, nil
}

// UpdateFASTA downloads a FASTA file from internet if the remove version is newer than the local version.
// name is the name of the species, taxonomy ID or one of the build in common names
// dir is the download directory
func UpdateFASTA(name string, dir string) error {
	url, err := UniprotURL(name)
	if err != nil {
		return err
	}
	msg, err := checkDownload(url, dir)
	if err != nil {
		return err
	}
	if msg != `` {
		fmt.Printf("Not downloading %s: %s\n", url, msg)
	} else {
		err = download(url, dir)
		if err != nil {
			return err
		}
	}
	return nil
}
