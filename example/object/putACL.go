package main

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"net/http"

	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/tencentyun/cos-go-sdk-v5/debug"
)

func logStatus(err error) {
	if err == nil {
		return
	}
	if cos.IsNotFoundError(err) {
		// WARN
		fmt.Println("WARN: Resource is not existed")
	} else if e, ok := cos.IsCOSError(err); ok {
		fmt.Printf("ERROR: Code: %v\n", e.Code)
		fmt.Printf("ERROR: Message: %v\n", e.Message)
		fmt.Printf("ERROR: Resource: %v\n", e.Resource)
		fmt.Printf("ERROR: RequestId: %v\n", e.RequestID)
		// ERROR
	} else {
		fmt.Printf("ERROR: %v\n", err)
		// ERROR
	}
}

func main() {
	u, _ := url.Parse("https://test-1253846586.cos.ap-guangzhou.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  os.Getenv("SECRETID"),
			SecretKey: os.Getenv("SECRETKEY"),
			Transport: &debug.DebugRequestTransport{
				RequestHeader:  true,
				RequestBody:    true,
				ResponseHeader: true,
				ResponseBody:   true,
			},
		},
	})

	opt := &cos.ObjectPutACLOptions{
		Header: &cos.ACLHeaderOptions{
			XCosACL: cos.ACL.Private,
		},
	}
	name := "test/hello.txt"
	_, err := c.Object.PutACL(context.Background(), name, opt)
	logStatus(err)

	// with body
	opt = &cos.ObjectPutACLOptions{
		Body: &cos.ACLXml{
			Owner: &cos.Owner{
				ID: "qcs::cam::uin/100000760461:uin/100000760461",
			},
			AccessControlList: []cos.ACLGrant{
				{
					Grantee: &cos.ACLGrantee{
						Type: "RootAccount",
						ID:   "qcs::cam::uin/100000760461:uin/100000760461",
					},

					Permission: "FULL_CONTROL",
				},
			},
		},
	}

	_, err = c.Object.PutACL(context.Background(), name, opt)
	logStatus(err)
}
