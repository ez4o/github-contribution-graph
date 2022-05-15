package server

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"path"
	"time"

	"server/model"
	"server/util"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/patrickmn/go-cache"
)

type Server struct {
	c           *cache.Cache
	GitHubToken string
}

func (s *Server) handleRender(w http.ResponseWriter, r *http.Request) {
	params := util.GetRequestParams(r.URL.Query())
	if params.Username == "" {
		fmt.Fprintln(w, "Url param 'username' is missing.")

		return
	}

	log.Println("New render request.")
	log.Printf("username=%s\n", params.Username)

	if params.ImgUrl == "" {
		params.ImgUrl = "https://images.unsplash.com/photo-1518791841217-8f162f1e1131?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=800&q=60"
	} else {
		imgUrl, err := url.QueryUnescape(params.ImgUrl)
		if err != nil {
			fmt.Fprintln(w, "Error while parsing image url.")

			return
		}

		params.ImgUrl = imgUrl
	}

	contributiondata, username, err := util.GetContributionData(s.c, params.Username, s.GitHubToken, params.LastNDays)
	if err != nil {
		fmt.Fprintln(w, "Error while parsing response from GitHub, please check all your request parameters.")

		return
	}

	image, err := util.GetImageFromUrl(s.c, params.ImgUrl)
	if err != nil {
		fmt.Fprintln(w, "Error while parsing image:", err.Error())

		return
	}

	tmpl, err := template.New("index.html").Funcs(
		template.FuncMap{
			"hydrateContributionData": util.HydrateContributionData,
		},
	).ParseFiles(
		path.Join("dist", "index.html"),
	)
	if err != nil {
		fmt.Fprintln(w, "Error while parsing template.")

		return
	}

	w.Header().Set("Content-Type", "text/html")

	tmpl.Execute(w, &model.TemplateData{
		ContributionData: contributiondata,
		Username:         username,
		ImgType:          image.MimeType,
		ImgBase64String:  image.Base64String,
	})
}

func (s *Server) sendSVG(svgContent string, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "image/svg+xml; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=600")
	w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; font-src 'self' data:; script-src 'self' 'unsafe-inline' 'unsafe-eval'")
	w.WriteHeader(http.StatusOK)

	fmt.Fprintln(w, svgContent)
	w.(http.Flusher).Flush()
}

func (s *Server) handleSVGInBackground(r *http.Request, b *rod.Browser) {
	log.Println("New SVG request in background.")

	page := b.MustPage(fmt.Sprintf("http://localhost:8687/?%s", r.URL.Query().Encode()))
	defer page.Close()
	page.MustWaitLoad()

	if page.MustHas("#svg-container") {
		svg := page.MustElement("#svg-container")
		svgContent := svg.MustHTML()

		s.c.Set("cachedSvgContent:"+r.URL.Query().Encode(), svgContent, cache.DefaultExpiration)
		s.c.Set("cachedSvgContent:old:"+r.URL.Query().Encode(), svgContent, 15778463000000000)
	}
}

func (s *Server) handleSVG(w http.ResponseWriter, r *http.Request, b *rod.Browser) {
	if svgContent, found := s.c.Get("cachedSvgContent:" + r.URL.Query().Encode()); found {
		log.Println("SVG found in cache.")

		s.sendSVG(svgContent.(string), w)

		return
	}

	if svgContent, found := s.c.Get("cachedSvgContent:old:" + r.URL.Query().Encode()); found {
		log.Println("SVG found in old cache.")

		s.c.Delete("cachedSvgContent:old:" + r.URL.Query().Encode())
		s.sendSVG(svgContent.(string), w)

		go s.handleSVGInBackground(r, b)

		return
	}

	log.Println("New SVG request.")

	page := b.MustPage(fmt.Sprintf("http://localhost:8687/?%s", r.URL.Query().Encode()))
	defer page.Close()
	page.MustWaitLoad()

	if page.MustHas("pre") {
		pre := page.MustElement("pre")

		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Cache-Control", "no-store, max-age=0")
		w.WriteHeader(http.StatusInternalServerError)

		fmt.Fprintln(w, pre.MustText())

		return
	} else if page.MustHas("#svg-container") {
		svg := page.MustElement("#svg-container")
		svgContent := svg.MustHTML()

		s.sendSVG(svgContent, w)

		s.c.Set("cachedSvgContent:"+r.URL.Query().Encode(), svgContent, cache.DefaultExpiration)
		s.c.Set("cachedSvgContent:old:"+r.URL.Query().Encode(), svgContent, 15778463000000000)

		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Cache-Control", "no-store, max-age=0")
	w.WriteHeader(http.StatusInternalServerError)

	fmt.Fprintln(w, "Error while parsing SVG.")
}

func (s *Server) Start() {
	s.c = cache.New(30*time.Minute, 1*time.Minute)

	ssrServer := http.NewServeMux()
	ssrServer.HandleFunc("/", s.handleRender)
	go http.ListenAndServe(":8687", ssrServer)

	l := launcher.New()
	l.Set("no-sandbox")
	headlessBrowserUrl := l.MustLaunch()
	browser := rod.New().ControlURL(headlessBrowserUrl).MustConnect()

	svgServer := http.NewServeMux()
	svgServer.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if len(r.URL.String()) > 2 {
			if r.URL.String()[1] == '?' {
				s.handleSVG(w, r, browser)

				return
			}
		}

		http.NotFound(w, r)
	})
	http.ListenAndServe(":8686", svgServer)
}
