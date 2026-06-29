package model

// ProblemDetails is an RFC 9457 "Problem Details for HTTP APIs" object.
//
// RFC 9457 (https://www.rfc-editor.org/info/rfc9457/) obsoletes the older
// RFC 7807. The wire format is unchanged; responses carrying this object
// should use the "application/problem+json" media type.
type ProblemDetails struct {
	// Type is a URI reference identifying the problem type. When there is no
	// semantics beyond the HTTP status code, RFC 9457 recommends "about:blank".
	Type string `json:"type"`
	// Title is a short, human-readable summary of the problem type.
	Title string `json:"title"`
	// Status is the HTTP status code generated for this occurrence.
	Status int `json:"status"`
	// Detail is a human-readable explanation specific to this occurrence.
	Detail string `json:"detail"`
	// Instance is a URI reference identifying the specific occurrence.
	Instance string `json:"instance"`
}
