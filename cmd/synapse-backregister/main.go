package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

var htmlTemplate *template.Template

func init() {
	var err error
	htmlTemplate, err = template.New("page").Parse(`
<!DOCTYPE html>
<html>
<head>
	<title>Synapse Registration</title>
</head>
<body>
<h1>Synapse Registration</h1>
{{ if .Notice }}
	<h2>{{.Notice}}</h2>
{{ end }}

	<br>
	<form method="POST">
		Username: <br>
		<input type="text" name="Username" /><br>
		Password: <br>
		<input type="password" name="Password" /><br>
		<input type="submit" name="Submit" />
	</form>

</body>
</html>
`)
	if err != nil {
		log.Fatalf("template didn't parse: %v", err)
	}
}

func main() {
	sharedSecret := os.Getenv("SYNAPSE_SECRET")
	if sharedSecret == "" {
		log.Fatal("must specify SYNAPSE_SECRET environment variable")
	}

	matrixServer := os.Getenv("SYNAPSE_SERVER")
	if matrixServer == "" {
		log.Fatal("must specify SYNAPSE_SERVER environment variable")
	}

	admin := false
	if os.Getenv("REGISTER_ADMINS") != "" {
		admin = true
	}

	http.ListenAndServe(":8000", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "POST" {
			uname := req.FormValue("Username")
			pass := req.FormValue("Password")

			if uname == "" {
				w.WriteHeader(400)
				logIfErr(htmlTemplate.Execute(w, map[string]string{"Notice": "Must enter a username"}))
				return
			}

			if len(pass) < 10 {
				w.WriteHeader(400)
				logIfErr(htmlTemplate.Execute(w, map[string]string{"Notice": "Password must be 10+ chars"}))
				return
			}

			hm := hmac.New(sha1.New, []byte(sharedSecret))

			hm.Write([]byte(uname))
			hm.Write([]byte{0})
			hm.Write([]byte(pass))
			hm.Write([]byte{0})
			if admin {
				hm.Write([]byte("admin"))
			} else {
				hm.Write([]byte("notadmin"))
			}
			hexDigest := hex.EncodeToString(hm.Sum(nil))

			synapseReqData := map[string]interface{}{
				"user":     uname,
				"password": pass,
				"mac":      hexDigest,
				"type":     "org.matrix.login.shared_secret",
				"admin":    admin,
			}

			reqJson, err := json.Marshal(synapseReqData)
			if err != nil {
				w.WriteHeader(400)
				logIfErr(htmlTemplate.Execute(w, map[string]string{"Notice": "OH NO INTERNAL JSON FAILURE!"}))
				return
			}

			serverLocation := strings.TrimRight(matrixServer, "/")
			regResp, err := http.Post(fmt.Sprintf("%s/_matrix/client/api/v1/register", serverLocation), "application/json", bytes.NewReader(reqJson))
			if err != nil {
				log.Printf("error: %v", err)
				w.WriteHeader(500)
				logIfErr(htmlTemplate.Execute(w, map[string]string{"Notice": "Error hitting registration server"}))
				return
			}
			if regResp.StatusCode > 400 {
				w.WriteHeader(regResp.StatusCode)
				logIfErr(htmlTemplate.Execute(w, map[string]string{"Notice": "Registration error :(!"}))
				return
			}

			w.WriteHeader(200)
			logIfErr(htmlTemplate.Execute(w, map[string]string{"Notice": "You're registered!"}))
			return
		} else {
			w.WriteHeader(200)
			logIfErr(htmlTemplate.Execute(w, map[string]string{}))
			return
		}
	}))
}

func logIfErr(err error) {
	if err != nil {
		log.Printf("Unexpected error: %v", err)
	}
}
