package api

import (
	"bytes"
	"fmt"
	"github.com/rs/zerolog/log"
	"html/template"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
)

func (a ApiController) Menu(res http.ResponseWriter, req *http.Request) {
	var (
		barUrls    []string
		tsUrls     []string
		mapUrls    []string
		statesUrls []string
	)

	barUrls = append(barUrls, "http://localhost:8080/oem-bar")
	barUrls = append(barUrls, "http://localhost:8080/test")
	tsUrls = append(tsUrls, "http://localhost:8080/oem-serie")

	for _, key := range a.keys {
		var (
			url_t  string
			params url.Values
		)
		params = url.Values{
			"key": {key},
		}
		url_t = "http://localhost:8080/oem-bar?" + params.Encode()
		barUrls = append(barUrls, url_t)
		url_t = "http://localhost:8080/oem-serie?" + params.Encode()
		tsUrls = append(tsUrls, url_t)
		url_t = "http://localhost:8080/map?" + params.Encode()
		mapUrls = append(mapUrls, url_t)
		url_t = "http://localhost:8080/state?" + params.Encode()
		statesUrls = append(statesUrls, url_t)
		log.Info().Msg("Stats")
		log.Info().Int("Keylen", len(key)).Msg("menu")
		fmt.Println("Menu Stats")
		fmt.Println("Keylen: ", len(key))
	}

	// Set the Content-Type header to text/html
	res.Header().Set("Content-Type", "text/html")

	// Write the HTML opening tags
	fmt.Fprintln(res, "<html><body>")
	fmt.Fprintln(res, "<h1>List of Bar plot </h1>")

	// Iterate over the URLs and create clickable links
	for _, url := range barUrls {
		fmt.Fprintf(res, `<a href="%s">%s</a><br>`, url, url)
	}

	fmt.Fprintln(res, "<h1>List of Time Series plot </h1>")
	// Iterate over the URLs and create clickable links
	for _, url := range tsUrls {
		fmt.Fprintf(res, `<a href="%s">%s</a><br>`, url, url)
	}

	fmt.Fprintln(res, "<h1>List of Map plot </h1>")
	// Iterate over the URLs and create clickable links
	for _, url := range mapUrls {
		fmt.Fprintf(res, `<a href="%s">%s</a><br>`, url, url)
	}

	fmt.Fprintln(res, "<h1>List of States plot </h1>")
	// Iterate over the URLs and create clickable links
	for _, url := range statesUrls {
		fmt.Fprintf(res, `<a href="%s">%s</a><br>`, url, url)
	}

	// Write the HTML closing tags
	fmt.Fprintln(res, "</body></html>")

}

func (a ApiController) MenuEcho(c echo.Context) error {
	var (
		barUrls    []string
		tsUrls     []string
		mapUrls    []string
		statesUrls []string
		res        *bytes.Buffer
	)

	log.Info().Msg("Menu Echo")
	fmt.Println("Menu Echo")

	res = new(bytes.Buffer)

	barUrls = append(barUrls, "http://localhost:8080/oem-bar")
	barUrls = append(barUrls, "http://localhost:8080/test")
	tsUrls = append(tsUrls, "http://localhost:8080/oem-serie")

	for _, key := range a.keys {
		var (
			url_t  string
			params url.Values
		)
		params = url.Values{
			"key": {key},
		}
		url_t = "http://localhost:8080/oem-bar?" + params.Encode()
		barUrls = append(barUrls, url_t)
		url_t = "http://localhost:8080/oem-serie?" + params.Encode()
		tsUrls = append(tsUrls, url_t)
		url_t = "http://localhost:8080/map?" + params.Encode()
		mapUrls = append(mapUrls, url_t)
		url_t = "http://localhost:8080/state?" + params.Encode()
		statesUrls = append(statesUrls, url_t)
		log.Info().Msg("Stats")
		log.Info().Int("Keylen", len(key)).Msg("menu")
		fmt.Println("Menu Stats")
		fmt.Println("Keylen: ", len(key))
	}

	// Set the Content-Type header to text/html
	// res.Header().Set("Content-Type", "text/html")

	// Write the HTML opening tags
	fmt.Fprintln(res, "<html><body>")
	fmt.Fprintln(res, "<h1>List of Bar plot </h1>")

	// Iterate over the URLs and create clickable links
	for _, url := range barUrls {
		fmt.Fprintf(res, `<a href="%s">%s</a><br>`, url, url)
	}

	fmt.Fprintln(res, "<h1>List of Time Series plot </h1>")
	// Iterate over the URLs and create clickable links
	for _, url := range tsUrls {
		fmt.Fprintf(res, `<a href="%s">%s</a><br>`, url, url)
	}

	fmt.Fprintln(res, "<h1>List of Map plot </h1>")
	// Iterate over the URLs and create clickable links
	for _, url := range mapUrls {
		fmt.Fprintf(res, `<a href="%s">%s</a><br>`, url, url)
	}

	fmt.Fprintln(res, "<h1>List of States plot </h1>")
	// Iterate over the URLs and create clickable links
	for _, url := range statesUrls {
		fmt.Fprintf(res, `<a href="%s">%s</a><br>`, url, url)
	}

	// Write the HTML closing tags
	fmt.Fprintln(res, "</body></html>")

	c.Blob(http.StatusOK, "text/html", res.Bytes())
	return nil

}

