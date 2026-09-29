package main

import (
	"context"
	"log"

	"github.com/cloudwego/eino-ext/components/embedding/ark"
)

func Newembedde(ctx context.Context, muti_apikey ark.APIType, apikey string, embeddermodel string) *ark.Embedder {
	apiKey := muti_apikey
	embedder, err := ark.NewEmbedder(ctx, &ark.EmbeddingConfig{
		APIKey:  apikey,
		Model:   embeddermodel,
		APIType: &apiKey,
	})
	if err != nil {
		log.Fatal(err)
	}
	return embedder
}
