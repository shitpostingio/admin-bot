package analysisadapter

import "errors"

var (

	// FingerprintError is returned when a fingerprint wasn't successful.
	FingerprintError = errors.New("analysis: unable to perform fingerprint")

	// NSFWError is returned when a NSFW analysis wasn't successful.
	NSFWError = errors.New("analysis: unable to perform NSFW analysis")

	// GibberishError is returned when a gibberish analysis wasn't successful.
	GibberishError = errors.New("analysis: unable to perform gibberish analysis")
)
