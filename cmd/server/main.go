package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"simple-rest-api/internal/config"
	"simple-rest-api/internal/database"
	httpserver "simple-rest-api/internal/handlers"
	"simple-rest-api/internal/repository"
	"simple-rest-api/pkg/slogpretty"
	"strconv"
	"sync"
	"time"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	ctx := context.Background()

	if err := run(ctx, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, w io.Writer) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	config := config.MustLoad()
	logger := setupLogger(config.Env, w)

	db, err := database.New(config.ConnectionString)
	if err != nil {
		logger.Error("database connection failed", slog.String("error", err.Error()))
		return err
	}
	defer db.Close()

	logger.Info("database connected successfully")

	studentRepository := repository.NewStudentRepository(db)
	teacherRepository := repository.NewTeacherRepository(db)
	courseRepository := repository.NewCourseRepository(db)
	enrollmentRepository := repository.NewEnrollmentRepository(db)

	logger.Info("repositories initialized")

	server := httpserver.NewServer(
		logger,
		studentRepository,
		teacherRepository,
		courseRepository,
		enrollmentRepository,
	)
	httpServer := &http.Server{
		Addr:         net.JoinHostPort(config.Http.Host, strconv.Itoa(config.Http.Port)),
		ReadTimeout:  config.Http.Timeout,
		WriteTimeout: config.Http.Timeout,
		Handler:      server,
	}

	go func() {
		logger.Info(
			"Server listening",
			slog.String("host", config.Http.Host),
			slog.Int("port", config.Http.Port),
		)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("error listening and serving", slog.Any("error", err))
		}
	}()

	var wg sync.WaitGroup
	wg.Go(func() {
		<-ctx.Done()
		shutdownCtx := context.Background()
		shutdownCtx, cancel := context.WithTimeout(shutdownCtx, 10*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			logger.Error("error shutting down http server", slog.Any("error", err))
		}
	})
	wg.Wait()

	return nil
}

func setupLogger(env string, w io.Writer) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(w, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(w, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
