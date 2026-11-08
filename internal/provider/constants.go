// Package provider implements the idgen Terraform provider and its data sources.
package provider

const (
	// MaxIDLength is the maximum allowed length for generated IDs.
	MaxIDLength = 1024

	// MinIDLength is the minimum allowed length for generated IDs.
	MinIDLength = 1

	// WarnIDLength is the threshold at which we warn about unusually long IDs.
	WarnIDLength = 128

	// defaultNanoIDLength is the NanoID length used when none is configured.
	defaultNanoIDLength = 21
)

// Schema attribute names shared by several data sources.
const (
	attrID        = "id"
	attrSeed      = "seed"
	attrLength    = "length"
	attrGroupSize = "group_size"
)

// Proquint layout: each word encodes 2 bytes as 5 characters, words are joined by one separator.
const (
	proquintBytesPerWord = 2
	proquintCharsPerWord = 6 // 5 letters plus the separator
	uint32Bytes          = 4
	uint64Bytes          = 8
)
