package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/cloudwego/eino-ext/components/embedding/ark"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	//model, err := ark.NewChatModel(ctx, &ark.
	//	ChatModelConfig{
	//	APIKey: os.Getenv("ARK_API_KEY"),
	//	Model:  os.Getenv("MODEL"),
	//})
	//if err != nil {
	//	panic(err)
	//}
	//
	//template := prompt.FromMessages(schema.FString,
	//	schema.SystemMessage("你是一个{role}"),
	//	&schema.Message{
	//		Role:    schema.User,
	//		Content: "请你帮帮我，影心,帮我解决{task}",
	//	},
	//)
	//
	//params := map[string]any{
	//	"role": "影心",
	//	"task": "帮我偷(刷)金币",
	//}
	//messages, err := template.Format(ctx, params)
	//if err != nil {
	//	panic(err)
	//}
	//
	////input := []*schema.Message{
	////	schema.SystemMessage("你是一个可爱的高中美少女"),
	////	schema.UserMessage("你好"),
	////}
	////Genert式生成
	//response, err := model.Generate(ctx, messages)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(response.Content)
	//
	////流式生成
	//reader, err := model.Stream(ctx, input)
	//if err != nil {
	//	panic(err)
	//}
	//defer reader.Close()
	//
	//for {
	//	chunk, err := reader.Recv()
	//	if err != nil {
	//		panic(err)
	//	}
	//	fmt.Print(chunk.Context)
	apiType := ark.APITypeMultiModal
	embedder, err := ark.NewEmbedder(ctx, &ark.EmbeddingConfig{
		APIKey:  os.Getenv("ARK_API_KEY"),
		Model:   os.Getenv("EMBEDDER"),
		APIType: &apiType,
	})
	if err != nil {
		panic(err)
	}

	vectors, err := embedder.EmbedStrings(ctx, []string{
		"bdzm",
		"yingxin"})
	if err != nil {
		panic(err)
	}
	log.Println(vectors)

	fmt.Printf("向量维度%s:", len(vectors))

}
