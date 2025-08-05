// router.go
package main

import(
   "net/http"
   "github.com/asccclass/sherryserver"
)

func NewRouter(srv *SherryServer.Server, documentRoot string)(*http.ServeMux) {
   router := http.NewServeMux()

   // Static File server
   staticfileserver := SherryServer.StaticFileServer{documentRoot, "index.html"}
   staticfileserver.AddRouter(router)

   router.HandleFunc("POST /reflash/bpm", refreshScreen)
/*
   // App router
   router.HandleFunc("GET /api/notes", GetAll)

   router.Handle("/homepage", oauth.Protect(http.HandlerFunc(Home)))
   router.Handle("/upload", oauth.Protect(http.HandlerFunc(Upload)))
*/	
   return router
}
