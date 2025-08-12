package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	pob "github.com/reaksa-maii/one_digital_grpc_getway/proto/podcast/v3"
	podcastv3 "github.com/reaksa-maii/one_digital_grpc_getway/proto/podcast/v3"
	"google.golang.org/protobuf/types/known/timestamppb"

)

type AuthorServer struct {
	pob.UnimplementedAuthorServiceServer
	restBaseURL string
	client      *http.Client
}

func NewAuthorServer(restBaseURL string) *AuthorServer {
	return &AuthorServer{
		restBaseURL: restBaseURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type restAuthor struct {
	ID          int `json:"id"`
	Name        string `json:"name"`
	Locale      string `json:"locale"`
	DocumentID  string `json:"documentId"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
	PublishedAt string `json:"publishedAt"`
}

func parseTimePtr(s string) *timestamppb.Timestamp {
	if s == "" {
		return nil
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return nil
	}
	return timestamppb.New(t)
}
func toProtoAuthor(r *restAuthor) *pob.Author {
	return &podcastv3.Author{
		Id:          fmt.Sprintf("%d", r.ID),
		Name:        r.Name,
		Locale:      r.Locale,
		DocumentId:  r.DocumentID,
		CreatedAt:   parseTimePtr(r.CreatedAt),
		UpdatedAt:   parseTimePtr(r.UpdatedAt),
		PublishedAt: parseTimePtr(r.PublishedAt),
	}
}

func (s *AuthorServer) ListAuthors(ctx context.Context, req *pob.ListAuthorsRequest) (*pob.ListAuthorsResponse, error) {
	url := fmt.Sprintf("%s/api/authors", s.restBaseURL)
	httpReq, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	resp, err := s.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("REST error %d: %s", resp.StatusCode, string(body))
	}

	var data struct {
		Data []restAuthor `json:"data"`
	}
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&data); err != nil {
		return nil, err
	}
	out := &pob.ListAuthorsResponse{}
	for _, r := range data.Data {
		out.Authors = append(out.Authors, toProtoAuthor(&r))
	}
	return out, nil
}
func (s *AuthorServer) GetAuthor(ctx context.Context, req *pob.GetAuthorRequest) (*pob.Author, error) {
	url := fmt.Sprintf("%s/api/authors/%s", s.restBaseURL, req.Id)
	httpReq, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	resp, err := s.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("REST error %d: %s", resp.StatusCode, string(body))
	}

	var data struct {
		Data restAuthor `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return toProtoAuthor(&data.Data), nil
}

func (s *AuthorServer) CreateAuthor(ctx context.Context, req *pob.CreateAuthorRequest) (*pob.Author, error) {
	url := fmt.Sprintf("%s/api/authors", s.restBaseURL)

	payload := map[string]interface{}{
		"name":       req.Name,
		"locale":     req.Locale,
		"documentId": req.DocumentId,
	}
	b, _ := json.Marshal(payload)
	httpReq, _ := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(b))
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("REST error %d: %s", resp.StatusCode, string(body))
	}

	var data struct {
		Data restAuthor `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return toProtoAuthor(&data.Data), nil
}
