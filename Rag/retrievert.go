package main

import (
	"context"

	"github.com/cloudwego/eino-ext/components/embedding/ark"
	milvus2 "github.com/cloudwego/eino-ext/components/retriever/milvus2"
	"github.com/cloudwego/eino-ext/components/retriever/milvus2/search_mode"
	"github.com/milvus-io/milvus/client/v2/milvusclient"
)

func Newretriever(ctx context.Context, addr string, username string, password string, dbname string, topk int, embedder *ark.Embedder) *milvus2.Retriever {

	retriever, err := milvus2.NewRetriever(ctx, &milvus2.RetrieverConfig{
		ClientConfig: &milvusclient.ClientConfig{
			Address:  addr,
			Username: username,
			Password: password,
			DBName:   dbname,
		},
		Collection: "test",
		TopK:       topk,
		SearchMode: search_mode.NewApproximate(milvus2.COSINE),
		Embedding:  embedder,
	})
	if err != nil {
		panic(err)
	}

	return retriever

}
