package main

import (
	"context"

	"github.com/cloudwego/eino-ext/components/embedding/ark"
	"github.com/cloudwego/eino-ext/components/indexer/milvus2"
	"github.com/milvus-io/milvus/client/v2/milvusclient"
)

func Newindexer(ctx context.Context, addr string, username string, password string, dbname string, embedder *ark.Embedder) *milvus2.Indexer {

	indexer, err := milvus2.NewIndexer(ctx, &milvus2.IndexerConfig{
		ClientConfig: &milvusclient.ClientConfig{
			Address:  addr,
			Username: username,
			Password: password,
			DBName:   dbname,
		},
		Collection: "test",
		Vector: &milvus2.VectorConfig{
			Dimension:  2048,
			MetricType: milvus2.COSINE,
		},
		Embedding: embedder,
	})
	if err != nil {
		panic(err)
	}
	return indexer

}
