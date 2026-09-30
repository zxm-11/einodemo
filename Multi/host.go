package main

import (
	"context"
	"os"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/flow/agent/multiagent/host"
)

func newHost(ctx context.Context) (*host.Host, error) {
	arkApiKey := os.Getenv("ARK_API_KEY")
	arkModelName := os.Getenv("MODEL")
	chatModel, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		APIKey: arkApiKey,
		Model:  arkModelName,
	})
	if err != nil {
		return nil, err
	}

	return &host.Host{
		ToolCallingModel: chatModel,
		SystemPrompt:     "你可以同时计算加法和减法,当用户要提出问题时,你需要请求add或sub专家得到答案回答问题",
	}, nil
}
