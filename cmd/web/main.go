package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/KurosawaAngel/kindler/internal/adapters/email"
	"github.com/KurosawaAngel/kindler/internal/application/interactors"
	"github.com/KurosawaAngel/kindler/internal/presentation/web"
	"github.com/KurosawaAngel/kindler/internal/presentation/web/static"
	"github.com/ilyakaznacheev/cleanenv"
	"gopkg.in/gomail.v2"
)

type config struct {
	Email email.Config
	Addr  string `env:"ADDR" env-default:":8012"`
}

func main() {
	var cfg config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic(err)
	}

	mux := http.NewServeMux()

	g := gomail.NewDialer(cfg.Email.Host, cfg.Email.Port, cfg.Email.Username, cfg.Email.Username)
	mailer := email.New(g)
	sendInteractor := interactors.NewSendFile(mailer)
	handler := web.NewHandler(sendInteractor)

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write(static.Index)
	})
	mux.HandleFunc("POST /send", handler.SendFile)

	server := &http.Server{
		Addr:              cfg.Addr,
		ReadHeaderTimeout: time.Second * 30,
		Handler:           mux,
	}

	runServer(server)
}

func runServer(server *http.Server) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	<-ctx.Done()

	dctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := server.Shutdown(dctx); err != nil {
		panic(err)
	}
}
