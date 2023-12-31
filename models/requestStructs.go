package models

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/lib/pq" // PostgreSQL driver
)

func init() {
	orm.RegisterModel(new(Users), new(Car), new(HomeSetting), new(EmailLogs))
}

// >>>>>>>>>>>>Models For tables Start from Here <<<<<<<<<<<<<<<<<<<<<<
type Users struct {
	Id          uint      `json:"user_id" orm:"pk;auto"`
	FirstName   string    `json:"first_name" orm:"column(first_name);null" valid:"MaxSize(255);Required"`
	LastName    string    `json:"last_name" orm:"null" valid:"MaxSize(255);Required"`
	Email       string    `json:"email" orm:"unique"`
	PhoneNumber string    `json:"phone_number" orm:"null"`
	Country     int       `json:"country_id"`
	Role        string    `json:"role"`
	Age         int       `json:"age"`
	Password    string    `json:"password"`
	Otp         string    `orm:"null"`
	Verified    string    `orm:"null"`
	CreatedAt   time.Time `orm:"null"`
	UpdatedAt   time.Time `orm:"null"`
	DeletedAt   time.Time `orm:"null"`
}

type CarType string

const (
	Sedan     CarType = "sedan"
	Hatchback CarType = "hatchback"
	SUV       CarType = "SUV"
)

type Car struct {
	Id          uint      `json:"car_id" orm:"pk;auto;column(id)"`
	CarName     string    `orm:"column(car_name)"`
	CarImage    string    `orm:"null;column(car_image)" form:"file" json:"file"`
	ModifiedBy  string    `orm:"column(modified_by)"`
	Model       string    `orm:"column(model)"`
	Type        CarType   `orm:"column(car_type);type(enum)"`
	CreatedDate time.Time `orm:"null;column(ctreated_date)"`
	UpdateDate  time.Time `orm:"null;column(updated_at)"`
}

type HomeSetting struct {
	Id          uint      `orm:"pk;auto;column(id);type(integer)"`
	Section     string    `orm:"column(section);type(character);size(255)"`
	Type        string    `orm:"column(type);type(character);size(255)"`
	Key         string    `orm:"column(key);type(character);size(255)"`
	Value       string    `orm:"column(value);type(character);size(255)"`
	Demo        string    `orm:"column(demo);type(text);size(255)"`
	CreatedDate time.Time `orm:"null;column(created_at);type(timestamptz)"`
	UpdateDate  time.Time `orm:"null;column(update_at);type(timestamptz)"`
}

type EmailLogs struct {
	Id      uint   `orm:"pk;auto;column(LogId)"`
	To      string `orm:"column(emailTo)"`
	Name    string `orm:"column(name)"`
	Subject string
	Body    string
	Status  string
}

//<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<End Table Models>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

type UserLoginRequest struct {
	Email    string `json:"username"`
	Password string `json:"password"`
}

type NewUserRequest struct {
	FirstName   string `json:"first_name" valid:"MaxSize(255);MinSize(3);Required;Match(/^[a-zA-Z][a-zA-Z0-9._-]{0,31}$/)"`
	LastName    string `json:"last_name" valid:"MaxSize(255);MinSize(3);Required;Match(/^[a-zA-Z][a-zA-Z0-9._-]{0,31}$/)"`
	Email       string `json:"email" valid:"MaxSize(255);Required;Email"`
	Country     int    `json:"country_id" valid:"Required"`
	Role        string `json:"role" valid:"MaxSize(255);Required"`
	Age         int    `json:"age" valid:"Range(1, 140);Required"`
	PhoneNumber string `json:"phone_number" valid:"Required;IsMobile"`
	Password    string `json:"password" valid:"MaxSize(25);MinSize(6);Required"`
}

type UpdateUserRequest struct {
	Id          uint   `json:"user_id" valid:"Required"`
	FirstName   string `json:"first_name" valid:"MaxSize(255);MinSize(3);Required;Match(/^[a-zA-Z][a-zA-Z0-9._-]{0,31}$/)"`
	LastName    string `json:"last_name" valid:"MaxSize(255);MinSize(3);Required;Match(/^[a-zA-Z][a-zA-Z0-9._-]{0,31}$/)"`
	Email       string `json:"email" valid:"MaxSize(255);Required;Email"`
	PhoneNumber string `json:"phone_number"  valid:"IsMobile;Required"`
	Country     int    `json:"country_id" valid:"Required"`
	Role        string `json:"role" valid:"MaxSize(255);Required"`
	Age         int    `json:"age"  valid:"Range(1, 140);Required"`
}

type ResetUserPassword struct {
	CurrentPass string `json:"current_password"`
	NewPass     string `json:"new_password"`
	ConfirmPass string `json:"confirm_password"`
}

type JwtClaim struct {
	Email string
	ID    int
	jwt.StandardClaims
}

type UserDetailsRequest struct {
	Id        uint   `json:"user_id"`
	FirstName string `json:"first_name" `
	LastName  string `json:"last_name" `
	Email     string `json:"email"`
	Age       int    `json:"age"`
	Country   int    `json:"country_id"`
	CreatedAt string `json:"CreatedDate"`
}

