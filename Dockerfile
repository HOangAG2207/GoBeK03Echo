# =========================
# Stage 1: Base (chuẩn bị môi trường Go + cache dependency)
# =========================
FROM golang:alpine AS base

# Tạo thư mục làm việc trong container
RUN mkdir -p /opt/app
WORKDIR /opt/app

# Cài thêm build-base
RUN apk add --no-cache build-base

# Copy file go.mod và go.sum trước để tận dụng cache layer
COPY go.mod go sum ./
# Download dependency (được cache nếu go.mod/go.sum không đổi)
RUN go mod download

# Copy toàn bộ source code vào container
COPY . .

# =========================
# Stage 2: Build binary
# =========================
FROM base AS build

# Build binary:
# - CGO_ENABLED=0: build static binary (chạy được trên scratch/alpine)
# - GOOS, GOARCH: build cho Linux amd64
# - -ldflags "-w -s": giảm size binary (strip debug info)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -tags musl -ldflags="-w -s" \
    -o /opt/app/GoBeK03-Echo cmd/api/main.go

# =========================
# Stage 3: Run test + generate coverage
# =========================
FROM base AS test-exec

# Thư mục output coverage (có thể override khi build)
ARG _outputdir="/tmp/coverage"
ARG COVERAGE_EXCLUDE

# Chạy test + generate coverage
RUN mkdir -p ${_outputdir} && \
    go test ./... -coverprofile=coverage.tmp -covermode=atomic -coverpkg=./... -p 1 && \
	grep -v -E "${COVERAGE_EXCLUDE}" coverage.tmp > ${_outputdir}/coverage.out && \
    go tool cover -html=${_outputdir}/coverage.out -o ${_outputdir}/coverage.html

    # =========================
# Stage 4: Export coverage (artifact only)
# =========================
FROM scratch AS test

ARG _outputdir="/tmp/coverage"
# Chỉ copy file coverage ra image cuối (dùng cho CI/CD)
COPY --from=test-exec ${_outputdir}/coverage.out /
COPY --from=test-exec ${_outputdir}/coverage.html /

# =========================
# Stage 5: Final runtime image
# =========================
FROM alpine AS final

# ARG app_name=app
# Set timezone
ENV TZ=Asia/Ho_Chi_Minh

WORKDIR /app

# Copy binary từ stage build
COPY --from=build /opt/app/GoBeK03-Echo /app/GoBeK03-Echo
# Copy docs (nếu cần swagger/static file)
COPY --from=build /opt/app/docs /app/docs

# Cài timezone data + thiết lập timezone
RUN apk add --no-cache tzdata && \
    ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && \
    echo $TZ > /etc/timezone

# Chạy binary
CMD ["/app/GoBeK03-Echo"]