package api

type CommonPaginationReq struct {
	Page int `p:"page" in:"query" d:"1"  v:"min:0#分页号码错误"     dc:"分页号码，默认1"`
	Size int `p:"size" in:"query" d:"10" v:"max:50#分页数量最大50条" dc:"分页数量，最大50"`
}
