package clusterCache

import (
	"encoding/json"

	"github.com/valyala/fasthttp"
)

type ClusterApiRequest struct {
	Key     string `json:"key"`
	Method  string `json:"method"`
	Message string `json:"message"`
}

func (c *Cluster) ClusterApiHandler(ctx *fasthttp.RequestCtx) {
	data := ClusterApiRequest{}
	err := json.Unmarshal(ctx.PostBody(), &data)
	if err != nil {
		ctx.Error("Invalid json", 400)
		return
	}

	if !c.checkIfApiKeyIsCorrect(data.Key) {
		ctx.Error("Invalid api key", 401)
		return
	}

	switch data.Method {
	case "heartbeat":
		ctx.Response.SetStatusCode(200)
		ctx.Response.SetBody([]byte("OK"))
		return
	default:
		ctx.Error("Invalid method", 400)
		return
	}
}

func (c *Cluster) checkIfApiKeyIsCorrect(key string) bool {
	return key == c.Conf.ClusterKey
}

func (c *Cluster) SendHeartbeat() {

}
