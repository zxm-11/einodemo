package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

// Add Tool的实现
type AddTool struct{}

func GetAddTool() tool.InvokableTool {
	return &AddTool{}
}

type AddParam struct {
	A int `json:"a"`
	B int `json:"b"`
}

func (t *AddTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "add",
		Desc: "add two numbers",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"a": {
				Type:     "number",
				Desc:     "First number",
				Required: true,
			},
			"b": {
				Type:     "number",
				Desc:     "Second number",
				Required: true,
			},
		}),
	}, nil
}

func (t *AddTool) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	p := &AddParam{}
	err := json.Unmarshal([]byte(argumentsInJSON), &p)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d", p.A+p.B), nil

}

// subTool的实现
type SubTool struct{}

func GetsubTool() tool.InvokableTool {
	return &SubTool{}
}

type SubParam struct {
	A int `json:"A"`
	B int `json:"B"`
}

func (t *SubTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "sub",
		Desc: "sub two nubers",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"a": {
				Type:     "number",
				Desc:     "First number",
				Required: true,
			},
			"b": {
				Type:     "number",
				Desc:     "Second number",
				Required: true,
			},
		}),
	}, nil
}

func (t *SubTool) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	p := &SubParam{}
	err := json.Unmarshal([]byte(argumentsInJSON), &p)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d", p.A-p.B), nil
}

type AnalyzeTool struct{}

func GetAnalyzeTool() tool.InvokableTool {
	return &AnalyzeTool{}
}

type AnalyzeParam struct {
	Content string `json:"content"`
}

func (t *AnalyzeTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "analyze",
		Desc: "analyze the difficulty of the content",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"content": {
				Type:     "string",
				Desc:     "content to be analyzed",
				Required: true,
			},
		}),
	}, nil
}

func (t *AnalyzeTool) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	//解析输入参数
	p := &AnalyzeParam{}
	err := json.Unmarshal([]byte(argumentsInJSON), &p)
	if err != nil {
		return "", err
	}
	//调用模型
	arkAPIKey := os.Getenv("ARK_API_KEY")
	arkModelName := os.Getenv("MODEL")
	arkModel, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		APIKey: arkAPIKey,
		Model:  arkModelName,
	})
	if err != nil {
		return "", err
	}
	AnalyzeInput := []*schema.Message{
		{
			Role:    schema.System,
			Content: "你是一名老师,你需要分析用户的问题,判断用户问题的难度,难度分为中等,简单,困难,你需要根据用户的问题给出一个难度的评分,评分范围为1-10,1为简单,10为困难",
		},
		{
			Role:    schema.User,
			Content: p.Content,
		},
	}
	response, err := arkModel.Generate(ctx, AnalyzeInput)
	if err != nil {
		return "", err
	}
	return response.Content, nil
}
