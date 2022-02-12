package smg

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"
)

const (
	baseURL     = "https://www.example.com"
	n           = 5
	letterBytes = "////abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

var (
	lenLetters = len(letterBytes)
	routes     []string
)

func buildRoutes(n int) {
	rand.Seed(time.Now().UnixNano())

	routes = make([]string, n)
	for i := range routes {
		routes[i] = randString(rand.Intn(40) + 10)
	}
}

func randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(lenLetters)]
	}
	return string(b)
}

// TestCompleteAction tests the whole sitemap-generator module with a semi-basic usage
func TestCompleteAction(t *testing.T) {
	buildRoutes(10)
	randNum := rand.Intn(900) + 100
	path := fmt.Sprintf("/tmp/sitemap_output_%d", randNum)

	smi := NewSitemapIndex(true)
	smi.SetCompress(false)
	smi.SetHostname(baseURL)
	smi.SetSitemapIndexName("bomt_sitemap")
	smi.SetOutputPath(path)
	now := time.Now().UTC()

	// Testing a list of named sitemaps
	a := []string{"test_sitemap1", "test_sitemap2", "test_sitemap3", "test_sitemap4", "test_sitemap5"}
	for _, name := range a {
		sm := smi.NewSitemap()
		sm.SetName(name)
		for _, route := range routes {
			err := sm.Add(&SitemapLoc{
				Loc:        route,
				LastMod:    &now,
				ChangeFreq: Always,
				Priority:   0.4,
			})
			if err != nil {
				t.Fatal("Unable to add SitemapLoc:", name, err)
			}
		}
	}

	// Testing another one with autogenerated name
	sm := smi.NewSitemap()
	for _, route := range routes {
		err := sm.Add(&SitemapLoc{
			Loc:        route,
			LastMod:    &now,
			ChangeFreq: Daily,
			Priority:   0.8,
		})
		if err != nil {
			t.Fatal("Unable to add 6th SitemapLoc:", err)
		}
	}

	err := smi.Save()
	if err != nil {
		t.Fatal("Unable to Save SitemapIndex:", err)
	}

	err = smi.PingSearchEngines()
	if err != nil {
		t.Fatal("Unable to Ping search engines:", err)
	}

	// Checking 5 named output files
	for _, name := range a {
		f, err := os.Stat(filepath.Join(path, name+".xml"))
		if os.IsNotExist(err) || f.IsDir() {
			t.Fatal("Final file does not exist or is directory:", name, err)
		}
		if f.Size() == 0 {
			t.Fatal("Final file has zero size:", name)
		}
	}

	// Checking the 6th sitemap which was no-name
	f, err := os.Stat(filepath.Join(path, "sitemap6.xml"))
	if os.IsNotExist(err) || f.IsDir() {
		t.Fatal("Final 6th file does not exist or is directory:", err)
	}
	if f.Size() == 0 {
		t.Fatal("Final 6th file has zero size")
	}

	// Removing the generated path and files
	err = os.RemoveAll(path)
	if err != nil {
		t.Fatal("Unable to remove tmp path after testing:", err)
	}
}
