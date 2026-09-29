package main

import (
	"context"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
)

type Game struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

// 定义入参结构体(模型需要给tool提供的json参数)
type InputParams struct {
	Name string `json:"name" jsonschema:"description=the name of game"`
}

// GetGame :本地函数
func GetGame(_ context.Context, Params *InputParams) (string, error) {
	GameSet := []Game{
		{Name: "博德之门", Url: "https://baldursgate3.game/"},
		{Name: "原神", Url: "https://www.yuanshen.com/"},
	}
	for _, game := range GameSet {
		if game.Name == Params.Name {
			return game.Url, nil
		}
	}
	return "", nil
}

// CreateTool 创建tool
func CreateTool() tool.InvokableTool {
	getGameTool := utils.NewTool(
		&schema.ToolInfo{
			Name: "get_name",
			Desc: "get a game url by name",
			ParamsOneOf: schema.NewParamsOneOfByParams(
				map[string]*schema.ParameterInfo{
					"name": &schema.ParameterInfo{
						Type:     schema.String,
						Desc:     "game's name",
						Required: true,
					},
				},
			),
		}, GetGame)
	return getGameTool
}
