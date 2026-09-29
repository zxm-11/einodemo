package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/cloudwego/eino-ext/components/embedding/ark"
	milvus2 "github.com/cloudwego/eino-ext/components/retriever/milvus2"
	"github.com/cloudwego/eino-ext/components/retriever/milvus2/search_mode"
	"github.com/joho/godotenv"
	"github.com/milvus-io/milvus/client/v2/milvusclient"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	api_key := os.Getenv("ARK_API_KEY")
	model := os.Getenv("EMBEDDER")
	addr := os.Getenv("MILVUS_ADDR")
	username := os.Getenv("MILVUS_USERNAME")
	password := os.Getenv("MILVUS_PASSWORD")
	apiType := ark.APITypeMultiModal
	ctx := context.Background()
	emb, err := ark.NewEmbedder(ctx, &ark.EmbeddingConfig{
		APIKey:  api_key,
		Model:   model,
		APIType: &apiType,
	})
	if err != nil {
		log.Fatal(err)
	}

	//创建retriever
	retriever, err := milvus2.NewRetriever(ctx, &milvus2.RetrieverConfig{
		ClientConfig: &milvusclient.ClientConfig{
			Address:  addr,
			Username: username,
			Password: password,
			DBName:   "awosomeEino",
		},
		Collection: "test",
		TopK:       10,
		SearchMode: search_mode.NewApproximate(milvus2.COSINE),
		Embedding:  emb,
	})
	if err != nil {
		log.Fatal(err)
	}

	documents, err := retriever.Retrieve(ctx, "博德之门3")
	if err != nil {
		log.Fatal(err)
	}
	for i, doc := range documents {
		fmt.Printf("Document %d:\n", i)
		fmt.Printf("  ID: %s\n", doc.ID)
		fmt.Printf("  Content: %s\n", doc.Content)
		fmt.Printf("  Score: %v\n", doc.Score())
	}

}
