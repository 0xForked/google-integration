package server

import (
	"context"
	"database/sql"
	"errors"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/0xForked/goca/server/user"
	"github.com/0xForked/goca/web"
	"github.com/gin-gonic/gin"
)

func Run(db *sql.DB, engine *gin.Engine) {
	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	registerRouteAndModule(db, engine)
	server := &http.Server{
		Addr:              ":8000",
		Handler:           engine,
		ReadHeaderTimeout: time.Second * 60,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(
			err, http.ErrServerClosed,
		) {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	<-ctx.Done()
	stop()
	timeToHandle := 10
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(timeToHandle)*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %s\n", err)
	}
	if err := db.Close(); err != nil {
		log.Printf("Database close: %s\n", err)
	}
	log.Println("Server exiting")
}

type embeddedFile struct {
	fs.File
}

func (f *embeddedFile) Close() error {
	return nil
}

func (f *embeddedFile) Seek(offset int64, whence int) (int64, error) {
	return f.File.(io.Seeker).Seek(offset, whence)
}

func registerRouteAndModule(db *sql.DB, router *gin.Engine) {
	router.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusTemporaryRedirect, "/fe")
	})
	router.StaticFS("/fe", http.FS(web.SPAAssets()))
	router.NoRoute(func(ctx *gin.Context) {
		if !strings.Contains(ctx.FullPath(), "/fe/") {
			ctx.String(http.StatusNotFound,
				"route that you are looking for is not found")
			return
		}
		file, err := web.SPAAssets().Open("index.html")
		if err != nil {
			ctx.String(http.StatusInternalServerError,
				"failed to open spa file: ", err.Error())
			return
		}
		defer func() { _ = file.Close() }()
		fileInfo, err := file.Stat()
		if err != nil {
			ctx.String(http.StatusInternalServerError,
				"failed to get spa file info: ", err.Error())
			return
		}
		http.ServeContent(
			ctx.Writer, ctx.Request, fileInfo.Name(),
			fileInfo.ModTime(), &embeddedFile{file})
	})
	apiRG := router.Group("/api/v1")
	user.NewUserModuleProvider(apiRG, db)
}
