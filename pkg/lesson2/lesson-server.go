package lesson2 //TODO package update

import (
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

//TODO gradually move in the best bits from lesson and tidy up as we go - TODO do these now!
//TODO do all of the auth and websocket message handling upfront. Wanna get into a good state before moving onto controller & evt loop.

type LessonServer struct {
	Opts struct {
		planIDStr string
		address string
		authClientID *string
		code string
		inactivityTimeout time.Duration
	}
	Controller Controller
	Register struct {
		pass *string
		requiredAuth bool
		passIsFallback bool
		authDomain *string
		authProvider *string
		timeout time.Duration
		complete bool
	}
	EventLoop EventLoop

	upgrader websocket.Upgrader
}

func NewLessonServer() * LessonServer {
	srv := &LessonServer{}

	// Opts flags
	flag.StringVar(&srv.Opts.planIDStr, "plan-id", "", "plan UUID to load")
	flag.StringVar(&srv.Opts.address, "address", "localhost:8080", "Address to bind on")
	flag.StringVar(srv.Opts.authClientID, "oauth-client-id", "", "Oauth2 client id")
	flag.StringVar(&srv.Opts.code, "room-code", "", "Room code")
	flag.DurationVar(&srv.Opts.inactivityTimeout, "inactivity-timeout", time.Minute * 60, "Time after which without activity game ends automatically")

	//Register flags
	flag.StringVar(srv.Register.pass, "room-pass", "", "Room passcode")
	flag.BoolVar(&srv.Register.requiredAuth, "required-auth", true, "Require authentication for access")
	flag.BoolVar(&srv.Register.passIsFallback, "pass-fallback", false, "Instead of requiring pass, use it as a fallback incase domain condition unmet.")
	flag.StringVar(srv.Register.authDomain, "auth-domain", "", "Require email @domain")
	flag.StringVar(srv.Register.authProvider, "auth-provider", "", "One of a preset list of ")
	flag.DurationVar(&srv.Register.timeout, "register-timeout", time.Second * 60, "Time for everybody to register")
	srv.Register.complete = false

	//Postgres connection flags
	pgHost := flag.String("pg-host", "localhost", "Postgresql database host")
	pgPort := flag.String("pg-port", "5432", "Postgresql database port")
	pgUser := flag.String("pg-user", "postgres", "Postgresql user")
	pgPass := flag.String("pg-password", "admin", "Postgresql user password")
	pgDb := flag.String("pg-database", "esl_games", "Postgresql database")

	flag.Parse()
	srv.validateFlags()

	// Open the database connection
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", *pgHost, *pgUser, *pgPass, *pgDb, *pgPort)
	var db *gorm.DB
	var err error
	if db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{}); err != nil {
		zap.L().Fatal("Unable to open postgres connection: " + err.Error())
	}

	// Parse the plan UUID and make various channels
	var planUUID uuid.UUID
	if planUUID, err = uuid.FromString(srv.Opts.planIDStr); err != nil {
		zap.L().Fatal("Unable to parse plan UUID: " + err.Error())
	}
	//TODO make the various channels

	//TODO
	srv.Controller = NewController(db, planUUID) //TODO extend this ctor.
	srv.EventLoop = NewEventLoop() //TODO make this

	srv.upgrader = websocket.Upgrader{}

	return srv
}

// This will go fatal if the flags aren't valid
func (srv *LessonServer) validateFlags() {
	// set any optional flags to nil
	if srv.Opts.authClientID != nil && len(*srv.Opts.authClientID) == 0 {
		srv.Opts.authClientID = nil
	}
	if srv.Register.pass != nil && len(*srv.Register.pass) == 0 {
		srv.Register.pass = nil
	}
	if srv.Register.authDomain != nil && len(*srv.Register.authDomain) == 0 {
		srv.Register.authDomain = nil
	}
	if srv.Register.authProvider != nil && len(*srv.Register.authProvider) == 0 {
		srv.Register.authProvider = nil
	}

	if srv.Register.requiredAuth {
		if srv.Opts.authClientID == nil {
			zap.L().Fatal("required-auth is set, but no oauth-client-id provided")
		}
	}

	if srv.Register.passIsFallback && srv.Register.pass == nil {
		zap.L().Fatal("pass-fallback is set but pass is empty")
	}

	if len(srv.Opts.code) == 0 {
		zap.L().Fatal("room code is empty")
	}

	if len(srv.Opts.planIDStr) == 0 {
		zap.L().Fatal("Plan ID is empty")
	}

}

func (srv *LessonServer) Run() error {
	if err := srv.Controller.LoadPlan(); err != nil {
		zap.L().Error("Unable to load plan due to " + err.Error())
		return err
	}
	go srv.Controller.Run()
	go srv.EventLoop.Run()
	return nil
}

//TODO let's move in: console, threatre (assume the other apps will handle endpoints for rest)