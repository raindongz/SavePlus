package apis

import (
	"fmt"

	"github.com/gin-gonic/gin"
	db "github.com/randongz/save_plus/db/sqlc"
	"github.com/randongz/save_plus/token"
	"github.com/randongz/save_plus/utils"
)

// Server serves all HTTP requests for our SavePlus service
type Server struct {
	config     utils.Config
	store      db.Store
	router     *gin.Engine
	tokenMaker token.Maker
}

// NewServer create a new HTTP server and setup routing.
func NewServer(config utils.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	server.setUpRouter()

	//add routes to router
	return server, nil
}

func (server *Server) setUpRouter() {
	router := gin.Default()
	router.Use(setTraceId())

	// post related route do not need authenticate
	router.GET("/post/infoNoAuth", server.getPostDetailInfoWithOutAuth)
	router.GET("/post/getInterestList", server.GetInterestList)
	router.GET("/post/list", server.getPostList)

	// User related operations(Authentication needed)
	userGroup := router.Group("/user")

	userGroup.Handle("GET", "/viewMyPurchaseHistory", server.viewMyPurchaseHistory)
	userGroup.Handle("POST", "/getUserInfo", server.updateUserInfo)
	userGroup.Handle("POST", "/login", server.userLogin)
	userGroup.Handle("GET", "/create", server.createUser)

	// Post related operations(Authentication needed)// below routes need authentication
	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	authRoutes.POST("/post", server.createNewPost)
	authRoutes.GET("/post/infoAuth", server.getPostDetailInfoWithAuth)
	authRoutes.POST("/post/update", server.updatePostInfo)
	authRoutes.POST("/post/delete", server.deletePostInfo)
	authRoutes.POST("/post/interest", server.InterestPost)
	server.router = router
}

// Start start runs http server on specified address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
