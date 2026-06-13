package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	goredis "github.com/redis/go-redis/v9"

	"sandbox-api-gin/internal/api/controller"
	"sandbox-api-gin/internal/api/middleware"
	"sandbox-api-gin/internal/api/router"
	fxusecase "sandbox-api-gin/internal/application/usecase/fx"
	"sandbox-api-gin/internal/application/usecase/fx/bardata"
	"sandbox-api-gin/internal/application/usecase/fx/country"
	"sandbox-api-gin/internal/application/usecase/fx/economicindicator"
	"sandbox-api-gin/internal/application/usecase/fx/economicindicatordata"
	"sandbox-api-gin/internal/application/usecase/fx/summertime"
	"sandbox-api-gin/internal/application/usecase/fx/symbol"
	zigzagusecase "sandbox-api-gin/internal/application/usecase/fx/zigzag"
	userusecase "sandbox-api-gin/internal/application/usecase/user"
	"sandbox-api-gin/internal/config"
	fxservice "sandbox-api-gin/internal/domain/service/fx"
	"sandbox-api-gin/internal/infrastructure/external"

	"sandbox-api-gin/internal/infrastructure/infradb"
	infradbfx "sandbox-api-gin/internal/infrastructure/infradb/fx"
	"sandbox-api-gin/internal/infrastructure/infraredis"
	"sandbox-api-gin/internal/security"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

