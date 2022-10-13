package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ilhammhdd/go-toolkit/restkit"
)

const callTraceFileTestCORSHttpServer = "/test_cors_http_server.go"

func testCORSHandler() {
	var callTraceFunc = fmt.Sprintf("%s#TestCORS", callTraceFileTestCORSHttpServer)
	http.Handle("/cors", &restkit.MethodRouting{
		GetHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			type response struct {
				Authorization string `json:"Authorization"`
				XPingOther    string `json:"X-PINGOTHER"`
			}
			w.Header().Set("X-PINGOTHER", "PONG")
			w.WriteHeader(http.StatusOK)
			resp := response{
				Authorization: r.Header.Get("Authorization"),
				XPingOther:    r.Header.Get("X-PINGOTHER"),
			}
			respJSON, err := json.Marshal(resp)
			if err != nil {
				log.Printf("error at %s: %s", callTraceFunc, err.Error())
			}
			w.Write(respJSON)
		}),
		MethodsCORSHeaderPolicy: &restkit.MethodsCORSHeaderPolicy{
			http.MethodGet: restkit.CORSHeaderPolicy{
				AccessControlAllowOrigin:      "http://localhost:7575",
				AccessControlAllowCredentials: true,
				AccessControlAllowMethods:     restkit.NewCaseSensitiveStrings(restkit.UpperCase, http.MethodGet, http.MethodOptions, http.MethodPost),
				AccessControlAllowHeaders:     restkit.NewCaseSensitiveStrings(restkit.LowerCase, "Authorization", "X-PINGOTHER"),
			},
		},
	})

	server := http.Server{Addr: ":8686"}
	err := server.ListenAndServe()
	if err != nil {
		log.Printf("error at %s: %s", callTraceFunc, err.Error())
	}
}
