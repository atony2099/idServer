package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

type IdArgs struct {
	BizTag string
}

type Arith int

func (t *Arith) GetId(args IdArgs, replay *int64) error {
	var (
		bizTag string
		err    error
		id     int64
	)

	bizTag = args.BizTag
	if bizTag == "" {
		err = errors.New("need biz_tag param")
		return err
	}

	for { // 跳过ID=0, 一般业务不支持ID=0
		if id, err = GAlloc.NextId(bizTag); err != nil {
			return err
		}
		if id == 0 {
			continue
		}
		*replay = id
		return nil
	}
}

type allocResponse struct {
	Errno int    `json:"errno"`
	Msg   string `json:"msg"`
	Id    int64  `json:"id"`
}

type healthResponse struct {
	Errno int    `json:"errno"`
	Msg   string `json:"msg"`
	Left  int64  `json:"left"`
}

func handleAlloc(w http.ResponseWriter, r *http.Request) {
	var (
		resp   allocResponse = allocResponse{}
		err    error
		bytes  []byte
		bizTag string
	)

	if err = r.ParseForm(); err != nil {
		goto RESP
	}

	if bizTag = r.Form.Get("biz_tag"); bizTag == "" {
		err = errors.New("need biz_tag param")
		goto RESP
	}

	for { // 跳过ID=0, 一般业务不支持ID=0
		if resp.Id, err = GAlloc.NextId(bizTag); err != nil {
			goto RESP
		}
		if resp.Id != 0 {
			break
		}
	}

RESP:
	if err != nil {
		resp.Errno = -1
		resp.Msg = fmt.Sprintf("%v", err)
		w.WriteHeader(500)
	} else {
		resp.Msg = "success"
	}
	if bytes, err = json.Marshal(&resp); err == nil {
		w.Write(bytes)
	} else {
		w.WriteHeader(500)
	}
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	var (
		resp   healthResponse = healthResponse{}
		err    error
		bizTag string
	)

	if err = r.ParseForm(); err != nil {
		goto RESP
	}

	if bizTag = r.Form.Get("biz_tag"); bizTag == "" {
		err = errors.New("need biz_tag param")
		goto RESP
	}

	resp.Left = GAlloc.LeftCount(bizTag)
	if resp.Left == 0 {
		err = errors.New("no available id ")
		goto RESP
	}

RESP:
	if err != nil {
		resp.Errno = -1
		resp.Msg = fmt.Sprintf("%v", err)
		w.WriteHeader(500)
	} else {
		resp.Msg = "success"
	}
	if bytes, err := json.Marshal(&resp); err == nil {
		w.Write(bytes)
	} else {
		w.WriteHeader(500)
	}
}

func StartServer() error {
	// mux := http.NewServeMux()
	http.HandleFunc("/alloc", handleAlloc)
	http.HandleFunc("/health", handleHealth)

	arith := new(Arith)
	rpc.Register(arith)

	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	return http.Serve(l, nil)

	// return http.ListenAndServe("localhost:"+strconv.Itoa(GConf.HttpPort), nil)
}
