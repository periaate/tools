package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/periaate/blob"
	"github.com/periaate/blume/clog"
	"github.com/periaate/blume/fsio"
	"github.com/periaate/blume/maps"
	"github.com/periaate/blume/options"
	"github.com/periaate/blume/util"
	"github.com/periaate/blume/x/hnet"
	"github.com/periaate/blume/x/hnet/auth"
)

// TODO: validate link and session hosts

func main() {
	man := auth.NewManager(func(a *auth.Manager) {
		a.Sessions.Hooks.Set = func(s string, ei maps.ExpItem[auth.Session]) (st string, exi maps.ExpItem[auth.Session], op maps.Operation) {
			clog.Info("setting session", "session", ei.Value)
			b := blob.Blob(fsio.Join(ei.Value.Host, ei.Value.Cookie))
			clog.Info("setting session", "session", ei.Value)
			nerr := b.Set(ei.Value.Reader(), blob.JSON)
			clog.Info("setting session", "session", ei.Value)
			if nerr != nil {
				clog.Error("error setting session", "err", nerr)
				return s, ei, maps.OP_NIL
			}
			clog.Info("session set", "session", ei.Value)
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
			duration := p.Duration("duration", options.AtLeast(time.Minute))
			uses := p.Int("uses", options.AtLeast(1))
			host := p.String("host", options.NotZero[string]())
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

			res := fsio.Join(hnet.URL(host), link)
			fmt.Fprintf(w, "%s", res)
		})
		clog.Info("starting inward server", "addr", "http://localhost:10008")
		http.ListenAndServe("localhost:10008", nil)
	}()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		sessKey, _ := r.Cookie("X-Session")
		if sessKey == nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		v, ok := man.Sessions.Get(sessKey.Value)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		clog.Info("session recognized", "label", v.Label)
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("GET /fw-auth/{host}/{hash...}", func(w http.ResponseWriter, r *http.Request) {
		bucket := r.PathValue("host")
		hash := r.PathValue("hash")

		sessKey, _ := r.Cookie("X-Session")
		if sessKey != nil {
			v, ok := man.Sessions.Get(sessKey.Value)
			clog.Info("attempting cookie", "cookie", sessKey.Value)
			if ok {
				clog.Info("session recognized", "label", v.Label)
				w.WriteHeader(http.StatusOK)
				return
			}
		}

		clog.Info("authenticating", "hash", hash, "host", bucket)
		if len(hash) != 44 {
			clog.Error("invalid hash", "hash", hash)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		clog.Info("using link")
		_, ok := man.UseLink(hash, w)
		if !ok {
			clog.Error("error using link")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		clog.Info("link used")

		http.Redirect(w, r, "/", http.StatusFound)
	})

	http.ListenAndServe("localhost:10010", hnet.Log(mux))
}
