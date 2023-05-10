package elasticsearch

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/davveo/go-toolkit/logger"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/elastic/go-elasticsearch/v8/esutil"
	"net"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	esClient                   *elasticsearch.Client
	bulkIndexer                esutil.BulkIndexer
	once                       sync.Once
	defaultMaxIdleConnsPerHost = 10
)

func InitElasticSearch(conf *Config) (err error) {
	if !conf.isUseElasticSearch {
		return nil
	}

	cfg := elasticsearch.Config{
		Addresses: conf.address,
		Username:  conf.username,
		Password:  conf.password,
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   defaultMaxIdleConnsPerHost,
			ResponseHeaderTimeout: time.Second,
			DialContext:           (&net.Dialer{Timeout: time.Second}).DialContext,
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
			},
		},
	}

	once.Do(func() {
		esClient, err = elasticsearch.NewClient(cfg)
		if err != nil {
			logger.FatalErr("fail create es client err", err)
			return
		}
	})

	res, infoErr := esClient.Info()
	if err != nil {
		logger.FatalErr("err get info from es", err)
		return infoErr
	}

	bulkConfig := esutil.BulkIndexerConfig{
		Client:        esClient,         // The Elasticsearch client
		NumWorkers:    runtime.NumCPU(), // The number of worker goroutines
		FlushBytes:    int(5e+6),        // The flush threshold in bytes
		FlushInterval: 30 * time.Second, // The periodic flush interval
	}
	once.Do(func() {
		bulkIndexer, err = esutil.NewBulkIndexer(bulkConfig)
		if err != nil {
			logger.FatalErr("Error creating BulkIndexer", err)
			return
		}
	})

	logger.InfoKV("es启动成功", logger.KV("info", res))

	return nil
}

func IndexLists() (interface{}, error) {
	req := esapi.CatIndicesRequest{}
	res, err := req.Do(context.Background(), esClient)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer res.Body.Close() // res.Body  = io.ReadCloser

	buf := new(bytes.Buffer) //new(Type)作用是为T类型分配并清零一块内存，并将这块内存地址作为结果返回
	_, _ = buf.ReadFrom(res.Body)
	return strings.Split(buf.String(), "\n"), nil
}

func IndexExist(index string) bool {
	req := esapi.IndicesExistsRequest{
		Index: []string{index},
	}
	res, err := req.Do(context.Background(), esClient)
	if err != nil {
		logger.Errorf("index exist err: %+v", err)
		return false
	}
	defer res.Body.Close()
	if res.IsError() {
		return false
	}
	return true
}

func IndexCreate(index string, indexInfo map[string]interface{}) (map[string]interface{}, error) {
	data, err := json.Marshal(indexInfo)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("index create fail, error "+
			"marshaling document, index:%s, error:%s", index, err))
	}
	ioReader := bytes.NewReader(data)
	req := esapi.IndicesCreateRequest{
		Index:   index,
		Body:    ioReader,
		Timeout: 30 * time.Second,
	}
	res, err := req.Do(context.Background(), esClient)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("IndexCreate err : %v", err))
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, errors.New(fmt.Sprintf("[%s] Error indexing document "+
			"Index=%v, %s", res.Status(), req.Index, res.String()))
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, errors.New(fmt.Sprintf("error parsing the response body: %v", err))
	}
	return r, nil
}

func IndexGetMapping(index string) (map[string]interface{}, error) {
	req := esapi.IndicesGetMappingRequest{
		Index: []string{index},
	}
	res, err := req.Do(context.Background(), esClient)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("req.Do : %v", err))
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, errors.New(fmt.Sprintf("[%s] Error res， "+
			"Index=%v , %s", res.Status(), req.Index, res.String()))
	}
	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, errors.New(fmt.Sprintf("error parsing the response body: %v", err))
	}
	return r, nil
}

// IndexPutMapping 文档：https://www.elastic.co/guide/en/elasticsearch/reference/master/indices-put-mapping.html
func IndexPutMapping(index string, mapping map[string]interface{}) (map[string]interface{}, error) {
	byteMapping, err := json.Marshal(mapping)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(byteMapping)
	req := esapi.IndicesPutMappingRequest{
		Index:             []string{index},
		Body:              reader,
		AllowNoIndices:    nil,
		ExpandWildcards:   "",
		IgnoreUnavailable: nil,
		MasterTimeout:     0,
		Timeout:           0,
		WriteIndexOnly:    nil,
		Pretty:            false,
		Human:             false,
		ErrorTrace:        false,
		FilterPath:        nil,
		Header:            nil,
	}
	res, err := req.Do(context.Background(), esClient)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("req.Do : %v", err))
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, errors.New(fmt.Sprintf("[%s] Error res， "+
			"Index=%v,%s", res.Status(), req.Index, res.String()))
	}
	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, errors.New(fmt.Sprintf("error parsing the response body: %v", err))
	}
	return r, nil
}

