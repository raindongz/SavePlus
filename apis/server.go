package apis

import (
	"fmt"

	"github.com/gin-gonic/gin"
	db "github.com/randongz/save_plus/db/sqlc"
	"github.com/randongz/save_plus/token"
	"github.com/randongz/save_plus/utils"
)

// Server serves all HTTP requests for our SavePlus service
type Server struct{
	config utils.Config
	store db.Store
	router *gin.Engine
	tokenMaker token.Maker
}

// NewServer create a new HTTP server and setup routing.
func NewServer(config utils.Config, store db.Store) (*Server, error){
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config: config,
		store:  store,
		tokenMaker: tokenMaker,
	}

	//register custom validater to gin
	// if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
	// 	v.RegisterValidation("currency", validateCurrency)
	// }

	server.setUpRouter()

	//add routes to router 
	return server, nil
}

func (server *Server) setUpRouter(){
	router := gin.Default()

	// below routes don't need authentication

	// User related operations(no need for authentication)

	// Post related operations(no need for authentication)
		//router.POST("/user/login", server.loginUser)

	// below routes need authentication
	//authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	router.POST("/post", server.createNewPost)
	router.GET("/post/list", server.getPostList)
	router.GET("/post/listInfoNoAuth", server.getPostDetailInfoWithOutAuth)
	router.GET("/post/listInfoAuth", server.getPostDetailInfoWithAuth)
	router.POST("/post/update", server.updatePostInfo)
	router.POST("/post/delete", server.deletePostInfo)
	router.GET("/post/getInterestList", server.GetInterestList)
	router.POST("/post/interest", server.InterestPost)
	router.POST("/post/unInterest", server.UnInterestPost)
	server.router = router
}

//start runs http server on specified address 
func (server *Server) Start(address string) error{
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}