package services

import (
	"bloodBankManagement/pojo"
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var SECRET_KEY string = os.Getenv("SECRET_KEY")

type SignedDetails struct {
	Email      string
	First_name string
	Last_name  string
	Uid        string
	jwt.StandardClaims
}

type Connection struct {
	Server      string
	Database    string
	Collection  string
	Collection2 string
	Collection3 string
	Collection4 string
}

var CollectionBlood *mongo.Collection
var CollectionDonor *mongo.Collection
var CollectionPatient *mongo.Collection
var CollectionLogin *mongo.Collection
var ctx = context.TODO()
var insertDocs int

func (c *Connection) Connect() {
	clientOptions := options.Client().ApplyURI(c.Server)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	CollectionBlood = client.Database(c.Database).Collection(c.Collection)
	CollectionDonor = client.Database(c.Database).Collection(c.Collection2)
	CollectionPatient = client.Database(c.Database).Collection(c.Collection3)
	CollectionLogin = client.Database(c.Database).Collection(c.Collection4)
}

func insertLoginData(email, pass, userId string) error {

	_, err := CollectionLogin.InsertOne(ctx, bson.M{
		"email":    email,
		"password": pass,
		"user_id":  userId,
	})

	if err != nil {
		return errors.New("Unable to create new record")
	}

	return nil
}

// ==========================================Donor detail======================================
func (e *Connection) SaveDonorDetails(reqBody pojo.DonorDetailRequest) (string, error) {
	saveData, err := SetValueInModel(reqBody)
	if err != nil {
		return "", errors.New("Unable to parse date")
	}
	data, err := CollectionDonor.InsertOne(ctx, saveData)
	if err != nil {
		log.Println(err)
		return "", errors.New("Unable to store data")
	}
	fmt.Println(data)
	str, err := saveBloodQuantityInBloodDetails(reqBody)
	if err != nil {
		log.Println(err)
		return "", err
	}
	fmt.Println(str)

	resultId := data.InsertedID
	email := reqBody.Email
	pass := reqBody.Password
	userId := fmt.Sprintf("%v", resultId)
	fmt.Println("userId:", userId)
	insertLoginData(email, pass, userId)

	return userId, nil
}

func SetValueInModel(req pojo.DonorDetailRequest) (pojo.DonorDetail, error) {
	var data pojo.DonorDetail
	data.DepositDate = time.Now()
	data.DOB = req.DOB
	data.Units = req.Units
	data.First_name = req.First_name
	data.Last_name = req.Last_name
	data.Age = req.Age
	data.AdharCardNo = req.AdharCardNo
	data.BloodGroup = req.BloodGroup
	data.Active = true
	data.Location = req.Location
	return data, nil
}

//==========================================Patient Details=====================================

func (c *Connection) SavePatientData(reqBody pojo.PatientDetailRequest) (string, error) {

	saveData, err := SetValueInPatientModel(reqBody)
	result, err := CollectionPatient.InsertOne(ctx, saveData)

	if err != nil {
		return "", errors.New("Unable to create new record")
	}

	resultId := result.InsertedID
	email := reqBody.Email
	pass := reqBody.Password
	userId := fmt.Sprintf("%v", resultId)
	fmt.Println("userId:", userId)
	insertLoginData(email, pass, userId)

	// fmt.Println("resultId:", resultId)
	return userId, nil
}

func (e *Connection) ApplyBloodPatientDetails(reqBody pojo.PatientDetailRequest, tokenId string) (string, error) {
	verifyToken, err := ValidateToken(tokenId)

	if verifyToken != "" {

		return verifyToken, err
	}

	deduct, err := deductOrAddBloodUnitsFromBloodDetails(reqBody.BloodGroup, reqBody.ApplyUnits, reqBody.Location, "Deduct")
	if err != nil {
		return "", err
	}
	fmt.Println(deduct)

	return deduct, nil
}

func SetValueInPatientModel(req pojo.PatientDetailRequest) (pojo.PatientDetail, error) {
	var data pojo.PatientDetail
	data.DOB = req.DOB
	data.First_name = req.First_name
	data.Last_name = req.Last_name
	data.Age = req.Age
	data.AdharCardNo = req.AdharCardNo
	data.BloodGroup = req.BloodGroup
	data.Active = true
	data.Location = req.Location
	data.CreatedAt = time.Now()
	data.ApplyUnits = req.ApplyUnits
	data.ApplyDate = time.Now()
	return data, nil
}

func convertDate(dateStr string) (time.Time, error) {

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		log.Println(err)
		return date, err
	}
	return date, nil
}

