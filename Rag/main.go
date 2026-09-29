package main

import (
	"context"
	"log"
	"os"

	"github.com/cloudwego/eino-ext/components/embedding/ark"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	addr := os.Getenv("MILVUS_ADDR")
	username := os.Getenv("MILVUS_USERNAME")
	password := os.Getenv("MILVUS_PASSWORD")
	apikey := os.Getenv("ARK_API_KEY")
	//chatmodel := os.Getenv("MODEL")
	embeddermodel := os.Getenv("EMBEDDER")

	muti_apikey := ark.APITypeMultiModal
	dbname := "awosomeEino"
	topk := 2
	ctx := context.Background()

	//初始化组件
	embedder := Newembedde(ctx, muti_apikey, apikey, embeddermodel)
	//indexer := Newindexer(ctx, addr, username, password, dbname, embedder)
	retriever := Newretriever(ctx, addr, username, password, dbname, topk, embedder)
	//splitter := NewTrans(ctx)

	//bs, err := os.ReadFile("docs/baldurs-gate.md")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//docs := []*schema.Document{
	//	{
	//		ID:      "doc1",
	//		Content: string(bs),
	//	},
	//}
	//
	////transformer组件分割文档
	//results, err := splitter.Transform(ctx, docs)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	////indexer组件写入向量库
	//for i, doc := range results {
	//	doc.ID = docs[0].ID + "_" + strconv.Itoa(i)
	//	fmt.Println(doc.ID)
	//}
	//
	//ids, err := indexer.Store(ctx, results)
	//
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println(ids)
	results, err := retriever.Retrieve(ctx, "博得之门的主要职业")
	if err != nil {
		log.Fatal(err)
	}
	for _, doc := range results {
		println(doc.ID)
		println(doc.Content)
		println("——————————————————")
	}

}
