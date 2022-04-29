package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

var (
	urlBaseInfoSociete = "https://www.societe.com/cgi-bin/search?champs=%s"
)

type infoSociete struct {
	Siret       string
	Category    string
	CompanySize int
}

func getInfoCompany(societeName string) []string {
	log.Printf("get societe info for %s\n", societeName)

	response := make([]string, 4)

	c := colly.NewCollector(
		colly.AllowedDomains(
			"https://www.societe.com/",
			"https://www.societe.com",
			"http://www.societe.com/",
			"www.societe.com/",
			"societe.com/",
			"www.societe.com",
			"societe.com"),
		colly.Async(true),
	)

	c.Limit(&colly.LimitRule{
		RandomDelay: 5 * time.Second,
	})

	c2 := c.Clone()

	log.Println(c2.AllowedDomains)

	c.OnHTML("div[id=result_deno_societe]", func(e *colly.HTMLElement) {

		firstLink := e.ChildAttr("a", "href")
		urlToVisit := fmt.Sprintf("https://www.societe.com%s", firstLink)
		log.Printf("visiting %s", urlToVisit)
		if err := c2.Visit(urlToVisit); err != nil {
			log.Fatal(err)
		}

	})

	c2.OnHTML("div[id=tabledir]", func(div *colly.HTMLElement) {
		div.ForEach("div table", func(j int, table *colly.HTMLElement) {
			id := table.Attr("id")
			match, err := regexp.MatchString(`^dir\d$`, id)
			if err != nil {
				log.Fatal(err)
			}

			if match {
				table.ForEach("tbody tr td", func(i int, td *colly.HTMLElement) {
					if i == 1 {
						directorName := strings.TrimSpace(td.Text)
						response[3] = directorName
					}
				})

			}
		})
	})

	c2.OnHTML("table[id=rensjur]", func(table *colly.HTMLElement) {

		table.ForEach("body tr", func(i int, tr *colly.HTMLElement) {
			switch i {
			case 9:
				siret := getSiret(tr)
				response[0] = siret
			case 13:
				category := getCategory(tr)
				response[1] = category
			case 21:
				companySize := getCompanySize(tr)
				response[2] = companySize
			}
		})
	})
	url := fmt.Sprintf(urlBaseInfoSociete, societeName)
	c.Visit(url)

	c.Wait()
	c2.Wait()
	return response
}

func getSiret(tr *colly.HTMLElement) string {
	siret := ""
	tr.ForEach("td", func(j int, tr *colly.HTMLElement) {
		switch j {
		case 0:
			if tr.Text != "Numéro SIRET (siège)" {
				log.Fatalf("should be Numéro SIRET (siège) but was %s", tr.Text)
			}

		case 1:
			npSpaces := strings.ReplaceAll(tr.Text, " ", "")
			noJumpLine := strings.ReplaceAll(npSpaces, "\n", "")

			//siret is 14 characters number source : https://fr.wikipedia.org/wiki/Syst%C3%A8me_d%27identification_du_r%C3%A9pertoire_des_%C3%A9tablissements
			siret = noJumpLine[:14]
		}
	})
	return siret
}

func getCategory(tr *colly.HTMLElement) string {

	noSpaces := strings.ReplaceAll(tr.Text, " ", "")

	arr := strings.Split(noSpaces, "\n")
	if len(arr) >= 2 {
		return arr[2]
	}

	return ""
}

func getCompanySize(tr *colly.HTMLElement) string {
	return tr.ChildText("#trancheeff-histo-description")
}
