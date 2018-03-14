package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
)

// Config type
type Config struct {
	ID       int `storm:"id,increment"`
	Height   float64
	Birthday time.Time
}

// Entry type
type Entry struct {
	ID       int `storm:"id,increment"`
	Date     time.Time
	Calories int
	Food     string
}

// Weight type
type Weight struct {
	ID     int `storm:"id,increment"`
	Date   time.Time
	Weight float64
}

func main() {
	db, err := storm.Open("test.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Config
	err = addConfig(db, 186, time.Now().AddDate(-30, 0, 0))
	if err != nil {
		log.Fatal(err)
	}
	var configs []Config
	err = db.All(&configs, storm.Limit(1), storm.Reverse())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(configs)

	// Weight
	err = addWeight(db, 86, time.Now())
	if err != nil {
		log.Fatal(err)
	}
	var weights []Weight
	err = db.All(&weights, storm.Reverse())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(weights)

	// Entries
	err = addEntry(db, "apple", 100, time.Now())
	if err != nil {
		log.Fatal(err)
	}
	err = addEntry(db, "bread", 300, time.Now().AddDate(0, 0, -2))
	if err != nil {
		log.Fatal(err)
	}

	var today []Entry
	err = db.Range("Date", time.Now().AddDate(0, 0, -1), time.Now().AddDate(0, 0, 1), &today)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Entries from Today:")
	fmt.Println(today)
	var twoDaysAgo []Entry
	query := db.Select(q.Gt("Date", time.Now().AddDate(0, 0, -3)), q.Lt("Date", time.Now().AddDate(0, 0, -1)))
	err = query.Find(&twoDaysAgo)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Entries from Two Days Ago:")
	fmt.Println(twoDaysAgo)

	var data []Entry
	var filters []q.Matcher
	filters = append(filters, q.Eq("Calories", 300))
	filters = append(filters, q.Eq("Food", "bread"))
	err = db.Select(filters...).Bucket("Entry").Find(&data)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Filtered bread using Select() will find data", data)

	var data2 []Entry
	var filters2 []q.Matcher
	filters2 = append(filters2, q.Eq("Calories", 50))
	filters2 = append(filters2, q.Eq("Food", "bread"))
	err = db.Select(filters2...).Bucket("Entry").Find(&data2)
	if err != nil && strings.Index(err.Error(), "not found") == -1 {
		log.Fatal(err)
	}
	log.Println("Filtered bread using Select() will not find data because the calories should match nothing", data2)
}

func addConfig(db *storm.DB, height float64, birthday time.Time) error {
	config := Config{Height: height, Birthday: birthday}
	err := db.Save(&config)
	if err != nil {
		return fmt.Errorf("could not save config, %v", err)
	}
	fmt.Println("config saved")
	return nil
}

func addWeight(db *storm.DB, value float64, date time.Time) error {
	weight := Weight{Weight: value, Date: date}
	err := db.Save(&weight)
	if err != nil {
		return fmt.Errorf("could not save weight, %v", err)
	}
	fmt.Println("weight saved")
	return nil
}

func addEntry(db *storm.DB, food string, calories int, date time.Time) error {
	entry := Entry{Food: food, Calories: calories, Date: date}
	err := db.Save(&entry)
	if err != nil {
		return fmt.Errorf("could not save entry, %v", err)
	}
	fmt.Println("entry saved")
	return nil
}
