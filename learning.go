package main
import (
    "github.com/gocolly/colly"
    "fmt"
    "net/http"
    "crypto/tls"
	"log"
    "strings"
)
type company struct {
    Name string `selector:".address_details h5"`
    Logo string `selector:".company-list-details-logo img" attr:"src"`
    Address string `selector:".address_details p"`
    Website string `selector:".compant_cnt_details a[href]" attr:"href"`
    Phone string `selector:".company_phone"`
    Email string `selector:".company_email"`
}
type job_details struct {
    Title string 
    Details string
    Email string
}
type job struct {
    Role     string `selector:"div:first-child"`
    Link     string `selector:"div:first-child a[href]" attr:"href"`
	CompanyName  string `selector:".jobs-comp-name"`
	LastDate string `selector:".job-date"`
    Details job_details
    Email string
    Company company
}
func main(){
    c := colly.NewCollector()
    var jobs []job
	var currentJob job
    c.OnRequest(func(r *colly.Request) {
        fmt.Println("Visiting", r.URL)
    })
    i :=0
    c.OnError(func(_ *colly.Response, err error) {
        log.Println("Something went wrong:", err)
    })
    c.OnHTML(".contents", func(e *colly.HTMLElement) {
        // Extract company details
        urlEndpoint := e.Request.URL.String()
        if urlEndpoint == "https://xyz.in/companies/job-search"{
            return
        }
        var currentCompany company
        if err := e.Unmarshal(&currentCompany); err != nil {
            fmt.Println("Error extracting company details:", err)
            return
        }
        currentJob.Company = currentCompany
        // Define dynamic tags
        title := fmt.Sprintf("strong:contains('%s')", currentJob.Role)
        details := fmt.Sprintf("p:contains('%s')+p", currentJob.Role)
        email :=  fmt.Sprintf("p:contains('%s') +p + p:contains('Email')", currentJob.Role)
        jb := job_details{e.ChildText(title), strings.Split(e.ChildText(details), "Email:")[0], e.ChildText(email)}

        currentJob.Details = jb
    })
    c.OnHTML(".company-list.joblist", func(e *colly.HTMLElement) {
        // Extract job information
		if err := e.Unmarshal(&currentJob); err != nil {
			fmt.Println("Error extracting job information:", err)
			return
		}
        if i <10 {
            i +=1
        }else {
            return
        }
        c.Visit(currentJob.Link)
		// Append the current job entry to the slice
		jobs = append(jobs, currentJob)
        currentJob = job {}
	})
    c.OnResponse(func(r *colly.Response) {
        fmt.Println("\n Visited", r.Request.URL)
    })
    c.OnHTML("tr td:nth-of-type(1)", func(e *colly.HTMLElement) {
        fmt.Println("First column of a table row:", e.Text)
    })

    c.OnXML("//h1", func(e *colly.XMLElement) {
        fmt.Println(e.Text)
    })

    c.OnScraped(func(r *colly.Response) {
        fmt.Println("Finished", r.Request.URL)
    })
    c.WithTransport(&http.Transport{
        TLSClientConfig:&tls.Config{InsecureSkipVerify: true},
    })
    c.Visit("https://xyz.in/companies/job-search")
    for _, j := range jobs {
		fmt.Printf("\n\n Role: %s\n", j.Role)
		//fmt.Printf("Company: %s\n", j.Company)
		fmt.Printf("\n\n Last Date: %s\n", j.LastDate)
		fmt.Printf("\n\n Link: %s\n", j.Link)
		fmt.Printf("\n\n Company: %s\n", j.Company)
		fmt.Printf("\n\n Details: %s\n", j.Details)
		fmt.Println(strings.Repeat("-", 30)) // Just for separating entries visually
	}
}
