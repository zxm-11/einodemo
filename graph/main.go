package main

import (
	"context"
	"fmt"
	"log"

	"github.com/cloudwego/eino/compose"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	ctx := context.Background()

	//注册图
	g := compose.NewGraph[string, string]()
	//编写节点
	lambda0 := compose.InvokableLambda(func(ctx context.Context, input string) (output string, err error) {
		if input == "1" {
			return "haomao", nil
		} else if input == "2" {
			return "耄耋", nil
		} else if input == "3" {
			return "device", nil
		}
		return "", nil
	})

	lambda1 := compose.InvokableLambda(func(ctx context.Context, input string) (output string, err error) {
		return "喵!", nil
	})
	lambda2 := compose.InvokableLambda(func(ctx context.Context, input string) (output string, err error) {
		return "哈!", nil
	})
	lambda3 := compose.InvokableLambda(func(ctx context.Context, input string) (output string, err error) {
		return "没有人类了!", nil
	})
	//加入节点
	err = g.AddLambdaNode("lambda0", lambda0)

	if err != nil {
		log.Fatal("Error adding lambda0")
	}
	err = g.AddLambdaNode("lambda1", lambda1)
	if err != nil {
		log.Fatal("Error adding lambda1")
	}
	err = g.AddLambdaNode("lambda2", lambda2)
	if err != nil {
		log.Fatal("Error adding lambda2")
	}
	err = g.AddLambdaNode("lambda3", lambda3)
	if err != nil {
		log.Fatal("Error adding lambda3")
	}
	//加入分支
	err = g.AddBranch("lambda0", compose.NewGraphBranch(func(ctx context.Context, in string) (endNode string, err error) {
		if in == "haomao" {
			return "lambda1", nil
		} else if in == "耄耋" {
			return "lambda2", nil
		} else if in == "device" {
			return "lambda3", nil
		}
		return compose.END, nil
	}, map[string]bool{"lambda1": true, "lambda2": true, "lambda3": true}))
	//连接节点
	err = g.AddEdge(compose.START, "lambda0")
	if err != nil {
		log.Fatal("Error adding Lambda0")
	}
	err = g.AddEdge("lambda1", compose.END)
	if err != nil {
		log.Fatal("Error adding Lambda1")
	}
	err = g.AddEdge("lambda2", compose.END)
	if err != nil {
		log.Fatal("Error adding Lambda2")
	}
	err = g.AddEdge("lambda3", compose.END)
	if err != nil {
		log.Fatal("Error adding Lambda3")
	}
	//编译运行
	r, err := g.Compile(ctx)
	if err != nil {
		log.Fatal("Error compile")
	}

	answer, err := r.Invoke(ctx, "3")
	if err != nil {
		log.Fatal("Error invoke")
	}
	fmt.Println(answer)

}
