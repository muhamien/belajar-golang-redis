# Catatan Belajar Golang Redis

Repositori ini berisi referensi lengkap dan hasil eksperimen implementasi struktur data Redis menggunakan Golang (`github.com/redis/go-redis/v9`). Berikut adalah daftar fitur Redis yang dieksplorasi di dalam berkas tes beserta penjelasan, _snippet_ kode, dan ekspektasi *Output*-nya.

## ⚙️ Persiapan Lingkungan

Proyek ini menggunakan Docker Compose agar Anda tidak perlu menginstal server Redis secara manual di *Host OS*.
```yaml
# docker-compose.yml
version: '3.5'

services:
  redis:
    container_name: redis
    image: redis:6
    ports:
      - 6379:6379
```
Jalankan server Redis di *background*: `docker-compose up -d`

---

## 📚 Struktur Data & Fitur Redis

### 1. Connection & Ping
Membuat _Client_ koneksi ke Redis _server_ lokal dan menguji koneksinya menggunakan perintah ping bawaan.
```go
var client = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
	DB:   0, // Menggunakan database indeks 0 (default)
})

// Eksekusi Ping
result, err := client.Ping(ctx).Result()
fmt.Println(result) 
// Output: "PONG"
```

### 2. String (Tipe Data Dasar)
Menyimpan pasangan *Key-Value* tunggal layaknya variabel biasa. Sangat umum digunakan untuk menyimpan hasil _caching array/struct_ (dalam format JSON) atau merekam *Session Token* bersama parameter kadaluarsa waktu (*TTL*).
```go
// Disimpan hanya selama 3 detik, jika lebih key ini akan lenyap (nil)
client.Set(ctx, "name", "Muhammad Amien", time.Second*3)
result, err := client.Get(ctx, "name").Result()

fmt.Println(result) 
// Output: "Muhammad Amien"
```

### 3. List
Mengadopsi antrean berurutan (_Sequential List_/Linked List). Sangat berguna ketika Anda mengimplementasikan antrean sederhana *Queue* FIFO (*First In First Out*) atau tumpukan Log *Stack* LIFO (*Last In First Out*).
- `RPush`: Memasukkan elemen dari kanan (belakang antrean).
- `LPop`: Menarik (mengeluarkan) elemen dari kiri (depan antrean), lalu menghapusnya agar elemen di belakangnya maju.
```go
client.RPush(ctx, "antrean_tugas", "Tugas A")
client.RPush(ctx, "antrean_tugas", "Tugas B")
client.RPush(ctx, "antrean_tugas", "Tugas C")

tugasPertama := client.LPop(ctx, "antrean_tugas").Val()
tugasKedua := client.LPop(ctx, "antrean_tugas").Val()

fmt.Println(tugasPertama) // Output: "Tugas A"
fmt.Println(tugasKedua)   // Output: "Tugas B"
```

### 4. Set
Tipe data _Collection_ yang hanya memelihara nilai-nilai "Unik" secara otomatis tanpa perlu pemeriksaan berulang di _Application Level_ (Golang). Jika kita memasukkan nilai duplikat, Redis sekadar mengabaikannya.
```go
client.SAdd(ctx, "tag_artikel", "Golang")
client.SAdd(ctx, "tag_artikel", "Redis")
client.SAdd(ctx, "tag_artikel", "Golang") // Diabaikan karena telah eksis

jumlahTags := client.SCard(ctx, "tag_artikel").Val()
semuaTags := client.SMembers(ctx, "tag_artikel").Val()

fmt.Println(jumlahTags) // Output: 2
fmt.Println(semuaTags)  // Output: ["Golang", "Redis"] (Urutan bisa acak)
```

### 5. Sorted Set
Koleksi setrik (unik) yang elemennya diperkaya beban penimbang berupa parameter _Score_ bertipe *Float/Integer* sehingga elemen secara inheren terurut dari yang terkecil sampai tertinggi. Fitur ini lazim diadopsi pada _Game Leaderboard_ (Papan Skor Teratas).
```go
client.ZAdd(ctx, "papan_skor", redis.Z{Score: 100, Member: "Muhammad"})
client.ZAdd(ctx, "papan_skor", redis.Z{Score: 85, Member: "Amien"})
client.ZAdd(ctx, "papan_skor", redis.Z{Score: 96, Member: "Rauf"})

// Mengambil rank 0 (terendah) hingga rank 2 (tertinggi ketiga)
rankBawahKeAtas := client.ZRange(ctx, "papan_skor", 0, 2).Val()
fmt.Println(rankBawahKeAtas) 
// Output: ["Amien", "Rauf", "Muhammad"] 
// Catatan: Amien memiliki score terendah (85) sehingga berada di index 0

// Mengambil (dan menghapus) satu partisipan dengan score tertingginya (Rank 1)
pemenangTertinggi := client.ZPopMax(ctx, "papan_skor").Val()
fmt.Println(pemenangTertinggi[0].Member) 
// Output: "Muhammad"
```

### 6. Hash
Menyimpan sekumpulan *Key-Value* turunan (seperti tipe data `Map` atau _struct_ Golang). Struktur ini superior ketimbang mengubah entitas jadi String berbasis JSON apabila kebutuhannya: "Hanya *merubah*/meng-ngupdate parameter `name` tanpa perlu mengunduh parameter `id` dan yang lainnya".
```go
// Menyimpan entity User dengan ID=1 dan Name="Amien" ke Key = user:1
client.HSet(ctx, "user:1", "id", "1", "name", "Amien")

// Download satu persatu atau HGetAll untuk mengambil utuh sebagai map
userMap := client.HGetAll(ctx, "user:1").Val() 

fmt.Println(userMap["id"])   // Output: "1"
fmt.Println(userMap["name"]) // Output: "Amien"
```

