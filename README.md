# Panduan & Dokumentasi Kode Lengkap REST API Card Bridge (`bab6`)
*Panduan Ramah Pemula & Orang Awam: Penjelasan Konsep, Kode Baris demi Baris, Variabel, Sintaks, dan Fungsi.*

---

## 💡 1. Konsep Dasar & Analogia Sederhana (Untuk Orang Awam)

Sebelum masuk ke kode, mari kita pahami apa yang dibuat oleh aplikasi ini dengan analogi sederhana: **Restoran**.

```
[ Pelanggan / Browser / Postman ]
               │
               ▼  1. Pesan/Update Makanan (HTTP Request: GET, POST, PUT, DELETE)
┌──────────────────────────────────────────────┐
│  A. ROUTER (Buku Menu & Kasir)              │ -> Menerima pesanan dan mengarahkan ke Pelayan yang tepat.
└──────────────────────────────────────────────┘
               │
               ▼  2. Meneruskan Pesanan
┌──────────────────────────────────────────────┐
│  B. CONTROLLER (Pelayan / Waiter)           │ -> Memeriksa format pesanan & menyajikan makanan (JSON Response).
└──────────────────────────────────────────────┘
               │
               ▼  3. Minta Bahan / Olah Data
┌──────────────────────────────────────────────┐
│  C. REPOSITORY (Koki / Chef di Dapur)        │ -> Membaca & Menulis data langsung ke Gudang (Database SQL).
└──────────────────────────────────────────────┘
               │
               ▼  4. Ambil / Simpan Data
┌──────────────────────────────────────────────┐
│  D. DATABASE & ENTITY (Gudang & Bentuk Wadah)│ -> PostgreSQL (Gudang Data) & Card Struct (Bentuk Wadah Kartu).
└──────────────────────────────────────────────┘
```

* **REST API**: Cara aplikasi saling berkomunikasi menggunakan protokol HTTP. Klien mengirim *Request* (permintaan) dan server memberikan *Response* (jawaban berupa teks JSON).
* **HTTP Methods**:
  * **`GET`**: Membaca / mengambil data (seperti melihat menu).
  * **`POST`**: Menambahkan data baru (seperti memesan makanan baru).
  * **`PUT`**: Memperbarui/mengubah data yang *sudah ada* (seperti mengganti isi pesanan yang salah).
  * **`DELETE`**: Menghapus data (seperti membatalkan pesanan).
* **Layered Architecture**: Memisah-misah kode ke beberapa folder (Router, Controller, Repository, Entity, Database) agar jika ada perubahan atau kerusakan, kita tahu persis folder mana yang harus diperbaiki tanpa mengganggu bagian lain.

---

## 📖 2. Glosarium Istilah Pemula

*Baca ini agar lebih mudah memahami istilah teknis di bawah:*

| Istilah | Penjelasan Sederhana untuk Pemula |
| :--- | :--- |
| **Golang / Go** | Bahasa pemrograman yang cepat, ringan, dan cocok untuk backend server. |
| **Gin Gonic** | Alat bantu (framework) untuk membuat web server di Go. |
| **POST** | Metode HTTP untuk **menambah** data baru. |
| **GET** | Metode HTTP untuk **mengambil/membaca** data. |
| **PUT** | Metode HTTP untuk **memperbarui/mengubah** data yang sudah ada (secara keseluruhan). |
| **DELETE** | Metode HTTP untuk **menghapus** data. |
| **Struct** | Fitur Go untuk membuat wadah tipe data sendiri (seperti objek). |
| **JSON** | Format teks pertukaran data di internet (`"key": "value"`). |

---

## 🛠️ 3. Langkah-Langkah Membuat Proyek Go Baru Menggunakan Framework Gin & Library Pendukung

Bagi pemula yang ingin membuat proyek Go baru dari nol, ikuti tahapan lengkap di bawah ini:

### Langkah 1: Install Golang
Pastikan Go sudah terinstal di komputer Anda. Cek melalui terminal:
```bash
go version
```

