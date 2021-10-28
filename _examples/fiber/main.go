package main

import (
	"log"

	syserrors "errors"

	"github.com/gofiber/fiber/v2"
	"github.com/ovsinc/errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func processError(e error) error {
	if err, ok := e.(validation.Error); ok {
		return errors.New(
			err.Message(),
			errors.SetErrorType("validation"),
			errors.SetOperations(errors.DefaultCaller().FuncName),
			errors.SetContextInfo(errors.CtxMap(err.Params())),
			errors.AppendContextInfo("code", err.Code()),
			errors.AppendContextInfo("func", errors.DefaultCaller()),
		)
	}

	return errors.New(e.Error(), errors.SetOperations(errors.HandlerCaller().FuncName))
}

type Post struct {
	Name   string
	Gender string
}

func (p *Post) Validate() error {
	return errors.Combine(
		processError(validation.Validate(&p.Name, validation.Required, validation.Length(5, 50))),
		processError(validation.Validate(&p.Gender, validation.In("Female", "Male"))),
	)
}

func getPosts() ([]Post, error) {
	je1 := errors.New(
		"hello1",
		errors.SetErrorType("not found"),
		errors.SetOperations("write"),
		errors.SetSeverity(errors.SeverityError),
		errors.SetContextInfo(errors.CtxMap{"hello": "world", "my": "name"}),
	)

	je2 := errors.New(
		"hello2",
		errors.SetErrorType("not found"),
		errors.SetOperations("read"),
		errors.SetSeverity(errors.SeverityError),
		errors.SetContextInfo(errors.CtxMap{"hello2": "world", "my2": "name"}),
	)

	je3 := errors.New(
		"hello3",
		errors.SetErrorType("not found"),
		errors.SetOperations("read"),
		errors.SetSeverity(errors.SeverityError),
		errors.SetContextInfo(errors.CtxMap{"hello3": "world", "my3": "name"}),
	)

	return nil, errors.Combine(
		je1, je2, je3,
	)
}

func mapErrors(err error) []string {
	es := errors.Errors(err)
	ss := make([]string, 0, len(es))
	for _, e := range es {
		ss = append(ss, e.Error())
	}
	return ss
}

func myPost(c *fiber.Ctx) error {
	log.Printf("HandlerCaller: %v\n", errors.HandlerCaller())
	log.Printf("DefaultCaller: %v\n", errors.DefaultCaller())

	post := new(Post)
	if err := c.BodyParser(post); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	if err := post.Validate(); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	posts, err := getPosts() // your logic
	if err != nil {
		return c.Status(404).JSON(&fiber.Map{
			"success": false,
			"error":   mapErrors(err),
			"err1":    syserrors.New("hello err"),
		})
	}

	return c.JSON(&fiber.Map{
		"success": true,
		"posts":   posts,
	})
}

func main() {
	errors.DefaultMultierrFormatFunc = errors.JSONMultierrFuncFormat

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Post("/api/posts", myPost)

	log.Fatal(app.Listen(":3000"))
}
