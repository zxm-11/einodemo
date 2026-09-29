package main

import (
	"context"
	"log"

	"github.com/cloudwego/eino-ext/components/model/ark"
)

func NewArkModel(ctx context.Context, apikey string, chatmodel string) *ark.ChatModel {

	model, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		APIKey: apikey,
		Model:  chatmodel,
	})
	if err != nil {
		log.Fatal(err)
	}

	return model
}
