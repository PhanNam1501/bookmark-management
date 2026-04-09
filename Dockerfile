# Khuyến nghị nên fix cứng phiên bản Go (VD: 1.22) thay vì để golang:alpine chung chung
FROM golang:1.25-alpine AS base
RUN mkdir -p /opt/app
WORKDIR /opt/app

# Chỉ cần thiết nếu bạn dùng CGO (ví dụ xài thư viện go-sqlite3). Nếu app Go thuần, có thể bỏ dòng này.
RUN apk add --no-cache build-base

# Viết gộp lại để tối ưu số lượng layer của Docker
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# ==========================================
FROM base AS build

# 1. Sửa lỗi gõ phím: amin.go -> main.go
# 2. Thêm CGO_ENABLED=0 để tạo ra file binary tĩnh (statically linked), chạy ổn định trên Alpine
RUN CGO_ENABLED=0 GOOS=linux go build -tags musl -ldflags="-w -s" -o bookmark_service cmd/api/main.go

# ==========================================
# ==========================================
FROM base AS test-exec

ARG _outputdir="/tmp/coverage"
ARG COVERAGE_EXCLUDE="mocks|main.go|docs|test"

# 1. Cài đặt Redis vào thẳng container dùng để test
RUN apk add --no-cache redis

# 2. Tạo thư mục coverage
# 3. Bật Redis chạy ngầm (--daemonize yes)
# 4. Đợi 2 giây (sleep 2) để đảm bảo Redis đã khởi động xong hoàn toàn
# 5. Chạy go test (lúc này code test gọi localhost:6379 sẽ ăn ngay vào con Redis vừa bật)
RUN mkdir -p ${_outputdir} && \
    redis-server --daemonize yes && \
    sleep 2 && \
    go test ./... -coverprofile=coverage.tmp -covermode=atomic -coverpkg=./... -p 1 && \
    grep -v -E "${COVERAGE_EXCLUDE}" coverage.tmp > ${_outputdir}/coverage.out && \
    go tool cover -html=${_outputdir}/coverage.out -o ${_outputdir}/coverage.html
# ==========================================
# ==========================================
FROM scratch AS test
ARG _outputdir="/tmp/coverage"
COPY --from=test-exec ${_outputdir}/coverage.out /
COPY --from=test-exec ${_outputdir}/coverage.html /

# ==========================================
FROM alpine:latest AS final

ENV TZ=Asia/Ho_Chi_Minh

# BẮT BUỘC SỬA: Alpine gốc không có package tzdata. 
# Phải cài tzdata thì lệnh setup timezone bên dưới mới tìm thấy thư mục /usr/share/zoneinfo/
# Cài thêm ca-certificates để app gọi được các API HTTPS bên ngoài.
RUN apk add --no-cache tzdata ca-certificates && \
    ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && \
    echo $TZ > /etc/timezone

WORKDIR /app

COPY --from=build /opt/app/bookmark_service /app/bookmark_service
COPY --from=build /opt/app/docs /app/docs

# Nên khai báo Port để người xem Dockerfile biết app chạy ở cổng nào (VD: 8080)
EXPOSE 8080

# Bổ sung lệnh chạy app mặc định
CMD ["/app/bookmark_service"]