package main

import (
	"context"
	"os"

	"github.com/cloudwego/eino-ext/components/document/transformer/splitter/markdown"
	"github.com/cloudwego/eino/schema"
)

func main() {
	ctx := context.Background()
	transformer, err := markdown.NewHeaderSplitter(ctx, &markdown.HeaderConfig{
		Headers: map[string]string{
			"#":   "h1",
			"##":  "h2",
			"###": "h3",
		},
		TrimHeaders: true,
	})
	if err != nil {
		panic(err)
	}
	content, err := os.OpenFile("docs/baldurs-gate.md", os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		panic(err)
	}
	defer content.Close()

	bs, err := os.ReadFile("docs/baldurs-gate.md")
	if err != nil {
		panic(err)
	}
	docs := []*schema.Document{
		{
			ID:      "doc1",
			Content: string(bs),
		},
	}

	splitDocs, err := transformer.Transform(ctx, docs)
	if err != nil {
		panic(err)
	}

	for i, doc := range splitDocs {
		println("片段", i+1, ":", doc.Content)
		println("标题层级:")
		for k, v := range doc.MetaData {
			if k == "h1" || k == "h2" || k == "h3" {
				println(" ", k, ":", v)
			}
		}
	}

}
