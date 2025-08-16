package main

import (
	"log"
	"net/http"
	"os"
	"sync"
	"time"
	"tvprog/pkg/filter"
	"tvprog/pkg/formatter"
	"tvprog/pkg/parser"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

type Server struct {
	programs filter.ProgramsByDate
	mu       sync.RWMutex
}

func NewServer() *Server {
	return &Server{
		programs: make(filter.ProgramsByDate),
	}
}

func (s *Server) updatePrograms() error {
	xmlData, err := parser.FetchXML()
	if err != nil {
		return err
	}

	tv, err := parser.ParseXML(xmlData)
	if err != nil {
		return err
	}

	filteredPrograms, err := filter.FilterPrograms(tv)
	if err != nil {
		return err
	}

	s.mu.Lock()
	s.programs = filteredPrograms
	s.mu.Unlock()

	return nil
}

func (s *Server) startDailyUpdates() {
	c := cron.New()
	c.AddFunc("0 4 * * *", func() {
		if err := s.updatePrograms(); err != nil {
			log.Printf("Failed to update programs: %v", err)

			time.Sleep(30 * time.Minute)
			if err := s.updatePrograms(); err != nil {
				log.Printf("Retry failed: %v", err)
			}
		}
	})
	c.Start()
}

func (s *Server) getProgramsOfTheDay() *filter.DayPrograms {
	paris, _ := time.LoadLocation("Europe/Paris")
	today := time.Now().In(paris).Format("2006-01-02")

	s.mu.RLock()
	todayPrograms := s.programs[today]
	s.mu.RUnlock()

	if todayPrograms == nil {
		return &filter.DayPrograms{
			Programs:     make(map[string]filter.FilteredProgram),
			ChannelOrder: []string{},
		}
	}

	return todayPrograms
}

func (s *Server) handleText(c *gin.Context) {
	dayPrograms := s.getProgramsOfTheDay()
	table := formatter.FormatTable(dayPrograms)

	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.String(http.StatusOK, table)
}

func (s *Server) handleJSON(c *gin.Context) {
	dayPrograms := s.getProgramsOfTheDay()
	c.JSON(http.StatusOK, dayPrograms.Programs)
}

func main() {
	server := NewServer()

	if err := server.updatePrograms(); err != nil {
		log.Printf("Warning: Failed to fetch initial programs: %v", err)
	}

	server.startDailyUpdates()

	app := gin.New()
	app.Use(gin.Recovery())
	if os.Getenv("GIN_MODE") != "release" {
		app.Use(gin.Logger())
	}

	app.GET("/", server.handleText)
	app.GET("/json", server.handleJSON)

	log.Println("API started ✓")
	log.Fatal(app.Run(":3000"))
}
