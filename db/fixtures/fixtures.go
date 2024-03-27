package fixtures

import (
	"context"
	"fmt"
	"hotel-reservation/db"
	"hotel-reservation/types"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddUser(store *db.Store, fn string, ln string, isAdmin bool) *types.User {

	user, err := types.NewUserFromParams(types.CreateUserParams{
		FirstName: fn,
		LastName:  ln,
		Email:     fmt.Sprintf("%s@%s.com", fn, ln),
		Password:  fmt.Sprintf("%s_%s", fn, ln),
	})

	if err != nil {
		log.Fatal(err)
	}

	user.IsAdmin = isAdmin

	insertedUser, err := store.User.InsertUser(context.TODO(), user)

	if err != nil {
		log.Fatal(err)
	}

	return insertedUser
}

func AddHotel(store *db.Store, name string, loc string, rating int, rooms []primitive.ObjectID) *types.Hotel {

	var roomsIDs = rooms

	if roomsIDs == nil {
		roomsIDs = []primitive.ObjectID{}
	}

	hotel := types.Hotel{
		Name:     name,
		Location: loc,
		Rooms:    roomsIDs,
		Rating:   rating,
	}

	insertedHotel, err := store.Hotel.Insert(context.TODO(), &hotel)

	if err != nil {
		log.Fatal(err)
	}

	return insertedHotel
}

func AddRoom(store *db.Store, size string, ss bool, price float64, hotelID primitive.ObjectID) *types.Room {

	room := &types.Room{
		Size:    size,
		Seaside: ss,
		Price:   price,
		HotelID: hotelID,
	}

	insertedRoom, err := store.Room.InsertRoom(context.Background(), room)

	if err != nil {
		log.Fatal(err)
	}

	return insertedRoom
}

func AddBooking(store *db.Store, userID primitive.ObjectID, roomID primitive.ObjectID, from time.Time, till time.Time, numPersons int) *types.Booking {

	booking := &types.Booking{
		UserID:     userID,
		RoomID:     roomID,
		FromDate:   from,
		TillDate:   till,
		NumPersons: numPersons,
	}

	b, err := store.Booking.InsertBooking(context.Background(), booking)

	if err != nil {
		log.Fatal(err)
	}

	return b
}
