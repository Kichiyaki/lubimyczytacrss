package internal

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	lubimyCzytacBaseURL = "https://lubimyczytac.pl"
)

var (
	ErrAuthorNotFound   = errors.New("author not found")
	ErrUnexpectedStatus = errors.New("unexpected http status was returned by lubimyczytac.pl")
)

type Book struct {
	Title string
	URL   string
}

type Author struct {
	Name             string
	ShortDescription string
	URL              string
	Books            []Book
}

type LubimyCzytacClient struct {
	http *http.Client
}

func NewLubimyCzytacClient(client *http.Client) *LubimyCzytacClient {
	return &LubimyCzytacClient{http: client}
}

func (c *LubimyCzytacClient) GetAuthor(ctx context.Context, authorID string) (Author, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/autor/%s/x", lubimyCzytacBaseURL, authorID),
		nil,
	)
	if err != nil {
		return Author{}, fmt.Errorf("http.NewRequestWithContext: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return Author{}, fmt.Errorf("httpClient.Do: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		_, _ = io.Copy(io.Discard, resp.Body)
		if resp.StatusCode == http.StatusNotFound {
			return Author{}, ErrAuthorNotFound
		}
		return Author{}, ErrUnexpectedStatus
	}

	p, err := newAuthorPageParser(resp.Body)
	if err != nil {
		_, _ = io.Copy(io.Discard, resp.Body)
		return Author{}, fmt.Errorf("newAuthorPageParser: %w", err)
	}

	return Author{
		Name:             p.name(),
		ShortDescription: p.shortDescription(),
		URL:              p.url(),
		Books:            p.books(),
	}, nil
}

type authorPageParser struct {
	doc *goquery.Document
}

func newAuthorPageParser(r io.Reader) (authorPageParser, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return authorPageParser{}, fmt.Errorf("goquery.NewDocumentFromReader: %w", err)
	}
	return authorPageParser{
		doc: doc,
	}, nil
}

func (p authorPageParser) name() string {
	return strings.TrimSpace(p.doc.Find("#author-info .title-container").Text())
}

func (p authorPageParser) url() string {
	return strings.TrimSpace(p.doc.Find(`meta[property="og:url"]`).AttrOr("content", ""))
}

func (p authorPageParser) shortDescription() string {
	return strings.TrimSpace(p.doc.Find(`meta[name="description"]`).AttrOr("content", ""))
}

func (p authorPageParser) books() []Book {
	booksSel := p.doc.Find("#authorBooks .authorAllBooks__single")
	books := make([]Book, booksSel.Length())
	booksSel.Each(func(i int, selection *goquery.Selection) {
		books[i] = Book{
			Title: strings.TrimSpace(selection.Find(".authorAllBooks__singleTextTitle").Text()),
			URL:   lubimyCzytacBaseURL + selection.Find(".authorAllBooks__singleTextTitle").AttrOr("href", ""),
		}
	})
	return books
}
