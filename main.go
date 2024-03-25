package main

import (
	"hotel-reservation/api"
	"hotel-reservation/api/middleware"
	"hotel-reservation/db"

	"context"
	"flag"
	"log"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var config = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {

	listenAddr := flag.String("listenAddr", ":5000", "The listening address of the API server")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	var (
		userStore  = db.NewMongoUserStore(client)
		hotelStore = db.NewMongoHotelStore(client)
		roomStore  = db.NewMongoRoomStore(client, hotelStore)

		store = &db.Store{
			User:  userStore,
			Room:  roomStore,
			Hotel: hotelStore,
		}

		hotelHandler = api.NewHotelHandler(store)
		userHandler  = api.NewUserHandler(userStore)
		authHandler  = api.NewAuthHandler(userStore)

		app   = fiber.New(config)
		auth  = app.Group("/api")
		apiV1 = app.Group("/api/v1", middleware.JWTAuthentication)
	)

	// Auth -------------------------------------------------------
	auth.Post("/auth", authHandler.HandleAuthenticate)

	// User -------------------------------------------------------
	apiV1.Post("/user", userHandler.HandlePostUser)
	apiV1.Get("/user", userHandler.HandleGetUsers)
	apiV1.Get("/user/:id", userHandler.HandleGetUserById)
	apiV1.Put("/user/:id", userHandler.HandlePutUser)
	apiV1.Delete("/user/:id", userHandler.HandleDeleteUser)

	// Hotel ------------------------------------------------------
	apiV1.Get("/hotel", hotelHandler.HandlerGetHotels)
	apiV1.Get("/hotel/:id", hotelHandler.HandleGetHotelById)
	apiV1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms)

	app.Listen(*listenAddr)
}
