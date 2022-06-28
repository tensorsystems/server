/*
  Copyright 2021 Kidus Tiliksew

  This file is part of Tensor EMR.

  Tensor EMR is free software: you can redistribute it and/or modify
  it under the terms of the version 2 of GNU General Public License as published by
  the Free Software Foundation.

  Tensor EMR is distributed in the hope that it will be useful,
  but WITHOUT ANY WARRANTY; without even the implied warranty of
  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
  GNU General Public License for more details.

  You should have received a copy of the GNU General Public License
  along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package server

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"net/http"
	_ "net/http/pprof"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"github.com/tensoremr/server/pkg/auth"
	"github.com/tensoremr/server/pkg/conf"
	"github.com/tensoremr/server/pkg/controller"
	"github.com/tensoremr/server/pkg/graphql/graph"
	"github.com/tensoremr/server/pkg/graphql/graph/generated"
	"github.com/tensoremr/server/pkg/middleware"
	"github.com/tensoremr/server/pkg/repository"
	"gorm.io/gorm"
)

// Server ...
type Server struct {
	Gin           *gin.Engine
	Config        *conf.Configuration
	DB            *gorm.DB
	ACLEnforcer   *casbin.Enforcer
	TestDB        *gorm.DB          // Database connection
	ModelRegistry *repository.Model // Model registry for migration
}

// NewServer will create a new instance of the application
func NewServer() *Server {
	server := &Server{}

	server.ModelRegistry = repository.NewModel()
	server.NewEnforcer()

	if err := server.ModelRegistry.OpenPostgres(); err != nil {
		log.Fatalf("gorm: could not connect to db %q", err)
	}

	server.DB = server.ModelRegistry.DB

	server.ModelRegistry.RegisterAllModels()
	server.ModelRegistry.AutoMigrateAll()
	server.ModelRegistry.SeedData()
	server.RegisterJobs()

	server.Gin = server.NewRouter()

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	return server
}

// NewTestServer will create a new test instance
func NewTestServer(config *conf.Configuration) (*Server, error) {
	server := &Server{}
	server.Config = config

	server.ModelRegistry = repository.NewModel()
	err := server.ModelRegistry.OpenWithConfigTest(config)

	if err != nil {
		return nil, err
	}

	server.TestDB = server.ModelRegistry.DB

	server.ModelRegistry.RegisterAllModels()
	server.ModelRegistry.AutoMigrateAll()

	return server, nil
}

// Defining the Playground handler
func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// RegisterJobs ...
func (s *Server) RegisterJobs() {
	c := cron.New()
	c.AddFunc("@hourly", func() {
		var patientQueue repository.PatientQueue
		if err := patientQueue.ClearExpired(); err != nil {
			fmt.Println(err)
		}
	})
	c.Start()
}

// Defining the Graphql handler
func graphqlHandler(server *Server) gin.HandlerFunc {
	// NewExecutableSchema and Config are in the generated.go file
	// Resolver is in the resolver.go file

	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{
		Config:        server.Config,
		AccessControl: server.ACLEnforcer,
	}}))

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// NewRouter ...
func (s *Server) NewRouter() *gin.Engine {
	r := gin.Default()
	//r.Use(cors.Default())
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.GinContextToContextMiddleware())

	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	r.Group("/public")
	{
		r.POST("/login", auth.Login())
		r.POST("/legacy-login", auth.LegacyLogin())
		r.POST("/signup", auth.Signup)
		r.GET("/userTypes", controller.GetUserTypes)
		r.GET("/patientQueues", controller.GetPatientQueues)
		r.GET("/organizationDetails", controller.GetOrganizationDetails)

		r.Static("/files", "./files")

		r.GET("/recreate-opthalmology-exam", controller.RecreateOpthalmologyExam)

		r.GET("/rxnorm-drugs", controller.GetDrugs)
		r.GET("/rxnorm-intractions", controller.GetDrugIntractions)
	}

	r.GET("/clean", controller.ClearPatientsRecord)

	r.Use(middleware.AuthMiddleware())
	r.GET("/api", playgroundHandler())
	r.POST("/query", graphqlHandler(s))

	return r
}

// GetDB returns gorm (ORM)
func (s *Server) GetDB() *gorm.DB {
	return s.DB
}

// GetConfig return the current app configuration
func (s *Server) GetConfig() *conf.Configuration {
	return s.Config
}

// GetModelRegistry returns the model registry
func (s *Server) GetModelRegistry() *repository.Model {
	return s.ModelRegistry
}

// Start the http server
func (s *Server) Start() error {
	port := os.Getenv("ADDRESS")

	log.Fatal(s.Gin.Run(":" + port))
	return nil
}

func (s *Server) NewEnforcer() error {
	var model string
	var policy string

	appMode := os.Getenv("APP_MODE")

	if appMode == "release" {
		model = "/model.conf"
		policy = "/policy.csv"
	} else {
		model = "pkg/conf/model.conf"
		policy = "pkg/conf/policy.csv"
	}

	e, err := casbin.NewEnforcer(model, policy)
	if err != nil {
		log.Fatal(err)
	}

	s.ACLEnforcer = e
	return nil
}

// GracefulShutdown Wait for interrupt signal
// to gracefully shutdown the server with a timeout of 5 seconds.
func (s *Server) GracefulShutdown() {
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)
	<-quit
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// close database connection
	if s.DB != nil {
		db, _ := s.DB.DB()
		db.Close()
	}
}

// ShutdownTest shuts down test server
func (s *Server) ShutdownTest() {
	// close database connection
	if s.TestDB != nil {
		s.ModelRegistry.DropAll()
		db, _ := s.TestDB.DB()
		db.Close()
	}
}