func deductOrAddBloodUnitsFromBloodDetails(bloodGroup, units, location, methodCall string) (string, error) {
	unitInt, err := convertUnitsStringIntoInt(units)
	if err != nil {
		fmt.Println(err)
		return "", nil
	}
	filter := bson.D{
		{Key: "$and",
			Value: bson.A{
				bson.D{{Key: "blood_group", Value: bloodGroup}},
				bson.D{{Key: "location", Value: location}},
			},
		},
	}
	fmt.Println(filter)
	data, err := CollectionBlood.Find(ctx, filter)
	finalData, err := convertDbResultIntoBloodStruct(data)
	fmt.Println(finalData)
	if err != nil {
		return "", nil
	}
	if finalData == nil {
		return "", errors.New("Data not present in Blood details according to given location and desposited date")
	}
	if methodCall == "Deduct" {
		unit := finalData[0].Units
		if !(unit >= unitInt) {
			return "", errors.New("Insufficient Blood!")
		}
		addUnit := unit - unitInt
		fmt.Println("Total Units:", addUnit)
		CollectionBlood.FindOneAndUpdate(ctx, filter, bson.D{{Key: "$set", Value: bson.D{{Key: "units", Value: addUnit}}}})
		return "Blood units Deduct Successfully", nil
	} else if methodCall == "Add" {
		unit := finalData[0].Units
		addUnit := unit + unitInt
		fmt.Println("Total Units:", addUnit)
		CollectionBlood.FindOneAndUpdate(ctx, filter, bson.D{{Key: "$set", Value: bson.D{{Key: "units", Value: addUnit}}}})
		return "Blood Units Added Successfully", nil
	}
	return "", nil
}

func convertUnitsStringIntoInt(units string) (int, error) {
	unitReplace := strings.ReplaceAll(units, "ml", "")
	unitInt, err := strconv.Atoi(unitReplace)
	if err != nil {
		fmt.Println(err)
		return 0, nil
	}
	return unitInt, nil
}

func (e *Connection) DeletePendingBloodPatientDetails(idStr string, tokenId string) (string, error) {
	verifyToken, err := ValidateToken(tokenId)

	if verifyToken != "" {

		return verifyToken, err
	}

	id, err := primitive.ObjectIDFromHex(idStr)

	if err != nil {
		return "", err
	}
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "active", Value: false}}}}
	CollectionPatient.FindOneAndUpdate(ctx, filter, update)
	data, err := CollectionPatient.Find(ctx, filter)
	if err != nil {
		return "", err
	}
	dataConv, err := convertDbResultIntoPatientStruct(data)
	if err != nil {
		return "", err
	}
	str, err := deductOrAddBloodUnitsFromBloodDetails(dataConv[0].BloodGroup, dataConv[0].ApplyUnits, dataConv[0].Location, "Add")
	if err != nil {
		return "", err
	}
	fmt.Println(str)
	return "Documents Deactivated Successfully", err
}

func convertDbResultIntoPatientStruct(fetchDataCursor *mongo.Cursor) ([]*pojo.PatientDetail, error) {
	var finaldata []*pojo.PatientDetail
	for fetchDataCursor.Next(ctx) {
		var data pojo.PatientDetail
		err := fetchDataCursor.Decode(&data)
		if err != nil {
			return finaldata, err
		}
		finaldata = append(finaldata, &data)
	}
	return finaldata, nil
}

