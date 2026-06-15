package links

import (
	"context"
	"errors"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	repoMocks "github.com/HOangAG2207/GoBeK03Echo/internal/repository/links/mocks"
	utilsMocks "github.com/HOangAG2207/GoBeK03Echo/pkg/utils/mocks"
)

func TestService_ShortenLink(t *testing.T) {
	t.Parallel()

	type fields struct {
		mockRepo func() *repoMocks.Repository
		mockCode func() *utilsMocks.CodeGenerator
	}

	type args struct {
		url        string
		codeLength int
		exptime    int64
	}

	tests := []struct {
		name          string
		fields        fields
		args          args
		expectedCode  string
		expectedError bool
	}{
		{
			name: "success - first attempt",
			fields: fields{
				mockRepo: func() *repoMocks.Repository {
					m := repoMocks.NewRepository(t)

					m.On("GetURL", mock.Anything, "abc123").
						Return("", redis.Nil).Once()

					m.On(
						"StoreURL",
						mock.Anything,
						"abc123",
						"https://example.com",
						int64(60),
					).Return(nil).Once()

					return m
				},
				mockCode: func() *utilsMocks.CodeGenerator {
					m := utilsMocks.NewCodeGenerator(t)

					m.On("GenerateCode", 6).
						Return("abc123", nil).Once()

					return m
				},
			},
			args: args{
				url:        "https://example.com",
				codeLength: 6,
				exptime:    60,
			},
			expectedCode:  "abc123",
			expectedError: false,
		},
		{
			name: "generate code failed",
			fields: fields{
				mockRepo: func() *repoMocks.Repository {
					return repoMocks.NewRepository(t)
				},
				mockCode: func() *utilsMocks.CodeGenerator {
					m := utilsMocks.NewCodeGenerator(t)

					m.On("GenerateCode", 6).
						Return("", errors.New("generate failed")).Once()

					return m
				},
			},
			args: args{
				url:        "https://example.com",
				codeLength: 6,
				exptime:    60,
			},
			expectedError: true,
		},
		{
			name: "store url failed should return error",
			fields: fields{
				mockRepo: func() *repoMocks.Repository {
					m := repoMocks.NewRepository(t)

					m.On("GetURL", mock.Anything, "abc123").
						Return("", redis.Nil).Once()

					m.On(
						"StoreURL",
						mock.Anything,
						"abc123",
						"https://example.com",
						int64(60),
					).Return(errors.New("store failed")).Once()

					return m
				},
				mockCode: func() *utilsMocks.CodeGenerator {
					m := utilsMocks.NewCodeGenerator(t)

					m.On("GenerateCode", 6).
						Return("abc123", nil).Once()

					return m
				},
			},
			args: args{
				url:        "https://example.com",
				codeLength: 6,
				exptime:    60,
			},
			expectedError: true,
		},
		{
			name: "retry once then success",
			fields: fields{
				mockRepo: func() *repoMocks.Repository {
					m := repoMocks.NewRepository(t)

					m.On("GetURL", mock.Anything, "exist01").
						Return("https://old.com", nil).Once()

					m.On("GetURL", mock.Anything, "new123").
						Return("", redis.Nil).Once()

					m.On(
						"StoreURL",
						mock.Anything,
						"new123",
						"https://example.com",
						int64(60),
					).Return(nil).Once()

					return m
				},
				mockCode: func() *utilsMocks.CodeGenerator {
					m := utilsMocks.NewCodeGenerator(t)

					m.On("GenerateCode", 6).
						Return("exist01", nil).Once()

					m.On("GenerateCode", 6).
						Return("new123", nil).Once()

					return m
				},
			},
			args: args{
				url:        "https://example.com",
				codeLength: 6,
				exptime:    60,
			},
			expectedCode:  "new123",
			expectedError: false,
		},
		{
			name: "redis error should fail immediately",
			fields: fields{
				mockRepo: func() *repoMocks.Repository {
					m := repoMocks.NewRepository(t)

					m.On("GetURL", mock.Anything, "abc123").
						Return("", errors.New("redis down")).Once()

					return m
				},
				mockCode: func() *utilsMocks.CodeGenerator {
					m := utilsMocks.NewCodeGenerator(t)

					m.On("GenerateCode", 6).
						Return("abc123", nil).Once()

					return m
				},
			},
			args: args{
				url:        "https://example.com",
				codeLength: 6,
				exptime:    60,
			},
			expectedError: true,
		},
		{
			name: "max retry exceeded",
			fields: fields{
				mockRepo: func() *repoMocks.Repository {
					m := repoMocks.NewRepository(t)

					m.On("GetURL", mock.Anything, mock.Anything).
						Return("exists", nil)

					return m
				},
				mockCode: func() *utilsMocks.CodeGenerator {
					m := utilsMocks.NewCodeGenerator(t)

					m.On("GenerateCode", 6).
						Return("dup123", nil)

					return m
				},
			},
			args: args{
				url:        "https://example.com",
				codeLength: 6,
				exptime:    60,
			},
			expectedError: true,
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := tc.fields.mockRepo()
			codeGen := tc.fields.mockCode()

			s := NewService(repo, codeGen)

			code, err := s.ShortenLink(
				context.Background(),
				tc.args.url,
				tc.args.codeLength,
				tc.args.exptime,
			)

			if tc.expectedError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.expectedCode, code)
		})
	}
}
