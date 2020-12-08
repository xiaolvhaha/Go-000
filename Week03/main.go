package main

import (
	"fmt"
	"github.com/golang/sync/errgroup"
	"golang.org/x/net/context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)


func main(){

	ctx := context.Background()

	ctx,cancel := context.WithTimeout(ctx,30*time.Second)
	g, ctx := errgroup.WithContext(ctx)
	server :=&http.Server{
		Addr:    "localhost:8080",
		Handler: nil,
	}
	g.Go(func() error {
		//开启http服务
		 err := server.ListenAndServe()
		 if err != http.ErrServerClosed {
			return err
		}
		return nil
	})

	//do sth...


	c := make(chan os.Signal, 3)
	signal.Notify(c, syscall.SIGHUP, os.Kill,syscall.SIGINT)

	g.Go(func() error {
		select {
		case <-c:
			//理论上需要注册一下这里可以打印出收到了哪种信号
			fmt.Println("收到信号")
			time.Sleep(10*time.Second)
			fmt.Println("退出")
			cancel()

		case <-ctx.Done():
			fmt.Println("退出2")
			server.Shutdown(ctx)
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		fmt.Println(err)
	}



}
