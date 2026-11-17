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

#### Penjelasan Sintaks, Variabel, & Kata Kunci:
1. **`package entities`**: Menandai bahwa kode ini milik kelompok (*package*) `entities`. File lain bisa memakai struktur ini dengan meng-import `"bab6/entities"`.
2. **`type Card struct`**:
   * **`type`**: Kata kunci Go untuk membuat tipe data baru.
   * **`Card`**: Nama dari tipe data yang kita buat. Menggunakan huruf awal **kapital (`C`)** artinya tipe ini bersifat **Public/Exported** (bisa dibaca oleh package lain luar `entities`).
   * **`struct`**: Koleksi kumpulan variabel (field) yang dibungkus menjadi satu kesatuan objek.
3. **`ID string`**:
   * **`ID`**: Nama properti/variabel di dalam `Card`. Diawali huruf kapital agar bisa dibaca oleh library JSON serializer dan package lain.
   * **`string`**: Tipe data teks (huruf/angka sebagai karakter). Sesuai dengan tipe kolom `varchar(20)` di database PostgreSQL.
4. **`` `json: "id"` ``**: Disebut **Struct Tag**. Ini petunjuk untuk converter JSON bawaan Go. Saat data `Card` diubah menjadi JSON, nama variabel `ID` akan ditulis sebagai `"id"` (huruf kecil).

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

#### Penjelasan Sintaks & Instruksi SQL:
1. **`-- +migrate Up`**: Komentar khusus direktif untuk library `sql-migrate`. Menandakan instruksi SQL di bawahnya dijalankan ketika proses migrasi aplikasi bergerak maju (*Up*).
2. **`-- +migrate StatementBegin` & `-- +migrate StatementEnd`**: Penanda batas awal dan akhir satu blok transaksi perintah SQL.
3. **`CREATE TABLE card (`**: Perintah SQL untuk membuat tabel baru bernama `card`.
4. **`id varchar(20)`**: Membuat kolom bernama `id` bertipe teks dinamis (`varchar`) dengan panjang maksimum 20 karakter.

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

#### Penjelasan Sintaks, Variabel, & Fungsi:
1. **`package database`**: Nama package pengelola database.
2. **`import (...)`**: Memuat library pendukung:
   * `"database/sql"`: Library standar Go untuk mengolah database SQL.
   * `"embed"`: Fitur canggih Go untuk memasukkan file fisik (seperti `.sql`) langsung ke dalam program executable binary.
   * `"fmt"`: Library untuk mencetak teks ke layar terminal.
   * `migrate "github.com/rubenv/sql-migrate"`: Library migrasi database otomatis yang di-alias menjadi nama `migrate`.
3. **`//go:embed sql_migrations/*.sql`**: **Compiler Directive** (perintah khusus compiler Go). Menginstruksikan Go membaca semua file `.sql` di folder `sql_migrations/` dan memasukkannya ke variabel di bawahnya.
4. **`var dbMigration embed.FS`**: Variabel bertipe `embed.FS` (*Embedded File System*) yang menampung isi file `.sql` di dalam memori program.
5. **`var DbConnection *sql.DB`**: Variabel global bertipe pointer database (`*sql.DB`). Menggunakan huruf awal kapital (`D`) agar variabel koneksi ini bisa dipanggil dari package lain (`database.DbConnection`).
6. **`func DBMigration(dbParam *sql.DB)`**: Fungsi migrasi database.
   * **Parameter `dbParam *sql.DB`**: Menerima objek koneksi PostgreSQL yang dikirim dari `main.go`.
   * **`migrations := &migrate.EmbedFileSystemMigrationSource{...}`**: Menyiapkan sumber file migrasi dari variabel `dbMigration` dalam folder `"sql_migrations"`.
   * **`n, errs := migrate.Exec(dbParam, "postgres", migrations, migrate.Up)`**: Eksekusi migrasi ke PostgreSQL.
     * `n`: Jumlah file migrasi SQL yang baru diterapkan.
     * `errs`: Menampung objek error jika migrasi gagal.
   * **`if errs != nil { panic(errs) }`**: Jika terjadi error (misal database belum ada atau SQL salah), `panic()` di-trigger untuk menghentikan program seketika dan menampilkan pesan kesalahannya.
   * **`DbConnection = dbParam`**: Menyimpan koneksi database yang valid ke variabel publik `DbConnection`.
   * **`fmt.Println(...)`**: Mencetak pesan konfirmasi ke terminal bahwa migrasi sukses.

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

#### Penjelasan Syntax untuk UPDATE (PUT):
* **`sql := "UPDATE card SET id = $1 WHERE id = $2"`**: Perintah SQL untuk memperbarui data. `SET` menentukan kolom yang diubah, `WHERE` menentukan data mana yang diubah.
* **`db.Exec(...)`**: Menjalankan perintah SQL tersebut ke database.

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

	return router
}
```

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

#### Penjelasan Sintaks, Variabel, & Fungsi:
1. **`package main` & `func main()`**: Syarat mutlak program Go yang dapat dieksekusi (*executable*). Titik awal eksekusi program dimulai dari sini.
2. **`_ "github.com/lib/pq"`**: Disebut **Side-Effect Import** (diawali tanda garis bawah `_`). Kita tidak memanggil fungsi package `pq` secara langsung, tetapi kita membutuhkan fungsi internal `init()` milik package `pq` agar driver database bernama `"postgres"` terdaftar secara resmi di package standar `database/sql`.
3. **Blok `const (...)`**: Mendefinisikan nilai variabel tetap (konstanta) untuk konfigurasi PostgreSQL:
   * `host = "localhost"`: Alamat komputer tempat PostgreSQL berada.
   * `port = 5432`: Nomor port standar PostgreSQL.
   * `user = "postgres"`: Username login PostgreSQL.
   * `password = "nickytirta54"`: Password user PostgreSQL.
   * `dbName = "mcsbab6"`: Nama database yang digunakan di PostgreSQL.
4. **Blok `var ( DB *sql.DB; err error )`**: Variabel global penampung koneksi SQL (`DB`) dan objek error (`err`).
5. **`var PORT = ":50000"`**: Variabel penyimpan nomor port server web HTTP tempat aplikasi akan berjalan (akses: `http://localhost:50000`).
6. **`psqlInfo := fmt.Sprintf(...)`**: Rangkaian teks *Connection String* (DSN) format PostgreSQL.
7. **`DB, err = sql.Open("postgres", psqlInfo)`**: Inisialisasi awal koneksi database PostgreSQL.
8. **`database.DBMigration(DB)`**: Memanggil fungsi migrasi dari package `database`. Memastikan tabel `card` otomatis dibuat di PostgreSQL.
9. **`defer DB.Close()`**: Menjamin koneksi database ditutup dengan aman saat program berhenti.
10. **`routers.StartServer().Run(PORT)`**: Memanggil fungsi `StartServer()` dari `routers` dan menyalakan server di port `:50000`.

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
  -d '{"id": "CARD-BARU"}'
  ```
* **Hasil Response (`200 OK`)**:
  ```json
  {
    "message": "Data Berhasil Diupdate"
  }
  ```
