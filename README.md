# Interactive Scraper - Dark Web İzleme Paneli

Bu proje, belirlediğimiz dark web kaynaklarını (onion siteleri vb.) tarayıp verileri toplayan ve bunları analiz etmemizi sağlayan bir web paneli.

Arka tarafta **Go (Fiber)** ve **PostgreSQL** çalışıyor. İzlediğimiz kaynakların kritiklik seviyelerine göre istatistikleri görüp içerikleri listeleyebiliyoruz.

## Nasıl Çalıştırırım?


```bash
docker-compose up --build
```

- **Adres:** `http://localhost:3000`
- **Giriş Bilgileri:** `admin` / `admin`

## Neler Var?
- **Dashboard:** Toplanan veri sayısı, kaynak sayısı ve risk grafikleri.
- **Scraper:** Arka planda çalışıp verileri çeken servis.
- **Web Arayüzü:** Verileri incelemek için basit bir panel.

