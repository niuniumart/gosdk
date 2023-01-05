package escli

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/olivere/elastic/v7"
)

var mapping = `{
		"mappings": {
			"properties": {
				"blockHeight":         { "type": "long" },
				"blockHash":      		{ "type": "text"},
				"createTime":        	{ "type": "date"},
				"blockSize": 		{ "type": "long" },
				"txNumber": 		{ "type": "long" }
			}
		}
	}`

type blockItem struct {
	BlockHeight int64     `json:"blockHeight"`
	BlockHash   string    `json:"blockHash"`
	CreateTime  time.Time `json:"createTime"`
	BlockSize   int64     `json:"blockSize"`
	TxNumber    int64     `json:"txNumber"`
}

var (
	id    int64 = 1
	idstr       = fmt.Sprintf("%d", id)
)

func TestEsFactory_CreateEsCli(t *testing.T) {
	cli, err := Factory.CreateEsCli("", "", []string{"http://127.0.0.1:9200"})
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	ok, err := cli.IndexExists("block").Do(ctx)
	if err != nil {
		t.Log(err)
		return
	}
	if !ok {
		// 创建create
		res, err := cli.CreateIndex("block").Body(mapping).Do(ctx)
		if err != nil {
			t.Log(err)
			return
		}
		t.Log("create index", res)
	}

	// 添加index
	idxres, err := cli.Index().Index("block").Id(idstr).BodyJson(&blockItem{
		BlockHeight: int64(id),
		BlockHash:   "2cf48b30c48c71caeb994cf1b5025628d4af32302fe662547a6e491f9d220dc8",
		CreateTime:  time.Now(),
		BlockSize:   123456,
		TxNumber:    281,
	}).Refresh("true").Do(ctx)
	if err != nil {
		t.Log(err)
		return
	}
	t.Log("create index", idxres)

	// term 查询
	term := elastic.NewTermQuery("blockHeight", 1)
	searchResult, err := cli.Search().Index("block").Query(term).Do(ctx)
	if err != nil {
		t.Log(err)
		return
	}
	if searchResult.TotalHits() > 0 {
		fmt.Printf("Found a total of %d tweets\n", searchResult.TotalHits())
		for _, hit := range searchResult.Hits.Hits {
			var item blockItem
			err := json.Unmarshal(hit.Source, &item)
			if err != nil {
				t.Log(err)
			}
			fmt.Printf("item %+v\n", item)
		}
	} else {
		fmt.Print("Found no tweets\n")
	}

	// 删除索引
	deleteIndex, err := cli.DeleteIndex("block").Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
	t.Log(deleteIndex)
}
