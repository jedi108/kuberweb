package application

import (
	"net/http"

	"github.com/carbocation/interpose"
	gorilla_mux "github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"

	"git.betfavorit.cf/vadim.tsurkov/kuberweb/handlers"
	"git.betfavorit.cf/vadim.tsurkov/kuberweb/kub/clientKub"
	"git.betfavorit.cf/vadim.tsurkov/kuberweb/kub/kubService"
	"git.betfavorit.cf/vadim.tsurkov/kuberweb/middlewares"
	"git.betfavorit.cf/vadim.tsurkov/kuberweb/redisService"
)

const kubInsecure = true

// New is the constructor for Application struct.
func New(config *viper.Viper) (*Application, error) {
	dsn := config.Get("dsn").(string)
	kubAddr := config.Get("kubernetes_address").(string)
	kubToken := config.Get("kubernetes_token").(string)
	redis1 := config.Get("redis1").(string)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}

	cookieStoreSecret := config.Get("cookie_secret").(string)

	app := &Application{}
	app.config = config
	app.dsn = dsn
	app.db = db
	app.sessionStore = sessions.NewCookieStore([]byte(cookieStoreSecret))
	app.serviceKubernetes = kubService.InitInstance(clientKub.NewRestClient(kubAddr, kubToken, kubInsecure))
	app.serviceRedis = redisService.NewRedisClients()
	app.serviceRedis.AddRedisAndInit(redis1)

	return app, err
}

// Application is the application object that runs HTTP server.
type Application struct {
	config            *viper.Viper
	dsn               string
	db                *sqlx.DB
	sessionStore      sessions.Store
	serviceKubernetes *kubService.ServiceKubernetes
	serviceRedis      *redisService.RedisCache
}

func (app *Application) MiddlewareStruct() (*interpose.Middleware, error) {
	middle := interpose.New()
	middle.Use(middlewares.SetDB(app.db))
	middle.Use(middlewares.SetSessionStore(app.sessionStore))
	middle.Use(middlewares.SetKubernetesService(app.serviceKubernetes))
	middle.Use(middlewares.SetRedisClient(app.serviceRedis))

	middle.UseHandler(app.mux())

	go app.serviceKubernetes.Start()

	return middle, nil
}

func (app *Application) mux() *gorilla_mux.Router {
	MustLogin := middlewares.MustLogin

	router := gorilla_mux.NewRouter()

	router.Handle("/", MustLogin(http.HandlerFunc(handlers.GetHome))).Methods("GET")

	router.HandleFunc("/signup", handlers.GetSignup).Methods("GET")
	router.HandleFunc("/signup", handlers.PostSignup).Methods("POST")
	router.HandleFunc("/login", handlers.GetLogin).Methods("GET")
	router.HandleFunc("/login", handlers.PostLogin).Methods("POST")
	router.HandleFunc("/logout", handlers.GetLogout).Methods("GET")

	router.HandleFunc("/kubernetes/pods", handlers.PagePods).Methods("GET")
	router.HandleFunc("/kubernetes/deployments", handlers.PageDeployments).Methods("GET", "POST")
	router.HandleFunc("/redis/list", handlers.PageRedis).Methods("GET", "POST")

	router.Handle("/users/{id:[0-9]+}", MustLogin(http.HandlerFunc(handlers.PostPutDeleteUsersID))).Methods("POST", "PUT", "DELETE")
	router.Handle("/kubernetes/{id:[0-9]+}", MustLogin(http.HandlerFunc(handlers.PostPutDeleteUsersID))).Methods("POST", "PUT", "DELETE")

	// Path of static files must be last!
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("static")))

	return router
}