### Langkah 2: Buat Folder Proyek Baru
Buka terminal/command prompt, buat folder proyek baru, lalu masuk ke dalam folder tersebut:
```bash
mkdir proyek-baru
cd proyek-baru
```

### Langkah 3: Inisialisasi Modul Go (`go mod`)
Ketik perintah berikut untuk membuat file `go.mod` (pengenal modul proyek):
```bash
go mod init proyek-baru
```
*(Ganti `proyek-baru` dengan nama modul/proyek Anda).*

### Langkah 4: Install Pustaka / Dependensi Eksternal
Untuk membangun REST API yang terhubung ke database dan migrasi otomatis, jalankan perintah `go get` berikut beserta penjelasan fungsinya:

1. **Install Framework Gin Gonic (Web Server Framework)**
   ```bash
   go get -u "github.com/gin-gonic/gin"
   ```
   * **Fungsi:** Framework HTTP berkinerja tinggi untuk menangani routing URL, parameter, middleware CORS, dan pengiriman response berformat JSON dengan cepat dan mudah.

2. **Install Driver PostgreSQL (`lib/pq`)**
   ```bash
   go get -u "github.com/lib/pq"
   ```
   * **Fungsi:** Driver penghubung murni antara bahasa Go dengan database PostgreSQL sehingga fungsi standar `database/sql` bisa berkomunikasi langsung dengan server PostgreSQL.

3. **Install Library Migrasi SQL (`sql-migrate`)**
   ```bash
   go get -u "github.com/rubenv/sql-migrate"
   ```
   * **Fungsi:** Library untuk menjalankan migrasi skema tabel database secara terstruktur menggunakan file `.sql` (seperti `1_initiate.sql`).

4. **Install Library Pembaca File Konfigurasi `.env` (`godotenv`)**
   ```bash
   go get -u "github.com/joho/godotenv"
   ```
   * **Fungsi:** Library untuk memuat variabel rahasia atau konfigurasi (seperti password database, port, host) dari file `.env` eksternal agar kode lebih aman dan mudah dikonfigurasi.

### Langkah 5: Buat File `main.go`
Buat file bernama `main.go` di dalam folder proyek Anda, lalu isi dengan kode sederhana web server Gin berikut:
```go
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Membuat router default Gin
	r := gin.Default()

	// 2. Membuat endpoint GET sederhana
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// 3. Menjalankan server di port 8080
	r.Run(":8080")
}
```

### Langkah 6: Jalankan Program
Ketik perintah berikut di terminal:
```bash
go run main.go
```
Buka browser atau Postman, akses **`http://localhost:8080/ping`**, dan Anda akan melihat response JSON:
```json
{
  "message": "pong"
}
```

---

## 📁 4. Struktur Direktori Proyek

```text
testbab6/
├── controllers/
│   └── card_bridge_controllers.go      # Layer Pelayan: Mengatur request HTTP & format output JSON
├── databases/
│   ├── sql_migrations/
│   │   └── 1_initiate.sql              # Skrip SQL: Resep pembuatan tabel database awal
│   └── database.go                     # Layer Database: Menghubungkan & menjalankan migrasi otomatis
├── entities/
│   └── card_bridge_entity.go           # Layer Model/Entity: Bentuk cetakan objek data 'Card'
├── repositories/
│   └── card_bridge_repo.go             # Layer Koki/Dapur: Tempat perintah SQL (SELECT, INSERT, UPDATE, DELETE)
├── routers/
│   └── card_bridge_routers.go          # Layer Router/Menu: Menentukan jalan alamat URL (Endpoint) & CORS
├── go.mod                              # Catatan Identitas Modul & Daftar Library Eksternal
├── go.sum                              # Catatan Keamanan Checksum Library
├── main.go                             # Pintu Utama (Entry Point): Program pertama kali dinyalakan di sini
└── README.md                           # Dokumentasi proyek
```

