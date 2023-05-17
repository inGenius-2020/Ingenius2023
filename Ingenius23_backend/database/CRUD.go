package database

import (
	"Ingenius23/communication"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	log "github.com/urishabh12/colored_log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func GetDatabaseConnection() (*gorm.DB, error) {
	dsn := goDotEnvVariable("DB_CONN")
	if db == nil { //If first time asking for database operations
		var err error
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Panic("Error creating a connection to databse!", err)
			return nil, err
		}
		db.AutoMigrate(Team{})
		db.AutoMigrate(Meals{})
		db.AutoMigrate(User{})
		db.Set("gorm:auto_preload", true)
	}
	return db, nil
}

func TestDocs(){
	db, err := GetDatabaseConnection()
	if err != nil {
		log.Println(err)
      return
	}
   where_cond := User{
      Phone : "82177803098",
   }
   db.Debug().Model(&where_cond).Update("Name","Test");
}

func CreateUserRecord(b communication.UserInitRequest) (string, int, bool) {
	db, err := GetDatabaseConnection()
	if err != nil {
		log.Println(err)
		return "Internal error", http.StatusInternalServerError, false
	}
	query_meal := Meals{Dinner1: false, Midnight1: false, Breakfast1: false, Lunch1: false, Coffee1: false, Coffee2: false, Coffee3: false}
	result := db.Create(&query_meal)
	if result.Error != nil {
		log.Println(err)

		return "Internal error", http.StatusInternalServerError, false
	}
	//Check is team existing : '

	team := Team{Team_id: b.Team_id}
	existing_team := Team{}
	db.First(&existing_team, &team)
	if existing_team.Team_id != team.Team_id {
		result = db.Create(&team)
		if result.Error != nil {
			log.Println(err)

			return "Internal error", http.StatusInternalServerError, false
		}
	} else {
		team.Table_no = existing_team.Table_no
	}

	query_user := User{
		Name:     b.Name,
		SRN:      b.SRN,
		Email:    b.Email,
		Phone:    b.Phone,
		Team_id:  b.Team_id,
		Team:     team,
		Role:     "participant",
		Present:  false,
		Checkin:  false,
		Checkout: false,
		Meal_id:  query_meal.Meal_id,
		Meals:    query_meal,
	}
	result = db.Create(&query_user)
	if result.Error != nil {
		log.Println(err)

		return "Internal error", http.StatusInternalServerError, false
	}

	return "User created", http.StatusOK, true
}

func CheckUserRecords(request communication.CheckInRequest) (bool, *User) {
	db, err := GetDatabaseConnection()
	if err != nil {
		return false, nil
	}
	query_user := User{
		SRN: request.SRN,
	}
	var existing_user User
	db.First(&existing_user, &query_user)
	log.Println(request, existing_user)
	if existing_user.SRN == request.SRN &&
		existing_user.Name == request.Name &&
		existing_user.Email == request.Email &&
		existing_user.Phone == request.Phone {
		//Need to do this check even though we have primary key as gorm add's it own primary key 'Id' making our entire primary key compostie and non uniuqe
		return true, &existing_user
	} else {
		return false, nil
	}
}

func SetCheckedInUser(user User) {
	db, err := GetDatabaseConnection()
	if err != nil {
		log.Panic("Something going wrong recording checkin attempt!")
		return
	}
	result := db.Model(&User{}).Where(&User{SRN: user.SRN}).Update("Checkin", true)
	if result.RowsAffected == 0 {
		log.Panic("Something going wrong recording checkin attempt!")
	}
}

func GetFullUserRecord(req jwt.MapClaims) (string, int, bool, *User) {
	db, err := GetDatabaseConnection()
	if err != nil {
		return "Internal Error", http.StatusInternalServerError, false, nil
	}
	query_user := User{
		SRN: req["SRN"].(string),
	}
	var existing_user User
	db.First(&existing_user, &query_user)
	var eager_load User
	db.Debug().Preload("Team").Preload("Meals").First(&eager_load, &existing_user)
	if existing_user.SRN == req["SRN"] {
		//Need to do this check even though we have primary key as gorm add's it own primary key 'Id' making our entire primary key compostie and non uniuqe
		return "Found user records", http.StatusOK, true, &eager_load
	} else {
		return "Internal Error", http.StatusInternalServerError, false, nil
	}
}

