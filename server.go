package main

import (
	"log"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/sebmartin/collabd/game"
	"github.com/sebmartin/collabd/games/connect4"
	"github.com/sebmartin/collabd/graph"
	"github.com/sebmartin/collabd/graph/generated"
)

const defaultPort = "8080"

// This is an example server that launches the game server with
// a test game
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv, err := game.NewServer("sqlite", "db/models.sqlite")
	if err != nil {
		log.Fatalf("Failed to initalize game server: %s", err)
	}

	connect4.Register()

	r := gin.Default()
	r.POST("/query", graphqlHandler(srv))
	r.GET("/", playgroundHandler())

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	r.Run(":" + port)
}

func graphqlHandler(s *game.Server) gin.HandlerFunc {
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: &graph.Resolver{GameServer: s},
	}))

	return func(c *gin.Context) {
		srv.ServeHTTP(c.Writer, c.Request)
	}
}

func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