---

## 📘 5. Penjelasan Rinci Kode Berdasarkan Urutan Alur (`Entity` ➔ `Migration` ➔ `Database` ➔ `Repository` ➔ `Controller` ➔ `Router` ➔ `Main`)

---

### 5.1. `entities/card_bridge_entity.go` (Wadah Data Kartu)

**Tugas File:** Mendefinisikan bentuk data kartu yang dipakai dalam program Go dan format teks JSON saat dikirim ke pemanggil API.

```go
package entities

type Card struct {
	ID string `json: "id"`
}
```

#### Penjelasan Detail Baris demi Baris (Gaya Pemula):
* ◦ `package entities` : Menandakan bahwa file ini tergabung dalam satu kelompok folder bernama `entities`. Layaknya sebuah laci khusus di lemari, ini membuat kode data terorganisir dengan rapi.
* ◦ `type Card struct` : Perintah untuk membuat cetakan data baru bernama `Card`. Huruf `C` besar di awal kata sangat penting dalam bahasa Go agar cetakan data ini bisa dibaca dan digunakan oleh file lain di luar folder ini.
* ◦ `ID string` : Kotak penyimpanan data di dalam cetakan `Card` bernama `ID` dengan jenis data teks (`string`).
* ◦ `json: "id"` : Petunjuk khusus untuk menterjemahkan variabel `ID` menjadi tulisan `"id"` berhuruf kecil ketika dikirim dalam bentuk format JSON kepada pengguna.

---

### 5.2. `databases/sql_migrations/1_initiate.sql` (Skrip Pembuat Tabel)

**Tugas File:** Membuat tabel `card` di PostgreSQL secara otomatis saat aplikasi pertama kali berjalan.

```sql
-- +migrate Up
-- +migrate StatementBegin

CREATE TABLE card (
    id varchar(20)
);

-- +migrate StatementEnd
```

#### Penjelasan Detail Baris demi Baris (Gaya Pemula):
* ◦ `-- +migrate Up` : Perintah rahasia (direktif) yang dibaca oleh alat migrasi database untuk memberitahu bahwa kode di bawah ini bertugas *membangun atau menambah* isi database.
* ◦ `-- +migrate StatementBegin` & `-- +migrate StatementEnd` : Batas awal dan akhir satu paket perintah SQL agar dijalankan secara utuh dan aman oleh sistem.
* ◦ `CREATE TABLE card (` : Perintah SQL dasar untuk membuat tabel baru di dalam database dengan nama `card`.
* ◦ `id varchar(20)` : Membuat sebuah kolom di dalam tabel bernama `id` dengan jenis teks dinamis (`varchar`) dengan batas maksimal 20 huruf/angka.

---

### 5.3. `databases/database.go` (Pengelola Koneksi DB & Migrasi Otomatis)

**Tugas File:** Membaca file `.sql`, mengeksekusi pembuatan tabel ke PostgreSQL secara otomatis, dan menyimpan koneksi database agar siap dipakai oleh fungsi lain.

```go
package database

import (
	"database/sql"
	"embed"
	"fmt"

	migrate "github.com/rubenv/sql-migrate"
)

//go:embed sql_migrations/*.sql
var dbMigration embed.FS
var DbConnection *sql.DB

func DBMigration(dbParam *sql.DB) {
	migrations := &migrate.EmbedFileSystemMigrationSource{
		FileSystem: dbMigration,
		Root:       "sql_migrations",
	}
	n, errs := migrate.Exec(dbParam, "postgres", migrations, migrate.Up)

	if errs != nil {
		panic(errs)
	}
	DbConnection = dbParam
	fmt.Println("Migrasi berhasil dilakukan", n, migrations)
}
```

