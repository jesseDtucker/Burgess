/*
 * This program takes the aggregated data and stores it to the archived database.
 *
 *	Author: James Finlay
 *	Date: March 6th, 2015j
 */

package priority

import (
    // "fmt"
    "time"
    "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
    "models"
)

var (
    c_arch *mgo.Collection
    sleepDuration time.Duration
)

func Init(session *mgo.Session, sleepDur time.Duration) {
    c_arch = session.DB("retailers").C("archived")
    c_employ = session.DB("retailers").C("employees")
    sleepDuration = sleepDur
}

/** Convert from Customer hash to Archived array **/
func FormatArchived(hash *map[string]*models.Customer) *[]models.Archived {
    result := make([]models.Archived, 0)
    for _,value := range *hash {
        result = append(result, *value.ToArchived())
    }
    return &result
}

/** Push aggregated data to the archive database **/
func StoreArchived(now time.Time, customers *map[string]*models.Customer) {
    data := FormatArchived(customers)
    if (len(*data) == 0) { return }

    err := c_arch.Insert(
    	bson.M{
    		"t": now,
    		"data": *data,
    	})
    if err != nil { panic(err) }
}