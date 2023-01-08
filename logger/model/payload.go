package model

import "time"

type Payload struct {
	Time           time.Time       `json:"time"`
	Severity       string          `json:"severity"`
	Message        string          `json:"message"`
	Trace          string          `json:"logging.googleapis.com/trace,omitempty"`
	SpanID         string          `json:"logging.googleapis.com/spanId,omitempty"`
	SourceLocation *SourceLocation `json:"logging.googleapis.com/sourceLocation,omitempty"`
}

type SourceLocation struct {
	File     string `json:"file"`
	Line     int    `json:"line,string"`
	Function string `json:"function"`
}
