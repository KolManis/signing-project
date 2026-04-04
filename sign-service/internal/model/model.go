package model

type Document struct {
	UserID string
	Data   []byte
}

type SignedDocument struct {
	DocumentID string
	Signature  string
}
