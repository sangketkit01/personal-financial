package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/sangketkit01/personal-financial/db/sqlc"
	"github.com/sangketkit01/personal-financial/token"
	"github.com/sangketkit01/personal-financial/util"
)

type Server struct {
	config util.Config
	router *gin.Engine
	store db.Store
	tokenMaker token.Maker
}

func NewServer(config util.Config, store db.Store, tokenMaker token.Maker) (*Server, error){
	server := &Server{
		config: config,
		store: store,
		tokenMaker: tokenMaker,
	}

	server.setupRoute()
	server.router.Run(":5315")

	return server, nil
}

func (server *Server) setupRoute(){
	router := gin.Default()

	router.POST("/create-user",server.createUser)
	router.POST("/login-user",server.LoginUser)

	authRoute := router.Group("/")
	authRoute.Use(server.authMiddleware(server.tokenMaker))
	authRoute.PUT("/update-password", server.UpdateUserPassword)

	authRoute.POST("/new-financial", server.AddNewFinancial)
	authRoute.GET("/my-financial", server.MyFinancial)

	financialRoute := authRoute.Group("/financial")
	financialRoute.Use(server.FinancialMiddleware())
	financialRoute.GET("/get/:id", server.GetFinancialById)
	financialRoute.PUT("/update/:id", server.UpdateFinancial)
	financialRoute.DELETE("/delete/:id", server.DeleteFinancial)

	summaryRoute := authRoute.Group("/summary")

	summaryRoute.GET("/current-month", server.SummaryCurrentMonth)
	summaryRoute.GET("/current-year", server.SummaryCurrentYear)
	summaryRoute.GET("/each-year", server.SummaryEachYear)
	
	summaryRoute.GET("/month", server.SummaryByMonthYear)
	summaryRoute.GET("/summary/month/year", server.SummaryByYear)

	summaryRoute.GET("/type/month-year", server.SummaryTypeByMonthYear)
	summaryRoute.GET("/type/year", server.SummaryTypeByYear)

	server.router = router
}
