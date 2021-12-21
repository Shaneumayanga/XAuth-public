package ratelimit

//NEED improvements , i know :-)
import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tomasen/realip"
)

type RateLimiter struct {
	activeclients  map[string]int
	blockedClients map[string]bool
	clientnew      chan string
	limit          int
}

func NewRateLimit(limit int) *RateLimiter {
	return &RateLimiter{
		activeclients:  make(map[string]int),
		blockedClients: make(map[string]bool),
		clientnew:      make(chan string),
		limit:          limit,
	}
}

func (rl *RateLimiter) Start() {
	for {
		select {
		case clientIP := <-rl.clientnew:
			fmt.Printf("clientIP: %v\n", clientIP)
			activeClientRequest, ok := rl.activeclients[clientIP]
			if !ok {
				rl.activeclients[clientIP] = 0
			}
			if activeClientRequest >= rl.limit {
				//block ip here
				rl.blockedClients[clientIP] = true
				delete(rl.activeclients, clientIP)
			}
			rl.activeclients[clientIP] += 1
		case <-time.After(5 * time.Minute):
			//delete all blocked clients
			for blockedClients := range rl.blockedClients {
				rl.activeclients[blockedClients] = 0
			}
			for k := range rl.blockedClients {
				delete(rl.blockedClients, k)
			}
			fmt.Println("All unblocked!")
		}
	}
}

func (rl *RateLimiter) IsBlocked(ip string) bool {
	_, isblocked := rl.blockedClients[ip]
	return isblocked

}

func (rl *RateLimiter) Use() func(c *gin.Context) {
	return func(c *gin.Context) {
		clientIP := realip.FromRequest(c.Request)
		fmt.Printf("clientIP: %v\n", clientIP)
		rl.clientnew <- clientIP
		if rl.IsBlocked(clientIP) {
			http.Redirect(c.Writer, c.Request, "/onblock", http.StatusSeeOther)
		}
		c.Next()
	}
}

func (rl *RateLimiter) UseREST() func(c *gin.Context) {
	return func(c *gin.Context) {
		clientIP := realip.FromRequest(c.Request)
		fmt.Printf("clientIP: %v\n", clientIP)
		rl.clientnew <- clientIP
		if rl.IsBlocked(clientIP) {
			c.AbortWithStatusJSON(http.StatusForbidden, "You were blocked for 5 mins or less, but thank you for caring about me :-))))")
		}
		c.Next()
	}
}
