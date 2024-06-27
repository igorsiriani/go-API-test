package main

import (
	"encoding/csv"
	// "fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"golang.org/x/crypto/bcrypt"
)
 
var db *gorm.DB
var err error
 
func main() {
    db, err = gorm.Open("sqlite3", "test.db")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
 
    db.AutoMigrate(&User{}, &Record{})
 
    r := gin.Default()
 
	r.POST("/register", Register)
	r.POST("/login", Login)
    r.GET("/search", Search)
    r.GET("/upload", Upload)
	r.GET("/user", AuthMiddleware(), GetUser)
	r.PUT("/user/:id", AuthMiddleware(), UpdateUser)
	r.DELETE("/users:id", AuthMiddleware(), DeleteUser)
 
	r.Run(":8080")
}

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenString := c.GetHeader("Authorization")
        if tokenString == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
            c.Abort()
            return
        }

        println("tokenString: ")
        println(tokenString)

        // Remove "Bearer " prefix if present
        if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
            tokenString = tokenString[7:]
        }

        claims := &Claims{}
        token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
            return jwtKey, nil
        })

        if err != nil || !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        var user User
        if err := db.Where("username = ?", claims.Username).First(&user).Error; err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
            c.Abort()
            return
        }

        c.Set("userID", user.ID)
        c.Next()
    }
}

var jwtKey = []byte("your_secret_key")
 
func Register(c *gin.Context) {
    var creds Credentials
    if err := c.BindJSON(&creds); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }
 
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
        return
    }
 
    user := User{Username: creds.Username, Password: string(hashedPassword)}
    if err := db.Create(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
        return
    }
 
    c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}
 
func Login(c *gin.Context) {
    var creds Credentials
    if err := c.BindJSON(&creds); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }

    var user User
    if err := db.Where("username = ?", creds.Username).First(&user).Error; err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
        return
    }

    expirationTime := time.Now().Add(5 * time.Minute)
    claims := &Claims{
        Username: creds.Username,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: expirationTime.Unix(),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(jwtKey)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
 
func loadCSV(filename string) error {
	layout := "01/02/2006"

    file, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer file.Close()
 
    reader := csv.NewReader(file)
    reader.Comma = ';'
    reader.LazyQuotes = true
    records, err := reader.ReadAll()
    if err != nil {
        return err
    }
 
    for _, row := range records[1:] { // skipping header row
		MpnId, _ := strconv.Atoi(row[6])
		ChargeStartDate, _ := time.Parse(layout, row[18])
		ChargeEndDate, _ := time.Parse(layout, row[19])
		UsageDate, _ := time.Parse(layout, row[20])
		UnitPrice, _ := strconv.ParseFloat(row[33], 64)
		Quantity, _ := strconv.ParseFloat(row[34], 64)
		BillingPreTaxTotal, _ := strconv.ParseFloat(row[36], 64)
		PricingPreTaxTotal, _ := strconv.ParseFloat(row[38], 64)
		EffectiveUnitPrice, _ := strconv.ParseFloat(row[44], 64)
		PCToBCExchangeRate, _ := strconv.Atoi(row[45])
		PCToBCExchangeRateDate, _ := time.Parse(layout, row[46])
		PartnerEarnedCreditPercentage, _ := strconv.Atoi(row[49])
		CreditPercentage, _ := strconv.Atoi(row[50])

        db.Create(&Record{
			PartnerId: row[0],
			PartnerName: row[1],
			CustomerId: row[2],
			CustomerName: row[3],
			CustomerDomainName: row[4],
			CustomerCountry: row[5],
			MpnId: MpnId,
			Tier2MpnId: row[7],
			InvoiceNumber: row[8],
			ProductId: row[9],
			SkuId: row[10],
			AvailabilityId: row[11],
			SkuName: row[12],
			ProductName: row[13],
			PublisherName: row[14],
			PublisherId: row[15],
			SubscriptionDescription: row[16],
			SubscriptionId: row[17],
			ChargeStartDate: ChargeStartDate,
			ChargeEndDate: ChargeEndDate,
			UsageDate: UsageDate,
			MeterType: row[21],
			MeterCategory: row[22],
			MeterId: row[23],
			MeterSubCategory: row[24],
			MeterName: row[25],
			MeterRegion: row[26],
			Unit: row[27],
			ResourceLocation: row[28],
			ConsumedService: row[29],
			ResourceGroup: row[30],
			ResourceURI: row[31],
			ChargeType: row[32],
			UnitPrice: UnitPrice,
			Quantity: Quantity,
			UnitType: row[35],
			BillingPreTaxTotal: BillingPreTaxTotal,
			BillingCurrency: row[37],
			PricingPreTaxTotal: PricingPreTaxTotal,
			PricingCurrency: row[39],
			ServiceInfo1: row[40],
			ServiceInfo2: row[41],
			Tags: row[42],
			AdditionalInfo: row[43],
			EffectiveUnitPrice: EffectiveUnitPrice,
			PCToBCExchangeRate: PCToBCExchangeRate,
			PCToBCExchangeRateDate: PCToBCExchangeRateDate,
			EntitlementId: row[47],
			EntitlementDescription: row[48],
			PartnerEarnedCreditPercentage: PartnerEarnedCreditPercentage,
			CreditPercentage: CreditPercentage,
			CreditType: row[51],
			BenefitOrderId: row[52],
			BenefitId: row[53],
			BenefitType: row[54],
        })
    }
 
    return nil
}

func Upload(c *gin.Context) {
    if err := loadCSV("dataset/data.csv"); err != nil {
        log.Fatal(err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating Database"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Database created successfully"})
}
 
func GetUser(c *gin.Context) {
    userID, _ := c.Get("userID")
    var user User
    if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }
    c.JSON(http.StatusOK, gin.H{ "ID": user.ID, "Username": user.Username})
}
 
func UpdateUser(c *gin.Context) {
    userID, _ := c.Get("userID")
    var user User
    if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }
    if err := c.BindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }
    db.Save(&user)
    c.JSON(http.StatusOK, user)
}
 
func DeleteUser(c *gin.Context) {
    userID, _ := c.Get("userID")
    var user User
    if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }
    db.Delete(&user)
    c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func Search(c *gin.Context) {
    var searchRequest SearchRequest
    if err := c.ShouldBindJSON(&searchRequest); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }
 
    var query string
    var args []interface{}
    query = ""
 
    var count = 0
    for _, criteria := range searchRequest.Criteria {
        count += 1

        if count > 1 {
            query += " AND "
        }

        if criteria.Value2 == "" {
            // Single value match
            query += criteria.Column + " = ?"
            args = append(args, criteria.Value1)
        } else {
            // Value interval
            query += criteria.Column + " BETWEEN ? AND ?"
            args = append(args, criteria.Value1, criteria.Value2)
        }
    }
 
    var records []Record
    db.Where(query, args...).Find(&records)
    c.JSON(http.StatusOK, records)
}