// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
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

package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
)

var serveCfg struct {
	Peers        []string
	Data         map[string]string
	StoragePath  string
	Leader       bool
	ReadConcern  string
	WriteConcern int
}

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start acting as our test server",
	Long:  `This starts our server which can act like a dumb database`,
	Run: func(cmd *cobra.Command, args []string) {
		serveCfg.Data = map[string]string{}

		http.HandleFunc("/read", read)
		http.HandleFunc("/write", write)

		// add/remove peer.
		http.HandleFunc("/add", addPeer)
		http.HandleFunc("/remove", removePeer)

		// check status.
		http.HandleFunc("/status", status)
		http.HandleFunc("/lead", leader)

		err := http.ListenAndServe(":8080", nil) // set listen port
		if err != nil {
			glog.Fatal("ListenAndServe: ", err)
		}
		fmt.Println("Serving on port 8080")
	},
}

func leader(w http.ResponseWriter, r *http.Request) {

	writeResponse(w, "no such key", http.StatusExpectationFailed)

}

func status(w http.ResponseWriter, r *http.Request) {

}

func writeResponse(w http.ResponseWriter, resp string, respCode int) error {
	b, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	w.WriteHeader(respCode)
	w.Write(b)
	return err
}

func read(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		writeResponse(w, "invalid", http.StatusNotFound)
		return
	}

	if val, ok := serveCfg.Data[key]; ok {
		writeResponse(w, val, http.StatusOK)
		return
	}

	writeResponse(w, "no such key", http.StatusExpectationFailed)
}

func write(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	val := r.URL.Query().Get("val")

	if key == "" || val == "" {
		writeResponse(w, "invalid", http.StatusExpectationFailed)
		return
	}

	serveCfg.Data[key] = val
	writeResponse(w, "success", http.StatusOK)
	return
}

func addPeer(w http.ResponseWriter, r *http.Request) {
	writeResponse(w, "success", http.StatusOK)
}

func removePeer(w http.ResponseWriter, r *http.Request) {
	writeResponse(w, "success", http.StatusOK)
}

func init() {
	RootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringSliceVar(&serveCfg.Peers, "peer-dns", []string{}, "Domain names for all peers including myself.")
	serveCmd.Flags().StringVar(&serveCfg.StoragePath, "storage-path", "", "Storage path where we write our data")
	serveCmd.Flags().StringVar(&serveCfg.ReadConcern, "read-concern", "local", "all|majority|local")
	serveCmd.Flags().IntVar(&serveCfg.WriteConcern, "write-concern", 0, "Number of writers that need to ack")
}
