package server

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"io"
	"io/ioutil"
	"strings"

	"fmt"

	"github.com/slotix/dataflowkit/extract"
	"github.com/slotix/dataflowkit/paginate"
	"github.com/slotix/dataflowkit/parser"
	"github.com/slotix/dataflowkit/scrape"
	"github.com/slotix/dataflowkit/splash"
)

// ParseService provides operations on strings.
type ParseService interface {
	//	GetResponse(req splash.Request) (*splash.Response, error)
	Fetch(req splash.Request) (interface{}, error)
	ParseData(payload []byte) (io.ReadCloser, error)
	//	CheckServices() (status map[string]string)
}

type parseService struct {
}

//Fetch returns splash.Request
func (parseService) Fetch(req splash.Request) (interface{}, error) {
	//logger.Println(req)
	fetcher, err := scrape.NewSplashFetcher()

	if err != nil {
		logger.Println(err)
	}
	res, err := fetcher.Fetch(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (parseService) ParseData(payload []byte) (io.ReadCloser, error) {
	p, err := parser.NewParser(payload)
	if err != nil {
		return nil, err
	}
	fetcher, err := scrape.NewSplashFetcher()
	if err != nil {
		logger.Println(err)
	}
	pieces := []scrape.Piece{}
	pl := p.Payloads[0]
	selectors := []string{}
	names := []string{}
	for _, f := range pl.Fields {
		var extractor scrape.PieceExtractor
		params := make(map[string]interface{})
		if f.Extractor.Params != nil {
			params = f.Extractor.Params.(map[string]interface{})
		}
		switch f.Extractor.Type {
		//For Link type by default Two pieces with different Text and Attr="href" extractors will be added for field selector.
		case "link":
			extractor, err = extract.FillParams("text", params)
			if err != nil {
				logger.Println(err)
			}
			fName := fmt.Sprintf("%s_text", f.Name)
			pieces = append(pieces, scrape.Piece{
				Name:      fName,
				Selector:  f.Selector,
				Extractor: extractor,
			})
			names = append(names, fName)

			params["Attr"] = "href"
			extractor, err = extract.FillParams("attr", params)
			if err != nil {
				logger.Println(err)
			}
			fName = fmt.Sprintf("%s_link", f.Name)
			pieces = append(pieces, scrape.Piece{
				Name:      fName,
				Selector:  f.Selector,
				Extractor: extractor,
			})
			names = append(names, fName)
			//Add selector just one time for link type
			selectors = append(selectors, f.Selector)
		
		//For image type by default Two pieces with different Attr="src" and Attr="alt" extractors will be added for field selector.
		case "image":
			params["Attr"] = "src"
			extractor, err = extract.FillParams("attr", params)
			if err != nil {
				logger.Println(err)
			}
			fName := fmt.Sprintf("%s_src", f.Name)
			pieces = append(pieces, scrape.Piece{
				Name:      fName,
				Selector:  f.Selector,
				Extractor: extractor,
			})
			names = append(names, fName)

			params["Attr"] = "alt"
			extractor, err = extract.FillParams("attr", params)
			if err != nil {
				logger.Println(err)
			}
			fName = fmt.Sprintf("%s_alt", f.Name)
			pieces = append(pieces, scrape.Piece{
				Name:      fName,
				Selector:  f.Selector,
				Extractor: extractor,
			})
			names = append(names, fName)
			//Add selector just one time for link type
			selectors = append(selectors, f.Selector)
		default:
			extractor, err = extract.FillParams(f.Extractor.Type, params)
			if err != nil {
				logger.Println(err)
			}

			pieces = append(pieces, scrape.Piece{
				Name:      f.Name,
				Selector:  f.Selector,
				Extractor: extractor,
			})

			selectors = append(selectors, f.Selector)
			names = append(names, f.Name)

		}
		//	if f.Extractor.Type == "link" {

		//	} else {

		//	}
	}

	paginator := pl.Paginator
	config := &scrape.ScrapeConfig{
		Fetcher: fetcher,
		//DividePage: scrape.DividePageBySelector(".p"),
		DividePage: scrape.DividePageByIntersection(selectors),
		Pieces:     pieces,
		//Paginator: paginate.BySelector(".next", "href"),
		Paginator: paginate.BySelector(paginator.Selector, paginator.Attribute),
		Opts:      scrape.ScrapeOptions{MaxPages: paginator.MaxPages, Format: p.Format},
	}
	scraper, err := scrape.New(config)
	if err != nil {
		return nil, err
	}
	req := splash.Request{URL: pl.URL}
	results, err := scraper.ScrapeWithOpts(req, config.Opts)
	if err != nil {
		return nil, err
	}
	//logger.Println(results.Results[0][0])
	var buf bytes.Buffer
	switch config.Opts.Format {
	case "json":
		json.NewEncoder(&buf).Encode(results)
	case "csv":
		includeHeader := true
		w := csv.NewWriter(&buf)
		for i, page := range results.Results {
			if i != 0 {
				includeHeader = false
			}
			err = encodeCSV(names, includeHeader, page, ",", w)
			if err != nil {
				logger.Println(err)
			}
		}
		w.Flush()
	}
	//	logger.Println(string(b))
	readCloser := ioutil.NopCloser(bytes.NewReader(buf.Bytes()))
	return readCloser, nil

	/*
		res, err := formatResults(results, config.Opts.Format)
		if err !=nil {
			return nil, err
		}
		return res, nil*/
}

/*
func formatResults(res *scrape.ScrapeResults, format string)(io.ReadCloser, error){
	var buf bytes.Buffer
	switch format {
	case "json":
		json.NewEncoder(&buf).Encode(res)
	case "csv":
	includeHeader := true
		w := csv.NewWriter(&buf)
		for i, page := range res.Results {
			if i != 0 {
				includeHeader = false
			}
			err = encodeCSV(names, includeHeader, page, ",", w)
			if err != nil {
				logger.Println(err)
			}
		}
		w.Flush()
	}
	//	logger.Println(string(b))
	readCloser := ioutil.NopCloser(bytes.NewReader(buf.Bytes()))
	return readCloser, nil

}
*/
//encodeCSV writes data to w *csv.Writee.
//header - headers for csv.
//includeHeader include headers or not.
//rows - csv records to be written.
func encodeCSV(header []string, includeHeader bool, rows []map[string]interface{}, comma string, w *csv.Writer) error {
	if comma == "" {
		comma = ","
	}
	w.Comma = rune(comma[0])
	//Add Header string to csv or no
	if includeHeader {
		if err := w.Write(header); err != nil {
			return err
		}
	}
	r := make([]string, len(header))
	for _, row := range rows {
		for i, column := range header {
			switch v := row[column].(type) {
			case string:
				r[i] = v
			case []string:
				r[i] = strings.Join(v, ";")
			case nil:
				r[i] = ""
			}
		}
		if err := w.Write(r); err != nil {
			return err
		}
	}
	return nil
}

/*
func (parseService) GetResponse(req splash.Request) (*splash.Response, error) {
	splashURL, err := splash.NewSplashConn(req)
	response, err := splash.GetResponse(splashURL)
	return response, err
}
*/

/*
func (parseService) Fetch(req splash.Request) (io.ReadCloser, error) {
	logger.Println(req)
	splashURL, err := splash.NewSplashConn(req)
	content, err := splash.Fetch(splashURL)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func (parseService) ParseData_old(payload []byte) (io.ReadCloser, error) {
	p, err := parser.NewParser(payload)
	if err != nil {
		return nil, err
	}
	res, err := p.MarshalData()
	if err != nil {
		logger.Println(res, err)
		return nil, err
	}
	return res, nil
}
*/
//func (parseService) CheckServices() (status map[string]string) {
//	return CheckServices() //, allAlive
//}

// ServiceMiddleware is a chainable behavior modifier for ParseService.
type ServiceMiddleware func(ParseService) ParseService