/////////////blood detail//////////////////

func (e *Connection) SearchFilterBloodDetails(search pojo.BloodDetailsRequest, tokenId string) ([]*pojo.BloodBankDetail, string, error) {
	var searchData []*pojo.BloodBankDetail

	verifyToken, err := ValidateToken(tokenId)

	if verifyToken != "" {

		return searchData, verifyToken, err
	}

	filter := bson.D{}

	if search.BloodGroup != "" {
		filter = append(filter, primitive.E{Key: "blood_group", Value: bson.M{"$regex": search.BloodGroup}})
	}
	if search.Location != "" {
		filter = append(filter, primitive.E{Key: "location", Value: bson.M{"$regex": search.Location}})
	}
	if search.DepositDate != "" {
		depositDate, err := convertDate(search.DepositDate)
		if err != nil {
			return searchData, "", err
		}
		filter = append(filter, primitive.E{Key: "deposit-date", Value: bson.M{"$regex": depositDate}})
	}
	result, err := CollectionBlood.Find(ctx, filter)
	if err != nil {
		return searchData, "", err
	}
	data, err := convertDbResultIntoBloodStruct(result)
	if err != nil {
		return searchData, "", err
	}

	return data, "", nil
}

func convertDbResultIntoBloodStruct(fetchDataCursor *mongo.Cursor) ([]*pojo.BloodBankDetail, error) {
	var finaldata []*pojo.BloodBankDetail
	for fetchDataCursor.Next(ctx) {
		var data pojo.BloodBankDetail
		err := fetchDataCursor.Decode(&data)
		if err != nil {
			return finaldata, err
		}
		finaldata = append(finaldata, &data)
	}
	return finaldata, nil
}

func saveBloodQuantityInBloodDetails(reqBody pojo.DonorDetailRequest) (string, error) {
	var finalData []*pojo.BloodBankDetail
	unitInt, err := convertUnitsStringIntoInt(reqBody.Units)
	if err != nil {
		fmt.Println(err)
		return "", nil
	}
	filter := bson.D{
		{Key: "$and",
			Value: bson.A{
				bson.D{primitive.E{Key: "location", Value: reqBody.Location}},
				bson.D{primitive.E{Key: "blood_group", Value: reqBody.BloodGroup}},
			},
		},
	}
	data, err := CollectionBlood.Find(ctx, filter)
	finalData, err = convertDbResultIntoBloodStruct(data)
	if err != nil {
		return "", nil
	}
	if finalData == nil {
		saved, err := createNewEntryIntoBloodDetails(reqBody, unitInt)
		if err != nil {
			return "", err
		}
		fmt.Println(saved)
	} else {
		unitDB := finalData[0].Units
		addUnit := unitDB + unitInt
		fmt.Println("Total Units:", addUnit)
		CollectionBlood.FindOneAndUpdate(ctx, filter, bson.D{{Key: "$set", Value: bson.D{{Key: "units", Value: addUnit}}}})
	}
	return "Blood Details Saved Successfully", nil
}

func createNewEntryIntoBloodDetails(reqBody pojo.DonorDetailRequest, unitInt int) (string, error) {
	var bloodDetails pojo.BloodBankDetail

	bloodDetails.Units = unitInt
	bloodDetails.Location = reqBody.Location
	bloodDetails.BloodGroup = reqBody.BloodGroup
	bloodDetails.DepositDate = time.Now()
	bloodDetails.CreatedDate = time.Now()
	_, err := CollectionBlood.InsertOne(ctx, bloodDetails)
	if err != nil {
		log.Println(err)
		return "", nil
	}
	return "New entry created in blood details", nil
}

