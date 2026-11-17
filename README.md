# Dokumentasi Lengkap REST API Card Bridge (`bab6`)
*Panduan Ramah Pemula: Bedah Tuntas Seluruh File, Library, Variabel, dan Fungsi Tanpa Terlewat.*

---

## 💡 1. Konsep Dasar & Analogia Sederhana (Untuk Orang Awam)

Aplikasi ini menggunakan **Layered Architecture**. Bayangkan Restoran:

```
[ Pelanggan ]
       │
       ▼  1. Pesan (HTTP Request)
┌──────────────────────┐
│  A. ROUTER (Kasir)   │ -> Menerima pesanan, mengarahkan ke Pelayan.
└──────────────────────┘
       │
       ▼  2. Pesanan
┌──────────────────────┐
│  B. CONTROLLER (Waiter)│ -> Periksa pesanan, sajikan hasil (JSON).
└──────────────────────┘
       │
       ▼  3. Olah Pesanan
┌──────────────────────┐
│  C. REPOSITORY (Koki)│ -> Olah data langsung di Dapur (Query SQL).
└──────────────────────┘
       │
       ▼  4. Bahan Baku
┌──────────────────────┐
│  D. DATABASE (Gudang)│ -> PostgreSQL (Tempat simpan data).
└──────────────────────┘
```

---

## 📖 2. Glosarium & Simbol Penting

| Simbol / Istilah | Penjelasan Sederhana |
| :--- | :--- |
| **`//go:embed`** | Perintah *compiler* untuk membungkus file luar (seperti `.sql`) ke dalam program. |
| **`_` (Blank ID)** | Variabel "tempat sampah" untuk mengabaikan hasil yang tidak terpakai. |
| **`*` (Pointer)** | Alamat lokasi memori asli data. |
| **`&` (Address-of)** | Mengambil lokasi memori variabel agar bisa diisi nilainya oleh database. |
| **`:=`** | Cara cepat membuat & mengisi variabel baru. |
| **`defer`** | Tunda eksekusi kode sampai fungsi selesai. |
| **CORS** | Keamanan agar website luar boleh akses API kita. |

---

## 🛠️ 3. Instalasi Proyek (Langkah demi Langkah)

1. **Install Go & Inisialisasi**:
   ```bash
   mkdir testbab6 && cd testbab6
   go mod init bab6
   ```

2. **Install Library**:
   ```bash
   go get github.com/gin-gonic/gin
   go get github.com/lib/pq
   go get github.com/rubenv/sql-migrate
   ```

---

## 📘 4. Bedah Rinci File & Fungsi (Sesuai Alur Sistem)

### 4.1. `entities/card_bridge_entity.go`
**Tugas:** Wadah data (Struct) untuk objek kartu.
* **Import:** Tidak ada.
* **Isi:**
  ```go
  type Card struct { ID string `json: "id"` }
  ```
  * `ID`: Variabel teks. `json:"id"`: Label agar saat jadi JSON namanya menjadi `id`.

---

### 4.2. `databases/sql_migrations/1_initiate.sql`
**Tugas:** Perintah SQL pembuat tabel.
* **Penjelasan:**
  * `-- +migrate Up`: Instruksi alat migrasi agar menjalankan perintah ini saat aplikasi mulai.
  * `CREATE TABLE card (id varchar(20));`: SQL murni untuk membuat tabel `card`.

---

### 4.3. `databases/database.go`
**Tugas:** Koneksi database & migrasi otomatis.
* **Import:** `database/sql`, `embed`, `fmt`, `sql-migrate`.
* **Variabel Penting:**
  * `//go:embed ...`: Perintah agar file SQL masuk ke dalam file program hasil *compile*.
  * `DbConnection *sql.DB`: Variabel global kunci koneksi database.
* **Fungsi: `DBMigration(dbParam *sql.DB)`**
  * `dbParam`: Parameter kunci koneksi dari `main.go`.
  * `migrate.Exec(...)`: Fungsi untuk menjalankan SQL migrasi.
  * `panic(errs)`: Berhenti paksa jika ada error konfigurasi.

---

### 4.4. `repositories/card_bridge_repo.go`
**Tugas:** Eksekusi perintah SQL (CRUD).
* **Import:** `bab6/entities`, `database/sql`.
* **Fungsi `GetCards(db *sql.DB)`**:
  * `db.Query(sql)`: Jalankan perintah baca data.
  * `rows.Next()` & `rows.Scan(&data.ID)`: Ambil data dari SQL ke variabel Go.
* **Fungsi `InsertCard`, `UpdateCard`, `DeleteCard`**:
  * `db.Exec(sql, ...)`: Jalankan perintah ubah data.
  * `_`: Mengabaikan *metadata* hasil eksekusi (karena kita hanya butuh error-nya).

---

### 4.5. `controllers/card_bridge_controllers.go`
**Tugas:** Menangani permintaan HTTP (masuk & keluar).
* **Import:** `database`, `entities`, `repositories`, `net/http`, `gin`.
* **Fungsi `GetCards(c *gin.Context)`**:
  * Memanggil repo, lalu kirim JSON.
* **Fungsi `InsertCard`, `UpdateCard`, `DeleteCard`**:
  * `c.Param("id")`: Ambil ID dari URL.
  * `c.ShouldBindJSON(...)`: Ambil data dari JSON body (khusus `PUT/POST`).
  * `c.JSON(http.StatusOK, ...)`: Kirim jawaban sukses.

---

### 4.6. `routers/card_bridge_routers.go`
**Tugas:** Mengarahkan alamat URL ke fungsi di `Controller`.
* **Import:** `controllers`, `cors`, `gin`.
* **Fungsi `StartServer()`**:
  * `cors.New(...)`: Mengatur izin domain agar website lain tidak diblokir.
  * `router.GET/POST/PUT/DELETE`: Mendaftarkan alamat URL ke controller.

---

### 4.7. `main.go`
**Tugas:** Titik masuk utama aplikasi (Main Program).
* **Import:** `database`, `routers`, `sql`, `fmt`, `log`, `lib/pq` (driver SQL).
* **Fungsi `main()`**:
  * `sql.Open(...)`: Membuka koneksi database.
  * `database.DBMigration(DB)`: Menjalankan migrasi.
  * `routers.StartServer().Run(PORT)`: Menyalakan server web.

---

## 🚀 5. Cara Menjalankan
1. `CREATE DATABASE mcsbab6;` (di PostgreSQL).
2. Sesuaikan password di `main.go`.
3. `go mod tidy`
4. `go run main.go`
