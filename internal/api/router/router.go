package router

import (
	"sandbox-api-gin/internal/api/controller"

	"github.com/gin-gonic/gin"
)

func Setup(
	engine *gin.Engine,
	jwtMiddleware gin.HandlerFunc,
	authMiddleware gin.HandlerFunc,
	authController *controller.AuthController,
	userController *controller.UserController,
	tradeSimulationController *controller.TradeSimulationController,
	masterListController *controller.MasterListController,
	symbolController *controller.SymbolController,
	countryController *controller.CountryController,
	summerTimeController *controller.SummerTimeController,
	barDataController *controller.BarDataController,
	economicIndicatorController *controller.EconomicIndicatorController,
	economicIndicatorDataController *controller.EconomicIndicatorDataController,
	zigZagController *controller.ZigZagController,
) {
	api := engine.Group("/api")

	// 認証不要エンドポイント（@PublicApi相当）
	v1Public := api.Group("/v1/fx/master-list")
	{
		v1Public.GET("/symbol/:symbolType", masterListController.Symbol)
		v1Public.GET("/country", masterListController.Country)
		v1Public.GET("/currency-pair", masterListController.CurrencyPair)
		v1Public.GET("/currency-index", masterListController.CurrencyIndex)
		v1Public.GET("/economic-indicator/:countryCode", masterListController.EconomicIndicator)
	}

	// 認証必須エンドポイント（JWT + Auth middleware）
	v1 := api.Group("/v1")
	v1.Use(jwtMiddleware)
	v1.Use(authMiddleware)
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authController.Login)
			auth.POST("/logout-api", authController.Logout)
		}
		user := v1.Group("/user")
		{
			user.GET("", userController.Profile)
			user.POST("", userController.Registration)
			user.PUT("/:userId", userController.Update)
		}
		fx := v1.Group("/fx")
		{
			fx.POST("/trade/simulation", tradeSimulationController.Simulation)

			symbol := fx.Group("/symbol")
			{
				symbol.GET("/currency-pair-list", symbolController.CurrencyPairList)
				symbol.GET("/currency-index-list", symbolController.CurrencyIndexList)
				symbol.POST("/search", symbolController.Search)
				symbol.POST("", symbolController.Add)
				symbol.GET("/:symbol", symbolController.Get)
				symbol.PUT("/:symbol", symbolController.Update)
			}

			country := fx.Group("/country")
			{
				country.POST("/search", countryController.Search)
				country.POST("", countryController.Add)
				country.GET("/:code", countryController.Get)
				country.PUT("/:code", countryController.Update)
			}

			summerTime := fx.Group("/summer-time")
			{
				summerTime.POST("/search", summerTimeController.Search)
				summerTime.POST("", summerTimeController.Add)
				summerTime.GET("/:targetYear", summerTimeController.Get)
				summerTime.PUT("/:targetYear", summerTimeController.Update)
			}

			barData := fx.Group("/bar-data")
			{
				barData.POST("", barDataController.Search)
				barData.GET("/:symbolType/:barType", barDataController.Status)
				barData.POST("/import-csv/:symbol/:barType/:skipLatest", barDataController.ImportCsv)
			}

			zigzag := fx.Group("/zigzag")
			{
				zigzag.POST("", zigZagController.Search)
				zigzag.POST("/status", zigZagController.Status)
				zigzag.POST("/generate", zigZagController.Generate)
				zigzag.POST("/bar-data", zigZagController.BarData)
			}

			economicIndicator := fx.Group("/economic-indicator")
			{
				economicIndicator.POST("/search", economicIndicatorController.Search)
				economicIndicator.POST("", economicIndicatorController.Add)
				economicIndicator.GET("/:countryCode/:id", economicIndicatorController.Get)
				economicIndicator.PUT("/:countryCode/:id", economicIndicatorController.Update)
			}

			economicIndicatorData := fx.Group("/economic-indicator-data")
			{
				economicIndicatorData.POST("/search", economicIndicatorDataController.Search)
				economicIndicatorData.POST("", economicIndicatorDataController.Add)
				economicIndicatorData.GET("/:economicIndicatorId/:publication", economicIndicatorDataController.Get)
				economicIndicatorData.PUT("/:economicIndicatorId/:publication", economicIndicatorDataController.Update)
				economicIndicatorData.POST("/import-text", economicIndicatorDataController.ImportText)
			}
		}
	}
}
