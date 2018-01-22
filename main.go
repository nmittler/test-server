// Copyright 2017 Zack Butcher.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

const defaultPort uint16 = 9000

func main() {
	var (
		servingPort, healthCheckPort uint16
		healthy                      bool
	)

	root := &cobra.Command{
		Use:   "server",
		Short: "Starts Mixer as a server",
		Run: func(cmd *cobra.Command, args []string) {
			if servingPort == healthCheckPort {
				http.HandleFunc("/echo", echo)
				http.HandleFunc("/health", health(&healthy))
				err := http.ListenAndServe(toAddress(servingPort), nil)
				fmt.Printf("%v\n", err)
				return
			}

			healths := http.NewServeMux()
			healths.HandleFunc("/health", health(&healthy))
			go func() {
				err := http.ListenAndServe(toAddress(healthCheckPort), healths)
				fmt.Printf("%v\n", err)
			}()
			echos := http.NewServeMux()
			echos.HandleFunc("/echo", echo)
			err := http.ListenAndServe(toAddress(servingPort), echos)
			fmt.Printf("%v\n", err)
		},
	}

	root.PersistentFlags().Uint16VarP(&servingPort, "server-port", "s", defaultPort, "Main port to serve on; always on /echo")
	root.PersistentFlags().Uint16VarP(&healthCheckPort, "health-port", "c", defaultPort, "Port to serve health checks on; always on /health")
	root.PersistentFlags().BoolVar(&healthy, "healthy", true, "If false, the health check will report unhealthy")

	if err := root.Execute(); err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(-1)
	}
}

func health(healthy *bool) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if *healthy {
			w.Write([]byte("healthy"))
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
		}
	}
}

func echo(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(body)
}

func toAddress(port uint16) string {
	return fmt.Sprintf(":%d", port)
}
