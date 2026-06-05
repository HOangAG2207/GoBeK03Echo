package links

import (
	"context"
	"errors"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/mock"

	"github.com/HOangAG2207/GoBeK03Echo/internal/repository/links/mocks"
)

func TestService_ShortenURL(t *testing.T) {
	t.Parallel()

	type fields struct {
		mockRepo func() *mocks.Repository
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
		expectedError bool
	}{
		{
			name: "success - first attempt",
			fields: fields{
				mockRepo: func() *mocks.Repository {
					m := mocks.NewRepository(t)
					m.On("GetURL", mock.Anything, mock.Anything).
						Return("", redis.Nil).Once()

					m.On("StoreURL", mock.Anything, mock.Anything, "https://example.com", int64(60)).
						Return(nil).Once()

					return m
				},
			},
			args: args{
				url:        "https://example.com",
				codeLength: 6,
				exptime:    60,
			},
			expectedError: false,
		},
		{
			name: "store url failed should return error",
			fields: fields{
				mockRepo: func() *mocks.Repository {
					m := mocks.NewRepository(t)

					m.On("GetURL", mock.Anything, mock.Anything).
						Return("", redis.Nil).Once()

					m.On("StoreURL", mock.Anything, mock.Anything, "https://example.com", int64(60)).
						Return(errors.New("store failed")).Once()

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
				mockRepo: func() *mocks.Repository {
					m := mocks.NewRepository(t)

					// first: exists
					m.On("GetURL", mock.Anything, mock.Anything).
						Return("exist", nil).Once()

					// second: not exists
					m.On("GetURL", mock.Anything, mock.Anything).
						Return("", redis.Nil).Once()

					m.On("StoreURL", mock.Anything, mock.Anything, "https://example.com", int64(60)).
						Return(nil).Once()

					return m
				},
			},
			args: args{
				url:        "https://example.com",
				codeLength: 6,
				exptime:    60,
			},
			expectedError: false,
		},
		{
			name: "redis error (not Nil) should fail immediately",
			fields: fields{
				mockRepo: func() *mocks.Repository {
					m := mocks.NewRepository(t)

					m.On("GetURL", mock.Anything, mock.Anything).
						Return("", errors.New("redis down")).Once()

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
				mockRepo: func() *mocks.Repository {
					m := mocks.NewRepository(t)

					// luôn trả về exists => loop mãi đến max retry
					m.On("GetURL", mock.Anything, mock.Anything).
						Return("exists", nil)

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
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := tc.fields.mockRepo()
			s := NewService(repo)

			_, err := s.ShortenURL(
				context.Background(),
				tc.args.url,
				tc.args.codeLength,
				tc.args.exptime,
			)

			if (err != nil) != tc.expectedError {
				t.Errorf("ShortenURL() error = %v, expectedError %v", err, tc.expectedError)
			}
		})
	}
}
