package pcloud

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

type favoritesResponse struct {
	Error
	Items []Metadata `json:"items"`
}

func (c *Client) ListFavorites(ctx context.Context) ([]Metadata, error) {
	var resp favoritesResponse
	if err := c.do(ctx, "getfavourites", url.Values{}, &resp); err != nil {
		return nil, fmt.Errorf("list favorites: %w", err)
	}
	return resp.Items, nil
}

func (c *Client) addFavorite(ctx context.Context, params url.Values) error {
	var resp Error
	return c.do(ctx, "addfavourite", params, &resp)
}

func (c *Client) AddFavorite(ctx context.Context, fileID uint64) error {
	if err := c.addFavorite(ctx, url.Values{"fileid": {strconv.FormatUint(fileID, 10)}}); err != nil {
		return fmt.Errorf("add favorite %d: %w", fileID, err)
	}
	return nil
}

func (c *Client) AddFavoriteByPath(ctx context.Context, path string) error {
	if err := c.addFavorite(ctx, url.Values{"path": {path}}); err != nil {
		return fmt.Errorf("add favorite %s: %w", path, err)
	}
	return nil
}

func (c *Client) removeFavorite(ctx context.Context, params url.Values) error {
	var resp Error
	return c.do(ctx, "removefavourite", params, &resp)
}

func (c *Client) RemoveFavorite(ctx context.Context, fileID uint64) error {
	if err := c.removeFavorite(ctx, url.Values{"fileid": {strconv.FormatUint(fileID, 10)}}); err != nil {
		return fmt.Errorf("remove favorite %d: %w", fileID, err)
	}
	return nil
}

func (c *Client) RemoveFavoriteByPath(ctx context.Context, path string) error {
	if err := c.removeFavorite(ctx, url.Values{"path": {path}}); err != nil {
		return fmt.Errorf("remove favorite %s: %w", path, err)
	}
	return nil
}
