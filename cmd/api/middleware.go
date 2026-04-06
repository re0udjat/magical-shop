package main

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/re0udjat/magic-shop/internal/data"
	"github.com/re0udjat/magic-shop/internal/validator"
	"golang.org/x/time/rate"
)

func (app *app) recover() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a deferred function (which will always be run in the event of a panic as Go unwind the stack)
		defer func() {
			// Builtin recover function to check if there has been a panic or not
			if err := recover(); err != nil {
				// Set a "Connection: close" header on the response to make Go's HTTP server
				// automatically close the current connection after a response has been sent
				c.Header("Connection", "close")
				app.serverErrorResponse(c, fmt.Errorf("%s", err))
			}
		}()
		c.Next()
	}
}

func (app *app) rateLimit() gin.HandlerFunc {
	// Define a client struct to hold the rate limiter and last seen time for each client
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	// Launch a background goroutine which removes old entries from the clients map once
	// every minute
	go func() {
		for {
			time.Sleep(time.Minute)

			// Lock the mutext to prevent any rate limiter checks from happening while the
			// cleanup is taking place
			mu.Lock()

			// Iterate through the clients map and remove any entries that haven't been seen
			// in the last 3 minutes
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}

			// Unlock the mutex so that other goroutines can access the map
			mu.Unlock()
		}
	}()

	return func(c *gin.Context) {
		// Extract the client's IP address from the request
		if app.config.limiter.enabled {
			ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
			if err != nil {
				app.serverErrorResponse(c, err)
				return
			}

			// Lock the mutex to prevent this code from being executed by more than one goroutine at a time
			mu.Lock()

			// Check if the IP address already exists in the map. If it doesn't, then initialize a new rate
			// a new rate limiter and add the IP address and limiter to the map
			if _, found := clients[ip]; !found {
				clients[ip] = &client{limiter: rate.NewLimiter(rate.Limit(app.config.limiter.rps), app.config.limiter.burst)}
			}

			// Update the last seen time for the client
			clients[ip].lastSeen = time.Now()

			// Call the Allow() method on the limiter for the current IP address
			// If the request isn't allowed, unlock the mutex and send a 429 Too Many Requests response
			if !clients[ip].limiter.Allow() {
				mu.Unlock()
				app.rateLimitExceededResponse(c)
				return
			}

			// Unlock the mutex so that other goroutines can access the map
			mu.Unlock()
		}
		c.Next()
	}
}

func (app *app) authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Add the "Vary: Authorization" header to the response
		// Indicates to any caches that the response may vary based on the value of the
		// Authorization header in the request
		c.Header("Vary", "Authorization")

		// Get value of Authorization header from the request
		authorizationHeader := c.GetHeader("Authorization")

		// If there is no Authorization header found, add the AnnonymousUser to the request context
		// Then call the next handler in the chain and return without executing any of the code below
		if authorizationHeader == "" {
			c.Request = app.contextSetUser(c.Request, data.AnnonymousUser)
			c.Next()
			return
		}

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			app.invalidAuthenticationTokenResponse(c)
			c.Abort()
			return
		}

		// Extract the token from the header
		token := headerParts[1]

		// Validate the token
		v := validator.New()

		if data.ValidateTokenPlaintext(v, token); !v.Valid() {
			app.invalidAuthenticationTokenResponse(c)
			c.Abort()
			return
		}

		// Retrieve the details of the user associated with the authentication token
		user, err := app.models.Users.GetForToken(data.ScopeAuthentication, token)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				app.invalidAuthenticationTokenResponse(c)
			default:
				app.serverErrorResponse(c, err)
			}
			c.Abort()
			return
		}

		// Add the user to the request context
		c.Request = app.contextSetUser(c.Request, user)

		// Call the next handler in the chain
		c.Next()
	}
}
