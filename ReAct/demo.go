package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	arkAPIKey := os.Getenv("ARK_API_KEY")
	model := os.Getenv("MODEL")

	ctx := context.Background()
	arkModel, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		APIKey: arkAPIKey,
		Model:  model,
	})
	if err != nil {
		log.Fatal(err)
	}

	addTool := GetAddTool()
	subTool := GetSubTool()
	analyzeTool := GetAnalyzeTool()
	persona := `#Character:你是一个老师,会同时判断题目难易程度,给出问题的答案`
	raAgent, err := react.NewAgent(ctx, &react.AgentConfig{
		ToolCallingModel: arkModel,
		ToolsConfig: compose.ToolsNodeConfig{
			Tools:               []tool.BaseTool{addTool, subTool, analyzeTool},
			ExecuteSequentially: false,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	chatmsg := []*schema.Message{
		{
			Role:    schema.System,
			Content: persona,
		},
		{
			Role:    schema.User,
			Content: "请告诉我183*92-66这道题的难易程度和答案",
		},
	}

	resp, err := raAgent.Generate(ctx, chatmsg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp.Content)
}
