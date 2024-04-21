package model

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"example.com/gin/database"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var (
	jwtSecret = []byte(os.Getenv("AUTHO_KEY"))
)

type Stock struct {
	ID          int
	Description string
	Unit_price  int
	Units       int
}

type User struct {
	ID       int
	Username string
	Password string
}

type Claims struct {
	Username string `json:"username"`
	jwt.Claims
}

// For Stock Model

func AddStock(ctx *gin.Context) {
	body := Stock{}
	data, err := ctx.GetRawData()
	if err != nil {
		ctx.AbortWithStatusJSON(400, "User is not defined")
		return
	}
	err = json.Unmarshal(data, &body)
	if err != nil {
		ctx.AbortWithStatusJSON(400, "Bad Input")
		return
	}
	_, err = database.Db.Exec("insert into stocks(description,unit_price,units) values ($1,$2,$3)", body.Description, body.Unit_price, body.Units)
	if err != nil {
		fmt.Println(err)
		ctx.AbortWithStatusJSON(400, "Couldn't create the stock.")
	} else {
		ctx.JSON(http.StatusOK, "Stock is successfully created.")
	}
}

func GetStock(ctx *gin.Context) {
	rows, err := database.Db.Query("SELECT * from stocks")
	if err != nil {
		fmt.Println("Error querying the database:", err)
		ctx.AbortWithStatusJSON(400, "Row not defined")
		return
	}
	defer rows.Close()
	var result []Stock
	for rows.Next() {
		var ID int
		var Description string
		var Unit_price int
		var Units int
		if err := rows.Scan(&ID, &Description, &Unit_price, &Units); err != nil {
			fmt.Println("Error scanning row:", err)
			ctx.AbortWithStatusJSON(400, "Error querying db")
			return
		}
		data := Stock{
			ID:          ID,
			Description: Description,
			Unit_price:  Unit_price,
			Units:       Units,
		}
		result = append(result, data)
	}
	ctx.JSON(http.StatusOK, result)

}
func UpdateStock(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid stock ID"})
		return
	}
	var updatedStock Stock
	if err := ctx.BindJSON(&updatedStock); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
		return
	}
	result, err := database.Db.Exec("UPDATE stocks SET description=$1, unit_price=$2, units=$3 WHERE id=$4",
		updatedStock.Description, updatedStock.Unit_price, updatedStock.Units, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update stock"})
		return
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get rows affected"})
		return
	}
	if rowsAffected > 0 {
		ctx.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Stock with ID %d updated successfully", id)})
	} else {
		ctx.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Stock with ID %d not found", id)})
	}
}

func DeleteStock(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid stock ID"})
		return
	}
	result, err := database.Db.Exec("DELETE FROM stocks WHERE id=$1", id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete stock"})
		return
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get rows affected"})
		return
	}
	if rowsAffected > 0 {
		ctx.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Stock with ID %d deleted successfully", id)})
	} else {
		ctx.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Stock with ID %d not found", id)})
	}
}

// For User model

func AuthMiddleware(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		c.Abort()
		return
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		c.Abort()
		return
	}

	if !token.Valid {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		c.Abort()
		return
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		c.Abort()
		return
	}

	c.Set("username", claims.Username) // Store username in context for use in routes
	c.Next()
}
