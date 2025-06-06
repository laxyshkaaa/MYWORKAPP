package home

import (
	"alaricode/go-fiber/internal/vacancy"
	"alaricode/go-fiber/pkg/tadapter"
	"alaricode/go-fiber/views"
	"math"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/rs/zerolog"
)

type HomeHandler struct {
	router       fiber.Router
	customLogger *zerolog.Logger
	repository   *vacancy.VacancyRepository
	store        *session.Store
}

func NewHandler(router fiber.Router, customLogger *zerolog.Logger, repository *vacancy.VacancyRepository, store *session.Store) {
	h := &HomeHandler{
		router:       router,
		customLogger: customLogger,
		repository:   repository,
		store:        store,
	}
	h.router.Get("/", h.home)
	h.router.Get("/login", h.login)
	h.router.Get("/404", h.error)
	h.router.Get("/register", h.register)
}

func (h *HomeHandler) home(c *fiber.Ctx) error {
	PAGE_ITEMS := 2
	page := c.QueryInt("page", 1)
	sess, err := h.store.Get(c)
	if err != nil {
		panic(err)
	}
	authenticated := sess.Get("userID") != nil

	count := h.repository.GetCount()
	vacancies, err := h.repository.GetAll(PAGE_ITEMS, (page-1)*PAGE_ITEMS)
	if err != nil {
		h.customLogger.Error().Msg(err.Error())
		return c.SendStatus(500)
	}

	component := views.Main(
		vacancies,
		int(math.Ceil(float64(count)/float64(PAGE_ITEMS))),
		page,
		authenticated,
	)
	return tadapter.Render(c, component)
}

func (h *HomeHandler) login(c *fiber.Ctx) error {
	sess, err := h.store.Get(c)
	if err != nil {
		panic(err)
	}
	authenticated := sess.Get("userID") != nil
	component := views.Login(authenticated)
	return tadapter.Render(c, component)
}

func (h *HomeHandler) register(c *fiber.Ctx) error {
	sess, err := h.store.Get(c)
	if err != nil {
		panic(err)
	}
	authenticated := sess.Get("userID") != nil
	component := views.Register(authenticated)
	return tadapter.Render(c, component)
}

func (h *HomeHandler) error(c *fiber.Ctx) error {
	h.customLogger.Info().
		Bool("isAdmin", true).
		Str("email", "a@a.ru").
		Int("id", 10).
		Msg("Инфо")
	return c.SendString("Error")
}
