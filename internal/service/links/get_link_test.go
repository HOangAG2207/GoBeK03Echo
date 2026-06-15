package links

import (
	"context"
	"errors"
	"testing"

	repoMocks "github.com/HOangAG2207/GoBeK03Echo/internal/repository/links/mocks"
	utilsMocks "github.com/HOangAG2207/GoBeK03Echo/pkg/utils/mocks"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_GetLink(t *testing.T) {
	t.Parallel()

	type args struct {
		inputCode string
	}

	type expected struct {
		url string
		err error
	}

	type fields struct {
		mockRepo func() *repoMocks.Repository
		mockCode func() *utilsMocks.CodeGenerator
	}

	testCases := []struct {
		name     string
		args     args
		fields   fields
		expected expected
	}{
		{
			name: "success - URL found",
			args: args{
				inputCode: "abc123",
			},
			fields: fields{
				mockRepo: func() *repoMocks.Repository {
					m := repoMocks.NewRepository(t)

					m.On("GetURL", mock.Anything, "abc123").
						Return("https://example.com", nil).
						Once()

					return m
				},
				mockCode: func() *utilsMocks.CodeGenerator {
					return utilsMocks.NewCodeGenerator(t)
				},
			},
			expected: expected{
				url: "https://example.com",
			},
		},
		{
			name: "not found - redis.Nil should return ErrCodeNotFound",
			args: args{
				inputCode: "notfound",
			},
			fields: fields{
				mockRepo: func() *repoMocks.Repository {
					m := repoMocks.NewRepository(t)

					m.On("GetURL", mock.Anything, "notfound").
						Return("", redis.Nil).
						Once()

					return m
				},
				mockCode: func() *utilsMocks.CodeGenerator {
					return utilsMocks.NewCodeGenerator(t)
				},
			},
			expected: expected{
				err: ErrCodeNotFound,
			},
		},
		{
			name: "empty string is valid value",
			args: args{
				inputCode: "empty",
			},
			fields: fields{
				mockRepo: func() *repoMocks.Repository {
					m := repoMocks.NewRepository(t)

					m.On("GetURL", mock.Anything, "empty").
						Return("", nil).
						Once()

					return m
				},
				mockCode: func() *utilsMocks.CodeGenerator {
					return utilsMocks.NewCodeGenerator(t)
				},
			},
			expected: expected{
				url: "",
			},
		},
		{
			name: "repo error - should return error directly",
			args: args{
				inputCode: "errorcase",
			},
			fields: fields{
				mockRepo: func() *repoMocks.Repository {
					m := repoMocks.NewRepository(t)

					m.On("GetURL", mock.Anything, "errorcase").
						Return("", errors.New("redis down")).
						Once()

					return m
				},
				mockCode: func() *utilsMocks.CodeGenerator {
					return utilsMocks.NewCodeGenerator(t)
				},
			},
			expected: expected{
				err: errors.New("redis down"),
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := tc.fields.mockRepo()
			codeGen := tc.fields.mockCode()

			s := NewService(repo, codeGen)

			url, err := s.GetLink(
				context.Background(),
				tc.args.inputCode,
			)

			assert.Equal(t, tc.expected.url, url)

			if tc.expected.err != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expected.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			repo.AssertExpectations(t)
			codeGen.AssertExpectations(t)
		})
	}
}
