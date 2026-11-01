package main

import (
    "fmt"
    "log"

    "github.com/gin-gonic/gin"

    "vyomtech-backend/internal/config"
    "vyomtech-backend/internal/db"
    "vyomtech-backend/internal/routes"
)

func main() {
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("config load failed: %v", err)
    }

    gin.SetMode(gin.DebugMode)
    r := gin.Default()

    r.Use(func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin", cfg.CORSOrigin)
        c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        c.Next()
    })

    database, err := db.Connect(cfg)
    if err != nil {
        log.Printf("database connection failed: %v", err)
        log.Println("continuing startup without DB connection for now; migrations will run when the database is available")
    } else {
        if err := db.AutoMigrate(database); err != nil {
            log.Fatalf("database migration failed: %v", err)
        }
        log.Println("database migration successful")
    }

    routes.Register(r, cfg)

    if err := r.Run(fmt.Sprintf(":%d", cfg.Port)); err != nil {
        log.Fatalf("server failed to start: %v", err)
    }
}
