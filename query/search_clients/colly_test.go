package search_clients

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testHTML = `<head>
	<meta charset="UTF-8"/>
	<meta http-equiv="X-UA-Compatible" content="IE=edge"/>
	<meta name="viewport" content="width=device-width, initial-scale=1.0"/>
	<meta http-equiv="Content-Type" content="text/html;charset=UTF-8"/>
	<title>CSX Labs: Launching ideas into cyberspace.</title>
  
	<!-- Google Meta Tags -->
	<link rel="canonical" href="https://csxlabs.org"/>
	<meta name="description" content="CSX Labs (Computer Science Exploration Laboratories) is a collection of open research and development laboratories dedicated to exploring, developing, and promoting unsaturated technologies for the advancement of humanity."/>
	<meta name="keywords" content="CSX Labs, Computer Science Exploration Laboratories, Open Research Laboratories, San Ramon, California, San Ramon Laboratories, Texas Computer Labs, Texas Research Centers, San Ramon Research Centers, Research Centers"/>
	<!-- Snip -->
	<style>
  
	</style>
  </head>
  <body>
	<nav class="navbar">
	<!-- Snip -->
	</nav>
	<main>
	  <div class="section hero">
	  <!-- Snip -->
	  </div>
	  <!-- About -->
	  <div class="section about" id="about">
		<h1 class="about-heading">About</h1>
		<p class="about-text">CSX Labs (Computer Science Exploration Laboratories) is a collection of open research and development laboratories dedicated to exploring, developing, and promoting unsaturated technologies for the advancement of humanity.</p>
	  </div>
	  <!-- Snip -->
	  <!-- Footer -->
	  <footer class="section footer" id="footer">
	  <!-- Snip -->
		<div class="footer-right">
		  <div class="footer-social">
		  <!-- Snip -->
			<a href="https://github.com/CSXL" target="_blank" class="footer-link">
			  <img src="assets/github.svg" alt="GitHub" class="footer-social-icon"/>
			</a>
		  </div>
		</div>
	  </footer>
	</main>
  </body>`
)

func TestNewScraper(t *testing.T) {
	NewScraper()
}

func TestScraper_ScrapePage(t *testing.T) {
	s := NewScraper()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, testHTML)
	}))
	page, err := s.ScrapePage(ts.URL)
	assert.Nil(t, err)
	assert.Equal(t, "CSX Labs: Launching ideas into cyberspace.", page.Title)
	assert.Equal(t, "text/html; charset=utf-8", page.MIME)
	assert.Equal(t, 1, len(page.Links))
	assert.Equal(t, "https://github.com/CSXL", page.Links[0])
}
