package validator

import "github.com/fluxynet/gocipe/asset"

type document struct {
	Basic
}

// Document validator with common formats txt, pdf, csv, html, odt, ods, odp, doc, xls, ppt, docx, xlsx, pptx
func Document(maxSize int) asset.Validator {
	return document{
		Basic: Basic{
			MaxSize: maxSize,
			AllowedMimes: []string{
				"text/plain",
				"application/pdf",
				"text/csv",
				"text/html",
				"application/vnd.oasis.opendocument.text",
				"application/vnd.oasis.opendocument.spreadsheet",
				"application/vnd.oasis.opendocument.presentation",
				"application/msword",
				"application/vnd.ms-excel",
				"application/vnd.ms-powerpoint",
				"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
				"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
				"application/vnd.openxmlformats-officedocument.presentationml.presentation",
			},
		},
	}
}
