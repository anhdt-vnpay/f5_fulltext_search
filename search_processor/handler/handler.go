package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/olivere/elastic/v7"
)

func Insert(es *elastic.Client, index string, id string, data string) error {
	ctx := context.Background()
	ind, err := es.Index().
		Index(index).
		Id(id).
		BodyJson(data).
		Do(ctx)
	if err != nil {
		return err
	}

	fmt.Println("[Elastic][Insert] Successful: ", *ind)
	return nil
}

func Update(es *elastic.Client, index string, id string, data string) error {
	ctx := context.Background()
	ind, err := es.Index(). // still use Index() to completely replace old document
				Index(index).
				Id(id).
				BodyJson(data).
				Do(ctx)
	if err != nil {
		return err
	}

	fmt.Println("[Elastic][Update] Successful: ", *ind)
	return nil
}

func Delete(es *elastic.Client, index string, id string, data string) error {
	ctx := context.Background()
	ind, err := es.Delete().
		Index(index).
		Id(id).
		Do(ctx)
	if err != nil {
		return err
	}

	fmt.Println("[Elastic][Delete] Successful: ", *ind)
	return nil
}

func Search(es *elastic.Client, query string) ([]byte, error) {
	ctx := context.Background()
	fields := []string{} // left empty to search in all fields
	q := elastic.NewMultiMatchQuery(query, fields...)

	ind, err := es.Search().
		Query(q).
		Do(ctx)
	if err != nil {
		return nil, err
	}

	rs, err := parseHit(*ind)

	if err != nil {
		return nil, err
	}

	return rs, nil
}

func parseHit(searchResult elastic.SearchResult) ([]byte, error) {
	var rs []map[string]interface{}

	if searchResult.Hits.TotalHits.Value > 0 {
		fmt.Printf("Found a total of %d \n", searchResult.Hits.TotalHits.Value)

		// Iterate through results
		for _, hit := range searchResult.Hits.Hits {
			// hit.Index contains the name of the index

			// Deserialize hit.Source into a struct (could also be just a map[string]interface{})
			var t map[string]interface{}

			err := json.Unmarshal(hit.Source, &t)
			if err != nil {
				// Deserialization failed
				return nil, err
			}

			rs = append(rs, t)
		}
	}

	if len(rs) == 0 {
		return []byte{}, nil
	}

	bRs, err := json.Marshal(rs)
	if err != nil {
		return nil, err
	}
	return bRs, nil
}
