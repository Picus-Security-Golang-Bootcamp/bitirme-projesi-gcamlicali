package main

import (
	"fmt"
	"github.com/gcamlicali/tradeshopExample/internal/auth"
	"github.com/gcamlicali/tradeshopExample/internal/cart"
	"github.com/gcamlicali/tradeshopExample/internal/cart_item"
	"github.com/gcamlicali/tradeshopExample/internal/category"
	"github.com/gcamlicali/tradeshopExample/internal/order"
	"github.com/gcamlicali/tradeshopExample/internal/product"
	"github.com/gcamlicali/tradeshopExample/pkg/config"
	db "github.com/gcamlicali/tradeshopExample/pkg/database"
	"github.com/gcamlicali/tradeshopExample/pkg/graceful"
	logger "github.com/gcamlicali/tradeshopExample/pkg/logging"
	mw "github.com/gcamlicali/tradeshopExample/pkg/middleware"

	"github.com/gin-gonic/gin"

	"log"
	"net/http"
	"time"
)

func main() {
	log.Println("Trading cart service starting...")

	// Set envs for local development
	cfg, err := config.LoadConfig("./pkg/config/config-local")
	if err != nil {
		log.Fatalf("LoadConfig: %v", err)
	}

	// Set global logger
	logger.NewLogger(cfg)
	defer logger.Close()

	// Connect DB
	// Use golang-migrate instead of gorm auto migrate
	//https://github.com/golang-migrate/migrate
	DB := db.Connect(cfg)

	gin.SetMode(gin.ReleaseMode)

	// Init Gin and start gin engine (Recovery MW: if you don't want to panic exit, recovery returns 500 ErrorCode[read inside comments])
	r := gin.Default()

	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// your custom format
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.ServerConfig.Port),
		Handler:      r,
		ReadTimeout:  time.Duration(cfg.ServerConfig.ReadTimeoutSecs * int64(time.Second)),
		WriteTimeout: time.Duration(cfg.ServerConfig.WriteTimeoutSecs * int64(time.Second)),
	}

	// Router group
	rootRouter := r.Group(cfg.ServerConfig.RoutePrefix)
	authRooter := rootRouter.Group("/user")
	productRouter := rootRouter.Group("/product")
	categoryRouter := rootRouter.Group("/category")
	cartRouter := rootRouter.Group("/cart")
	orderRouter := rootRouter.Group("/order")

	//MW Control
	cartRouter.Use(mw.AuthMiddleware(cfg.JWTConfig.SecretKey))
	orderRouter.Use(mw.AuthMiddleware(cfg.JWTConfig.SecretKey))

	// Category Repository
	categoryRepo := category.NewCategoryRepository(DB)
	categoryRepo.Migration()
	productRepo := product.NewProductRepository(DB)
	productRepo.Migration()

	categoryService := category.NewCategoryService(categoryRepo, productRepo)
	category.NewCategoryHandler(categoryRouter, categoryService, cfg)

	//// Product Repository

	productService := product.NewProductService(productRepo, categoryRepo)
	product.NewProductHandler(productRouter, productService, cfg)

	cartItemRepo := cart_item.NewCartItemRepository(DB)
	cartItemRepo.Migration()

	cartRepo := cart.NewCartRepository(DB)
	cartRepo.Migration()
	cartService := cart.NewCartService(cartRepo, cartItemRepo, productRepo)
	cart.NewCartHandler(cartRouter, cartService)

	authRepo := auth.NewAuthRepository(DB)
	authRepo.Migration()
	authService := auth.NewAuthService(*authRepo, *cartRepo, cfg)
	authService.FillAdminData()
	auth.NewAuthHandler(authRooter, authService)

	orderRepo := order.NewOrderRepository(DB)
	orderRepo.Migration()
	orderService := order.NewOrderService(orderRepo, cartRepo, cartItemRepo, productRepo)
	order.NewOrderHandler(orderRouter, orderService)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	log.Println("Trading backend service started")
	graceful.ShutdownGin(srv, time.Duration(cfg.ServerConfig.TimeoutSecs*int64(time.Second)))
}
