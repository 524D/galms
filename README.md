# galms

**G**o **A**pp and **L**ibrary for **M**ass **S**pectrometry

![CodeQL](https://github.com/524D/galms/actions/workflows/codeql-analysis.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/524D/galms)](https://goreportcard.com/report/github.com/524D/galms)
[![Coverage Status](https://coveralls.io/repos/github/524D/galms/badge.svg?branch=main)](https://coveralls.io/github/524D/galms?branch=main)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/524D/galms/blob/master/LICENSE)

# *WARNING*

*GALMS is "work in progress".*
*Don't use for any serious work. Until version 1.0 is reached any part of the software is likely to change in ways that break all
dependencies.*

## Building

Minimum required Go version: 1.16

```bash
go install github.com/524D/galms@latest
```

This puts the `galms` (for Windows users: `galms.exe`) app into directory `${GOPATH}/bin`

## Library

The library in galms contains Go packages to enable simple creation of efficient MS software tools:

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

* TODO: galms isotopes: Compute isotopes
* TODO: galms decoy: Create decoy databases
* TODO: galms translate: Translate nucleotide sequence into peptide sequence

## Web server

TODO: galms contains a build-in webserver to give easy to use network access to the full toolbox without installation.

## Go advocacy

[This is why we use Go](whygo.md)
