package query

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"golang.org/x/exp/slog"

	"ragembeddings"
	"ragembeddings/bedrock"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	be "github.com/megaproaktiv/bedrockembedding/titan"
	"github.com/pgvector/pgvector-go"
)

func Query(c *gin.Context) {
	var req ragembeddings.QueryRequest

	err := c.BindJSON(&req)
	if err != nil {
		slog.Error("Error loading input parameter", "error", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	question := req.Question
	log.Println("Question received", question)
	// log.Println("Category", req.Category)
	// log.Println("Version", req.Version)

	ctx := context.Background()
	pgUser := os.Getenv("PGUSER")
	pgPassword := os.Getenv("PGPASSWORD")
	pgHost := os.Getenv("PGHOST")
	pgPort := os.Getenv("PGPORT")
	pgDatabase := os.Getenv("PGDATABASE")

	// Construct the connection string
	pgConnString := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=require",
		pgUser, pgPassword, pgHost, pgPort, pgDatabase,
	)

	conn, err := pgx.Connect(ctx, pgConnString)
	if err != nil {
		panic(err)
	}
	defer conn.Close(ctx)

	embedding, err := be.FetchEmbedding(question)
	if err != nil {
		panic(err)
	}

	// Todo refactor out
	rows, err := conn.Query(ctx, "SELECT id, content,context  FROM documents ORDER BY embedding <=> $1 LIMIT 10", pgvector.NewVector(embedding))
	// rows, err = conn.Query(ctx, "SELECT id, content FROM documents WHERE id != $1 ORDER BY embedding <=> (SELECT embedding FROM documents WHERE id = $1) LIMIT 5", documentId)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	promptTemplate, err := os.ReadFile("prompt.tmpl")
	var templateStr string
	templateStr = string(promptTemplate)
	content_separator := os.Getenv("CONTENT_SEPARATOR")
	if content_separator == "" {
		content_separator = "document"
	}
	preExcerpt := fmt.Sprintf("<%v>\n", content_separator)
	postExcerpt := fmt.Sprintf("</%v>\n", content_separator)

	documentExcerpts := ""
	Documents := make([]ragembeddings.RagDocument, 0)
	for rows.Next() {
		var id int64
		var content string
		var context string
		err = rows.Scan(&id, &content, &context)
		if err != nil {
			panic(err)
		}
		fmt.Println(id, content)
		documentExcerpts += preExcerpt
		documentExcerpts += content + "\n"
		documentExcerpts += postExcerpt
		Documents = append(Documents, ragembeddings.RagDocument{
			Content: content,
			Context: context,
		})

	}
	tmpl, err := template.New("Prompt").Parse(templateStr)
	if err != nil {
		log.Fatal("Error parsing template:", err)
	}
	data := ragembeddings.TemplateData{
		Question: question,
		Document: documentExcerpts,
	}
	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, data)
	if err != nil {
		log.Fatal("Error executing template:", err)
	}

	// Extract the string from the buffer
	prompt := buffer.String()

	answer := bedrock.Chat(prompt)
	response := ragembeddings.Response{
		Answer:    answer,
		Documents: Documents,
	}
	c.JSON(200, response)
}
