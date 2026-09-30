package main

import (
	"context"
	"fmt"
	"log"
	"os"

	ccb "github.com/cloudwego/eino-ext/callbacks/cozeloop"
	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/coze-dev/cozeloop-go"
	"github.com/joho/godotenv"
)

type State struct {
	History map[string]any
}

func genFunc(ctx context.Context) *State {
	return &State{History: make(map[string]any)}
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	model := os.Getenv("MODEL")
	apikey := os.Getenv("ARK_API_KEY")

	ctx := context.Background()

	client, err := cozeloop.NewClient()
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close(ctx)

	handler := ccb.NewLoopHandler(client)
	callbacks.AppendGlobalHandlers(handler)

	//注册graph
	outsideGraph := compose.NewGraph[map[string]string, string]()
	outlambda1 := compose.InvokableLambda(func(ctx context.Context, input map[string]string) (output map[string]string, err error) {
		//直接返回输出
		return input, nil
	})

	err = outsideGraph.AddLambdaNode("outlambda1", outlambda1)
	if err != nil {
		log.Fatal(err)
	}
	err = outsideGraph.AddEdge(compose.START, "outlambda1")
	if err != nil {
		log.Fatal(err)
	}

	insideGraph := compose.NewGraph[map[string]string, *schema.Message](
		compose.WithGenLocalState(genFunc),
	)

	lambda := compose.InvokableLambda(func(ctx context.Context, input map[string]string) (output map[string]string, err error) {

		_ = compose.ProcessState[*State](ctx, func(ctx context.Context, state *State) (err error) {
			state.History["tsundere_action"] = "我喜欢你"
			state.History["cute_action"] = "摸摸头"
			return nil
		})

		if input["role"] == "tsundere" {
			return map[string]string{
				"role":    "tsundere",
				"content": input["content"],
			}, nil
		}
		if input["role"] == "cute" {
			return map[string]string{
				"role":    "cute",
				"content": input["content"],
			}, nil
		}
		return map[string]string{"role": "user", "content": input["content"]}, nil
	})

	TsundereLambda := compose.InvokableLambda(func(ctx context.Context, input map[string]string) (output []*schema.Message, err error) {
		return []*schema.Message{
			{
				Role:    schema.System,
				Content: "你是一个高冷傲娇的大小姐,每次都会用傲娇的语气来回答我的问题",
			},
			{
				Role:    schema.User,
				Content: input["content"],
			},
		}, nil
	})
	CuteLambda := compose.InvokableLambda(func(ctx context.Context, input map[string]string) (output []*schema.Message, err error) {
		return []*schema.Message{
			{
				Role:    schema.System,
				Content: "你是一个可爱的小女孩,每次都会用可爱的语气来回答我的问题",
			},
			{
				Role:    schema.User,
				Content: input["content"],
			},
		}, nil
	})

	//state
	TsunderePreHandler := func(ctx context.Context, input map[string]string, state *State) (output map[string]string, err error) {
		input["content"] = input["content"] + state.History["tsundere_action"].(string)
		return input, nil
	}

	CutePreHandler := func(ctx context.Context, input map[string]string, state *State) (output map[string]string, err error) {
		input["content"] = input["content"] + state.History["cute_action"].(string) //any->类型断言
		return input, nil
	}

	chatmodel, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		APIKey: apikey,
		Model:  model,
	})
	if err != nil {
		log.Fatal(err)
	}
	err = insideGraph.AddLambdaNode("lambda", lambda)
	if err != nil {
		log.Fatal(err)
	}
	err = insideGraph.AddLambdaNode("tsundereNode", TsundereLambda, compose.WithStatePreHandler(TsunderePreHandler))
	if err != nil {
		log.Fatal(err)
	}
	err = insideGraph.AddLambdaNode("cuteNode", CuteLambda, compose.WithStatePreHandler(CutePreHandler))
	if err != nil {
		log.Fatal(err)
	}
	err = insideGraph.AddChatModelNode("model", chatmodel)
	if err != nil {
		log.Fatal(err)
	}
	//编写分支
	insideGraph.AddBranch("lambda", compose.NewGraphBranch(func(ctx context.Context, in map[string]string) (endNode string, err error) {
		if in["role"] == "tsundere" {
			return "tsundereNode", err
		}
		if in["role"] == "cute" {
			return "cuteNode", err
		}
		return "TsundereNode", nil
	}, map[string]bool{
		"tsundereNode": true, "cuteNode": true,
	}))
	//添加边(连接节点)
	err = insideGraph.AddEdge(compose.START, "lambda")
	if err != nil {
		log.Fatal(err)
	}
	err = insideGraph.AddEdge("tsundereNode", "model")
	if err != nil {
		log.Fatal(err)
	}
	err = insideGraph.AddEdge("cuteNode", "model")
	if err != nil {
		log.Fatal(err)
	}
	err = insideGraph.AddEdge("model", compose.END)
	if err != nil {
		log.Fatal(err)
	}

	witelambda := compose.InvokableLambda(func(ctx context.Context, input *schema.Message) (output string, err error) {
		f, err := os.OpenFile("orc_graph_withgraph.md", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644) //追加写 + 没有就创建 + 只写
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		if _, err := f.WriteString(input.Content + "\n---\n"); err != nil {
			return "", err
		}
		return "已经写入文件,请前往文件内查看内容", nil
	})

	//嵌套图: 子图要用 AddGraphNode, AddLambdaNode 只收 *Lambda
	err = outsideGraph.AddGraphNode("insideGraph", insideGraph)
	if err != nil {
		log.Fatal(err)
	}
	err = outsideGraph.AddLambdaNode("writelambda", witelambda)
	if err != nil {
		log.Fatal(err)
	}

	err = outsideGraph.AddEdge("outlambda1", "insideGraph")
	if err != nil {
		log.Fatal(err)
	}
	err = outsideGraph.AddEdge("insideGraph", "writelambda")
	if err != nil {
		log.Fatal(err)
	}
	err = outsideGraph.AddEdge("writelambda", compose.END)
	if err != nil {
		log.Fatal(err)
	}
	//编译
	//	r, err := insideGraph.Compile(ctx)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//
	//	answer, err := r.Invoke(ctx, map[string]string{
	//		"role":    "tsundere",
	//		"content": "你好啊",
	//	}, compose.WithCallbacks(genCallback()))
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	fmt.Println(answer.Content)
	//

	r, err := outsideGraph.Compile(ctx)
	if err != nil {
		log.Fatal(err)
	}
	result, err := r.Invoke(ctx, map[string]string{
		"role":    "tsundere",
		"content": "嘻嘻",
	})
	fmt.Println(result)

}

func genCallback() callbacks.Handler {
	handler := callbacks.NewHandlerBuilder().OnStartFn(func(ctx context.Context, info *callbacks.RunInfo, input callbacks.CallbackInput) context.Context {
		fmt.Printf("当前%s输入:%s\n", info.Component, input)
		return ctx
	}).OnEndFn(func(ctx context.Context, info *callbacks.RunInfo, output callbacks.CallbackOutput) context.Context {
		fmt.Printf("当前%s节点输出:%s\n", info.Component, output)
		return ctx
	}).Build()
	return handler
}
