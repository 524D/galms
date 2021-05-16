// Copyright 2021 Rob Marissen
// SPDX-License-Identifier: MIT

package cmd

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/524D/galms/digest"
	"github.com/524D/galms/fasta"

	"github.com/spf13/cobra"
)

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

		fastas := make([]fasta.Fasta, 0)
		for _, n := range args {
			upd, err := cmd.Flags().GetBool("update")
			if err != nil {
				log.Fatalf("Getstrings 'update' flag failed: %v", err)
			}
			if upd {
				err = fasta.UpdateFASTA(n, dataDir)
			}
			if err != nil {
				log.Fatalf("%v", err)
			}

			fn, err := fasta.FuzzyNameFile(n, dataDir)
			if err != nil {
				log.Fatalf("%s is not a valid filename nor a fuzzy name: %v", n, err)
			}
			file, err := os.Open(fn)
			if err != nil {
				log.Fatalf("Can't open file %s: %v", args[0], err)
			}
			defer file.Close()

			f, err := fasta.Read(file)
			if err != nil {
				log.Fatal("fasta Read failed")
			}
			fastas = append(fastas, f)
		}

		missing, err := cmd.Flags().GetString("missing")
		if err != nil {
			log.Fatalf("Getstrings 'missing' flag failed: %v", err)
		}
		contains, err := cmd.Flags().GetString("contains")
		if err != nil {
			log.Fatalf("Getstrings 'contains' flag failed: %v", err)
		}
		enzymeName, err := cmd.Flags().GetString("enzyme")
		if err != nil {
			log.Fatalf("Getstrings 'enzyme' flag failed: %v", err)
		}
		enzyme, err := digest.NamedEnzyme(enzymeName)
		if err != nil {
			log.Fatalf("%s: %v", enzymeName, err)
		}
		an, err := cmd.Flags().GetBool("analyse")
		if err != nil {
			log.Fatalf("Getstrings 'analyse' flag failed: %v", err)
		}
		pt, err := cmd.Flags().GetBool("proteotypic")
		if err != nil {
			log.Fatalf("Getstrings 'proteotypic' flag failed: %v", err)
		}

		for _, f := range fastas {
			if an {
				fasta.Analyse(f, enzyme)
			}

			if contains != `` || missing != `` {
				for _, p := range f.Prots() {
					if contains == `` || strings.Contains(p.Sequence(), contains) {

						if missing == `` || !strings.Contains(p.Sequence(), missing) {
							fmt.Printf("Length: %d %s %s\n", len(p.Sequence()), p.ID(), p.Description())
						}
					}
				}
			}
		}

		if pt {
			pts := fasta.ProteotypicPeps(fastas, enzyme)
			fasta.WriteProteotypicPeps(os.Stdout, pts)
		}

	},
}

func init() {
	rootCmd.AddCommand(fastaCmd)

	fastaCmd.PersistentFlags().StringP("missing", "m", "", "List proteins which don't contain the specified sequence")
	fastaCmd.PersistentFlags().StringP("contains", "c", "", "List proteins which contain the specified sequence")
	fastaCmd.PersistentFlags().StringP("enzyme", "e", "trypsin", "Use the specified cleavage enzyme {Trypsin,Trypsin_Simple, Trypsin/P,Lys_C,PepsinA,Chymotrypsin}")
	fastaCmd.PersistentFlags().BoolP("update", "u", false, "Update FASTA file")
	fastaCmd.PersistentFlags().BoolP("analyse", "a", false, "Analyse proteins")
	fastaCmd.PersistentFlags().BoolP("proteotypic", "p", false, "List proteotypic peptides")

}