// run はサーバー起動のエントリポイント。
// deferを使って確実にリソースを解放できるよう、main()から分離している。
func run() error {
	// APP_ENVに応じた設定ファイルを読み込む
	// APP_ENV未設定 → .env（production）
	// APP_ENV=local  → .env.local
	// APP_ENV=docker → .env.docker
	appEnv := os.Getenv("APP_ENV")
	envFile := ".env"
	if appEnv != "" {
		envFile = ".env." + appEnv
	}
	if err := godotenv.Load(envFile); err != nil {
		slog.Info("設定ファイルが見つかりません（環境変数を使用します）", "file", envFile)
	}

	cfg := config.Load()

	if len(cfg.JWTIssuers) == 0 {
		return fmt.Errorf("JWT_ISSUER1が設定されていません")
	}

	// MySQL接続
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Asia%%2FTokyo&clientFoundRows=true",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBSchema,
	)
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("MySQL接続エラー: %w", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			slog.Error("MySQL切断エラー", "error", err)
		}
	}()

	if err := db.PingContext(context.Background()); err != nil {
		return fmt.Errorf("MySQL疎通確認エラー: %w", err)
	}
	slog.Info("MySQL接続成功", "host", cfg.DBHost, "port", cfg.DBPort)

	// Redis接続
	redisClient := goredis.NewClient(&goredis.Options{
		Addr: fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
	})
	defer func() {
		if err := redisClient.Close(); err != nil {
			slog.Error("Redis切断エラー", "error", err)
		}
	}()
	slog.Info("Redis接続設定完了", "host", cfg.RedisHost, "port", cfg.RedisPort)

	// JWTプロバイダ
	jwtProvider, err := security.NewJwtProvider(cfg.JWTIssuers[0], cfg.JWTAudiences)
	if err != nil {
		return fmt.Errorf("JWTプロバイダ初期化エラー: %w", err)
	}
	slog.Info("JWTプロバイダ初期化完了", "issuer", cfg.JWTIssuers[0])

	// リポジトリ
	sessionRepo := infraredis.NewRedisSessionRepository(redisClient, cfg.SessionTTL)
	userRepo := infradb.NewMySQLUserRepository(db)

	// FXリポジトリ・サービス
	symbolRepo := infradbfx.NewMySQLSymbolRepository(db)
	countryRepo := infradbfx.NewMySQLCountryRepository(db)
	economicIndicatorRepo := infradbfx.NewMySQLEconomicIndicatorRepository(db)
	economicIndicatorDataRepo := infradbfx.NewMySQLEconomicIndicatorDataRepository(db)
	summerTimeRepo := infradbfx.NewMySQLSummerTimeRepository(db)
	barDataRepo := infradbfx.NewMySQLBarDataRepository(db)
	gaitameService := external.NewGaitameRateService(symbolRepo, cfg.FxRateURL, redisClient)
	tradeSimulationRepo := infradbfx.NewMySQLTradeSimulationRepository(db, redisClient, gaitameService)
	calculator := fxservice.NewFxTradeCalculator()

	// ユースケース
	loginUseCase := userusecase.NewLoginUseCase(userRepo, sessionRepo)
	logoutUseCase := userusecase.NewLogoutUseCase(sessionRepo)
	getProfileUseCase := userusecase.NewGetProfileUseCase(userRepo)
	registerUserUseCase := userusecase.NewRegisterUserUseCase(userRepo)
	updateUserUseCase := userusecase.NewUpdateUserUseCase(userRepo)
	tradeSimulationUseCase := fxusecase.NewTradeSimulationUseCase(tradeSimulationRepo, calculator)
	getMasterUseCase := fxusecase.NewGetMasterUseCase(symbolRepo, countryRepo, economicIndicatorRepo)
	searchSymbolUseCase := symbol.NewSearchSymbolUseCase(symbolRepo)
	addSymbolUseCase := symbol.NewAddSymbolUseCase(symbolRepo)
	getSymbolUseCase := symbol.NewGetSymbolUseCase(symbolRepo)
	updateSymbolUseCase := symbol.NewUpdateSymbolUseCase(symbolRepo)
	searchCountryUseCase := country.NewSearchCountryUseCase(countryRepo)
	addCountryUseCase := country.NewAddCountryUseCase(countryRepo)
	getCountryUseCase := country.NewGetCountryUseCase(countryRepo)
	updateCountryUseCase := country.NewUpdateCountryUseCase(countryRepo)
	searchSummerTimeUseCase := summertime.NewSearchSummerTimeUseCase(summerTimeRepo)
	addSummerTimeUseCase := summertime.NewAddSummerTimeUseCase(summerTimeRepo)
	getSummerTimeUseCase := summertime.NewGetSummerTimeUseCase(summerTimeRepo)
	updateSummerTimeUseCase := summertime.NewUpdateSummerTimeUseCase(summerTimeRepo)
	searchBarDataUseCase := bardata.NewSearchBarDataUseCase(barDataRepo)
	statusBarDataUseCase := bardata.NewStatusBarDataUseCase(barDataRepo)
	searchEconomicIndicatorUseCase := economicindicator.NewSearchEconomicIndicatorUseCase(economicIndicatorRepo)
	getEconomicIndicatorUseCase := economicindicator.NewGetEconomicIndicatorUseCase(economicIndicatorRepo)
	addEconomicIndicatorUseCase := economicindicator.NewAddEconomicIndicatorUseCase(economicIndicatorRepo)
	updateEconomicIndicatorUseCase := economicindicator.NewUpdateEconomicIndicatorUseCase(economicIndicatorRepo)
	searchEconomicIndicatorDataUseCase := economicindicatordata.NewSearchEconomicIndicatorDataUseCase(economicIndicatorDataRepo)
	getEconomicIndicatorDataUseCase := economicindicatordata.NewGetEconomicIndicatorDataUseCase(economicIndicatorDataRepo)
	addEconomicIndicatorDataUseCase := economicindicatordata.NewAddEconomicIndicatorDataUseCase(economicIndicatorDataRepo)
	updateEconomicIndicatorDataUseCase := economicindicatordata.NewUpdateEconomicIndicatorDataUseCase(economicIndicatorDataRepo)
	importEconomicIndicatorDataUseCase := economicindicatordata.NewImportEconomicIndicatorDataUseCase(
		economicIndicatorDataRepo, economicIndicatorRepo, countryRepo,
		cfg.StorageBucket, cfg.StorageFX, cfg.IndicatorExcludeList,
	)

	// ZigZag
	zigZagRepo := infradbfx.NewMySQLZigZagRepository(db)
	zigZagDomainService := fxservice.NewZigZagDomainService()
	searchZigZagUseCase := zigzagusecase.NewSearchZigZagUseCase(zigZagRepo)
	getZigZagStatusUseCase := zigzagusecase.NewGetZigZagStatusUseCase(zigZagRepo)
	generateZigZagUseCase := zigzagusecase.NewGenerateZigZagUseCase(zigZagRepo, zigZagDomainService)
	getZigZagBarDataUseCase := zigzagusecase.NewGetZigZagBarDataUseCase(zigZagRepo)

	// コントローラ
	authController := controller.NewAuthController(loginUseCase, logoutUseCase)
	userController := controller.NewUserController(getProfileUseCase, registerUserUseCase, updateUserUseCase)
	tradeSimulationController := controller.NewTradeSimulationController(tradeSimulationUseCase)
	masterListController := controller.NewMasterListController(getMasterUseCase)
	symbolController := controller.NewSymbolController(searchSymbolUseCase, addSymbolUseCase, getSymbolUseCase, updateSymbolUseCase)
	countryController := controller.NewCountryController(searchCountryUseCase, addCountryUseCase, getCountryUseCase, updateCountryUseCase)
	summerTimeController := controller.NewSummerTimeController(searchSummerTimeUseCase, addSummerTimeUseCase, getSummerTimeUseCase, updateSummerTimeUseCase)
	barDataController := controller.NewBarDataController(searchBarDataUseCase, statusBarDataUseCase)
	economicIndicatorController := controller.NewEconomicIndicatorController(
		searchEconomicIndicatorUseCase, getEconomicIndicatorUseCase,
		addEconomicIndicatorUseCase, updateEconomicIndicatorUseCase,
	)
	economicIndicatorDataController := controller.NewEconomicIndicatorDataController(
		searchEconomicIndicatorDataUseCase, getEconomicIndicatorDataUseCase,
		addEconomicIndicatorDataUseCase, updateEconomicIndicatorDataUseCase,
		importEconomicIndicatorDataUseCase,
	)
	zigZagController := controller.NewZigZagController(
		searchZigZagUseCase, getZigZagStatusUseCase, generateZigZagUseCase, getZigZagBarDataUseCase,
	)

	// ミドルウェア
	jwtMw := middleware.JwtMiddleware(jwtProvider, sessionRepo)
	authMw := middleware.AuthMiddleware()

	// Ginエンジン
	gin.SetMode(cfg.GINMode)
	engine := gin.Default()

	// CORS設定（JWT_ORIGIN1/2で許可オリジンを指定）
	if len(cfg.JWTOrigins) > 0 {
		engine.Use(cors.New(cors.Config{
			AllowOrigins:     cfg.JWTOrigins,
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}))
		slog.Info("CORS設定完了", "origins", cfg.JWTOrigins)
	}

	// ルーター設定
	router.Setup(engine, jwtMw, authMw,
		authController, userController, tradeSimulationController,
		masterListController, symbolController,
		countryController, summerTimeController, barDataController,
		economicIndicatorController, economicIndicatorDataController,
		zigZagController,
	)

	// Graceful Shutdown
	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: engine,
	}

	serverErrCh := make(chan error, 1)
	go func() {
		slog.Info("サーバー起動", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrCh <- err
		}
		close(serverErrCh)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErrCh:
		return fmt.Errorf("サーバーエラー: %w", err)
	case sig := <-quit:
		slog.Info("シャットダウン開始", "signal", sig)
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("graceful shutdownエラー: %w", err)
	}

	slog.Info("シャットダウン完了")
	return nil
}
