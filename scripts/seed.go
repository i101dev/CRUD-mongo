package main

import (
	"context"
	"hotel-reservation/db"
	"hotel-reservation/types"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client     *mongo.Client
	roomStore  db.RoomStore
	hotelStore db.HotelStore
	userStore  db.UserStore
	ctx        = context.Background()
)

func seedUser(fName string, lName string, email string, password string, isAdmin bool) {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		FirstName: fName,
		LastName:  lName,
		Email:     email,
		Password:  password,
	})

	if err != nil {
		log.Fatal(err)
	}

	user.IsAdmin = isAdmin

	_, err = userStore.InsertUser(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}
}

func SeedHotel(name string, location string, rating int) {

	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}

	rooms := []types.Room{
		{
			Size:  "small",
			Price: 100.1,
		},
		{
			Size:  "normal",
			Price: 200.2,
		},
		{
			Size:  "kingsize",
			Price: 300.3,
		},
	}

	insertedHotel, err := hotelStore.Insert(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println("Inserted hotel - ", insertedHotel)

	for _, room := range rooms {

		room.HotelID = insertedHotel.ID

		_, err := roomStore.InsertRoom(ctx, &room)

		// fmt.Println("Inserted room - ", r)

		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	SeedHotel("Sloop", "Asgard", 25)
	SeedHotel("Hoptop", "Exemy", 90)
	seedUser("Jimmy", "Jackme", "jimmy@mail.com", "fuqyourcouch", false)
	seedUser("Fuq", "Yocouch", "fyc@mail.com", "fuqyourcouch", false)
	seedUser("admin", "admin", "admin@mail.com", "iamtheadmin", true)
}

func init() {

	var err error

	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	hotelStore = db.NewMongoHotelStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
	userStore = db.NewMongoUserStore(client)
}
