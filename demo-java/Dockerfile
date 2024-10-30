# Sử dụng image Maven với JDK 21 từ Eclipse Temurin để build
FROM maven:3.9.4-eclipse-temurin-21 AS build

# Đặt thư mục làm việc
WORKDIR /app

# Copy file cấu hình và source code vào container
COPY pom.xml ./
COPY src ./src

# Build ứng dụng
RUN mvn clean package -DskipTests

# Sử dụng JDK 21 runtime image nhẹ từ Eclipse Temurin cho runtime
FROM eclipse-temurin:21-jdk-jammy

# Đặt thư mục làm việc
WORKDIR /app

# Copy file jar từ giai đoạn build vào image cuối
COPY --from=build /app/target/*.jar app.jar

# Expose port của ứng dụng (thường là 8080)
EXPOSE 8091

# Lệnh để khởi chạy ứng dụng
ENTRYPOINT ["java", "-jar", "app.jar"]
