package pcloud

import (
	"context"
	"net/url"
	"strconv"
)

type favoritesResponse struct {
	Error
	Items []Metadata `json:"items"`
}

// ListFavorites returns all files and folders marked as favorites.
func (c *Client) ListFavorites(ctx context.Context) ([]Metadata, error) {
	var resp favoritesResponse
	if err := c.do(ctx, "getfavourites", url.Values{}, &resp); err != nil {
		return nil, err
	}
	return resp.Items, nil
}

func (c *Client) addFavorite(ctx context.Context, params url.Values) error {
	var resp Error
	return c.do(ctx, "addfavourite", params, &resp)
}

// AddFavorite marks a file as a favorite by numeric ID.
func (c *Client) AddFavorite(ctx context.Context, fileID uint64) error {
	return c.addFavorite(ctx, url.Values{"fileid": {strconv.FormatUint(fileID, 10)}})
}

// AddFavoriteByPath marks a file as a favorite by path.
func (c *Client) AddFavoriteByPath(ctx context.Context, path string) error {
	return c.addFavorite(ctx, url.Values{"path": {path}})
}

func (c *Client) removeFavorite(ctx context.Context, params url.Values) error {
	var resp Error
	return c.do(ctx, "removefavourite", params, &resp)
}

// RemoveFavorite removes a file from favorites by numeric ID.
func (c *Client) RemoveFavorite(ctx context.Context, fileID uint64) error {
	return c.removeFavorite(ctx, url.Values{"fileid": {strconv.FormatUint(fileID, 10)}})
}

// RemoveFavoriteByPath removes a file from favorites by path.
func (c *Client) RemoveFavoriteByPath(ctx context.Context, path string) error {
	return c.removeFavorite(ctx, url.Values{"path": {path}})
}
