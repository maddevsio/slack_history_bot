package service

import (
	"fmt"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/index/store/goleveldb"
	"github.com/blevesearch/blevex/lang/ru"
	"github.com/gen1us2k/log"
)

type SearchService struct {
	BaseService
	sh         *SlackHistoryBot
	logger     log.Logger
	index      bleve.Index
	batchCount int
	batch      *bleve.Batch
}

func (ss *SearchService) Name() string {
	return "search_service"
}
func (ss *SearchService) Run() error {
	return nil
}
func (ss *SearchService) Init(sh *SlackHistoryBot) error {
	ss.sh = sh
	ss.logger = log.NewLogger(ss.Name())
	indexName := "history.bleve"
	index, err := bleve.Open(indexName)
	if err == bleve.ErrorIndexPathDoesNotExist {
		mapping := ss.buildMapping()
		kvStore := goleveldb.Name
		kvConfig := map[string]interface{}{
			"create_if_missing": true,
			//		"write_buffer_size":         536870912,
			//		"lru_cache_capacity":        536870912,
			//		"bloom_filter_bits_per_key": 10,
		}

		index, err = bleve.NewUsing(indexName, mapping, "upside_down", kvStore, kvConfig)
	}
	if err != nil {
		return err
	}
	ss.index = index
	ss.batch = index.NewBatch()
	return nil
}
func (ss *SearchService) buildMapping() *bleve.IndexMapping {
	ruFieldMapping := bleve.NewTextFieldMapping()
	ruFieldMapping.Analyzer = ru.AnalyzerName

	eventMapping := bleve.NewDocumentMapping()
	eventMapping.AddFieldMappingsAt("message", ruFieldMapping)

	mapping := bleve.NewIndexMapping()
	mapping.DefaultMapping = eventMapping
	mapping.DefaultAnalyzer = ru.AnalyzerName
	return mapping
}

func (ss *SearchService) IndexMessage(data IndexData) error {
	ss.batch.Index(data.ID, data)
	if ss.batch.Size() > 100 {
		err := ss.index.Batch(ss.batch)
		if err != nil {
			return err
		}
		ss.batchCount += ss.batch.Size()
		ss.batch = ss.index.NewBatch()
	}
	if ss.batch.Size() > 0 {
		ss.index.Batch(ss.batch)
		ss.batchCount += ss.batch.Size()
	}
	return nil
}

func (ss *SearchService) Search(query, channel string) (*bleve.SearchResult, error) {
	stringQuery := fmt.Sprintf("/.*%s.*/", query)
	ss.logger.Info(query)
	ch := bleve.NewTermQuery(channel)
	mq := bleve.NewMatchPhraseQuery(query)
	rq := bleve.NewRegexpQuery(query)
	qsq := bleve.NewQueryStringQuery(stringQuery)
	q := bleve.NewDisjunctionQuery([]bleve.Query{ch, mq, rq, qsq})
	search := bleve.NewSearchRequest(q)
	search.Fields = []string{"username", "message", "channel", "timestamp"}
	return ss.index.Search(search)
}
