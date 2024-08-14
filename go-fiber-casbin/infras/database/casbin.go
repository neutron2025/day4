package database

import (
	"fmt"
	"log"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
)

func Casbin() *casbin.Enforcer {
	// Initialize casbin adapter
	adapter, err := gormadapter.NewAdapterByDB(adminDB)
	if err != nil {
		panic(fmt.Sprintf("failed to initialize casbin adapter: %v", err))
	}

	// Load model configuration file and policy store adapter
	e, err := casbin.NewEnforcer("config/restful_rbac_model.conf", adapter)
	if err != nil {
		panic(fmt.Sprintf("failed to create casbin enforcer: %v", err))
	}

	// Add policy - One-time run
	// Add policy for admin - One-time run
	hasPolicy, err := e.HasPolicy("admin", "/api/admin/*", "(GET)|(POST)|(PUT)|(DELETE)")
	if err != nil {
		// Handle error
		log.Fatal("Error checking policy for admin:", err)
	}
	if !hasPolicy {
		success, err := e.AddPolicy(
			"admin",
			"/api/admin/*",
			"(GET)|(POST)|(PUT)|(DELETE)",
		)
		if err != nil {
			// Handle error when adding policy
			log.Fatal("Error adding policy for admin:", err)
		}
		if !success {
			log.Fatal("Failed to add policy for admin")
		}
	}

	// Add policy for user - One-time run
	hasPolicy, err = e.HasPolicy("user", "/api/users/:id/*", "(GET)|(PUT)") // 使用 err 而不是 errs
	if err != nil {
		// Handle error
		log.Fatal("Error checking policy for user:", err)
	}
	if !hasPolicy {
		success, err := e.AddPolicy(
			"user",
			"/api/users/:id/*",
			"(GET)|(PUT)",
		)
		if err != nil {
			// Handle error when adding policy
			log.Fatal("Error adding policy for user:", err)
		}
		if !success {
			log.Fatal("Failed to add policy for user")
		}
	}

	// Load the policies
	e.LoadPolicy()
	return e
}
