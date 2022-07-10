package server

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

const (
	ASYNC = "async"
	SYNC  = "sync"
)

type Server struct {
	app         *fiber.App
	taskManager *TaskManager
	done        chan struct{}
}

func NewServer() *Server {
	taskManager := NewTaskManager()
	app := fiber.New()

	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path} body:${body}\n",
	}))

	app.Get("/schedule", func(c *fiber.Ctx) error {
		return c.JSON(taskManager.Shedule())
	})
	app.Get("/time", func(c *fiber.Ctx) error {
		return c.SendString(taskManager.Time().String())
	})

	app.Post("/add", func(c *fiber.Ctx) error {
		mod := c.Query("mod")

		dur, err := time.ParseDuration(string(c.Body()))
		if err != nil {
			log.Println(err)
			return c.Status(fiber.StatusBadRequest).SendString("invalud body format")
		}

		task := Task{
			Dur:      dur,
			TaskDone: make(chan struct{}, 1),
		}

		switch mod {
		case ASYNC:
			taskManager.AddTask(task)
			return c.SendString("task added")
		case SYNC:
			<-taskManager.AddTask(task).TaskDone
			return c.SendString("task done")
		default:
			return c.Status(fiber.StatusBadRequest).SendString("invalid mod")
		}
	})

	return &Server{
		app:         app,
		taskManager: taskManager,
		done:        make(chan struct{}, 1000),
	}
}

func (s *Server) Run() error {
	go func() {
		for {
			select {
			case <-s.taskManager.newTask:
				log.Println("got new task")
				s.taskManager.ProccessTask()
				log.Println("proccess end")
			case <-s.done:
				return
			}
		}
	}()
	return s.app.Listen(":5000")
}

func (s *Server) Shutdown() {
	s.app.Shutdown()
	s.done <- struct{}{}
}
