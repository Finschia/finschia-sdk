package cli

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/line/link/contrib/load_test/loadgenerator"
	"github.com/spf13/cobra"
)

func RunSlaveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run-slave",
		Short: "Run slave client (load generator)",
		RunE:  runSlave,
	}
	cmd.Flags().Int(FlagPort, 8000, "port of load generator")
	return cmd
}

func runSlave(cmd *cobra.Command, args []string) error {
	port, err := cmd.Flags().GetInt(FlagPort)
	if err != nil {
		return err
	}
	router := mux.NewRouter()
	lg := loadgenerator.NewLoadGenerator()

	loadgenerator.RegisterHandlers(lg, router)
	srv := &http.Server{
		Handler: router,
		Addr:    fmt.Sprintf(":%d", port),
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Panicf("listen: %s\n", err)
		}
	}()
	log.Print("Server Started")
	<-done
	log.Print("Server Stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Panicf("Server Shutdown Failed:%+v", err)
	}
	log.Print("Server Exited Properly")
	return nil
}