func SetUserAttendance(req jwt.MapClaims) (string, int, bool) {
	db, err := GetDatabaseConnection()
	if err != nil {
		return "Internal Error", http.StatusInternalServerError, false
	}
	update_user := User{
		SRN:        req["SRN"].(string),
		Present:    true,
		Entry_time: time.Now(),
	}
	result := db.Updates(&update_user)
	if result.RowsAffected == 0 {
		return "Attendance already recorded.", http.StatusForbidden, false
		//We are sure users exists cus token wont be valid if not. Unless we mess up our databse
	} else {
		return "Attendance recorded", http.StatusOK, true
	}
}

func SetUserCheckout(req jwt.MapClaims) (string, int, bool) {
	db, err := GetDatabaseConnection()
	if err != nil {
		return "Internal Error", http.StatusInternalServerError, false
	}
	update_user := User{
		SRN:       req["SRN"].(string),
		Checkout:  true,
		Exit_time: time.Now(),
	}
	result := db.Updates(&update_user)
	if result.RowsAffected == 0 {
		return "User already left, invalid action!", http.StatusForbidden, false
		//We are sure users exists cus token wont be valid if not. Unless we mess up our databse
	} else {
		return "Checkout Recorded!", http.StatusOK, true
	}
}

func SetFoodStatus(req jwt.MapClaims, FoodString string) (string, int, bool, *User) {
	message, httpstatus, status, fulluserrecord := GetFullUserRecord(req)
	if status == false {
		return message, httpstatus, status, fulluserrecord
	} else {
		result := db.Table("meals").Where(&Meals{Meal_id: fulluserrecord.Meal_id}).Update(FoodString, true)
		if result.RowsAffected == 0 {
			log.Println(result.Error)
			if result.Error != nil {
				return "Internal Error", http.StatusInternalServerError, false, nil
			}
			return "Invalid Food ID.", http.StatusBadRequest, false, nil
		} else {
			return "Updated food record successfully.", http.StatusOK, true, fulluserrecord //This full user record is NOT having the updated just made to meals
		}
		// if FoodString == "Dinner1" {
		// 	result := db.Table("meals").Where(&Meals{Meal_id: fulluserrecord.Meal_id}).Update("Dinner1", true)
		// 	if result.RowsAffected == 0 {
		// 		log.Println(result.Error)
		// 		return "Internal sdferver error.", http.StatusInternalServerError, false, nil
		// 	}
		// 	return "Updates user food records", http.StatusOK, true, fulluserrecord
		// }
		// if FoodString == "Midnight1" {
		// 	fulluserrecord.Meals.Midnight1 = true
		// 	result := db.Updates(fulluserrecord)
		// 	if result.RowsAffected == 0 {
		// 		return "Internal server error.", http.StatusInternalServerError, false, nil
		// 	}
		// 	return "Updates user food records", http.StatusOK, true, fulluserrecord
		// }
		// if FoodString == "Coffee1" {
		// 	fulluserrecord.Meals.Coffee1 = true
		// 	result := db.Updates(fulluserrecord)
		// 	if result.RowsAffected == 0 {
		// 		return "Internal server error.", http.StatusInternalServerError, false, nil
		// 	}
		// 	return "Updates user food records", http.StatusOK, true, fulluserrecord
		// }
		// if FoodString == "Coffee2" {
		// 	fulluserrecord.Meals.Coffee2 = true
		// 	result := db.Updates(fulluserrecord)
		// 	if result.RowsAffected == 0 {
		// 		return "Internal server error.", http.StatusInternalServerError, false, nil
		// 	}
		// 	return "Updates user food records", http.StatusOK, true, fulluserrecord
		// }
		// if FoodString == "Coffee3" {
		// 	fulluserrecord.Meals.Coffee3 = true
		// 	result := db.Updates(fulluserrecord)
		// 	if result.RowsAffected == 0 {
		// 		return "Internal server error.", http.StatusInternalServerError, false, nil
		// 	}
		// 	return "Updates user food records", http.StatusOK, true, fulluserrecord
		// }
		// if FoodString == "Breakfast1" {
		// 	fulluserrecord.Meals.Breakfast1 = true
		// 	result := db.Updates(fulluserrecord)
		// 	if result.RowsAffected == 0 {
		// 		return "Internal server error.", http.StatusInternalServerError, false, nil
		// 	}
		// 	return "Updates user food records", http.StatusOK, true, fulluserrecord
		// }
		// if FoodString == "Lunch1" {
		// 	fulluserrecord.Meals.Lunch1 = true
		// 	result := db.Updates(fulluserrecord)
		// 	if result.RowsAffected == 0 {
		// 		return "Internal server error.", http.StatusInternalServerError, false, nil
		// 	}
		// 	return "Updates user food records", http.StatusOK, true, fulluserrecord
		// }
	}
}
