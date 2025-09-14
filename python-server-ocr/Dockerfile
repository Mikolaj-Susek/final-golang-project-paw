# Używamy oficjalnego, lekkiego obrazu Python 3.9
FROM python:3.12-slim

# Ustawiamy katalog roboczy wewnątrz kontenera
WORKDIR /app

# Kopiujemy plik z zależnościami
COPY requirements.txt .

# Instalujemy zależności Pythona
# --no-cache-dir zapobiega zapisywaniu cache, co zmniejsza rozmiar obrazu
RUN pip install --no-cache-dir -r requirements.txt

# Kopiujemy resztę plików aplikacji serwera do kontenera
COPY . .

# Wystawiamy port 50051, na którym nasłuchuje serwer gRPC
EXPOSE 50051

# Komenda do uruchomienia serwera przy starcie kontenera
CMD ["python", "server.py"]