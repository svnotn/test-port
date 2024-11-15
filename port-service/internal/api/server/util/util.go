package util

import "github.com/valyala/fasthttp"

type QueryParam struct {
	Key, Value []byte
}

func GetQueryParams(args *fasthttp.Args) []QueryParam {
	queryParams := make([]QueryParam, 0, args.Len())
	args.VisitAll(func(key, value []byte) { queryParams = append(queryParams, QueryParam{Key: key, Value: value}) })
	return queryParams
}
