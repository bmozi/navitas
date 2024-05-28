package navitas

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/CloudyKit/jet/v6"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/johnobriggs/navitas/render"
	"github.com/johnobriggs/navitas/session"
	"github.com/joho/godotenv"
)

const version = "1.0.0"

// Navitas is the overall type for the Navitas package. Members that are exported in this type
// are available to any application that uses it.
type Navitas struct {
	AppName  string
	Debug    bool
	Version  string
	ErrorLog *log.Logger
	InfoLog  *log.Logger
	RootPath string
	Routes   *chi.Mux
	Render   *render.Render
	Session  *scs.SessionManager
	JetViews *jet.Set
	config   config
}

type config struct {
	port        string
	renderer    string
	cookie      cookieConfig
	sessionType string
}

// New reads the .env file, creates our application config, populates the Navitas type with settings
// based on .env values, and creates necessary folders and files if they don't exist
func (n *Navitas) New(rootPath string) error {
	pathConfig := initPaths{
		rootPath:    rootPath,
		folderNames: []string{"handlers", "migrations", "views", "data", "public", "tmp", "logs", "middleware"},
	}

	err := n.Init(pathConfig)
	if err != nil {
		return err
	}

	err = n.checkDotEnv(rootPath)
	if err != nil {
		return err
	}

	// read .env
	err = godotenv.Load(rootPath + "/.env")
	if err != nil {
		return err
	}

	// create loggers
	infoLog, errorLog := n.startLoggers()
	n.InfoLog = infoLog
	n.ErrorLog = errorLog
	n.Debug, _ = strconv.ParseBool(os.Getenv("DEBUG"))
	n.Version = version
	n.RootPath = rootPath
	n.Routes = n.routes().(*chi.Mux)

	n.config = config{
		port:     os.Getenv("PORT"),
		renderer: os.Getenv("RENDERER"),
		cookie: cookieConfig{
			name:     os.Getenv("COOKIE_NAME"),
			lifetime: os.Getenv("COOKIE_LIFETIME"),
			persist:  os.Getenv("COOKIE_PERSISTS"),
			secure:   os.Getenv("COOKIE_SECURE"),
			domain:   os.Getenv("COOKIE_DOMAIN"),
		},
		sessionType: os.Getenv("SESSION_TYPE"),
	}

	// create session
	sess := session.Session{
		CookieLifetime: n.config.cookie.lifetime,
		CookiePersist:  n.config.cookie.persist,
		CookieName:     n.config.cookie.name,
		SessionType:    n.config.sessionType,
		CookieDomain:   n.config.cookie.domain,
	}
	n.Session = sess.InitSession()

	var views = jet.NewSet(
		jet.NewOSFileSystemLoader(fmt.Sprintf("%s/views", rootPath)),
		jet.InDevelopmentMode(),
	)
	n.JetViews = views

	n.createRenderer()
	return nil
}

func (n *Navitas) Init(p initPaths) error {
	root := p.rootPath
	for _, path := range p.folderNames {
		// create folder if it doesn't exist
		err := n.CreateDirIfNotExist(root + "/" + path)
		if err != nil {
			return err
		}
	}
	return nil
}

// ListenAndServe starts the web server
func (n *Navitas) ListenAndServe() {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", os.Getenv("PORT")),
		ErrorLog:     n.ErrorLog,
		Handler:      n.Routes,
		IdleTimeout:  30 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 600 * time.Second,
	}

	n.InfoLog.Printf("Listening on port %s", os.Getenv("PORT"))
	err := srv.ListenAndServe()
	n.ErrorLog.Fatal(err)
}

func (n *Navitas) checkDotEnv(path string) error {
	err := n.CreateFileIfNotExists(fmt.Sprintf("%s/.env", path))
	if err != nil {
		return err
	}
	return nil
}

func (n *Navitas) startLoggers() (*log.Logger, *log.Logger) {
	var infoLog *log.Logger
	var errorLog *log.Logger

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	return infoLog, errorLog
}

func (n *Navitas) createRenderer() {
	myRenderer := render.Render{
		Renderer: n.config.renderer,
		RootPath: n.RootPath,
		Port:     n.config.port,
		JetViews: n.JetViews,
	}
	n.Render = &myRenderer
}
