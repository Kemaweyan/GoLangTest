package main

import(
        "github.com/gin-gonic/gin"
        "io/ioutil"
        "strings"
        "encoding/json"
        "net/http"
)

func main() {
    r := gin.Default()

    r.POST("/checkText", func(c *gin.Context) {
        //read a request body
        x, err := ioutil.ReadAll(c.Request.Body)

        if err != nil {
            //oops! can't read the request body
            //it looks like an internal server error
            c.Writer.WriteHeader(http.StatusInternalServerError)
            return
        }

        //a request struct
        type Request struct {
            Site []string
            SearchText string
        }

        //create a JSON decoder for the request
        dec := json.NewDecoder(strings.NewReader(string(x)))

        var req Request //here the requst data will be stored

        //decode the request
        if err := dec.Decode(&req); err != nil {
            //can't decode the request
            //probably the request is incorrect
            c.Writer.WriteHeader(http.StatusBadRequest)
            return
        }

        //a response structure
        type Response struct {
            FoundAtSite string
        }

        for _, site := range req.Site {
            //get a page from each site in the request
            site_resp, err := http.Get(site)

            if err != nil || site_resp.StatusCode != http.StatusOK {
                //something went wrong, the site is unavailable
                //just skip this site
                continue
            }

            //read the page
            page, err := ioutil.ReadAll(site_resp.Body)
            //close the response body immediately
            //there is no reason to defer this call
            site_resp.Body.Close()

            if err != nil {
                //can't read the page
                //just skip this site
                continue
            }

            //check the SearchText on the page
            if strings.Contains(string(page), req.SearchText) {
                //create a response and send it
                resp := Response{FoundAtSite: site}
                c.JSON(http.StatusOK, resp)
                //the text has been found, get out from the function
                return
            }
        }
        //the SearchText has not been found on specified sites
        //send the NoContent status
        c.Writer.WriteHeader(http.StatusNoContent)
    })
    r.Run()
}

