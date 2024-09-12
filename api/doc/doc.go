// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package doc

import (
	"context"

	"github.com/gq-leon/go-pandoc/api/doc/v1"
)

type IDocV1 interface {
	Doc(ctx context.Context, req *v1.DocReq) (res *v1.DocRes, err error)
}
