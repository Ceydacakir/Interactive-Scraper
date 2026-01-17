# Interactive Scraper - Dark Web İzleme Paneli

Bu proje, belirlediğimiz dark web kaynaklarını (onion siteleri vb.) tarayıp verileri toplayan ve bunları analiz etmemizi sağlayan bir web paneli.

Arka tarafta **Go (Fiber)** ve **PostgreSQL** çalışıyor. İzlediğimiz kaynakların kritiklik seviyelerine göre istatistikleri görüp içerikleri listeleyebiliyoruz.

## Nasıl Çalıştırırım?

Her şeyi Dockerize ettim, o yüzden kurulumla uğraşmana gerek yok. Şu komutu yapıştırman yeterli:

```bash
docker-compose up --build
```

Biraz bekledikten sonra tarayıcıdan panele ulaşabilirsin:
- **Adres:** `http://localhost:3000`
- **Giriş Bilgileri:** `admin` / `admin`

## Neler Var?
- **Dashboard:** Toplanan veri sayısı, kaynak sayısı ve risk grafikleri.
- **Scraper:** Arka planda çalışıp verileri çeken servis.
- **Web Arayüzü:** Verileri incelemek için basit bir panel.
