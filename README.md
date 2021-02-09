# galms

**G**o **A**pp and **L**ibrary for **M**ass **S**pectrometry

## Library

The library in galms contains Go modules to enable simple creation of efficient MS software tools:

* Read/write common files in common formats (mzML, mzID, mzXML, pepXML, FASTA)
* Compute masses and isotopic distributions
* Convert various representations of molecules into a molecular formula (amino acids, glycans, ...).
* Digest proteins into peptides
* Predict various LC/MS experiment values (retention times, fragmentation patterns, ionization efficiency)
* Conversion of nucleotide sequence into peptide sequence
* Use web services and obtain data from EBI EMBL
* General purpose functions (binning, ...)

## Tools

The app provides various MS tools. These serve to
demonstrate the use of the library, but are also useful
as general purpose MS tools.

All tools are accessed as sub commands of galms:

* galms isotopes: Compute isotopes
* galms decoy: Create decoy databases
* galms translate: Translate nucleotide sequence into peptide sequence

## Web server

galms contains a build-in webserver to give easy to use network access to the full toolbox without installation.

## Why another Mass Spectrometry library?

[Go](https://golang.org/) is a programming language that makes it easy to build simple, reliable, and efficient software. This allows fast creation of high quality software, an important goal in modern (academic) environments.

MS libraries already exists for other programming languages, some carrying decades of development effort. However, these have been created in programming languages (in many cases also decades old) that more complex or less efficient. We think the advantages of the Go programming language outweighs the completeness of available libraries.


