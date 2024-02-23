package apis

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/randongz/save_plus/db/sqlc"
	"github.com/randongz/save_plus/token"
	"github.com/randongz/save_plus/utils"
	"github.com/stretchr/testify/require"
)

func setAuthHeader(t *testing.T, req *http.Request, tokenMaker token.Maker, authorizationType string, id int64, duration time.Duration) {
	token, err := tokenMaker.CreateToken(strconv.Itoa(int(id)), duration)
	require.NoError(t, err)

	authorizationHeader := fmt.Sprintf("%s %s", authorizationType, token)
	req.Header.Set(authorizationHeaderKey, authorizationHeader)
}

func TestServer_createNewPost(t *testing.T) {
	type TestCreatePostRequest struct {
		Title        string `json:"title"`
		Content      string `json:"content"`
		TotalPrice   string `json:"total_price"`
		DeliveryType int16  `json:"delivery_type"`
		Area         string `json:"area"`
		ItemNum      int32  `json:"item_num"`
		PostStatus   int16  `json:"post_status"`
		Negotiable   int16  `json:"negotiable"`
		Images       string `json:"images"`
	}

	tests := []struct {
		name          string
		req           any
		method        string
		path          string
		isAuth        bool
		setAuthHeader func(t *testing.T, req *http.Request, token token.Maker, authorzationType string, id int64, duration time.Duration)
		checkResponse func(t *testing.T, req *http.Request, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			method: "POST",
			path:   "/post",
			isAuth: true,
			req: TestCreatePostRequest{
				Title:        "this is title",
				Content:      "this is content",
				TotalPrice:   "13",
				DeliveryType: 0,
				Area:         "Boston",
				ItemNum:      2,
				PostStatus:   0,
				Negotiable:   0,
				Images:       "www.baidu.com",
			},
			setAuthHeader: setAuthHeader,
			checkResponse: func(t *testing.T, req *http.Request, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := SetUpRouter()
			server := NewTestServer(t)
			if tt.isAuth {
				router.POST(tt.path, authMiddleware(server.tokenMaker), server.createNewPost)
			} else {
				router.POST(tt.path, server.createNewPost)
			}

			marshalledReq, err := json.Marshal(tt.req)
			require.NoError(t, err)
			req, err := http.NewRequest(tt.method, tt.path, bytes.NewBuffer(marshalledReq))
			tt.setAuthHeader(t, req, server.tokenMaker, authorizationTypeBearer, 25, server.config.AccessTokenDuration)
			require.NoError(t, err)
			rsp := httptest.NewRecorder()
			router.ServeHTTP(rsp, req)
			//call functions in test case
			tt.checkResponse(t, req, rsp)

		})
	}
}

func TestServer_getPostList(t *testing.T) {
	type fields struct {
		config     utils.Config
		store      db.Store
		router     *gin.Engine
		tokenMaker token.Maker
	}
	type args struct {
		ctx *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &Server{
				config:     tt.fields.config,
				store:      tt.fields.store,
				router:     tt.fields.router,
				tokenMaker: tt.fields.tokenMaker,
			}
			server.getPostList(tt.args.ctx)
		})
	}
}

func TestServer_checkIfUserIdExists(t *testing.T) {
	type fields struct {
		config     utils.Config
		store      db.Store
		router     *gin.Engine
		tokenMaker token.Maker
	}
	type args struct {
		ctx *gin.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &Server{
				config:     tt.fields.config,
				store:      tt.fields.store,
				router:     tt.fields.router,
				tokenMaker: tt.fields.tokenMaker,
			}
			got, err := server.checkIfUserIdExists(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Server.checkIfUserIdExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Server.checkIfUserIdExists() = %v, want %v", got, tt.want)
			}
		})
	}
}
