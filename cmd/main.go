package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/apex/gateway"
	"github.com/gin-gonic/gin"
)

func inLambda() bool {
	if lambdaTaskRoot := os.Getenv("LAMBDA_TASK_ROOT"); lambdaTaskRoot != "" {
		return true
	}
	return false
}

func getFib(n int) int {
	f := make([]int, n+1)
	f[0] = 0
	f[1] = 1
	for i := 2; i <= n; i++ {
		f[i] = f[i-1] + f[i-2]
	}
	return f[n]
}

func setupRouter() *gin.Engine {

	r := gin.Default()

	r.GET("/fib/n=:n", func(c *gin.Context) {
		num := c.Param("n")
		n, err := strconv.Atoi(num)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "n needs to be a valid positive integer."})
			return
		} else if n < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "n cannot be negative."})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": getFib(n)})
	})

	return r
}

func main() {
	if inLambda() {
		fmt.Println("running aws lambda in aws")
		log.Fatal(gateway.ListenAndServe(":8080", setupRouter()))
	} else {
		fmt.Println("running aws lambda in local")
		log.Fatal(http.ListenAndServe(":8080", setupRouter()))
	}
}
