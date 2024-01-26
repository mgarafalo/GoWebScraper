package main

import (
	"log"
	"strconv"

	"github.com/gocolly/colly"
	"github.com/gofiber/fiber/v2"
)

type IDPAClub struct {
	ClubName           string `json:"clubName"`
	PhoneNumber        string `json:"phoneNumber"`
	Location           string `json:"location"`
	Website            string `json:"website"`
	PrimaryClubContact string `json:"primaryClubContact"`
}

func ScrapeIdpa(c *colly.Collector, r *[]IDPAClub, pageNumber int) {
	c.OnHTML("table#results-table", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(_ int, el *colly.HTMLElement) {
			idpaClub := IDPAClub{
				ClubName:           el.ChildText("td > a > h3"),
				Location:           el.ChildText("td:nth-child(2) > p:nth-of-type(1)"),
				PhoneNumber:        el.ChildText("td > p:nth-of-type(2)"),
				Website:            el.ChildText("td:nth-child(3) > p > a"),
				PrimaryClubContact: "",
			}
			*r = append(*r, idpaClub)
		})
	})

	url := "https://www.idpa.com/clubs/page/" + strconv.Itoa(pageNumber) + "/?search-location&search-radius&search-word&search-id&search_country&search_state"
	c.Visit(url)
}

func main() {
	app := fiber.New()

	app.Get("/scrape", func(c *fiber.Ctx) error {
		var response []IDPAClub
		cly := colly.NewCollector()

		for i := 1; i <= 48; i++ {
			ScrapeIdpa(cly, &response, i)
		}

		return c.JSON(fiber.Map{"reults": response})
	})

	log.Fatal(app.Listen(":8080"))
}
