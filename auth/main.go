package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"time"

	"github.com/periaate/blume/clog"
	"github.com/periaate/blume/fsio"
	"github.com/periaate/blume/gen"
	"github.com/periaate/blume/x/hnet"
	"github.com/periaate/blume/x/hnet/auth"
)

func main() {
	args := fsio.Args()
	if len(args) == 0 {
		clog.Fatal("not enough args")
	}

	tar, ok := gen.GetShift(0, args)
	if !ok {
		clog.Fatal("not enough args", "reason", "no proxy target address")
	}
	out, ok := gen.GetShift(1, args)
	if !ok {
		clog.Fatal("not enough args", "reason", "no proxy host address")
	}
	link, _ := gen.GetShift(2, args)
	link = hnet.URL(link, hnet.Opt_HTTPS)

	uri, err := url.Parse(hnet.URL(tar))
	if err != nil {
		clog.Fatal("couldn't parse url", "err", err)
	}

	man := auth.NewManager()
	clog.Info("starting proxy", "target", tar)
	proxy := httputil.NewSingleHostReverseProxy(uri)
	go func() {
		http.HandleFunc("GET /gen/{time}/{uses}", func(w http.ResponseWriter, r *http.Request) {
			times := r.PathValue("time")
			uses := r.PathValue("uses")

			if len(times) == 0 || len(uses) == 0 {
				return
			}

			itime, err := strconv.Atoi(times)
			if err != nil {
				return
			}

			iuses, err := strconv.Atoi(uses)
			if err != nil {
				return
			}

			if iuses == 0 {
				return
			}

			if len(out) == 0 {
				fmt.Fprintf(w, "%s", man.NewLink(iuses, time.Duration(itime)*time.Minute))
				return
			}
			res := fsio.Join(link, man.NewLink(iuses, time.Duration(itime)*time.Minute))
			fmt.Fprintf(w, "%s", res)
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

	clog.Info("starting outwards server", "addr", "http://"+out)
	http.ListenAndServe(out, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
