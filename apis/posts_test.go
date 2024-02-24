package apis

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/randongz/save_plus/token"
	"github.com/stretchr/testify/require"
)

func setAuthHeader(t *testing.T, req *http.Request, tokenMaker token.Maker, authorizationType string, id int64, duration time.Duration) {
	token, err := tokenMaker.CreateToken(strconv.Itoa(int(id)), duration)
	require.NoError(t, err)

	authorizationHeader := fmt.Sprintf("%s %s", authorizationType, token)
	req.Header.Set(authorizationHeaderKey, authorizationHeader)
}

func Test_createNewPost(t *testing.T) {
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
		function      func(server *Server) func(ctx *gin.Context)
		setAuthHeader func(t *testing.T, req *http.Request, token token.Maker, authorzationType string, id int64, duration time.Duration)
		checkResponse func(t *testing.T, req *http.Request, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			method: "POST",
			path:   "/post",
			function: func(server *Server) func(ctx *gin.Context) {
				return server.createNewPost
			},
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

			marshalledReq, err := json.Marshal(tt.req)
			require.NoError(t, err)
			req, err := http.NewRequest(tt.method, tt.path, bytes.NewBuffer(marshalledReq))
			if tt.isAuth {
				router.POST(tt.path, authMiddleware(server.tokenMaker), tt.function(server))
				tt.setAuthHeader(t, req, server.tokenMaker, authorizationTypeBearer, 25, server.config.AccessTokenDuration)
			} else {
				router.POST(tt.path, tt.function(server))
			}
			require.NoError(t, err)
			rsp := httptest.NewRecorder()
			router.ServeHTTP(rsp, req)
			//call functions in test case
			tt.checkResponse(t, req, rsp)

		})
	}
}

func Test_GetPostList(t *testing.T) {
	tests := []struct {
		name          string
		req           any
		method        string
		path          string
		function      func(server *Server) func(ctx *gin.Context)
		isAuth        bool
		setAuthHeader func(t *testing.T, req *http.Request, token token.Maker, authorzationType string, id int64, duration time.Duration)
		checkResponse func(t *testing.T, req *http.Request, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			method: "GET",
			path:   "/post/list",
			function: func(server *Server) func(ctx *gin.Context) {
				return server.getPostList
			},
			isAuth:        false,
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

			// marshalledReq, err := json.Marshal(tt.req)
			// require.NoError(t, err)
			// req, err := http.NewRequest(tt.method, tt.path, bytes.NewBuffer(marshalledReq))
			queryParams := url.Values{}
			queryParams.Set("page_size", "2")
			queryParams.Set("page_num", "1")
			url := tt.path + "?" + queryParams.Encode()

			req, err := http.NewRequest(tt.method, url, nil)
			if tt.isAuth {
				router.GET(tt.path, authMiddleware(server.tokenMaker), tt.function(server))
				tt.setAuthHeader(t, req, server.tokenMaker, authorizationTypeBearer, 25, server.config.AccessTokenDuration)
			} else {
				router.GET(tt.path, tt.function(server))
			}
			require.NoError(t, err)
			rsp := httptest.NewRecorder()
			router.ServeHTTP(rsp, req)
			//call functions in test case
			tt.checkResponse(t, req, rsp)

		})
	}
}

func Test_GetPostDetailWithOutAuth(t *testing.T) {
	tests := []struct {
		name          string
		req           any
		method        string
		path          string
		function      func(server *Server) func(ctx *gin.Context)
		isAuth        bool
		setAuthHeader func(t *testing.T, req *http.Request, token token.Maker, authorzationType string, id int64, duration time.Duration)
		checkResponse func(t *testing.T, req *http.Request, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			method: "GET",
			path:   "/post/infoNoAuth",
			function: func(server *Server) func(ctx *gin.Context) {
				return server.getPostDetailInfoWithOutAuth
			},
			isAuth:        false,
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

			// marshalledReq, err := json.Marshal(tt.req)
			// require.NoError(t, err)
			// req, err := http.NewRequest(tt.method, tt.path, bytes.NewBuffer(marshalledReq))
			queryParams := url.Values{}
			queryParams.Set("id", "6")
			url := tt.path + "?" + queryParams.Encode()

			req, err := http.NewRequest(tt.method, url, nil)
			if tt.isAuth {
				router.GET(tt.path, authMiddleware(server.tokenMaker), tt.function(server))
				tt.setAuthHeader(t, req, server.tokenMaker, authorizationTypeBearer, 25, server.config.AccessTokenDuration)
			} else {
				router.GET(tt.path, tt.function(server))
			}
			require.NoError(t, err)
			rsp := httptest.NewRecorder()
			router.ServeHTTP(rsp, req)
			//call functions in test case
			tt.checkResponse(t, req, rsp)

		})
	}
}

