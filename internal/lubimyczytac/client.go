package lubimyczytac

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	defaultBaseURL = "https://lubimyczytac.pl"
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
	ID               string
	Name             string
	ShortDescription string
	URL              string
	Books            []Book
}

type Client struct {
	http    *http.Client
	baseURL string
}

type ClientOption func(c *Client)

func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

func NewClient(client *http.Client, opts ...ClientOption) *Client {
	c := &Client{http: client, baseURL: defaultBaseURL}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Client) GetAuthor(ctx context.Context, id string) (Author, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/autor/%s/x", c.baseURL, id),
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
		ID:               id,
		Name:             p.name(),
		ShortDescription: p.shortDescription(),
		URL:              p.url(),
		Books:            p.books(),
	}, nil
}

type authorPageParser struct {
	doc     *goquery.Document
	baseURL *url.URL
}

func newAuthorPageParser(r io.Reader) (authorPageParser, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return authorPageParser{}, fmt.Errorf("goquery.NewDocumentFromReader: %w", err)
	}
	baseURL, err := url.Parse(doc.Find("head base").AttrOr("href", ""))
	if err != nil {
		return authorPageParser{}, fmt.Errorf("url.Parse: %w", err)
	}
	return authorPageParser{
		doc:     doc,
		baseURL: baseURL,
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
		bookUrl := url.URL{
			Scheme: p.baseURL.Scheme,
			Host:   p.baseURL.Host,
			Path:   selection.Find(".authorAllBooks__singleTextTitle").AttrOr("href", ""),
		}
		books[i] = Book{
			Title: strings.TrimSpace(selection.Find(".authorAllBooks__singleTextTitle").Text()),
			URL:   bookUrl.String(),
		}
	})
	return books
}
