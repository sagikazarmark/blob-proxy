package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"syscall"

	"github.com/oklog/run"
	"github.com/spf13/pflag"
	"gocloud.dev/blob"
	_ "gocloud.dev/blob/fileblob"
	_ "gocloud.dev/blob/s3blob"
	"gocloud.dev/gcerrors"
)

// Provisioned by ldflags
// nolint: gochecknoglobals
var (
	version    string
	commitHash string
	buildDate  string
)

func main() {
	flags := pflag.NewFlagSet("blob-proxy", pflag.ExitOnError)

	addr := flags.String("addr", ":8000", "Listen address")

	_ = flags.Parse(os.Args[1:])

	log.Println("starting application version", version, fmt.Sprintf("(%s)", commitHash), "built on", buildDate)

	b, err := blob.OpenBucket(context.Background(), os.Getenv("BUCKET"))
	if err != nil {
		log.Fatal(err)
	}
	defer b.Close()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uri := r.RequestURI

		if r.Method == http.MethodPut {
			bw, err := b.NewWriter(r.Context(), uri, nil)
			if err != nil {
				log.Println(err)

				http.Error(w, "failed write data", http.StatusInternalServerError)

				return
			}

			_, err = bw.ReadFrom(r.Body)
			if err != nil {
				log.Println(err)

				http.Error(w, "failed write data", http.StatusInternalServerError)

				return
			}

			err = bw.Close()
			if err != nil {
				log.Println(err)

				http.Error(w, "failed write data", http.StatusInternalServerError)

				return
			}

			w.WriteHeader(http.StatusAccepted)
		} else if r.Method == http.MethodGet {
			br, err := b.NewReader(r.Context(), uri, nil)
			if err != nil {
				switch gcerrors.Code(err) {
				case gcerrors.NotFound:
					http.Error(w, "key not found", http.StatusNotFound)
				default:
					log.Println(err)

					http.Error(w, "failed to retrieve data", http.StatusInternalServerError)
				}

				return
			}
			defer br.Close()

			_, err = io.Copy(w, br)
			if err != nil {
				log.Println(err)
			}
		}
	})

	httpServer := &http.Server{
		Addr:    *addr,
		Handler: handler,
	}

	log.Println("listening on", *addr)

	httpLn, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatal(err)
	}

	var group run.Group

	group.Add(
		func() error { return httpServer.Serve(httpLn) },
		func(err error) { _ = httpServer.Shutdown(context.Background()) },
	)
	defer httpServer.Close()

	// Setup signal handler
	group.Add(run.SignalHandler(context.Background(), syscall.SIGINT, syscall.SIGTERM))

	err = group.Run()
	if err != nil {
		if _, ok := err.(run.SignalError); ok {
			log.Println(err)

			return
		}

		// Fatal error
		// We don't use fatal, so deferred functions can do their jobs.
		log.Println(err)
	}
}
