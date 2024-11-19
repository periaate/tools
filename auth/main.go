package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/periaate/blume/clog"
	"github.com/periaate/blume/fsio"
	"github.com/periaate/blume/x/hnet/auth"
)

var tar = ""

func main() {
	args := fsio.Args()
	if len(args) == 0 {
		clog.Fatal("not enough args")
	}

	tar = args[0]
	uri, err := url.Parse(tar)
	if err != nil {
		clog.Fatal("couldn't parse url", "err", err)
	}

	man := auth.NewManager()
	clog.Info("starting proxy", "target", tar)
	proxy := httputil.NewSingleHostReverseProxy(uri)
	go func() {
		http.HandleFunc("GET /gen/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "%s", man.NewLink(20, 30*time.Minute))
		})
		clog.Info("starting inward server", "addr", "http://localhost:8099")
		http.ListenAndServe("localhost:8099", nil)
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{hash}", func(w http.ResponseWriter, r *http.Request) {
		hash := r.PathValue("hash")
		fmt.Println(hash)
		fmt.Println(len(hash))
		if len(hash) != 44 {
			clog.Error("invalid hash", "hash", hash)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		fmt.Println(hash)
		_, err := man.UseLink(hash, w)
		if err != nil {
			clog.Error("error using link", "err", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	})

	clog.Info("starting outwards server", "addr", "http://"+args[1])
	http.ListenAndServe(args[1], http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get cookie key
		sessKey, err := r.Cookie("X-Session")
		if err != nil {
			clog.Error("no session cookie", "err", err)
			mux.ServeHTTP(w, r)
			return
		}

		man.Stringify()

		ok := man.IsValidSession(sessKey.Value)
		if !ok {
			clog.Error("session not valid")
			mux.ServeHTTP(w, r)
			return
		}

		clog.Info("proxying")
		proxy.ServeHTTP(w, r)
	}))
}
