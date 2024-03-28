package main

import (
	"context"
	"fmt"
	"hotel-reservation/api"
	"hotel-reservation/db"
	"hotel-reservation/db/fixtures"
	"log"
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	ctx := context.Background()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	hotelStore := db.NewMongoHotelStore(client)

	store := db.Store{
		User:    db.NewMongoUserStore(client),
		Booking: db.NewMongoBookingStore(client),
		Room:    db.NewMongoRoomStore(client, hotelStore),
		Hotel:   hotelStore,
	}

	user := fixtures.AddUser(&store, "jimbo", "fatback", false)
	admin := fixtures.AddUser(&store, "admin", "hotel", true)
	hotel := fixtures.AddHotel(&store, "Hotel Alpha", "Narnia", 77, nil)
	room := fixtures.AddRoom(&store, "large", true, 345, hotel.ID)

	for i := 0; i < 100; i++ {
		name := fmt.Sprintf("Random-Hotel-%d", i)
		loc := fmt.Sprintf("Random-Location-%d", i)
		fixtures.AddHotel(&store, name, loc, rand.Intn(5)+1, nil)
	}

	var (
		from       = time.Now().AddDate(0, 0, 2)
		till       = time.Now().AddDate(0, 0, 4)
		numPersons = 3
	)

	booking := fixtures.AddBooking(&store, user.ID, room.ID, from, till, numPersons)

	fmt.Printf("\n*** >>> [admin] -> %+v\n", api.CreateTokenFromUser(admin))
	fmt.Printf("\n*** >>> [user] - > %+v\n", api.CreateTokenFromUser(user))
	fmt.Printf("\n*** >>> [hotel] %+v\n", hotel)
	fmt.Printf("\n*** >>> [room] %+v\n", room)
	fmt.Printf("\n*** >>> [booking] %+v\n", booking)

	fmt.Println("\n*** >>> SEED COMPLETE")
	fmt.Print("-----------------------------------------------------------------------------\n")
}
