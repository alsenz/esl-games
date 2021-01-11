// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"flag"
	"github.com/alsenz/esl-games/pkg/lesson"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"html/template"
	"log"
	"net/http"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options

//TODO use https://github.com/pascaldekloe/jwt
//TODO OR use oauth2 library
//TODO also github.com/coreos/go-oidc/v3/oidc
//TODO use https://github.com/unrolled/secure

//TODO write a little mock http server that just returns a Set-Cookie with some grants in return for setting whatever you want
//TODO will take 10 seconds and will nicely

func (lesson *Lesson) play(w http.ResponseWriter, r *http.Request) {

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		//TODO migrate to zap log
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	//TODO we wanna generate a client token on register
	if NewConsoleAcceptMessage().write(c) != nil {

	}
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func watch(w http.ResponseWriter, r *http.Request) {

}

func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/echo")
}

func main() {
	//TODO let's get a bunch of flags together here... these are as important
	var lessonCode, authDomain, authProvider, authClientID string
	//TODO we YET need an auth provider here as an argument
	var requireAuth bool
	flag.StringVar(&lessonCode, "lesson-code", "", "The lesson code for access (default empty string = no lesson code required)")
	flag.StringVar(&authDomain, "auth-domain", "", "If require-auth is set, only accept ")
	flag.StringVar(&authProvider, "auth-provider", "", "Oidc auth provider")
	flag.StringVar(&authClientID, "auth-client-id", "", "Oidc oauth2 client ID")
	//TODO other params finishing off https://github.com/coreos/go-oidc/blob/v3/example/idtoken/app.go
	flag.BoolVar(&requireAuth)
	flag.Parse()

	provider, err := oidc.NewProvider(context.Background(), "https://accounts.google.com")
	if err != nil {
		zap.L().Fatal(err.Error())
	}
	oidcConfig := &oidc.Config{
		ClientID: authClientID,
	}
	verifier := provider.Verifier(oidcConfig)
	//TODO most of this stuff needs to get set in the handle...
	oauth2TheatreConfig := oauth2.Config{
		ClientID:     authClientID,
		ClientSecret: clientSecret, //TODO add this as a flag
		Endpoint:     provider.Endpoint(),
		RedirectURL:  "TOOD WE NEED TO KNOW OUR OWN REDIRECT! Which might be different? May be theatre or otherwise?",
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}
	oauth2ConsoleConfig := oauth2.Config{
		ClientID:     authClientID,
		ClientSecret: clientSecret, //TODO add this as a flag
		Endpoint:     provider.Endpoint(),
		RedirectURL:  "TOOD WE NEED TO KNOW OUR OWN REDIRECT! Which might be different? May be theatre or otherwise?",
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	//TODO move this through to Lesson and then pick up from there.


	//TODO plan needs to be asynchronously set
	//
	register := lesson.Register{
		Timeout:      0,
		Done:         false,
		RequireLogin: false,
		OptDomain:    nil,
		LessonCode:   lessonCode,
	}
	//TODO
	lessonSrv := lesson.NewLesson(register)

	//TODO secure wrapper round these -- TODO let's add the flags now and copy the pattern.
	http.HandleFunc("/play", lesson.PlayHanlder)
	http.HandleFunc("/watch", lesson.WatchHandler)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

//TODO delete this once we've used it.
var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>  
window.addEventListener("load", function(evt) {

    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;

    var print = function(message) {
        var d = document.createElement("div");
        d.textContent = message;
        output.appendChild(d);
    };

    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
        }
        ws = new WebSocket("{{.}}");
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
            print("RESPONSE: " + evt.data);
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };

    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input.value);
        ws.send(input.value);
        return false;
    };

    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };

});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<p>Click "Open" to create a connection to the server, 
"Send" to send a message to the server and "Close" to close the connection. 
You can change the message and send multiple times. #
<p>
<form>
<button id="open">Open</button>
<button id="close">Close</button>
<p><input id="input" type="text" value="Hello world!">
<button id="send">Send</button>
</form>
</td><td valign="top" width="50%">
<div id="output"></div>
</td></tr></table>
</body>
</html>
`))
