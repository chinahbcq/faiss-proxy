package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"testing"
	"time"

	"encoding/binary"
	"encoding/json"

	"log"

	gw "github.com/chinahbcq/faiss-proxy/proto"
	"github.com/stretchr/testify/assert"
)

const (
	Timeout = 2
	Host    = "http://0.0.0.0:3839"
	DBName  = "testDB"
	Dim     = 128
)

func httpPostWithTimeout(url string, timeout int, data []byte) ([]byte, bool) {
	client := &http.Client{
		Timeout: time.Second * time.Duration(timeout),
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		log.Fatal(err.Error())
		return nil, false
	}
	req.Close = true
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return []byte(err.Error()), false
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	return body, true
}
func Test_Faiss(t *testing.T) {
	req := gw.PingRequest{}
	req.Payload = "ping"

	js, _ := json.Marshal(req)

	url := Host + "/faiss/1.0/ping"
	body, ok := httpPostWithTimeout(url, Timeout, js)

	assert.Equal(t, true, ok)
	if !ok {
		return
	}
	var res gw.PingResponse
	err := json.Unmarshal(body, &res)
	assert.Equal(t, nil, err)
	assert.Equal(t, "pong", res.Payload)

	//dbNew
	dbNewReq := gw.DbNewRequest{}
	dbNewReq.DbName = DBName
	dbNewReq.MaxSize = 100000
	dbNewReq.Model = "index_trained.faissindex"
	dbNewReq.RequestId = fmt.Sprintf("%d", time.Now().Unix())
	js, _ = json.Marshal(dbNewReq)
	url = Host + "/faiss/1.0/db/new"
	body, ok = httpPostWithTimeout(url, Timeout, js)
	assert.Equal(t, true, ok)
	log.Println(string(body))
	var dbRes gw.EmptyResponse
	err = json.Unmarshal(body, &dbRes)
	assert.Equal(t, nil, err)
	assert.Equal(t, dbNewReq.RequestId, dbRes.RequestId)
	assert.Equal(t, int64(0), dbRes.ErrorCode)

	//dbList
	dbListReq := gw.DbListRequest{}
	dbListReq.RequestId = fmt.Sprintf("%d", time.Now().Unix())
	js, _ = json.Marshal(dbListReq)
	url = Host + "/faiss/1.0/db/list"
	body, ok = httpPostWithTimeout(url, Timeout, js)
	assert.Equal(t, true, ok)
	log.Println(string(body))
	var dbListRes gw.DbListResponse
	err = json.Unmarshal(body, &dbListRes)
	assert.Equal(t, nil, err)
	assert.Equal(t, dbListReq.RequestId, dbListRes.RequestId)
	assert.Equal(t, int64(0), dbListRes.ErrorCode)
	assert.Equal(t, true, len(dbListRes.DbStatus) > 0)

	//HSet
	HSetReq := gw.HSetRequest{}
	HSetReq.RequestId = fmt.Sprintf("%d", time.Now().Unix())
	HSetReq.DbName = DBName
	HSetReq.Feature = genFeature(Dim)
	js, _ = json.Marshal(HSetReq)
	url = Host + "/faiss/1.0/hset"
	body, ok = httpPostWithTimeout(url, Timeout, js)
	assert.Equal(t, true, ok)
	log.Println(string(body))
	var HSetRes gw.HSetResponse
	err = json.Unmarshal(body, &HSetRes)
	assert.Equal(t, nil, err)
	assert.Equal(t, HSetReq.RequestId, HSetRes.RequestId)
	assert.Equal(t, int64(0), HSetRes.ErrorCode)
	assert.Equal(t, uint64(1), HSetRes.Id)

	//HGet
	HGetReq := gw.HGetDelRequest{}
	HGetReq.RequestId = fmt.Sprintf("%d", time.Now().Unix())
	HGetReq.DbName = DBName
	HGetReq.Id = HSetRes.Id
	js, _ = json.Marshal(HGetReq)
	url = Host + "/faiss/1.0/hget"
	body, ok = httpPostWithTimeout(url, Timeout, js)
	assert.Equal(t, true, ok)
	log.Println(string(body))
	var HGetRes gw.HGetResponse
	err = json.Unmarshal(body, &HGetRes)
	assert.Equal(t, nil, err)
	assert.Equal(t, HGetReq.RequestId, HGetRes.RequestId)
	assert.Equal(t, int64(0), HGetRes.ErrorCode)
	assert.Equal(t, uint64(Dim), HGetRes.Dimension)
	assert.Equal(t, HSetReq.Feature, HGetRes.Feature)

	//HSearch
	HSearchReq := gw.HSearchRequest{}
	HSearchReq.RequestId = fmt.Sprintf("%d", time.Now().Unix())
	HSearchReq.DbName = DBName
	HSearchReq.Feature = genFeature(Dim)
	js, _ = json.Marshal(HSetReq)
	url = Host + "/faiss/1.0/hsearch"
	body, ok = httpPostWithTimeout(url, Timeout, js)
	assert.Equal(t, true, ok)
	log.Println(string(body))
	var HSearchRes gw.HSearchResponse
	err = json.Unmarshal(body, &HSearchRes)
	assert.Equal(t, nil, err)
	assert.Equal(t, HSearchReq.RequestId, HSearchRes.RequestId)
	assert.Equal(t, int64(0), HSearchRes.ErrorCode)
	assert.Equal(t, true, len(HSearchRes.Results) > 0)

	//HDel
	HDelReq := gw.HGetDelRequest{}
	HDelReq.RequestId = fmt.Sprintf("%d", time.Now().Unix())
	HDelReq.DbName = DBName
	HDelReq.Id = HSetRes.Id
	js, _ = json.Marshal(HDelReq)
	url = Host + "/faiss/1.0/hdel"
	body, ok = httpPostWithTimeout(url, Timeout, js)
	assert.Equal(t, true, ok)
	log.Println(string(body))
	var HDelRes gw.EmptyResponse
	err = json.Unmarshal(body, &HDelRes)
	assert.Equal(t, nil, err)
	assert.Equal(t, HDelReq.RequestId, HDelRes.RequestId)
	assert.Equal(t, int64(0), HDelRes.ErrorCode)

	//dbDel
	dbDelReq := gw.DbDelRequest{}
	dbDelReq.RequestId = fmt.Sprintf("%d", time.Now().Unix())
	dbDelReq.DbName = DBName
	js, _ = json.Marshal(dbDelReq)
	url = Host + "/faiss/1.0/db/del"
	body, ok = httpPostWithTimeout(url, Timeout, js)
	assert.Equal(t, true, ok)
	var dbDelRes gw.EmptyResponse
	err = json.Unmarshal(body, &dbDelRes)
	assert.Equal(t, nil, err)
	assert.Equal(t, dbDelRes.RequestId, dbDelRes.RequestId)
	assert.Equal(t, int64(0), dbDelRes.ErrorCode)
}

func genFeature(dim int) []byte {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var f = make([]float32, dim)
	for i := 0; i < dim; i++ {
		f[i] = r.Float32()
	}
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, f)

	return buf.Bytes()
}
