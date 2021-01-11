# gomztal

**Go** ***m***/***z*** **T**ools **A**nd **L**ibrary is a Mass Spectrometry (MS) toolkit and library in [Go](https://golang.org/) 

## Why?

Go is a programming language that makes it easy to build simple, reliable, and efficient software. This allows fast creation of high quality software, an important goal in modern (academic) environments.

For other programming languages, many MS libraries already exists, some carrying decades of development effort. However, these have been created in programming languages (in many cases also decades old) that have a different level of complexity or efficiency. We think the advantages of the Go programming language outweighs the completeness of available libraries.

## Library

The library in gomztal contains Go modules to enable simple creation of efficient MS software tools:

* Read/write common files in common formats (mzML, mzID, mzXML, pepXML, FASTA)
* Compute masses and isotopic distributions
* Convert various representations of molecules into a molecular formula (amino acids, glycans, ...).
* Digest proteins into peptides
* Predict various LC/MS experiment values (retention times, fragmentation patterns, ionization efficiency)
* Conversion of nucleotide sequence into peptide sequence
* Use web services and obtain data from EBI EMBL
* General purpose functions (binning, ...)

## Tools

Various MS tools are provided, both to demonstrate the use of the library, and as useful general purpose tools.

All tools are accessed as sub commands of gomztal:

* gomztal isotopes: Compute isotopes
* gomztal decoy: Create decoy databases
* gomztal translate: Translate nucleotide sequence into peptide sequence

## Web server

gomztal contains a build-in webserver to give easy to use network access to the full toolbox without installation.

