package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/periaate/blob"
	"github.com/periaate/blume/clog"
	"github.com/periaate/blume/fsio"
	"github.com/periaate/blume/hnet"
	"github.com/periaate/blume/hnet/auth"
	"github.com/periaate/blume/maps"
	"github.com/periaate/blume/str"
	"github.com/periaate/blume/util"
)

// TODO: validate link and session hosts

func main() {
	args, _ := fsio.Args()
	serveAt := args[0]

	man := auth.NewManager(func(a *auth.Manager) {
		a.Sessions.Hooks.Set = func(s string, ei maps.ExpItem[auth.Session]) (st string, exi maps.ExpItem[auth.Session], op maps.Operation) {
			b := blob.Blob(fsio.Join(ei.Value.Host, ei.Value.Cookie))
			nerr := b.Set(ei.Value.Reader(), blob.JSON)
			if nerr != nil {
				clog.Error("error setting session", "err", nerr)
				return s, ei, maps.OP_NIL
			}
			return s, ei, maps.OP_NIL
		}
	})

	del := []blob.Blob{}

	blob.SetIndex("/blume/etc/auth/")

	for b := range blob.I.Sync.Keys() {
		r, _, _ := b.Get()
		s, err := util.From[auth.Session](r)
		if err != nil {
			clog.Error("error decoding session", "err", err)
			continue
		}

		ok := man.Register(s)
		if !ok {
			clog.Error("error registering session", "session", s, "blob", b)
			del = append(del, b)
			continue
		}
		clog.Info("registering cookie", "cookie", s.Cookie, "until", time.Until(s.T).String())
	}

	for _, b := range del {
		b.Del()
	}

	go func() {
		http.HandleFunc("GET /gen/{host}/{label}/{duration}/{uses}", func(w http.ResponseWriter, r *http.Request) {
			p := hnet.PathValue(r)
			duration := p.Duration("duration", util.AtLeast(time.Minute))
			uses := p.Int("uses", util.AtLeast(1))
			host := p.String("host", util.NotZero[string]())
			label := p.String("label")

			if len(p.Nerrors) > 0 {
				for _, v := range p.Nerrors {
					clog.Error("found error parsing path values", "err", v.Error())
				}
				p.Nerrors[0].Respond(w)
				return
			}

			t := duration * time.Minute
			link, _ := man.NewLink(uses, label, host, t)

			res := fsio.Join(string(hnet.URL(host).Format()), string(link))
			fmt.Fprintf(w, "%s", res)
		})
		clog.Info("starting inward server", "addr", hnet.URL(serveAt))
		http.ListenAndServe(serveAt, nil)
	}()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		sessKey, _ := r.Cookie("X-Session")
		if sessKey == nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		_, ok := man.Sessions.Get(sessKey.Value)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("OPTIONS /", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("GET /fw-auth/{host}/{hash...}", func(w http.ResponseWriter, r *http.Request) {
		bucket := r.PathValue("host")
		hash := r.PathValue("hash")

		b := str.SplitWithAll(bucket, false, ".")
		if len(b) >= 3 {
			bucket = b[len(b)-2] + "." + b[len(b)-1]
		}
		clog.Debug("bucket", "bucket", bucket)
		if len(bucket) == 0 {
			clog.Error("invalid hash", "hash", hash)
			w.WriteHeader(http.StatusUnauthorized)
		}

		sessKey, _ := r.Cookie("X-Session")
		if sessKey != nil {
			v, ok := man.Sessions.Get(sessKey.Value)
			clog.Debug("attempting cookie", "cookie", sessKey.Value)
			if ok {
				clog.Debug("session recognized", "label", v.Label)
				w.WriteHeader(http.StatusOK)
				return
			}
		}

		clog.Debug("authenticating", "hash", hash, "host", bucket)
		if len(hash) != 44 {
			clog.Error("invalid hash", "hash", hash, "remote", r.RemoteAddr)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		clog.Debug("using link")
		_, ok := man.UseLink(hash, w)
		if !ok {
			clog.Error("error using link")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		clog.Debug("link used")

		http.Redirect(w, r, "/", http.StatusFound)
	})

	fmt.Println(os.Getenv("FW_AUTH_ADDR"))
	addr := ":10010"
	clog.Info("starting fwauth server", "addr", hnet.URL(addr).Format())
	http.ListenAndServe(":10010", hnet.CORS{}.Handler(mux))
	// http.ListenAndServe(os.Getenv("FW_AUTH_ADDR"), hnet.CORS{}.Handler(mux))
}