type SendOtpData struct {
	Username string `json:"username"`
}

type ResetUserPasswordOtp struct {
	Email   string `json:"email"`
	Otp     string `json:"otp"`
	NewPass string `json:"new_password"`
}

type VerifyEmailOTPRequest struct {
	Username string `json:"username"`
	Otp      string `json:"otp"`
}

type GetCountryWiseUserRequest struct {
	Country string `json:"country"`
	Count   int    `json:"count"`
}

/// Car request structs

type GetNewCarRequest struct {
	CarName    string  `json:"car_name" form:"car_name" valid:"MaxSize(255);MinSize(3);Required;Match(/^[a-zA-Z][a-zA-Z0-9._-]{0,31}$/)"`
	CarImage   string  `json:"car_imag" form:"file"`
	ModifiedBy string  `json:"modified_by" form:"modified_by" valid:"MaxSize(255);MinSize(3);Required;Match(/^[a-zA-Z][a-zA-Z0-9._-]{0,31}$/)"`
	Model      string  `json:"model" form:"model" valid:"MaxSize(255);MinSize(3);Required;Match(/^[a-zA-Z][a-zA-Z0-9._-]{0,31}$/)"`
	Type       CarType `json:"type" form:"type"`
}

type UpdateCarRequest struct {
	Id         uint    `json:"car_id" form:"car_id"`
	CarName    string  `json:"car_name" form:"car_name" valid:"MaxSize(255);MinSize(3);Required;Match(/^[a-zA-Z][a-zA-Z0-9._-]{0,31}$/)"`
	CarImage   string  `json:"car_imag" form:"file"`
	ModifiedBy string  `json:"modified_by" form:"modified_by" valid:"MaxSize(255);MinSize(3);Required;Match(/^[a-zA-Z][a-zA-Z0-9._-]{0,31}$/)"`
	Model      string  `json:"model" form:"model" valid:"MaxSize(255);MinSize(3);Required;Match(/^[a-zA-Z][a-zA-Z0-9._-]{0,31}$/)"`
	Type       CarType `json:"type" form:"type"`
}

type GetcarRequest struct {
	Id uint `json:"car_id"`
}
type OutgoingCallerID struct {
	PhoneNumber string `json:"phone_number"`
}

type SearchRequest struct {
	Search string `json:"search"`
}

type CarDetailsRequest struct {
	CarName    string  `json:"car_name"`
	CarImage   string  `json:"car_imag"`
	ModifiedBy string  `json:"modified_by"`
	Model      string  `json:"model"`
	Type       CarType `json:"type"`
}

// Home Setting reuests

type InserNewHomeSettingRequest struct {
	Section string `json:"section" form:"section" valid:"MaxSize(255);MinSize(3);Required;Match(/^[a-zA-Z][a-zA-Z0-9._-]{0,31}$/)"`
	Type    string `json:"type" form:"type" valid:"MaxSize(255);MinSize(3);Required;Match(/^[a-zA-Z][a-zA-Z0-9._-]{0,31}$/)"`
	Key     string `json:"key" form:"key"`
	Value   string `json:"value" form:"value"`
}

type UpdateHomeSetingRequest struct {
	Id      uint   `json:"home_seting_id" form:"home_seting_id"`
	Section string `json:"section" form:"section" valid:"MaxSize(255);MinSize(3);Required;Match(/^[a-zA-Z][a-zA-Z0-9._-]{0,31}$/)"`
	Type    string `json:"type" form:"type" valid:"MaxSize(255);MinSize(3);Required;Match(/^[a-zA-Z][a-zA-Z0-9._-]{0,31}$/)"`
	Key     string `json:"key" form:"key"`
	Value   string `json:"value" form:"value"`
}

type GetHomeSettingRequest struct {
	Id uint `json:"home_seting_id" form:"home_seting_id"`
}

type UserWiseHomeRequest struct {
	Section   string `json:"section"`
	Type      string `json:"type" `
	Key       string `json:"key" `
	Value     string `json:"value"`
	FirstName string `json:"first_name" `
	LastName  string `json:"last_name" `
}

// db syncronization

