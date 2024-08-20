# 構建階段
FROM node:20 as build

WORKDIR /app

# 複製 package.json 和 package-lock.json（如果可用）
COPY package*.json ./

# 安裝項目依賴
RUN npm install

# 複製項目文件並構建應用
COPY . .
RUN npm run build

# 生產階段
FROM nginx:alpine

# 複製構建文件到 Nginx 服務目錄
COPY --from=build /app/build /usr/share/nginx/html

# 暴露端口 80
EXPOSE 80

# 啟動 Nginx 服務
CMD ["nginx", "-g", "daemon off;"]