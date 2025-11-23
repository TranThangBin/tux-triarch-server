package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

var (
	ORION_URL = os.Getenv("ORION_URL")
)

func main() {
	app := pocketbase.New()

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		se.Router.Any("/orion/{path...}", func(re *core.RequestEvent) error {
			path := re.Request.PathValue("path")
			if !strings.HasPrefix(path, "/") {
				path = "/" + path
			}

			orionEndpoint, err := url.Parse(ORION_URL + path)
			if err != nil {
				return re.JSON(500, map[string]string{"error": "invalid ORION_URL"})
			}

			orionEndpoint.RawQuery = re.Request.URL.RawQuery

			body := re.Request.Body
			if re.Request.Method == "GET" {
				body = nil
			}

			req, err := http.NewRequest(
				re.Request.Method,
				orionEndpoint.String(),
				body,
			)
			if err != nil {
				return re.JSON(500, map[string]string{"error": err.Error()})
			}

			for k, vals := range re.Request.Header {
				if k == "Authorization" {
					continue
				}
				for _, v := range vals {
					req.Header.Add(k, v)
				}
			}

			client := &http.Client{}

			resp, err := client.Do(req)
			if err != nil {
				return re.JSON(502, map[string]string{"error": err.Error()})
			}
			defer resp.Body.Close()

			for k, vals := range resp.Header {
				for _, v := range vals {
					re.Response.Header().Add(k, v)
				}
			}

			re.Response.WriteHeader(resp.StatusCode)

			_, err = io.Copy(re.Response, resp.Body)
			return err
		}).Bind(apis.RequireAuth())

		return se.Next()
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
