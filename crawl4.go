package main

import (
	"database/sql"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func ConnDB() (*sql.DB, error) {

	//access database
	//Note : adjust your user and password access. for this code, we use user : root, password : password
	db, err := sql.Open("mysql", "root:stasiun@tcp(localhost:3306)/dropship")

	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func main() {

	// get connection database
	db, err := ConnDB()
	if err != nil {
		log.Panic(err)
	}

	//prepare query to insert data
	stmt, err := db.Prepare("replace INTO source (sumber,username,namatoko,link,level,feedback,alamat,joindate,lastlogin,tolak,deliveryrespon,pelanggan,deskripsi,catatanpelapak)VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?);")
	if err != nil {
		log.Fatal(err)
	}

	c := colly.NewCollector()

	// Allow requests only to store.xkcd.com
	c.AllowedDomains = []string{"www.bukalapak.com"}

	detail := c.Clone()

	// Extract product details
	c.OnHTML(".user__name", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.ChildAttr("a", "href"))
		detail.Visit(link)
	})

	detail.OnHTML("section[id=display_user]", func(e *colly.HTMLElement) {

		stringDropship := ""
		e.DOM.Find("#user_term_condition").Each(func(index int, item *goquery.Selection) {
			stringDropship = stringDropship + " " + item.Text()
		})
		e.DOM.Find(".user-header__description").Each(func(index int, item *goquery.Selection) {
			stringDropship = stringDropship + " " + item.Text()
		})
		if strings.Contains(strings.ToLower(stringDropship), "dropship") || strings.Contains(strings.ToLower(stringDropship), "reseller") || strings.Contains(strings.ToLower(stringDropship), "grosir") {

			username := strings.TrimSpace(e.DOM.Find(".user-description .user__username strong").Text())
			namatoko := e.DOM.Find(".user-description h5 a").Text()
			link := "https://bukalapak.com/u/" + username
			level := e.DOM.Find(".user-description .user__level").Text()
			feedback := e.DOM.Find(".user-description .user-feedback-summary").Text()
			alamat := strings.TrimSpace(e.DOM.Find(".user-location .user-address").Text())
			join := e.DOM.Find(".user-meta-join-at").Text()
			lastlogin := strings.TrimSpace(e.DOM.Find(".user-meta-last-login").Text())
			tolak := strings.TrimSpace(e.DOM.Find(".user-meta-rejection-rate").Text())
			deliveryrespon := strings.TrimSpace(e.DOM.Find(".user-meta-delivery-response").Text())
			pelanggan := strings.Replace(strings.TrimSpace(e.DOM.Find(".user-meta-subscribers-total").Text()), "Memiliki ", "", 1)
			pelanggan = strings.Replace(pelanggan, " pelanggan", "", 1)
			description := strings.TrimSpace(e.DOM.Find(".user-header__description").Text())
			catatanpelapak := strings.TrimSpace(e.DOM.Find(".c-seller-tnc").Text())

			fmt.Printf("Name:%s, Username:%s \n", namatoko, username)
			_, err = stmt.Exec("Bukalapak", &username, &namatoko, &link, &level, &feedback, &alamat, &join, &lastlogin, &tolak, &deliveryrespon, &pelanggan, &description, &catatanpelapak)

			if err != nil {
				print(err.Error())
			}
		}
	})
	detail.OnRequest(func(r *colly.Request) {
	})

	// Find and visit next page links
	c.OnHTML(`.next_page`, func(e *colly.HTMLElement) {
		fmt.Print("Visiting ", e.Attr("href"), "\n")
		e.Request.Visit(e.Attr("href"))
	})

	c.Visit("https://www.bukalapak.com/products/s")

	// Display collector's statistics
	log.Println(c)
}