func Test_GetPostDetailWithAuth(t *testing.T) {
	tests := []struct {
		name          string
		req           any
		method        string
		path          string
		function      func(server *Server) func(ctx *gin.Context)
		isAuth        bool
		setAuthHeader func(t *testing.T, req *http.Request, token token.Maker, authorzationType string, id int64, duration time.Duration)
		checkResponse func(t *testing.T, req *http.Request, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			method: "GET",
			path:   "/post/infoAuth",
			function: func(server *Server) func(ctx *gin.Context) {
				return server.getPostDetailInfoWithAuth
			},
			isAuth:        true,
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

			// marshalledReq, err := json.Marshal(tt.req)
			// require.NoError(t, err)
			// req, err := http.NewRequest(tt.method, tt.path, bytes.NewBuffer(marshalledReq))
			queryParams := url.Values{}
			queryParams.Set("post_id", "6")
			url := tt.path + "?" + queryParams.Encode()

			req, err := http.NewRequest(tt.method, url, nil)
			if tt.isAuth {
				router.GET(tt.path, authMiddleware(server.tokenMaker), tt.function(server))
				tt.setAuthHeader(t, req, server.tokenMaker, authorizationTypeBearer, 25, server.config.AccessTokenDuration)
			} else {
				router.GET(tt.path, tt.function(server))
			}
			require.NoError(t, err)
			rsp := httptest.NewRecorder()
			router.ServeHTTP(rsp, req)
			//call functions in test case
			tt.checkResponse(t, req, rsp)

		})
	}
}

