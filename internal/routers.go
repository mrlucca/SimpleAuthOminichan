package internal

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func StartServer() {
	app := fiber.New()

	app.Post("/users", createUser)
	app.Post("/users/validate", checkUserIsValid)
	app.Get("/users", listUsers)

	app.Listen("0.0.0.0:80")
}

func createUser(c *fiber.Ctx) error {
	user := new(UserCreate)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}
	log.Println("new user")
	log.Println(user)

	if user.Email == "" || user.Pass == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Email and password are required",
		})
	}

	if UserExistsFromEmail(user.Email) {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": "User already exists",
		})
	}

	CreateUser(*user)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User created successfully",
	})
}

func checkUserIsValid(c *fiber.Ctx) error {
	userValidation := new(UserValidation)
	if err := c.BodyParser(userValidation); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	if UserIsValid(userValidation.Email, userValidation.Password) {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "User is valid",
		})
	}

	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"message": "Invalid email or password",
	})
}
func listUsers(c *fiber.Ctx) error {
	var users []User = GetAllUsers()
	return c.JSON(users)
}