// IndexReindex 复制索引，(es不能修改直接索引名称，或修改mapping字段，需通过复制原始数据到新的索引来实现；或者添加索引别名）
func IndexReindex(sourceIndex, destIndex string) (map[string]interface{}, error) {
	data := map[string]interface{}{
		"source": map[string]string{
			"index": sourceIndex,
		},
		"dest": map[string]string{
			"index": destIndex,
		},
	}
	m, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	body := bytes.NewReader(m)
	req := esapi.ReindexRequest{
		Body: body,
	}
	res, err := req.Do(context.Background(), esClient)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("req.Do : %v", err))
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, errors.New(fmt.Sprintf("[%s] Error res,%s", res.Status(), res.String()))
	}
	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, errors.New(fmt.Sprintf("error parsing the response body: %v", err))
	}
	return r, nil
}

func IndexClose(index string) (bool, error) {
	if index == "" {
		return false, nil
	}
	req := esapi.IndicesCloseRequest{
		Index: []string{index},
	}
	res, err := req.Do(context.Background(), esClient)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return false, errors.New(fmt.Sprintf("Error response: %s", res.String()))
	}
	return true, nil
}

func IndexDelete(indexName string) (bool, error) {
	req := esapi.IndicesDeleteRequest{
		Index: []string{indexName},
	}
	res, err := req.Do(context.Background(), esClient)
	if err != nil {
		return false, errors.New(fmt.Sprintf("Index Delete Error getting response: %s", err))
	}
	defer res.Body.Close()
	if res.IsError() {
		return false, errors.New(fmt.Sprintf("Index Delete Error response: %s", res.String()))
	}
	return true, nil
}

// IndexAlias index别名操作（添加）
func IndexAlias(index, alias string) (map[string]interface{}, error) {
	items := make(map[string]map[string]string)
	items["add"] = map[string]string{
		"index": index,
		"alias": alias,
	}

	var actions []map[string]map[string]string
	actions = append(actions, items)
	data := map[string]interface{}{
		"actions": actions,
	}
	marshal, _ := json.Marshal(data)
	body := bytes.NewReader(marshal)
	req := esapi.IndicesPutAliasRequest{
		Index: []string{index},
		Body:  body,
	}
	res, err := req.Do(context.Background(), esClient)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("req.Do : %v", err))
	}
	defer res.Body.Close()
	if res.IsError() {
		return nil, errors.New(fmt.Sprintf("[%s] Error res,%s", res.Status(), res.String()))
	}
	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, errors.New(fmt.Sprintf("error parsing the response body: %v", err))
	}
	return r, nil
}

func IndexAliasLists(index string) (map[string]interface{}, error) {
	req := esapi.IndicesGetAliasRequest{
		Index: []string{index},
	}
	res, err := req.Do(context.Background(), esClient)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("req.Do : %v", err))
	}
	defer res.Body.Close()
	if res.IsError() {
		return nil, errors.New(fmt.Sprintf("[%s] Error res,%s", res.Status(), res.String()))
	}
	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, errors.New(fmt.Sprintf("error parsing the response body: %v", err))
	}
	return r, nil

}

func IndexIsClose(indexName string) (bool, error) {
	req := esapi.IndicesGetSettingsRequest{
		Index: []string{indexName},
	}
	res, err := req.Do(context.Background(), esClient)
	if err != nil {
		return false, errors.New(fmt.Sprintf("Error getting response: %s", err))
	}
	defer res.Body.Close()
	if res.IsError() {
		return false, errors.New(fmt.Sprintf("Error response: %s", res.String()))
	}
	var data map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return false, err
	}
	if info, ok := data[indexName].(map[string]interface{}); ok {
		if settings, ok := info["settings"].(map[string]interface{}); ok {
			if index, ok := settings["index"].(map[string]interface{}); ok {
				if verifiedBeforeClose, ok := index["verified_before_close"].(string); ok {
					close, err := strconv.ParseBool(verifiedBeforeClose)
					if err != nil {
						return false, err
					}
					if close {
						return true, nil
					}
				}
			}
		}
	}
	return false, nil
}

