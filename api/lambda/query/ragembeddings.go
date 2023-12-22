package ragembeddings

type QueryRequest struct {
	Question string `json:"question"`
	// Category string `json:"category"`
	// Version  string `json:"version"`
}

type RagDocument struct {
	Content string `json:"content"`
	Context string `json:"context"`
}

type Response struct {
	Answer    string        `json:"answer"`
	Documents []RagDocument `json:"documents"`
}

type TemplateData struct {
	Question string
	Document string
}
