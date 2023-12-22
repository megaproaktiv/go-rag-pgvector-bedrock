package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"

	be "github.com/megaproaktiv/bedrockembedding/titan"

	"github.com/jackc/pgx/v5"
	"github.com/pgvector/pgvector-go"
	dl "github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/textsplitter"
)

func main() {
	filename := flag.String("f", "testdata/sample.pdf", "Filename to process")
	dropTable := flag.Bool("d", false, "Drop table before import")
	// Parse the flags.
	flag.Parse()

	// Postgres
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
	_, err = conn.Exec(ctx, "CREATE EXTENSION IF NOT EXISTS vector")
	if err != nil {
		panic(err)
	}

	if *dropTable {
		log.Println("Drop table, Create table")
		_, err = conn.Exec(ctx, "DROP TABLE IF EXISTS documents")
		if err != nil {
			panic(err)
		}

		_, err = conn.Exec(ctx, "CREATE TABLE documents (id bigserial PRIMARY KEY, content text, context text, embedding vector(1536))")
		if err != nil {
			panic(err)
		}

		_, err = conn.Exec(ctx, "CREATE INDEX embedding_idx ON documents USING ivfflat(embedding)")
		if err != nil {
			panic(err)
		}
	}
	fmt.Printf("Start import %v\n", *filename)
	f, err := os.Open(*filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	finfo, err := f.Stat()
	if err != nil {
		panic(err)
	}
	p := dl.NewPDF(f, finfo.Size())

	split := textsplitter.NewRecursiveCharacter()
	split.ChunkSize = 300
	split.ChunkOverlap = 50

	docs, err := p.LoadAndSplit(context.Background(), split)
	if err != nil {
		panic(err)
	}
	// Sentence-window retrieval

	for i, doc := range docs {
		re := regexp.MustCompile(`(\w{2,})-(\w{2,})`)
		docs[i].PageContent = re.ReplaceAllString(doc.PageContent, "$1$2")

	}
	for i, doc := range docs {
		content := doc.PageContent
		singleEmbedding, err := be.FetchEmbedding(content)
		if err != nil {
			panic(err)
		}
		context := content
		if i > 0 && i < len(docs)-1 {
			context = docs[i-1].PageContent + docs[i].PageContent + docs[i+1].PageContent
		}

		_, err = conn.Exec(ctx, "INSERT INTO documents (content, context, embedding) VALUES ($1, $2, $3)", content, context, pgvector.NewVector(singleEmbedding))
		if err != nil {
			panic(err)
		}
		fmt.Printf("Cluster %v: {%v}\n / [%v]\n", i+1, content, context)
	}

}