### 7. GeoPoint
Mengelola kumpulan koordinat lokasi geografis (*Longitude* & *Latitude*). Dapat mengevaluasi batas jarak yang melingkari pusat lokasi pengguna dalam satuan unit Kilometer (km) maupun Meter (m).
```go
// Insert Kordinat
client.GeoAdd(ctx, "lokasi_toko", &redis.GeoLocation{Name: "Toko A", Longitude: 110.408456, Latitude:  -7.739936})
client.GeoAdd(ctx, "lokasi_toko", &redis.GeoLocation{Name: "Toko B", Longitude: 110.4144109, Latitude: -7.7420133})

// Mengukur Jarak Antara dua Titik (Dalam Kilometer)
jarakTokoAKeB := client.GeoDist(ctx, "lokasi_toko", "Toko A", "Toko B", "km").Val()
fmt.Println(jarakTokoAKeB) 
// Output: 0.6963 (Jaraknya hanya sekitar 696 Meter)

// Query: Berikan Toko apa saja yang berada di dalam Radius Titik Pusat saya (15 Kilometer)
sekitarSaya := client.GeoSearch(ctx, "lokasi_toko", &redis.GeoSearchQuery{
	Longitude:  110.410960,
	Latitude:   -7.740996,
	Radius:     15,
	RadiusUnit: "km",
}).Val()

fmt.Println(sekitarSaya) 
// Output: ["Toko A", "Toko B"]
```

### 8. HyperLogLog
Alat spesialis untuk menghitung _Unique Cardinality_ (Perkiraan perhitungan nilai yang unik) skala luar biasa masif (Misal 1 milyar penonton TV). Algoritmanya luar biasa hemat memori namun menghasilkan _Margin of error_ (± 0.81%). Biasa dipakai mengakumulasi "_Unique Visitor / IP Addresses_ / Views".
```go
client.PFAdd(ctx, "ip_pengunjung", "Muhammad", "Amien", "Rauf")
client.PFAdd(ctx, "ip_pengunjung", "Joko", "Gibs", "Prabs")
client.PFAdd(ctx, "ip_pengunjung", "Gibs", "Amien", "Prabs") // Redundan

totalPengunjung := client.PFCount(ctx, "ip_pengunjung").Val()

fmt.Println(totalPengunjung) 
// Output: 6 (Menghitung elemen unik saja secara probabilistik)
```

### 9. Pipeline & Transaction
Mekanisme membungkus rentetan *Command* untuk dikirim secara berbarengan hanya dalam *Satu Request* (1 kali jalan network ping). Sangat menghemat latensi TCP/IP.
- **Pipelined:** Mem-paketkan query (Cukup Cepat), namun instruksinya belum tentu berjalan rapat, karena masih memungkinkan *query* app Client orang lain menyelip masuk di ruang jeda antaranya.
- **TxPipelined (Transaction):** Modus Ekstra (Rapat/Atomic). Menggunakan kunci konvensi perintah internal (`MULTI` & `EXEC`). Aplikasi klien lain diblokir sementara agar set *Commands* kita dipastikan rampung berurutan tanpa disela dari luar.
```go
// Transaction Mode (Rekomendasi)
_, err := client.TxPipelined(ctx, func(pipeliner redis.Pipeliner) error {
	// Antrean eksekusi belum dijalankan secara sungguhan pada tahap ini.
	pipeliner.SetEx(ctx, "nama", "Joko", 5*time.Second)
	pipeliner.SetEx(ctx, "kota", "Cirebon", 5*time.Second)
	return nil
}) 
// Setelah return, redis Go barulah 'Menembakkan' paket Command via jaringan.

fmt.Println(client.Get(ctx, "nama").Val()) // Output: "Joko"
fmt.Println(client.Get(ctx, "kota").Val()) // Output: "Cirebon"
```

### 10. Streaming (Pub/Sub Modern)
Struktur tipe *Message Queue* terdepan Redis (*append-only log*). Lebih persisten dibanding Pub/Sub primitif. Data log tersimpan (bisa ditelusuri ulang jika Subscriber terputus).
```go
// XAdd: Menambah log entri ke dalam Stream channel bernama "member_baru"
err := client.XAdd(ctx, &redis.XAddArgs{
	Stream: "member_baru",
	Values: map[string]interface{}{
		"nama":    "Eko",
		"negara": "Indonesia",
	},
}).Err()

// Channel Stream akan menampung Dictionary/Map: {"nama": "Eko", "negara": "Indonesia"}
```

### 11. Pub/Sub (Publish/Subscribe)
Sistem Penyiaran Pesan Real-time tradisional. Tanpa penyimpanan (Memori = _Zero_). Saat ada event "A", semua fungsi yang ikut _Listening / Subscribe_ akan menerimanya bersamaan. Jika node Golang sedang *Crashed/Down*, instruksi tadi lenyap terbawa angin (*Fire and Forget*).
```go
// Node 1 (Subscriber yang memantau telinga pada channel-1)
subscriber := client.Subscribe(ctx, "channel-1")
message, err := subscriber.ReceiveMessage(ctx)
fmt.Println(message.Payload) 
// Output (Reaktif): "Pesan Masuk Ke-0", lalu "Pesan Masuk Ke-1" ...

// Node 2 (Publisher yang menyiarkan)
for i := 0; i < 2; i++ {
	client.Publish(ctx, "channel-1", fmt.Sprintf("Pesan Masuk Ke-%d", i))
}
```