// Login is the api used to tget a single user
func (c *Connection) Login(data pojo.SignInInputRequest) (string, error) {

	var err error
	var foundUser *pojo.SignInInput

	cursor, err := CollectionLogin.Find(ctx, bson.D{primitive.E{Key: "email", Value: data.Email}})

	if err != nil {
		return "", errors.New("No record found in db")
	}

	for cursor.Next(ctx) {
		var e pojo.SignInInput
		err := cursor.Decode(&e)
		if err != nil {
			return "", err
		}
		foundUser = &e

	}

	if foundUser == nil {
		return "", errors.New("No data present in db for given email")
	}

	passwordIsValid := passwordVerify(data)
	if passwordIsValid != nil {
		return "", errors.New("login or passowrd is incorrect")
	}

	// str, err := verifyUserId(data)

	// if err != nil {
	// 	return "", errors.New("User Id is invalid")
	// }

	// userId := str
	// token, refreshToken, _ := GenerateAllTokens(foundUser.Email, foundUser.First_name, foundUser.Last_name)
	token, _ := GenerateAllTokens(foundUser.Email)

	return token, err
}
func passwordVerify(data pojo.SignInInputRequest) error {
	var passData *pojo.SignInInput
	var err error
	cursor, err := CollectionLogin.Find(ctx, bson.D{primitive.E{Key: "password", Value: data.Password}})
	fmt.Println("cursor:", cursor)
	if err != nil {
		return errors.New("login or passowrd is incorrect")
	}

	for cursor.Next(ctx) {
		var e pojo.SignInInput
		err := cursor.Decode(&e)
		if err != nil {
			return err
		}
		passData = &e
		fmt.Println("foundUser:", passData)
	}

	if passData == nil {
		return errors.New("No data present in db for given password")
	}
	return err
}

func verifyUserId(data pojo.SignInInputRequest) (string, error) {
	var donorData *pojo.DonorDetail
	var recordData *pojo.PatientDetail

	// var cursor *mongo.Cursor
	var err error
	var str = ""
	userId, err := primitive.ObjectIDFromHex(data.UserId)
	fmt.Println("userId:", userId)
	if err != nil {
		return "", err
	}
	if data.UserType == "donor" {
		cursor, err := CollectionDonor.Find(ctx, bson.D{primitive.E{Key: "_id", Value: userId}})
		fmt.Println("cursor:", cursor)
		if err != nil {
			return "", errors.New("User Id invalid")
		}

		for cursor.Next(ctx) {
			var e pojo.DonorDetail
			err := cursor.Decode(&e)
			if err != nil {
				return "", err
			}
			donorData = &e
			fmt.Println("foundUser:", donorData)
		}

		if donorData == nil {
			return "", errors.New("No data present in db for given user Id")
		}
		str = getDonorIdRecord(donorData)
		fmt.Println("str:", str)

	} else {
		cursor, err := CollectionPatient.Find(ctx, bson.D{primitive.E{Key: "_id", Value: userId}})
		fmt.Println("cursor:", cursor)
		if err != nil {
			return "", errors.New("User Id invalid")
		}

		for cursor.Next(ctx) {
			var e pojo.PatientDetail
			err := cursor.Decode(&e)
			if err != nil {
				return "", err
			}
			recordData = &e

		}

		if recordData == nil {
			return "", errors.New("No data present in db for given user Id")
		}
		str = getPatientIdRecord(recordData)
	}
	return str, err
}

func getDonorIdRecord(donorData *pojo.DonorDetail) string {
	userId := donorData.ID.Hex()

	return userId
}

func getPatientIdRecord(recordData *pojo.PatientDetail) string {
	userId := recordData.ID.Hex()

	return userId
}

func GenerateAllTokens(email string) (signedToken string, err error) {
	claims := &SignedDetails{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))

	return token, err
}

// ValidateToken validates the jwt token
func ValidateToken(signedToken string) (string, error) {
	var msg = ""
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)

	if err != nil {
		msg = err.Error()
		return msg, err
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = fmt.Sprintf("the token is invalid")
		msg = err.Error()
		return msg, err
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = fmt.Sprintf("token is expired")
		msg = err.Error()
		return msg, err
	}

	return msg, err
}