func SynchronizeModelWithDB(table_name string, modelType reflect.Type) error {
	// Get a database connection
	db, err := sql.Open("postgres", "user=postgres password=root dbname=mydb sslmode=disable")
	if err != nil {
		return err
	}
	defer db.Close()

	// Query the database to retrieve column information
	rows, err := db.Query("SELECT column_name, data_type FROM information_schema.columns WHERE table_name = '" + table_name + "'")
	if err != nil {
		return err
	}
	defer rows.Close()

	// Map to store column names and data types from the database
	dbColumns := make(map[string]string)

	// Iterate over the result set and populate the map
	for rows.Next() {
		var columnName, dataType string
		if err := rows.Scan(&columnName, &dataType); err != nil {
			return err
		}
		dbColumns[columnName] = dataType
	}

	if modelType.Kind() != reflect.Struct {
		return errors.New("tableModelType must be a struct type")
	}

	// Compare with the Beego model
	modelColumns := make(map[string]string)
	// model := new(table_modle)
	// modelType := reflect.TypeOf(model).Elem()

	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		ormTag := field.Tag.Get("orm")
		columnName := extractColumnName(ormTag)
		dataTypes := extractDatatype(ormTag)
		modelColumns[columnName] = dataTypes
	}
	// Identify differences and perform modifications
	// Example: remove columns not in the model
	for columnName := range dbColumns {
		if _, exists := modelColumns[columnName]; !exists {
			// Drop the column from the database
			_, err := db.Exec(fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s", table_name, columnName))
			if err != nil {
				return err
			}
			log.Printf("Dropped column: %s", columnName)
		}
	}

	for columnName, modelType := range modelColumns {
		dbType, exists := dbColumns[columnName]

		// Column exists but data type is different
		if exists && dbType != modelType {
			if modelType == "integer" && dbType == "varchar" || dbType == "char" || dbType == "text" {
				_, err := db.Exec(fmt.Sprintf("UPDATE %s SET %s = NULL WHERE %s !~ E'^\\d+$';", table_name, columnName, columnName))
				if err != nil {
					return err
				}
				_, err = db.Exec(fmt.Sprintf("ALTER TABLE %s ALTER COLUMN %s TYPE %s USING %s::%s", table_name, columnName, modelType, columnName, modelType))
				if err != nil {
					return err
				}
				log.Printf("Changed data type of column %s to %s", columnName, modelType)

			} else if modelType == "integer" && dbType == "float" {
				_, err := db.Exec(fmt.Sprintf("ALTER TABLE %s ALTER COLUMN %s TYPE %s USING ROUND(%s::numeric)", table_name, columnName, modelType, columnName))
				if err != nil {
					return err
				}
				log.Printf("Changed data type of column %s to %s", columnName, modelType)
			} else if modelType == "float" && dbType == "varchar" || dbType == "char" || dbType == "text" {
				_, err := db.Exec(fmt.Sprintf("UPDATE %s SET %s = NULL WHERE %s !~ E'^\\d+$';", table_name, columnName, columnName))
				if err != nil {
					return err
				}
				_, err = db.Exec(fmt.Sprintf("ALTER TABLE %s ALTER COLUMN %s TYPE %s USING NULLIF(%s, '')::%s", table_name, columnName, modelType, columnName, modelType))
				if err != nil {
					return err
				}
			} else {
				_, err = db.Exec(fmt.Sprintf("UPDATE %s SET %s = NULL", table_name, columnName))
				if err != nil {
					if modelType == "varchar" || modelType == "char" || modelType == "text" {
						_, err := db.Exec(fmt.Sprintf("UPDATE %s SET %s = ''", table_name, columnName))
						if err != nil {
							return err
						}
					}
					if modelType == "int" || modelType == "integer" || modelType == "float" {
						_, err := db.Exec(fmt.Sprintf("UPDATE %s SET %s = 0", table_name, columnName))
						if err != nil {
							return err
						}
					}
				}
				if modelType == "varchar" || modelType == "char" || modelType == "text" {
					_, err = db.Exec(fmt.Sprintf("ALTER TABLE %s ALTER COLUMN %s TYPE %s(255)", table_name, columnName, modelType))
					if err != nil {
						return err
					}
					log.Printf("Changed data type of column %s to %s", columnName, modelType)
				} else {
					// Alter the column type in the database
					_, err := db.Exec(fmt.Sprintf("ALTER TABLE %s ALTER COLUMN %s TYPE %s", table_name, columnName, modelType))
					if err != nil {
						return err
					}
					log.Printf("Changed data type of column %s to %s", columnName, modelType)
				}

			}
		}
	}

	// Update Beego model if needed
	// Example: add columns from the model to the database
	for columnName, dataType := range modelColumns {
		if _, exists := dbColumns[columnName]; !exists {
			// Add the column to the database
			num, err := db.Exec(fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", table_name, columnName, dataType))
			log.Printf(">>>>>>>>>>>>>>>>>>>>>> %s", num)
			if err != nil {
				return err
			}
			log.Printf("Added column: %s", columnName)
		}
	}

	return nil
}

func extractColumnName(ormTag string) string {
	if strings.Contains(ormTag, "column") {
		parts := strings.Split(ormTag, ";")
		for _, part := range parts {
			if strings.HasPrefix(part, "column") {
				// Extract the column name from "column(column_name)"
				return strings.Trim(strings.Split(part, "(")[1], ")")
			}
		}
	}
	return ""
}

func extractDatatype(ormTag string) string {
	if strings.Contains(ormTag, "type") {
		parts := strings.Split(ormTag, ";")
		for _, part := range parts {
			if strings.HasPrefix(part, "type") {
				// Extract the column name from "column(column_name)"
				return strings.Trim(strings.Split(part, "(")[1], ")")
			}
		}
	}
	return ""
}
