/*
Copyright © 2021 Rob Marissen

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/524D/galms/fasta"

	"github.com/spf13/cobra"
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
		URL:         `https://ftp.expasy.org/databases/uniprot/current_release/knowledgebase/reference_proteomes/Eukaryota/UP000005640_9606.fasta.gz`,
	},
	{
		commonNames: []string{`ecoli`, `e.coli`, `e-coli`, `escherichia coli`, `K12`},
		taxonomyID:  83333,
		URL:         `https://ftp.expasy.org/databases/uniprot/current_release/knowledgebase/reference_proteomes/Bacteria/UP000000625_83333.fasta.gz`,
	},
	{
		commonNames: []string{`yeast`, `s. cerevisiae`, `saccharomyces cerevisiae`},
		taxonomyID:  83333,
		URL:         `https://ftp.expasy.org/databases/uniprot/current_release/knowledgebase/reference_proteomes/Eukaryota/UP000002311_559292.fasta.gz`,
	},
	{
		commonNames: []string{`mouse`, `mus musculus`, `house mouse`},
		taxonomyID:  10090,
		URL:         `https://ftp.expasy.org/databases/uniprot/current_release/knowledgebase/reference_proteomes/Eukaryota/UP000000589_10090.fasta.gz`,
	},
	{
		commonNames: []string{`rat`, `rattus norvegicus`, `norway rat`, `brown rat`},
		taxonomyID:  10116,
		URL:         `https://ftp.expasy.org/databases/uniprot/current_release/knowledgebase/reference_proteomes/Eukaryota/UP000002494_10116.fasta.gz`,
	},
	{
		commonNames: []string{`zebra fish`, `danio rerio`, `zebrafish`},
		taxonomyID:  7955,
		URL:         `https://ftp.expasy.org/databases/uniprot/current_release/knowledgebase/reference_proteomes/Eukaryota/UP000000437_7955.fasta.gz`,
	},
	{
		commonNames: []string{`cow`, `bos taurus`, `dairy cow`, `domestic cow`},
		taxonomyID:  9913,
		URL:         `https://ftp.expasy.org/databases/uniprot/current_release/knowledgebase/reference_proteomes/Eukaryota/UP000009136_9913.fasta.gz`,
	},
	{
		commonNames: []string{`fruit fly`, `drosophila melanogaster`},
		taxonomyID:  7227,
		URL:         `https://ftp.expasy.org/databases/uniprot/current_release/knowledgebase/reference_proteomes/Eukaryota/UP000000803_7227.fasta.gz`,
	},
	{
		commonNames: []string{`Nematode worm`, `Caenorhabditis elegans`},
		taxonomyID:  6239,
		URL:         `https://ftp.expasy.org/databases/uniprot/current_release/knowledgebase/reference_proteomes/Eukaryota/UP000001940_6239.fasta.gz`,
	},
	{
		commonNames: []string{`SARS-CoV-2`, `corona`, `corona virus`, `2019-nCoV`, `covid 19`},
		taxonomyID:  2697049,
		URL:         `https://ftp.expasy.org/databases/uniprot/current_release/knowledgebase/reference_proteomes/Viruses/UP000464024_2697049.fasta.gz`,
	},
	{
		commonNames: []string{`Lambda phage`, `coliphage λ`, `Escherichia virus Lambda`},
		taxonomyID:  10710,
		URL:         `https://ftp.expasy.org/databases/uniprot/current_release/knowledgebase/reference_proteomes/Viruses/UP000001711_10710.fasta.gz`,
	},
	{
		commonNames: []string{`Phi X 174`, `ΦX174`, `Escherichia virus phiX174`, `Bacteriophage phi-X174`},
		taxonomyID:  10847,
		URL:         ``, // Not listed!
	},
	{
		commonNames: []string{`SV40`, `simian vacuolating virus 40`, `simian virus 40`, `Macaca mulatta polyomavirus 1`},
		taxonomyID:  1891767,
		URL:         `https://ftp.expasy.org/databases/uniprot/current_release/knowledgebase/reference_proteomes/Viruses/UP000007705_1891767.fasta.gz`,
	},
	{
		commonNames: []string{`Herpes simplex virus 1`, `Herpes simplex 1`, `HSV-1`, `Human herpesvirus 1`},
		taxonomyID:  10298,
		URL:         `https://ftp.expasy.org/databases/uniprot/current_release/knowledgebase/reference_proteomes/Viruses/UP000110586_10298.fasta.gz`,
	},
	{
		commonNames: []string{`Herpes simplex virus 2`, `Herpes simplex 2`, `HSV-2`, `Human herpesvirus 2`, `Human alphaherpesvirus 2`},
		taxonomyID:  10310,
		URL:         `https://ftp.expasy.org/databases/uniprot/current_release/knowledgebase/reference_proteomes/Viruses/UP000270953_10310.fasta.gz`,
	},
	{
		commonNames: []string{`Escherichia virus T4`, `T4 phage`},
		taxonomyID:  10665,
		URL:         `https://ftp.expasy.org/databases/uniprot/current_release/knowledgebase/reference_proteomes/Viruses/UP000009087_10665.fasta.gz`,
	},
	{
		commonNames: []string{`Tobacco mosaic virus`, ``, ``},
		taxonomyID:  12242,
		URL:         ``, // Not listed!
	},
	// {
	// 	commonNames: []string{``, ``, ``},
	// 	taxonomyID:  ,
	// 	URL:         ``,
	// },
}

// fastaCmd represents the fasta command
var fastaCmd = &cobra.Command{
	Use:   "fasta",
	Short: "process a FASTA file according to options",
	Long: `The 'fasta' subcommand read fasta files

	The FASTA file is specified by the last argument(s). This can be either a
	plain filename, a symbolic identifier e.g. human, yeast, swiss,
	or a numeric taxonomy ID e.g. 9606.

	When a symbolic identifier is specified and the file is not present in the  
	FASTA directory, the file is retrieved from uniprot.org`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Fatal("Last argument must be name of FASTA file")
		}

		usr, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}
		dataDir := filepath.Join(usr.HomeDir, `data`, `fasta`)
		// Make dir in case it does not exist
		os.MkdirAll(dataDir, os.ModePerm)

		for _, n := range args {
			upd, err := cmd.Flags().GetBool("update")
			if err != nil {
				log.Fatalf("Getstrings 'update' flag failed: %v", err)
			}
			if upd {
				err = updateFASTA(n, dataDir)
			}

			fn, err := fuzzyNameFile(n, dataDir)
			if err != nil {
				log.Fatalf("%s is not a valid filename nor a fuzzy name", n)
			}
			file, err := os.Open(fn)
			if err != nil {
				log.Fatalf("Can't open file %s", args[0])
			}
			defer file.Close()

			f, err := fasta.Read(file)
			if err != nil {
				log.Fatal("fasta Read failed")
			}
			missing, err := cmd.Flags().GetString("missing")
			if err != nil {
				log.Fatalf("Getstrings 'missing' flag failed: %v", err)
			}
			contains, err := cmd.Flags().GetString("contains")
			if err != nil {
				log.Fatalf("Getstrings 'contains' flag failed: %v", err)
			}

			for _, p := range f.Prots() {
				if contains == `` || strings.Contains(p.Sequence(), contains) {

					if missing == `` || !strings.Contains(p.Sequence(), missing) {
						fmt.Printf("Length: %d %s %s\n", len(p.Sequence()), p.ID(), p.Description())
					}
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(fastaCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// fastaCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// fastaCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	fastaCmd.PersistentFlags().StringP("missing", "m", "", "List proteins which don't contain the specified sequence")
	fastaCmd.PersistentFlags().StringP("contains", "c", "", "List proteins which contain the specified sequence")
	fastaCmd.PersistentFlags().BoolP("update", "u", false, "Update FASTA file")

	// Flag 'analyse':
	// - print length/mass distribution distribution
	// - print number of proteins in certain class,
	//   e.g.:
	//     with motif for export from cell, import golgi, zinc finger, multiple motifs in one prot
	//       zinc-binding motif HEXXHXXGXXH
	//       N-terminal secretory signal peptide and a prodomain with a conserved PRCGXPD motif
	//     isoforms
	//     phosphorilation/glycosylation patterns
	// - print statistically significant sequence features, e.g. missing subsequence/regex, abundant subsequence/regex
	// - digest (selectable, default tryptic)
	//   for each peptides, print distribution of number of occurrences (unique)
	//   print prots with no unique peptides (merging I and L)
	//   print prots with no peptides that can be measured or uniquely identified, e.g. because of mass range, ionizability
	//

	// In Utils:
	// Search MS1 molecule: search a molecules isotopic pattern in MS1. E.g. search different glycoforms of a glycosilated peptide
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

// Convert a download URL to the pathname there we will store the file locally
func urlToPath(url string, dir string) string {
	fn := url[strings.LastIndex(url, "/")+1:]
	// Files are uncompressed before storing, remove .gz extention
	fn = strings.TrimSuffix(fn, `.gz`)
	pn := filepath.Join(dir, fn)
	return pn
}

func fuzzyNameFile(name string, dir string) (string, error) {
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
// FIXME: download to temporary name and rename when finished!
func download(url string, dir string) error {
	pn := urlToPath(url, dir)
	out, err := os.Create(pn)
	if err != nil {
		return err
	}
	defer out.Close()
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
	_, err = io.Copy(out, r)
	if err != nil {
		return err
	}
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

func updateFASTA(name string, dir string) error {
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
