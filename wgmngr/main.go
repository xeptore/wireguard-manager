package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	router := gin.Default()
	router.Use(cors.Default())
	store := cookie.NewStore([]byte(`|k'%#C85Y]L>tDuI=MaoqAHk+Gu_P|>(JuJ~=2XGNRBrr=/M##+_j6Ea'|HJ?&k2`), []byte(`T$W4k{_ep26|{,.Za+AK{KbbY5:f<dWR`))

	router.Use(sessions.Sessions("mysession", store))

	router.GET("/incr", func(c *gin.Context) {
		session := sessions.Default(c)
		var count int
		v := session.Get("sss")
		fmt.Printf("%+#v\n", session.ID())
		if v == nil {
			count = 0
		} else {
			count = v.(int)
			count++
		}
		session.Set("sss", count)
		session.Save()
		c.JSON(200, gin.H{"count": count})
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	log.Println("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")

}