func Test_UpdatePostInfo(t *testing.T) {
	type TestUpdatePostInfo struct {
		PostId       int64  `json:"post_id"`
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
		function      func(server *Server) func(ctx *gin.Context)
		setAuthHeader func(t *testing.T, req *http.Request, token token.Maker, authorzationType string, id int64, duration time.Duration)
		checkResponse func(t *testing.T, req *http.Request, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			method: "POST",
			path:   "/post/update",
			function: func(server *Server) func(ctx *gin.Context) {
				return server.updatePostInfo
			},
			isAuth: true,
			req: TestUpdatePostInfo{
				PostId:       6,
				Title:        "this is new title",
				Content:      "this is new content",
				TotalPrice:   "12",
				DeliveryType: 0,
				Area:         "Boston",
				ItemNum:      5,
				PostStatus:   0,
				Negotiable:   0,
				Images:       "www.baidu.com",
			},
			setAuthHeader: setAuthHeader,
			checkResponse: func(t *testing.T, req *http.Request, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:   "Unauthorized",
			method: "POST",
			path:   "/post/update",
			function: func(server *Server) func(ctx *gin.Context) {
				return server.updatePostInfo
			},
			isAuth: true,
			req: TestUpdatePostInfo{
				PostId:       2,
				Title:        "this is new title",
				Content:      "this is new content",
				TotalPrice:   "12",
				DeliveryType: 0,
				Area:         "Boston",
				ItemNum:      5,
				PostStatus:   0,
				Negotiable:   0,
				Images:       "www.baidu.com",
			},
			setAuthHeader: setAuthHeader,
			checkResponse: func(t *testing.T, req *http.Request, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := SetUpRouter()
			server := NewTestServer(t)

			marshalledReq, err := json.Marshal(tt.req)
			require.NoError(t, err)
			req, err := http.NewRequest(tt.method, tt.path, bytes.NewBuffer(marshalledReq))
			if tt.isAuth {
				router.POST(tt.path, authMiddleware(server.tokenMaker), tt.function(server))
				tt.setAuthHeader(t, req, server.tokenMaker, authorizationTypeBearer, 25, server.config.AccessTokenDuration)
			} else {
				router.POST(tt.path, tt.function(server))
			}
			require.NoError(t, err)
			rsp := httptest.NewRecorder()
			router.ServeHTTP(rsp, req)
			//call functions in test case
			tt.checkResponse(t, req, rsp)

		})
	}
}

// func Test_DeletePost(t *testing.T) {
// 	type TestDeletePostRequest struct {
// 		PostId int64 `json:"post_id"`
// 	}
// 	tests := []struct {
// 		name          string
// 		req           any
// 		method        string
// 		path          string
// 		isAuth        bool
// 		function      func(server *Server) func(ctx *gin.Context)
// 		setAuthHeader func(t *testing.T, req *http.Request, token token.Maker, authorzationType string, id int64, duration time.Duration)
// 		checkResponse func(t *testing.T, req *http.Request, recorder *httptest.ResponseRecorder)
// 	}{
// 		{
// 			name:   "OK",
// 			method: "POST",
// 			path:   "/post/delete",
// 			function: func(server *Server) func(ctx *gin.Context) {
// 				return server.deletePostInfo
// 			},
// 			isAuth: true,
// 			req: TestDeletePostRequest{
// 				PostId: 1,
// 			},
// 			setAuthHeader: setAuthHeader,
// 			checkResponse: func(t *testing.T, req *http.Request, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusOK, recorder.Code)
// 			},
// 		},
// 		{
// 			name:   "Unauthorized",
// 			method: "POST",
// 			path:   "/post/delete",
// 			function: func(server *Server) func(ctx *gin.Context) {
// 				return server.deletePostInfo
// 			},
// 			isAuth: true,
// 			req: TestDeletePostRequest{
// 				PostId: 2,
// 			},
// 			setAuthHeader: setAuthHeader,
// 			checkResponse: func(t *testing.T, req *http.Request, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusUnauthorized, recorder.Code)
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			router := SetUpRouter()
// 			server := NewTestServer(t)

// 			marshalledReq, err := json.Marshal(tt.req)
// 			require.NoError(t, err)
// 			req, err := http.NewRequest(tt.method, tt.path, bytes.NewBuffer(marshalledReq))
// 			if tt.isAuth {
// 				router.POST(tt.path, authMiddleware(server.tokenMaker), tt.function(server))
// 				tt.setAuthHeader(t, req, server.tokenMaker, authorizationTypeBearer, 25, server.config.AccessTokenDuration)
// 			} else {
// 				router.POST(tt.path, tt.function(server))
// 			}
// 			require.NoError(t, err)
// 			rsp := httptest.NewRecorder()
// 			router.ServeHTTP(rsp, req)
// 			//call functions in test case
// 			tt.checkResponse(t, req, rsp)

// 		})
// 	}
// }

func Test_GetPostInterestList(t *testing.T) {
	tests := []struct {
		name          string
		req           any
		method        string
		path          string
		function      func(server *Server) func(ctx *gin.Context)
		isAuth        bool
		setAuthHeader func(t *testing.T, req *http.Request, token token.Maker, authorzationType string, id int64, duration time.Duration)
		checkResponse func(t *testing.T, req *http.Request, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			method: "GET",
			path:   "/post/getInterestList",
			function: func(server *Server) func(ctx *gin.Context) {
				return server.GetInterestList
			},
			isAuth:        false,
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

			// marshalledReq, err := json.Marshal(tt.req)
			// require.NoError(t, err)
			// req, err := http.NewRequest(tt.method, tt.path, bytes.NewBuffer(marshalledReq))
			queryParams := url.Values{}
			queryParams.Set("post_id", "1")
			url := tt.path + "?" + queryParams.Encode()

			req, err := http.NewRequest(tt.method, url, nil)
			if tt.isAuth {
				router.GET(tt.path, authMiddleware(server.tokenMaker), tt.function(server))
				tt.setAuthHeader(t, req, server.tokenMaker, authorizationTypeBearer, 25, server.config.AccessTokenDuration)
			} else {
				router.GET(tt.path, tt.function(server))
			}
			require.NoError(t, err)
			rsp := httptest.NewRecorder()
			router.ServeHTTP(rsp, req)
			//call functions in test case
			tt.checkResponse(t, req, rsp)

		})
	}
}

func Test_Interest(t *testing.T) {
	tests := []struct {
		name          string
		req           any
		method        string
		path          string
		isAuth        bool
		function      func(server *Server) func(ctx *gin.Context)
		setAuthHeader func(t *testing.T, req *http.Request, token token.Maker, authorzationType string, id int64, duration time.Duration)
		checkResponse func(t *testing.T, req *http.Request, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			method: "POST",
			path:   "/post/interest",
			function: func(server *Server) func(ctx *gin.Context) {
				return server.InterestPost
			},
			isAuth: true,
			req: InterestPostRequest{
				PostId: 6,
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

			marshalledReq, err := json.Marshal(tt.req)
			require.NoError(t, err)
			req, err := http.NewRequest(tt.method, tt.path, bytes.NewBuffer(marshalledReq))
			if tt.isAuth {
				router.POST(tt.path, authMiddleware(server.tokenMaker), tt.function(server))
				tt.setAuthHeader(t, req, server.tokenMaker, authorizationTypeBearer, 25, server.config.AccessTokenDuration)
			} else {
				router.POST(tt.path, tt.function(server))
			}
			require.NoError(t, err)
			rsp := httptest.NewRecorder()
			router.ServeHTTP(rsp, req)
			//call functions in test case
			tt.checkResponse(t, req, rsp)

		})
	}
}
