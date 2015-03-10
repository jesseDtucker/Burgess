/*
 * This program performs the priority calculation for customer help.
 *
 *	Author: James Finlay
 *	Date: March 6th, 2015
 */

package priority

import (
	"models"
	"time"
	"gopkg.in/mgo.v2"
)

var (
    Customers = make(map[string]*models.Customer, 0)        // active Customers
    Employees = make(map[string]*models.Employee, 0)        // active Employees

	c_employ *mgo.Collection
    EmployeesAll = make(map[string]*models.Employee, 0)     // all Employees
    EmployeePullTime time.Time

    userExpiration = 20 * time.Second			// Time until user considered dead
	interactionExpiration =  10*time.Second 	// Time until interaction considered dead
	interactionDistance = float32(4)			// 4 metres, boundary box
)

/*
 *	Pull employee data into 'EmployeesAll'
 */
func PullEmployeeData() {
    var result []models.Employee
    err := c_employ.Find(nil).All(&result)
    if err != nil { panic(err) }

    for _,value := range result {
        EmployeesAll[value.MAC] = &value
    }
}

/*
 *	Find employee by MAC address.
 */
func FindEmployee(MAC string) *models.Employee {
    if Employees[MAC] != nil {
        return Employees[MAC]
    }
    if EmployeesAll[MAC] != nil {
        Employees[MAC] = EmployeesAll[MAC]
        Employees[MAC].FirstSeen = time.Now()
        return Employees[MAC]
    }
    return nil
}

/*
 *	Update employees & customers with new positions and time
 */
func UpdateUsers(data *map[string]*models.Position) {

    for _,value := range *data {

        // Employee?
        employee := FindEmployee(value.Wifi)
        if employee != nil {
            employee.LastSeen = time.Now()
            employee.Position = *value
        } else if Customers[value.Wifi] != nil {
            Customers[value.Wifi].LastSeen = time.Now()
            Customers[value.Wifi].Position = *value
        } else {
            Customers[value.Wifi] = &models.Customer{value.Wifi, time.Now(),
                time.Now(), *value, time.Now(), 0, nil}
        }
    }

    // Kill expired users
    for key,value := range Employees {
        if time.Since(value.LastSeen) > userExpiration {
            delete(Employees, key)
        }
    }
    for key,value := range Customers {
        if time.Since(value.LastSeen) > userExpiration {
            delete(Customers, key)
        }
    }

	// Kill expired interactions
	for _,employee := range Employees {
		for i,interaction := range employee.Interactions {
			if time.Since(interaction.LastTime) > interactionExpiration {
				// TODO::JF Should store this event for analytics
				eTime := interaction.GetPriorityTime()
				if eTime.UnixNano() > interaction.Customer.ExpiryTime.UnixNano() {
					interaction.Customer.ExpiryTime = eTime
				}
				// kill the interaction
				interaction.Customer.RemoveInteraction(interaction)
				employee.Interactions = append(employee.Interactions[:i],
					employee.Interactions[i+1:]...)
			}
		}
	}
}

/*
 *	Update interactions between Customers and Employees
 */
func UpdateInteractions() {
	for _,customer := range Customers {
		for _,employee := range Employees {
			if employee.Position.X > customer.Position.X + interactionDistance {continue}
			if employee.Position.Y > customer.Position.Y + interactionDistance {continue}
			if employee.Position.X < customer.Position.X - interactionDistance {continue}
			if employee.Position.Y < customer.Position.Y - interactionDistance {continue}

			interaction := models.FindByCustomer(employee.Interactions, customer)
			if interaction != nil {
				interaction.LastTime = time.Now()
			} else {
				interaction = &models.Interaction{employee, customer, time.Now(), time.Now()}
				customer.Interactions = append(customer.Interactions, interaction)
				employee.Interactions = append(employee.Interactions, interaction)
			}
		}

		// Priority
		if len(customer.Interactions) == 0 {
			dt := float32(customer.ExpiryTime.UnixNano() -
				time.Now().UnixNano() + int64(sleepDuration))
			if dt == 0 {
				customer.Priority = 1
			} else {
				customer.Priority += (1-customer.Priority)*float32(sleepDuration)/dt
			}
			if customer.Priority > 1 { customer.Priority = 1}
			if customer.Priority < 0 { customer.Priority = 0}
		}
	}
}

/*
 *	Super function for updating priorities. This also helps track analytics-
 *	based data.
 */
func UpdatePriorities(data *map[string]*models.Position) *map[string]*models.Customer {
    // TODO::JF - Filter by retailer

    // If it's been a while, update our EmployeeAll data
    if time.Since(EmployeePullTime) > time.Hour {
        PullEmployeeData()
        EmployeePullTime = time.Now()
    }

	UpdateUsers(data)
	UpdateInteractions()

	return &Customers
}