#### Penjelasan Detail Baris demi Baris (Gaya Pemula):
* ◦ `package database` : Menandakan file ini adalah bagian dari kelompok pengelola database.
* ◦ `import (...)` : Memanggil alat-alat bantu dari luar (pustaka standar Go maupun pustaka tambahan) agar bisa digunakan di file ini.
* ◦ `//go:embed sql_migrations/*.sql` : Perintah ajaib di Go untuk membungkus dan menyembunyikan file-file `.sql` langsung ke dalam hasil akhir program, jadi kita tidak perlu membawa folder file `.sql` terpisah saat aplikasi dijalankan di komputer lain.
* ◦ `var dbMigration embed.FS` : Variabel penampung file `.sql` yang sudah dibungkus tadi ke dalam memori sistem file virtual.
* ◦ `var DbConnection *sql.DB` : Variabel global yang memegang kunci koneksi database agar bisa dipanggil kapan saja oleh file lain (`controllers` atau `repositories`).
* ◦ `func DBMigration(dbParam *sql.DB)` : Fungsi utama untuk menjalankan migrasi database. Menerima kiriman kunci koneksi database (`dbParam`) dari file `main.go`.
* ◦ `migrations := &migrate.EmbedFileSystemMigrationSource{...}` : Mengatur sumber folder file migrasi agar membaca dari variabel `dbMigration` di dalam folder `sql_migrations`.
* ◦ `n, errs := migrate.Exec(dbParam, "postgres", migrations, migrate.Up)` : Perintah utama untuk mengeksekusi migrasi ke database. `n` (angka integer) berisi jumlah file migrasi yang sukses dijalankan, dan `errs` (tipe error) akan berisi pesan kesalahan jika proses migrasi gagal.
* ◦ `"postgres"` : Memberitahu library migrasi bahwa jenis database yang kita tuju adalah PostgreSQL.
* ◦ `migrate.Up` : Arah migrasi maju (membangun/membuat tabel). Kebalikannya adalah `migrate.Down` untuk menghapus kembali tabel.
* ◦ `if errs != nil { panic(errs) }` : Pengaman jika migrasi gagal (misal koneksi salah), program akan langsung dihentikan secara paksa (`panic`) agar kita tahu letak kesalahannya.
* ◦ `DbConnection = dbParam` : Menyimpan koneksi database yang sudah berhasil diuji ke variabel global `DbConnection`.
* ◦ `fmt.Println(...)` : Mencetak tulisan sukses migrasi ke layar terminal.

---

### 5.4. `repositories/card_bridge_repo.go` (Perintah Query SQL Langsung)

**Tugas File:** Tempat berkumpulnya fungsi yang berhubungan langsung dengan database PostgreSQL (menjalankan query `SELECT`, `INSERT`, `UPDATE`, `DELETE`).

```go
package repositories

import (
	"bab6/entities"
	"database/sql"
)

// ... (fungsi GetCards, InsertCard, DeleteCard sudah dijelaskan sebelumnya)

// Contoh fungsi UPDATE yang akan ditambahkan:
func UpdateCard(db *sql.DB, card entities.Card) (err error) {
	sql := "UPDATE card SET id = $1 WHERE id = $2" // Contoh jika ingin update ID
	_, err = db.Exec(sql, card.ID, "id_lama") // Logika perlu disesuaikan dengan kebutuhan
	return err
}
```

#### Penjelasan Detail Baris demi Baris (Gaya Pemula):
* ◦ `package repositories` : Menandakan file ini berada di kelompok khusus repositori data (tempat koki meracik perintah SQL).
* ◦ `import (...)` : Memuat cetakan data `entities` dan pustaka SQL Go.
* ◦ `func GetCards(db *sql.DB) (result []entities.Card, err error)` : Fungsi untuk mengambil semua data kartu. Mengembalikan kumpulan data kartu (`result`) dan status error (`err`).
* ◦ `db.Query(sql)` : Perintah untuk menjalankan query pencarian/pembacaan data `SELECT` ke database.
* ◦ `defer rows.Close()` : Perintah pengaman agar setelah selesai membaca data, lembaran hasil bacaan langsung ditutup kembali agar RAM komputer tidak bocor.
* ◦ `rows.Next()` & `rows.Scan(&data.ID)` : Perulangan membaca baris data satu per satu dari database dan menyalinnya ke dalam variabel `data.ID`.
* ◦ `db.Exec(sql, card.ID)` : Perintah eksekusi untuk `INSERT`, `UPDATE`, atau `DELETE` data ke database. Menggunakan parameter `$1` agar aman dari serangan peretas (*SQL Injection*).

