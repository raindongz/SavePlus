package apis

// import (
// 	"bytes"
// 	"encoding/json"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"github.com/randongz/save_plus/token"
// 	"github.com/stretchr/testify/require"
// )

// func Test_View_My_Interest(t *testing.T) {
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
// 			path:   "/post/viewMyInterestList",
// 			function: func(server *Server) func(ctx *gin.Context) {
// 				return server.viewMyInterestList
// 			},
// 			isAuth:        true,
// 			req:           viewMyInterestReq{},
// 			setAuthHeader: setAuthHeader,
// 			checkResponse: func(t *testing.T, req *http.Request, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusOK, recorder.Code)
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
