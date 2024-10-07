package types


type Document struct {
	ObjectKey string
	FileName string
}


type DocumentText struct {
	Name   string
	Chunks []string
}

type TitanEmbeddingInput struct {
	// must use camelcase for AWS here
	InputText string `json:"inputText"`
}