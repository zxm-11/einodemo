package main

import (
	"context"
	"log"
	"os"

	"github.com/cloudwego/eino-ext/components/embedding/ark"
	milvus2 "github.com/cloudwego/eino-ext/components/indexer/milvus2"
	"github.com/cloudwego/eino/schema"
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

	//创建embedding模型
	emb, err := ark.NewEmbedder(ctx, &ark.EmbeddingConfig{
		APIKey:  api_key,
		Model:   model,
		APIType: &apiType,
	})
	if err != nil {
		log.Fatal(err)
	}

	//创建索引器
	indexer, err := milvus2.NewIndexer(ctx, &milvus2.IndexerConfig{
		ClientConfig: &milvusclient.ClientConfig{
			Address:  addr,
			Username: username,
			Password: password,
			DBName:   "awosomeEino",
		},
		Collection: "test",
		Vector: &milvus2.VectorConfig{
			Dimension:  2048,
			MetricType: milvus2.COSINE,
		},
		Embedding: emb,
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Indexer created successfully")

	docs := []*schema.Document{
		{
			ID:      "doc1",
			Content: "你说的对,但原神是一款开放大世界二次元游戏",
			MetaData: map[string]any{
				"auth": "EchoEcho!!",
			},
		},
		{
			ID:      "doc2",
			Content: "你说的对,但博德之门3是一款开放世界DND游戏",
			MetaData: map[string]any{
				"auth": "zxm-11",
			},
		},
	}

	ids, err := indexer.Store(ctx, docs)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Stored %v documents successful", ids)

}
