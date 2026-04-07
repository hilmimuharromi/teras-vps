package routes

import (
	"teras-vps/backend/controllers"
	"teras-vps/backend/middleware"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Setup configures all application routes
func Setup(app *fiber.App, db *gorm.DB, redis *redis.Client) {
	// Initialize controllers
	authController := controllers.NewAuthController(db)
	vmController := controllers.NewVMController(db)
	billingController := controllers.NewBillingController(db)
	userController := controllers.NewUserController(db)
	adminController := controllers.NewAdminController(db)
	supportController := controllers.NewSupportController(db)

	// Initialize middleware
	authMiddleware := middleware.Auth(db)
	adminMiddleware := middleware.AdminOnly(db)

	// API v1
	api := app.Group("/api/v1")

	// Health check (public)
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "TerasVPS API is running",
			"version": "1.0.0",
		})
	})

	// Auth routes (public)
	auth := api.Group("/auth")
	auth.Post("/register", authController.Register)
	auth.Post("/login", authController.Login)
	auth.Post("/logout", authController.Logout)
	auth.Get("/me", authMiddleware, authController.GetMe)

	// Protected routes (require authentication)
	protected := api.Group("")
	protected.Use(authMiddleware)

	// VM routes
	vm := protected.Group("/vms")
	vm.Get("", vmController.ListVMs)
	vm.Post("", vmController.CreateVM)
	vm.Get("/:id", vmController.GetVM)
	vm.Put("/:id", vmController.UpdateVM)
	vm.Delete("/:id", vmController.DeleteVM)
	vm.Post("/:id/start", vmController.StartVM)
	vm.Post("/:id/stop", vmController.StopVM)
	vm.Post("/:id/reboot", vmController.RebootVM)
	vm.Get("/:id/stats", vmController.GetVMStats)
	vm.Post("/:id/backup", vmController.CreateBackup)
	vm.Get("/:id/backups", vmController.ListBackups)

	// Billing routes
	billing := protected.Group("/billing")
	billing.Get("/invoices", billingController.ListInvoices)
	billing.Get("/invoices/:id", billingController.GetInvoice)
	billing.Post("/invoices/:id/pay", billingController.PayInvoice)
	billing.Get("/plans", billingController.ListPlans)

	// SSH Keys routes
	sshKeys := protected.Group("/ssh-keys")
	sshKeys.Get("", userController.ListSSHKeys)
	sshKeys.Post("", userController.AddSSHKey)
	sshKeys.Delete("/:id", userController.DeleteSSHKey)

	// User routes
	user := protected.Group("/user")
	user.Get("/profile", userController.GetProfile)
	user.Put("/profile", userController.UpdateProfile)
	user.Put("/password", userController.ChangePassword)

	// Support tickets routes
	support := protected.Group("/support")
	support.Get("/tickets", supportController.ListTickets)
	support.Post("/tickets", supportController.CreateTicket)
	support.Get("/tickets/:id", supportController.GetTicket)
	support.Post("/tickets/:id/messages", supportController.AddMessage)
	support.Post("/tickets/:id/close", supportController.CloseTicket)

	// Admin routes (require admin role)
	admin := api.Group("/admin")
	admin.Use(authMiddleware)
	admin.Use(adminMiddleware)
	admin.Get("/users", adminController.ListUsers)
	admin.Get("/vms", adminController.ListAllVMs)
	admin.Get("/stats", adminController.GetStats)
	admin.Post("/suspend-user/:id", adminController.SuspendUser)
	admin.Post("/unsuspend-user/:id", adminController.UnsuspendUser)
}
