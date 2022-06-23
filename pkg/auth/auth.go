/*
  Copyright 2021 Kidus Tiliksew

  This file is part of Tensor EMR.

  Tensor EMR is free software: you can redistribute it and/or modify
  it under the terms of the version 2 of GNU General Public License as published by
  the Free Software Foundation.

  Tensor EMR is distributed in the hope that it will be useful,
  but WITHOUT ANY WARRANTY; without even the implied warranty of
  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
  GNU General Public License for more details.

  You should have received a copy of the GNU General Public License
  along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package auth

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/tensoremr/server/pkg/jwt"
	"github.com/tensoremr/server/pkg/repository"
)

// LoginPayload login body
type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginPayload login body
type LegacyLoginPayload struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse token response
type LoginResponse struct {
	Token string `json:"token"`
}

// Login logs users in
func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload LoginPayload
		var user repository.User

		err := c.ShouldBindJSON(&payload)
		if err != nil {
			c.JSON(400, gin.H{
				"message": "invalid json",
			})
			c.Abort()
			return
		}

		// Check if user exists
		eErr := user.GetByEmail(payload.Email)
		if eErr != nil {
			c.JSON(401, gin.H{
				"message": "Invalid user credentials",
			})
			c.Abort()
			return
		}

		// Check password validity
		pErr := user.CheckPassword(payload.Password)
		if pErr != nil {
			log.Println(err)
			c.JSON(401, gin.H{
				"message": "Invalid user credentials",
			})
			c.Abort()
			return
		}

		// Check if user is active
		if user.Active == false {
			c.JSON(401, gin.H{
				"message": "Your account is inactive",
			})
			c.Abort()
			return
		}

		jwtSecret := os.Getenv("JWT_SECRET")
		jwtIssuer := os.Getenv("JWT_ISSUER")

		// Generate JWT Token
		jwtWrapper := jwt.Wrapper{
			SecretKey:       jwtSecret,
			Issuer:          jwtIssuer,
			ExpirationHours: 24,
		}

		signedToken, err := jwtWrapper.GenerateToken(user)
		if err != nil {
			log.Println(err)
			c.JSON(500, gin.H{
				"msg": "error signing token",
			})
			c.Abort()
			return
		}

		tokenResponse := LoginResponse{
			Token: signedToken,
		}

		c.JSON(200, tokenResponse)

		return
	}
}

// Legacy Login logs users in
func LegacyLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload LegacyLoginPayload
		var user repository.User

		err := c.ShouldBindJSON(&payload)
		if err != nil {
			c.JSON(400, gin.H{
				"message": "invalid json",
			})
			c.Abort()
			return
		}

		// Check if user is legacy
		cErr := user.CheckIfUserLegacy(payload.Username)
		if cErr == nil {
			c.JSON(401, gin.H{
				"message": "User account is not legacy",
			})
			c.Abort()
			return
		}

		// Check if user exists
		eErr := user.GetByOldUserName(payload.Username)
		if eErr != nil {
			c.JSON(401, gin.H{
				"message": "Invalid user credentials",
			})
			c.Abort()
			return
		}

		// Check password validity
		pErr := user.CheckPassword(payload.Password)
		if pErr != nil {
			c.JSON(401, gin.H{
				"message": "Invalid user credentials",
			})
			c.Abort()
			return
		}

		user.Email = payload.Email
		if err := user.Update(nil); err != nil {
			c.JSON(500, gin.H{
				"message": "Sever error",
			})
			c.Abort()
			return
		}

		// Check if user is active
		if user.Active == false {
			c.JSON(401, gin.H{
				"message": "Your account is inactive",
			})
			c.Abort()
			return
		}

		jwtSecret := os.Getenv("JWT_SECRET")
		jwtIssuer := os.Getenv("JWT_ISSUER")

		// Generate JWT Token
		jwtWrapper := jwt.Wrapper{
			SecretKey:       jwtSecret,
			Issuer:          jwtIssuer,
			ExpirationHours: 24,
		}

		signedToken, err := jwtWrapper.GenerateToken(user)
		if err != nil {
			log.Println(err)
			c.JSON(500, gin.H{
				"msg": "error signing token",
			})
			c.Abort()
			return
		}

		tokenResponse := LoginResponse{
			Token: signedToken,
		}

		c.JSON(200, tokenResponse)

		return
	}
}

// Signup creates a user in db
func Signup(c *gin.Context) {
	var user repository.User

	err := c.ShouldBindJSON(&user)
	if err != nil {
		log.Println(err)

		c.JSON(400, gin.H{
			"msg": "invalid json",
		})
		c.Abort()

		return
	}

	err = user.HashPassword()
	if err != nil {
		log.Println(err.Error())

		c.JSON(500, gin.H{
			"msg": "error hashing password",
		})
		c.Abort()

		return
	}

	err = user.Save(nil)
	if err != nil {
		log.Println(err)

		c.JSON(500, gin.H{
			"msg": "error creating user",
		})
		c.Abort()

		return
	}

	c.JSON(200, user)
}
