package lib

import (
	"io/ioutil"

	"https://github.com/IminaDanzi/Docker-Project/rqctx"
)

func postToDatabaseHandler(ctx *rqctx.Context) error {
	requestBody, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		return err
	}
	return ctx.RequestBody.Insert(requestBody)
}