func (a ApiController) testMenuEcho(c echo.Context) error {
	var (
		barUrls    []string
		tsUrls     []string
		mapUrls    []string
		statesUrls []string
		res        *bytes.Buffer
	)

	//Refresh the key is more event has been added
	a.keys = a.srv.GetKeys()

	log.Info().Msg("testMenuEcho Echo")
	fmt.Println("testMenuEcho Echo")
	log.Info().Msg("Stats")
	log.Info().Int("Keylen", len(a.keys)).Msg("menu")
	fmt.Println("Menu Stats")
	fmt.Println("Keylen: ", len(a.keys))

	res = new(bytes.Buffer)

	barUrls = append(barUrls, "http://localhost:8080/oem-bar")
	barUrls = append(barUrls, "http://localhost:8080/test")
	tsUrls = append(tsUrls, "http://localhost:8080/oem-serie")

	for _, key := range a.keys {
		var (
			url_t  string
			params url.Values
		)
		params = url.Values{
			"key": {key},
		}
		url_t = "http://localhost:8080/oem-bar?" + params.Encode()
		barUrls = append(barUrls, url_t)
		url_t = "http://localhost:8080/oem-serie?" + params.Encode()
		tsUrls = append(tsUrls, url_t)
		url_t = "http://localhost:8080/map?" + params.Encode()
		mapUrls = append(mapUrls, url_t)
		url_t = "http://localhost:8080/state?" + params.Encode()
		statesUrls = append(statesUrls, url_t)
	}

	// Create an instance of TableData with the URL lists
	data := struct {
		BarUrls    []string
		TsUrls     []string
		MapUrls    []string
		StatesUrls []string
	}{
		BarUrls:    barUrls,
		TsUrls:     tsUrls,
		MapUrls:    mapUrls,
		StatesUrls: statesUrls,
	}

	// Define the HTML template
	tmpl := template.Must(template.New("table").Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Table with URL Lists</title>
    <style>
        table {
            width: 50%;
            border-collapse: collapse;
            margin: 20px auto;
            font-family: Arial, sans-serif;
        }

        th, td {
            border: 1px solid #dddddd;
            text-align: left;
            padding: 8px;
            vertical-align: top;
        }

        th {
            background-color: #f2f2f2;
        }

        tr:nth-child(even) {
            background-color: #f9f9f9;
        }

        tr:hover {
            background-color: #e0e0e0;
        }

        ul {
            padding-left: 20px;
            margin: 0;
        }

        li {
            margin-bottom: 5px;
        }
    </style>
</head>
<body>

<table>
    <tr>
        <th>Bar URLs</th>
        <th>TS URLs</th>
    </tr>
    <tr>
        <td>
            <ul>
                {{range .BarUrls}}
                <li><a href="{{.}}">{{.}}</a></li>
                {{end}}
            </ul>
        </td>
        <td>
            <ul>
                {{range .TsUrls}}
                <li><a href="{{.}}">{{.}}</a></li>
                {{end}}
            </ul>
        </td>
    </tr>
    <tr>
        <th>Map URLs</th>
        <th>States URLs</th>
    </tr>
    <tr>
        <td>
            <ul>
                {{range .MapUrls}}
                <li><a href="{{.}}">{{.}}</a></li>
                {{end}}
            </ul>
        </td>
        <td>
            <ul>
                {{range .StatesUrls}}
                <li><a href="{{.}}">{{.}}</a></li>
                {{end}}
            </ul>
        </td>
    </tr>
</table>

</body>
</html>
`))

	// Execute the template with the data
	if err := tmpl.Execute(res, data); err != nil {
		panic(err)
	}

	// Print the generated HTML (or you can write it to a file)
	// For this example, we'll just print it to the console
	c.Blob(http.StatusOK, "text/html", res.Bytes())
	return nil

}
