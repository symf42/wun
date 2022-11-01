package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthorizationHeader struct {
	Basic string `header:"Authorization"`
}

func AuthorizationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		authorizationHeader := AuthorizationHeader{}

		err := c.ShouldBindHeader(&authorizationHeader)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if !strings.HasPrefix(authorizationHeader.Basic, "Basic ") {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		basicAuthString, err := base64.StdEncoding.DecodeString(authorizationHeader.Basic[6:])
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		basicAuthParts := strings.Split(string(basicAuthString), ":")

		conn, err := dbConnect()
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		defer conn.Close()

		stmt, err := conn.Prepare("SELECT `id`, `password` FROM `user` WHERE `email` = ? AND `activated_at` IS NOT NULL;")
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		var userId int
		var passwd string
		err = stmt.QueryRow(basicAuthParts[0]).Scan(&userId, &passwd)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(passwd), []byte(basicAuthParts[1])); err != nil {
			fmt.Println(err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		} else {
			c.Set("userId", userId)
			c.Next()
			return
		}

	}
}
