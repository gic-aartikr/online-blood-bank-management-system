package pojo

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BloodBankDetail struct {
	Id          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	BloodGroup  string             `bson:"blood_group,omitempty" json:"blood_group,omitempty"`
	Units       int                `bson:"units,omitempty" json:"units,omitempty"`
	Location    string             `bson:"location,omitempty" json:"location,omitempty"`
	DepositDate time.Time          `bson:"deposit_date,omitempty" json:"deposit_date,omitempty"`
	CreatedDate time.Time          `bson:"created_date,omitempty" json:"created_date,omitempty"`
}

type DonorDetail struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	First_name  string             `bson:"first_name,omitempty" json:"first_name" validate:"required,min=2,max=100"`
	Last_name   string             `bson:"last_name,omitempty" json:"last_name" validate:"required,min=2,max=100"`
	Age         string             `bson:"age,omitempty" json:"age,omitempty"`
	AdharCardNo string             `bson:"adhar_card_no,omitempty" json:"adhar_card_no,omitempty"`
	Phone       string             `bson:"phone,omitempty" json:"phone" validate:"required"`
	DOB         string             `bson:"dob,omitempty" json:"dob,omitempty"`
	BloodGroup  string             `bson:"blood_group,omitempty" json:"blood_group,omitempty"`
	Units       string             `bson:"units,omitempty" json:"units,omitempty"`
	DepositDate time.Time          `bson:"deposit_date,omitempty" json:"deposit_date,omitempty"`
	Location    string             `bson:"location,omitempty" json:"location,omitempty"`
	Active      bool               `bson:"active,omitempty" json:"active,omitempty"`
	CreatedAt   time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt   time.Time          `bson:"updated_at,omitempty" json:"updated_at" bson:"updated_at,omitempty"`
}

type DonorDetailRequest struct {
	First_name  string `bson:"first_name,omitempty" json:"first_name" validate:"required,min=2,max=100"`
	Last_name   string `bson:"last_name,omitempty" json:"last_name" validate:"required,min=2,max=100"`
	Age         string `bson:"age,omitempty" json:"age,omitempty"`
	AdharCardNo string `bson:"adhar_card_no,omitempty" json:"adhar_card_no,omitempty"`
	Email       string `bson:"email,omitempty" json:"email" validate:"email,required"`
	Password    string `bson:"password,omitempty" json:"password" validate:"required"`
	Phone       string `bson:"phone,omitempty" json:"phone" validate:"required"`
	DOB         string `bson:"dob,omitempty" json:"dob,omitempty"`
	BloodGroup  string `bson:"blood_group,omitempty" json:"blood_group,omitempty"`
	Units       string `bson:"units,omitempty" json:"units,omitempty"`
	DepositDate string `bson:"deposit_date,omitempty" json:"deposit_date,omitempty"`
	Location    string `bson:"location,omitempty" json:"location,omitempty"`
}

type PatientDetail struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	First_name  string             `bson:"first_name,omitempty" json:"first_name" validate:"required,min=2,max=100"`
	Last_name   string             `bson:"last_name,omitempty" json:"last_name" validate:"required,min=2,max=100"`
	Age         string             `bson:"age,omitempty" json:"age,omitempty"`
	DOB         string             `bson:"dob,omitempty" json:"dob,omitempty"`
	AdharCardNo string             `bson:"adhar_card_no,omitempty" json:"adhar_card_no,omitempty"`
	Phone       string             `bson:"phone,omitempty" json:"phone" validate:"required"`
	BloodGroup  string             `bson:"blood_group,omitempty" json:"blood_group,omitempty"`
	Location    string             `bson:"location,omitempty" json:"location,omitempty"`
	Active      bool               `bson:"active,omitempty" json:"active,omitempty"`
	ApplyUnits  string             `bson:"apply_units,omitempty" json:"apply_units,omitempty"`
	ApplyDate   time.Time          `bson:"apply_date,omitempty" json:"apply_date,omitempty"`
	GivenDate   time.Time          `bson:"given_date,omitempty" json:"given_date,omitempty"`
	CreatedAt   time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt   time.Time          `bson:"updated_at,omitempty" json:"updated_at" bson:"updated_at,omitempty"`
}

type PatientDetailRequest struct {
	First_name  string    `bson:"first_name,omitempty" json:"first_name" validate:"required,min=2,max=100"`
	Last_name   string    `bson:"last_name,omitempty" json:"last_name" validate:"required,min=2,max=100"`
	Age         string    `bson:"age,omitempty" json:"age,omitempty"`
	DOB         string    `bson:"dob,omitempty" json:"dob,omitempty"`
	AdharCardNo string    `bson:"adhar_card_no,omitempty" json:"adhar_card_no,omitempty"`
	Email       string    `bson:"email,omitempty" json:"email" validate:"email,required"`
	Password    string    `bson:"password,omitempty" json:"password" validate:"required"`
	Phone       string    `bson:"phone,omitempty" json:"phone" validate:"required"`
	BloodGroup  string    `bson:"blood_group,omitempty" json:"blood_group,omitempty"`
	GivenDate   time.Time `bson:"given_date,omitempty" json:"given_date,omitempty"`
	Location    string    `bson:"location,omitempty" json:"location,omitempty"`
	Active      bool      `bson:"active,omitempty" json:"active,omitempty"`
	ApplyUnits  string    `bson:"apply_units,omitempty" json:"apply_units,omitempty"`
}
type SignInInput struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Email    string             `bson:"email,omitempty" json:"email,omitempty"`
	Password string             `bson:"password,omitempty" json:"password,omitempty"`
	UserId   string             `bson:"user_id,omitempty" json:"user_id"`
}

type SignInInputRequest struct {
	Email    string `bson:"email,omitempty" json:"email,omitempty"`
	Password string `bson:"password,omitempty" json:"password,omitempty"`
	UserId   string `bson:"user_id,omitempty" json:"user_id"`
	UserType string `bson:"user_type,omitempty" json:"user_type"`
}

// type SignUpInput struct {
// 	Name            string    `json:"name" bson:"name" binding:"required"`
// 	Email           string    `json:"email" bson:"email" binding:"required"`
// 	Password        string    `json:"password" bson:"password" binding:"required,min=8"`
// 	PasswordConfirm string    `json:"passwordConfirm" bson:"passwordConfirm,omitempty" binding:"required"`
// 	Role            string    `json:"role" bson:"role"`
// 	Verified        bool      `json:"verified" bson:"verified"`
// 	CreatedAt       time.Time `json:"created_at" bson:"created_at"`
// 	UpdatedAt       time.Time `json:"updated_at" bson:"updated_at"`
// }

type Response struct {
	Success    string      `json:"success,omitempty"`
	SucessCode string      `json:"successCode,omitempty"`
	Response   interface{} `json:"response,omitempty"`
}

type BloodDetailsRequest struct {
	BloodGroup  string `bson:"blood_group,omitempty" json:"blood_group,omitempty"`
	Units       string `bson:"units,omitempty" json:"units,omitempty"`
	Location    string `bson:"location,omitempty" json:"location,omitempty"`
	DepositDate string `bson:"deposit_date,omitempty" json:"deposit_date,omitempty"`
}