---

### 5.5. `controllers/card_bridge_controllers.go` (Pengatur Alur Request & Response HTTP)

**Tugas File:** Menerima input dari pengguna, meminta repository memproses ke database, dan mengembalikan jawaban dalam format teks JSON.

```go
// ... (fungsi GetCards, InsertCard, DeleteCard sudah dijelaskan sebelumnya)

// Contoh handler UPDATE (PUT) yang akan ditambahkan:
func UpdateCard(c *gin.Context) {
	var card entities.Card
	idLama := c.Param("id") // Mendapatkan ID yang ingin diubah
	
    // Menerima data baru dari body JSON (opsional, tergantung implementasi)
	if err := c.ShouldBindJSON(&card); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := repositories.UpdateCard(database.DbConnection, card) // Memanggil repositori

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Data Berhasil Diupdate"})
}
```

#### Penjelasan Detail Baris demi Baris (Gaya Pemula):
* ◦ `package controllers` : Menandakan file ini bertugas sebagai pelayan utama penghubung antara pengguna dan sistem dapur.
* ◦ `c *gin.Context` : Objek penghubung bawaan framework Gin yang membawa seluruh data permintaan (request) dari pengguna sekaligus alat untuk mengirim jawaban (response).
* ◦ `c.Param("id")` : Perintah untuk mengambil potongan teks penting dari alamat URL (misal mengambil ID kartu dari `/card/delete/CARD-01`).
* ◦ `c.ShouldBindJSON(&card)` : Perintah untuk membaca data berformat JSON yang dikirim oleh pengguna melalui Body request (biasanya dipakai pada metode `POST` atau `PUT`).
* ◦ `c.JSON(http.StatusOK, ...)` : Perintah untuk mengirim balasan kembali kepada pengguna berupa data berformat JSON dengan status sukses (`200 OK`).

---

### 5.6. `routers/card_bridge_routers.go` (Pendaftaran Alamat URL & CORS)

**Tugas File:** Mengatur daftar alamat URL (endpoint), menentukan HTTP Method (`GET`, `POST`, `PUT`, `DELETE`), dan mengatur akses izin domain eksternal (*CORS*).

```go
package routers

import (
	"bab6/controllers"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func StartServer() *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "DELETE", "PUT", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	router.GET("/cards", controllers.GetCards)
	router.POST("/card/input/:id", controllers.InsertCard)
	router.PUT("/card/update/:id", controllers.UpdateCard)
	router.DELETE("/card/delete/:id", controllers.DeleteCard)

	router.GET("/cards", controllers.GetCards)
	router.POST("/card/input/:id", controllers.InsertCard)
	router.PUT("/card/update/:id", controllers.UpdateCard)
	router.DELETE("/card/delete/:id", controllers.DeleteCard)

	return router
}
```

#### Penjelasan Detail Baris demi Baris (Gaya Pemula):
* ◦ `package routers` : Kelompok file yang bertugas sebagai buku daftar alamat rute jalan di aplikasi.
* ◦ `router := gin.Default()` : Perintah untuk menyalakan mesin pengatur alamat web berbasis Gin (lengkap dengan fitur pemantau log terminal dan pengaman agar server tidak gampang mati).
* ◦ `router.Use(cors.New(...))` : Pengaturan izin lintas domain (CORS) agar aplikasi frontend (seperti React/Vue atau Postman) bisa mengakses API ini tanpa diblokir oleh browser.
* ◦ `router.GET(...)`, `router.POST(...)`, `router.PUT(...)`, `router.DELETE(...)` : Peta jalur utama penghubung alamat URL ke fungsi pelayan (controller) yang tepat.

