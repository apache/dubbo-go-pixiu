package seata

import (
	"bytes"
	"encoding/json"
	"net/http"

	"vimagination.zapto.org/byteio"

	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

var (
	VarHost        = "host"
	VarQueryString = "queryString"
)

type RequestContext struct {
	ActionContext map[string]string
	Headers       http.Header
	Body          []byte
	Trailers      http.Header
}

func (ctx *RequestContext) Encode() ([]byte, error) {
	var (
		err               error
		actionContextData []byte
		headersData       []byte
		trailersData      []byte
		b                 bytes.Buffer
	)

	w := byteio.BigEndianWriter{Writer: &b}

	if ctx.ActionContext == nil || len(ctx.ActionContext) == 0 {
		w.WriteUint32(0)
	} else {
		actionContextData, err = json.Marshal(ctx.ActionContext)
		if err != nil {
			return nil, err
		}

		w.WriteUint32(uint32(len(actionContextData)))
		w.Write(actionContextData)
	}

	if ctx.Headers == nil || len(ctx.Headers) == 0 {
		w.WriteUint32(0)
	} else {
		headersData, err = json.Marshal(ctx.Headers)
		if err != nil {
			return nil, err
		}
		w.WriteUint32(uint32(len(headersData)))
		w.Write(headersData)
	}

	if ctx.Trailers == nil || len(ctx.Trailers) == 0 {
		w.WriteUint32(0)
	} else {
		trailersData, err = json.Marshal(ctx.Trailers)
		if err != nil {
			return nil, err
		}
		w.WriteUint32(uint32(len(trailersData)))
		w.Write(trailersData)
	}

	w.WriteUint32(uint32(len(ctx.Body)))
	b.Write(ctx.Body)

	return b.Bytes(), nil
}

func (ctx *RequestContext) Decode(b []byte) {
	var (
		length32          uint32 = 0
		err               error
		actionContextData []byte
		headersData       []byte
		bodyData          []byte
		trailersData      []byte
	)
	r := byteio.BigEndianReader{Reader: bytes.NewReader(b)}

	length32, _, _ = r.ReadUint32()
	if length32 > 0 {
		actionContextData = make([]byte, length32, length32)
		r.Read(actionContextData)
	}

	length32, _, _ = r.ReadUint32()
	if length32 > 0 {
		headersData = make([]byte, length32, length32)
		r.Read(headersData)
	}

	length32, _, _ = r.ReadUint32()
	if length32 > 0 {
		trailersData = make([]byte, length32, length32)
		r.Read(trailersData)
	}

	length32, _, _ = r.ReadUint32()
	if length32 > 0 {
		bodyData = make([]byte, length32, length32)
		r.Read(bodyData)
	}

	if actionContextData != nil {
		err = json.Unmarshal(actionContextData, &(ctx.ActionContext))
		if err != nil {
			logger.Errorf("unmarshal action context failed, %v", err)
		}
	}

	if headersData != nil {
		err = json.Unmarshal(headersData, &(ctx.Headers))
		if err != nil {
			logger.Errorf("unmarshal headers failed, %v", err)
		}
	}

	if trailersData != nil {
		err = json.Unmarshal(trailersData, &(ctx.ActionContext))
		if err != nil {
			logger.Errorf("unmarshal trailers failed, %v", err)
		}
	}

	ctx.Body = bodyData
}