// DocumentEntity Document
type DocumentEntity struct {
	Id   string
	Data *map[string]interface{}
}

type Pager struct {
	pageNumber int
	pageSize   int
	totalCount int
	data       []interface{}
}

func (p *Pager) GetPageNumber() int {
	return p.pageNumber
}

func (p *Pager) GetPageSize() int {
	return p.pageSize
}

func (p *Pager) GetTotalCount() int {
	return p.totalCount
}

func (p *Pager) GetData() []interface{} {
	return p.data
}

func DocumentBatchSave(index string, docs []*DocumentEntity) error {
	for _, doc := range docs {
		data, err := json.Marshal(doc.Data)
		if err != nil {
			logger.Errorf("Cannot encode doc %s: %s", doc.Id, err)
			return err
		}
		err = bulkIndexer.Add(
			context.Background(),
			esutil.BulkIndexerItem{
				Index: index,
				// Action field configures the operation to perform (index, create, delete, update)
				Action: "index",
				// DocumentID is the (optional) document ID
				DocumentID: doc.Id,
				// Body is an `io.Reader` with the payload
				Body: bytes.NewReader(data),
				// OnSuccess is called for each successful operation
				OnSuccess: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem) {
					logger.Infof("batch save success! index:%s, id:%s", res.Index, res.DocumentID)
				},
				// OnFailure is called for each failed operation
				OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
					info, _ := json.Marshal(res)
					if err != nil {
						logger.Errorf("batch save has fail! info:%s ERROR: %v", info, err)
					} else {
						logger.Infof("batch save has fail! info:%s ERROR: %s: %s", info, res.Error.Type, res.Error.Reason)
					}
				},
			},
		)
		if err != nil {
			return errors.New(fmt.Sprintf("batch save has fail! index:%s Unexpected error: %v", index, err))
		}
	}
	return nil
}

func DocumentSave(index string, doc DocumentEntity) error {
	if index == "" {
		return fmt.Errorf("document save fail, index can not be empty")
	}
	if doc.Id == "" || doc.Data == nil || len(*doc.Data) == 0 {
		return fmt.Errorf("document save fail, param doc invalid. index:%s", index)
	}
	data, err := json.Marshal(doc.Data)
	if err != nil {
		return fmt.Errorf("document save fail, error marshaling document, index:%s, error:%s", index, err)
	}
	req := esapi.IndexRequest{
		Index:      index,
		DocumentID: doc.Id,
		Body:       bytes.NewReader(data),
		Timeout:    30 * time.Second,
	}
	res, err := req.Do(context.Background(), esClient)
	if err != nil {
		return fmt.Errorf("error getting response: %v", err)
	}
	defer res.Body.Close()
	if res.IsError() {
		return fmt.Errorf("[%s] Error indexing document ID=%s, Index=%v", res.Status(), req.DocumentID, req.Index)
	} else {
		// Deserialize the response into a map.
		var r map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			return fmt.Errorf("error parsing the response body: %v", err)
		} else {
			// Print the response status and indexed document version.
			return fmt.Errorf("[%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
		}
	}

}

func DocumentFind(req esapi.SearchRequest) (*Pager, error) {
	var p = Pager{
		pageNumber: 1,
		pageSize:   20,
		totalCount: 0,
		data:       nil,
	}
	res, err := req.Do(context.Background(), esClient)
	if err != nil {
		logger.Errorf("error getting response: %v", err)
		return nil, err
	}
	defer res.Body.Close()
	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			logger.Errorf("Error parsing the response body: %v", err)
			return nil, err
		} else {
			var errorType string
			var errorReason string
			if errorInfo, ok := e["error"].(map[string]interface{}); ok {
				if v, ok := errorInfo["type"].(string); ok {
					errorType = v
				}
				if v, ok := errorInfo["reason"].(string); ok {
					errorReason = v
				}
			}
			logger.Warnf("Search index fail! type:%s, reason:%s", errorType, errorReason)
		}
		return nil, errors.New("search request fail")
	}
	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		logger.Errorf("Error parsing the response body: %v", err)
		return nil, err
	}
	if hits1, ok := r["hits"].(map[string]interface{}); ok {
		if total, ok := hits1["total"].(map[string]interface{}); ok {
			if totalCount, ok := total["value"].(float64); ok {
				p.totalCount = int(totalCount)
			}
		}
		if data, ok := hits1["hits"].([]interface{}); ok {
			p.data = data
			p.pageSize = len(data)
		}
	}
	return &p, nil
}