---

### 5.7. `main.go` (Pintu Masuk Utama Aplikasi)

**Tugas File:** File pertama yang dijalankan oleh komputer. Menyiapkan parameter database, membuka koneksi PostgreSQL, memanggil fungsi migrasi, dan menyalakan web server HTTP.

```go
package main

import (
	database "bab6/databases"
	"bab6/routers"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "nickytirta54" // SESUAIKAN DENGAN PASSWORD POSTGRE YANG TELAH DIDAFTARKAN
	dbName   = "mcsbab6"      // SESUAIKAN DENGAN NAMA DATABASE YANG DIBUAT
)

var (
	DB  *sql.DB
	err error
)

func main() {
	var PORT = ":50000" // ISI DENGAN PORT BEBAS (DAPAT DICARI DI GOOGLE DENGAN KEYWORD "UNUSED PORT RANGE")

	psqlInfo := fmt.Sprintf(
		`host=%s port=%d user=%s password=%s dbname=%s sslmode=disable`,
		host, port, user, password, dbName,
	)

	DB, err = sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Fatal("Error Open DB: ", err)
	}

	database.DBMigration(DB)

	defer DB.Close()

	routers.StartServer().Run(PORT)
	fmt.Printf("Success Connected")
}
```

#### Penjelasan Detail Baris demi Baris (Gaya Pemula):
* ◦ `package main` : Kelompok utama wajib bagi file Go yang bisa dijalankan langsung sebagai program mandiri.
* ◦ `_ "github.com/lib/pq"` : Panggilan khusus (side-effect import) dengan tanda garis bawah `_` di depannya. Berfungsi mendaftarkan sopir/driver PostgreSQL ke dalam sistem Go agar komputer tahu cara mengobrol dengan database PostgreSQL.
* ◦ `const (...)` : Kotak penyimpanan data tetap (konstanta) yang berisi alamat host, port, username, password, dan nama database PostgreSQL Anda.
* ◦ `sql.Open("postgres", psqlInfo)` : Perintah untuk mulai membuka pintu komunikasi/koneksi pertama kali ke server PostgreSQL menggunakan data koneksi (`psqlInfo`) yang sudah dirangkai.
* ◦ `database.DBMigration(DB)` : Memerintahkan program untuk masuk ke folder database dan menjalankan migrasi otomatis pembuatan tabel.
* ◦ `defer DB.Close()` : Perintah pengaman agar koneksi database ditutup secara otomatis saat program dimatikan.
* ◦ `routers.StartServer().Run(PORT)` : Perintah pamungkas untuk menyalakan server web HTTP di port `:50000` agar siap menerima kunjungan dari pengguna.

---

## 🚀 8. Cara Menjalankan Proyek Utama (`bab6`)

1. **Buat Database PostgreSQL**
   ```sql
   CREATE DATABASE mcsbab6;
   ```
2. **Sesuaikan Konfigurasi** di `main.go` (`password` dan `dbName`).
3. **Unduh Dependensi**
   ```bash
   go mod tidy
   ```
4. **Jalankan Aplikasi**
   ```bash
   go run main.go
   ```

---

## 🧪 9. Pengujian Endpoint API (Termasuk PUT)

Base URL Server: **`http://localhost:50000`**

### Tambahan: Memperbarui Data Kartu (`PUT`)
* **URL Endpoint**: `http://localhost:50000/card/update/:id`
* **Metode**: `PUT`
* **Perintah cURL** (dengan body JSON baru):
  ```bash
  curl -X PUT http://localhost:50000/card/update/CARD-001 \
  -H "Content-Type: application/json" \
  -d '{"id": "CARD-BARU"}`
  ```
* **Hasil Response (`200 OK`)**:
  ```json
  {
    "message": "Data Berhasil Diupdate"
  }
  ```
