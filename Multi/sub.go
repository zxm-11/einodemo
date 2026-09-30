package main

import (
	"context"
	"os"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent"
	"github.com/cloudwego/eino/flow/agent/multiagent/host"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"
)

func newSubSpecialist(ctx context.Context) (*host.Specialist, error) {

	arkApiKey := os.Getenv("ARK_API_KEY")
	arkModelName := os.Getenv("MODEL")
	chatMode, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		APIKey: arkApiKey,
		Model:  arkModelName,
	})
	if err != nil {
		return nil, err
	}

	subTool := GetsubTool()
	raAgent, err := react.NewAgent(ctx, &react.AgentConfig{
		ToolCallingModel: chatMode,
		ToolsConfig: compose.ToolsNodeConfig{
			Tools: []tool.BaseTool{subTool},
		},
	})
	if err != nil {
		return nil, err
	}
	return &host.Specialist{
		AgentMeta: host.AgentMeta{
			Name:        "sub_specialist",
			IntendedUse: "sub Two numbers and return the result",
		},
		Invokable: func(ctx context.Context, input []*schema.Message, opts ...agent.AgentOption) (output *schema.Message, err error) {
			return raAgent.Generate(ctx, input, opts...)
		},
	}, nil

}
