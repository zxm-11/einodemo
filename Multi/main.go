package main

import (
	"context"
	"fmt"
	"log"

	"github.com/cloudwego/eino/flow/agent/multiagent/host"
	"github.com/cloudwego/eino/schema"
	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//路由
	h, err := newHost(ctx)
	if err != nil {
		panic(err)
	}

	//两个专家
	adder, err := newAddSpecialist(ctx)
	if err != nil {
		panic(err)
	}
	suber, err := newSubSpecialist(ctx)
	if err != nil {
		panic(err)
	}

	hostMA, err := host.NewMultiAgent(ctx, &host.MultiAgentConfig{
		Host: *h,
		Specialists: []*host.Specialist{
			adder,
			suber,
		},
		Summarizer: &host.Summarizer{
			ChatModel:    h.ToolCallingModel,
			SystemPrompt: "请总结一下各个专家的回答",
		},
	})
	if err != nil {
		panic(err)
	}
	input := []*schema.Message{
		{
			Role:    schema.User,
			Content: "帮我计算17899+114514-25922",
		},
	}

	resp, err := hostMA.Generate(ctx, input)
	if err != nil {
		panic(err)
	}

	fmt.Println(resp.Content)

}
