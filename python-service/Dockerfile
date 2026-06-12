# базовый образ
FROM python:3.11-slim

# рабочая директория
WORKDIR /app

# копируем зависимости
COPY requirements.txt .

# ставим зависимости
RUN pip install --no-cache-dir -r requirements.txt

# копируем код
COPY . .

EXPOSE 50051

# переменные окружения
ENV PYTHONUNBUFFERED=1

# запуск
CMD ["python", "server.py"]
