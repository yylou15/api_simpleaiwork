package service

import (
	"context"
	"errors"
	"strings"

	"api/biz/say_right/dal/query"
)

var ErrProRequired = errors.New("pro required")
var ErrTemplateNotFound = errors.New("template not found")

type TemplateItem struct {
	ID          int32    `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	IsPro       bool     `json:"is_pro"`
	IsLocked    bool     `json:"is_locked"`
}

type CategoryWithTemplates struct {
	ID          int32          `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Icon        string         `json:"icon"`
	Templates   []TemplateItem `json:"templates"`
}

type TemplateDetailResult struct {
	ID            int32    `json:"id"`
	CategoryID    int32    `json:"category_id"`
	CategoryName  string   `json:"category_name"`
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	Tags          []string `json:"tags"`
	IsPro         bool     `json:"is_pro"`
	ReplySoft     string   `json:"reply_soft"`
	ReplyNeutral  string   `json:"reply_neutral"`
	ReplyFirm     string   `json:"reply_firm"`
	WhenNotToUse  string   `json:"when_not_to_use"`
	BestPractices []string `json:"best_practices"`
}

type TemplateService interface {
	ListTemplatesByCategory(ctx context.Context, userID int32) ([]CategoryWithTemplates, error)
	GetTemplateDetail(ctx context.Context, userID int32, templateID int32) (*TemplateDetailResult, error)
}

type templateService struct {
	q *query.Query
}

func NewTemplateService() TemplateService {
	return &templateService{
		q: query.Q,
	}
}

func (s *templateService) ListTemplatesByCategory(ctx context.Context, userID int32) ([]CategoryWithTemplates, error) {
	user, err := s.q.User.WithContext(ctx).Where(s.q.User.ID.Eq(userID)).First()
	if err != nil {
		return nil, err
	}

	categories, err := s.q.Category.WithContext(ctx).
		Where(s.q.Category.IsActive.Eq(1)).
		Order(s.q.Category.SortOrder, s.q.Category.ID).
		Find()
	if err != nil {
		return nil, err
	}

	templates, err := s.q.Template.WithContext(ctx).
		Where(s.q.Template.IsActive.Eq(1)).
		Order(s.q.Template.CategoryID, s.q.Template.SortOrder, s.q.Template.ID).
		Find()
	if err != nil {
		return nil, err
	}

	templatesByCategory := make(map[int32][]TemplateItem)
	for _, t := range templates {
		isPro := t.IsPro != 0
		item := TemplateItem{
			ID:          t.ID,
			Title:       t.Title,
			Description: t.Description,
			Tags:        splitTags(t.TagsText),
			IsPro:       isPro,
			IsLocked:    isPro && user.IsPro == 0,
		}
		templatesByCategory[t.CategoryID] = append(templatesByCategory[t.CategoryID], item)
	}

	result := make([]CategoryWithTemplates, 0, len(categories))
	for _, c := range categories {
		result = append(result, CategoryWithTemplates{
			ID:          c.ID,
			Name:        c.Name,
			Description: c.Description,
			Icon:        c.Icon,
			Templates:   templatesByCategory[c.ID],
		})
	}

	return result, nil
}

func (s *templateService) GetTemplateDetail(ctx context.Context, userID int32, templateID int32) (*TemplateDetailResult, error) {
	user, err := s.q.User.WithContext(ctx).Where(s.q.User.ID.Eq(userID)).First()
	if err != nil {
		return nil, err
	}

	template, err := s.q.Template.WithContext(ctx).
		Where(s.q.Template.ID.Eq(templateID)).
		Where(s.q.Template.IsActive.Eq(1)).
		First()
	if err != nil {
		return nil, ErrTemplateNotFound
	}

	isPro := template.IsPro != 0
	if isPro && user.IsPro == 0 {
		return nil, ErrProRequired
	}

	detail, err := s.q.TemplateDetail.WithContext(ctx).
		Where(s.q.TemplateDetail.TemplateID.Eq(templateID)).
		First()
	if err != nil {
		return nil, ErrTemplateNotFound
	}

	category, err := s.q.Category.WithContext(ctx).
		Where(s.q.Category.ID.Eq(template.CategoryID)).
		First()
	if err != nil {
		return nil, err
	}

	title := template.Title
	if detail.Headline != "" {
		title = detail.Headline
	}

	result := &TemplateDetailResult{
		ID:            template.ID,
		CategoryID:    template.CategoryID,
		CategoryName:  category.Name,
		Title:         title,
		Description:   detail.Summary,
		Tags:          splitTags(template.TagsText),
		IsPro:         isPro,
		ReplySoft:     detail.ReplySoft,
		ReplyNeutral:  detail.ReplyNeutral,
		ReplyFirm:     detail.ReplyFirm,
		WhenNotToUse:  detail.WhenNotToUse,
		BestPractices: splitLines(detail.BestPractices),
	}

	return result, nil
}

func splitTags(value string) []string {
	parts := strings.FieldsFunc(value, func(r rune) bool {
		return r == ',' || r == '，'
	})
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		if v := strings.TrimSpace(p); v != "" {
			result = append(result, v)
		}
	}
	return result
}

func splitLines(value string) []string {
	parts := strings.Split(strings.ReplaceAll(value, "\r\n", "\n"), "\n")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		if v := strings.TrimSpace(p); v != "" {
			result = append(result, v)
		}
	}
	return result
}
