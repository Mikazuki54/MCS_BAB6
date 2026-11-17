# Panduan & Dokumentasi Kode Lengkap REST API Card Bridge (`bab6`)
*Panduan Ramah Pemula & Orang Awam: Penjelasan Konsep, Kode Baris demi Baris, Variabel, Sintaks, dan Fungsi.*

---

## 💡 1. Konsep Dasar & Analogia Sederhana (Untuk Orang Awam)

Sebelum masuk ke kode, mari kita pahami apa yang dibuat oleh aplikasi ini dengan analogi sederhana: **Restoran**.

```
[ Pelanggan / Browser / Postman ]
               │
               ▼  1. Pesan Makanan (HTTP Request: GET, POST, DELETE)
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
  * **`DELETE`**: Menghapus data (seperti membatalkan pesanan).
* **Layered Architecture**: Memisah-misah kode ke beberapa folder (Router, Controller, Repository, Entity, Database) agar jika ada perubahan atau kerusakan, kita tahu persis folder mana yang harus diperbaiki tanpa mengganggu bagian lain.

---

## 📁 2. Struktur Direktori Proyek

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
│   └── card_bridge_repo.go             # Layer Koki/Dapur: Tempat perintah SQL (SELECT, INSERT, DELETE)
├── routers/
│   └── card_bridge_routers.go          # Layer Router/Menu: Menentukan jalan alamat URL (Endpoint) & CORS
├── go.mod                              # Catatan Identitas Modul & Daftar Library Eksternal
├── go.sum                              # Catatan Keamanan Checksum Library
├── main.go                             # Pintu Utama (Entry Point): Program pertama kali dinyalakan di sini
└── README.md                           # Panduan & Dokumentasi
```

---

## 📘 3. Penjelasan Rinci Kode Baris demi Baris & Sintaks

---

### 3.1. `entities/card_bridge_entity.go` (Wadah Data Kartu)

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

### 3.2. `databases/sql_migrations/1_initiate.sql` (Skrip Pembuat Tabel)

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

### 3.3. `databases/database.go` (Pengelola Koneksi DB & Migrasi Otomatis)

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

### 3.4. `repositories/card_bridge_repo.go` (Perintah Query SQL Langsung)

**Tugas File:** Tempat berkumpulnya fungsi yang berhubungan langsung dengan database PostgreSQL (menjalankan query `SELECT`, `INSERT`, `DELETE`).

```go
package repositories

import (
	"bab6/entities"
	"database/sql"
)

func GetCards(db *sql.DB) (result []entities.Card, err error) {
	sql := "SELECT * FROM card"
	rows, err := db.Query(sql)

	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var data entities.Card
		err = rows.Scan(&data.ID)
		if err != nil {
			return
		}
		result = append(result, data)
	}
	return
}

func InsertCard(db *sql.DB, card entities.Card) (err error) {
	sql := "INSERT INTO card(id) values($1)"
	_, err = db.Exec(sql, card.ID)
	return err
}

func DeleteCard(db *sql.DB, card entities.Card) (err error) {
	sql := "DELETE FROM card WHERE id = $1"
	_, err = db.Exec(sql, card.ID)
	return err
}
```

#### Penjelasan Sintaks, Variabel, & Fungsi:

##### 1. Fungsi `GetCards`:
* **Header `func GetCards(db *sql.DB) (result []entities.Card, err error)`**:
  * `db *sql.DB`: Parameter koneksi database.
  * `(result []entities.Card, err error)`: **Named Return Values**. Variabel `result` (kumpulan data `entities.Card` dalam bentuk slice `[]`) dan `err` (objek error) disiapkan sejak awal. Di akhir fungsi, kita cukup memanggil `return` tanpa menulis ulang nama variabelnya.
* **`sql := "SELECT * FROM card"`**: Perintah SQL untuk mengambil seluruh baris data dari tabel `card`.
* **`rows, err := db.Query(sql)`**: Menjalankan query `SELECT`. `rows` menyimpan penunjuk baris-baris data hasil pembacaan.
* **`defer rows.Close()`**: Kata kunci `defer` menjamin perintah `rows.Close()` dijalankan paling akhir saat fungsi `GetCards` selesai, sehingga koneksi memori database tidak bocor (*memory leak*).
* **Perulangan `for rows.Next()`**: Perulangan membaca baris demi baris data dari database.
  * `var data entities.Card`: Variabel sementara wadah 1 kartu.
  * `err = rows.Scan(&data.ID)`: Memindahkan/menyalin isi nilai kolom dari database ke variabel alamat memori `&data.ID`.
  * `result = append(result, data)`: Memasukkan data kartu yang baru dibaca ke dalam daftar list `result`.

##### 2. Fungsi `InsertCard`:
* **Header `func InsertCard(db *sql.DB, card entities.Card) (err error)`**: Menambahkan kartu baru.
* **`sql := "INSERT INTO card(id) values($1)"`**: Query memasukkan data. Simbol `$1` adalah **Placeholder Parameter** PostgreSQL. Ini mencegah celah keamanan **SQL Injection**.
* **`_, err = db.Exec(sql, card.ID)`**: Menjalankan perintah SQL `INSERT`. Simbol `_` (blank identifier) digunakan untuk mengabaikan nilai return metadata eksekusi SQL karena kita hanya butuh variabel status error-nya (`err`).

##### 3. Fungsi `DeleteCard`:
* **Header `func DeleteCard(db *sql.DB, card entities.Card) (err error)`**: Menghapus kartu berdasarkan ID.
* **`sql := "DELETE FROM card WHERE id = $1"`**: Perintah SQL menghapus baris tabel yang memiliki `id` sama dengan nilai `$1`.
* **`_, err = db.Exec(sql, card.ID)`**: Menjalankan proses penghapusan data kartu.

---

### 3.5. `controllers/card_bridge_controllers.go` (Pengatur Alur Request & Response HTTP)

**Tugas File:** Menerima input dari pengguna (melalui URL/HTTP), meminta repository memproses ke database, dan mengembalikan jawaban dalam format teks JSON.

```go
package controllers

import (
	database "bab6/databases"
	"bab6/entities"
	"bab6/repositories"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetCards(c *gin.Context) {
	var result gin.H

	card, err := repositories.GetCards(database.DbConnection)

	if err != nil {
		result = gin.H{
			"result": err.Error(),
		}
	} else {
		result = gin.H{
			"result": card,
		}
	}
	c.JSON(http.StatusOK, result)
}

func InsertCard(c *gin.Context) {
	var card entities.Card
	idCard := c.Param("id")
	card.ID = idCard

	err := repositories.InsertCard(database.DbConnection, card)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, card)
}

func DeleteCard(c *gin.Context) {
	var card entities.Card
	idCard := c.Param("id")
	card.ID = idCard

	err := repositories.DeleteCard(database.DbConnection, card)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Data Berhasil Dihapus", "id": idCard})
}
```

#### Penjelasan Sintaks, Variabel, & Fungsi:

##### 1. Fungsi `GetCards(c *gin.Context)`:
* **Parameter `c *gin.Context`**: Objek konteks dari framework Gin. Memuat semua data request yang masuk dan menyediakan metode untuk memberikan response balik.
* **`var result gin.H`**: `gin.H` adalah bentuk cepat (*shortcut*) buatan Gin untuk struktur data Map `map[string]interface{}` (digunakan untuk membentuk format objek JSON).
* **`card, err := repositories.GetCards(database.DbConnection)`**: Memanggil fungsi `GetCards` di layer repository dengan mengirimkan pointer database global.
* **Percabangan `if err != nil`**:
  * Jika ada error: `result` diisi pesan error `{ "result": "pesan error" }`.
  * Jika sukses: `result` diisi daftar data kartu `{ "result": [ {"id": "CARD01"} ] }`.
* **`c.JSON(http.StatusOK, result)`**: Mengirim response HTTP status `200 OK` dengan payload data JSON dari variabel `result`.

##### 2. Fungsi `InsertCard(c *gin.Context)`:
* **`var card entities.Card`**: Membuat wadah objek kartu kosong.
* **`idCard := c.Param("id")`**: Membaca isi parameter dinamis `:id` dari path URL (contoh URL `/card/input/C100` -> `idCard` bernilai `"C100"`).
* **`card.ID = idCard`**: Mengisikan teks ID kartu ke field `ID` struct `card`.
* **`err := repositories.InsertCard(...)`**: Memanggil fungsi repository untuk menyimpan kartu ke PostgreSQL.
* **`if err != nil`**: Jika gagal simpan, kirim HTTP Status Code `500 Internal Server Error` (`http.StatusInternalServerError`) dalam bentuk JSON `{ "error": "pesan error" }` dan langsung `return` (keluar dari fungsi).
* **`c.JSON(http.StatusOK, card)`**: Jika sukses simpan, kirim HTTP status `200 OK` dan tampilkan data kartu yang baru dibuat.

##### 3. Fungsi `DeleteCard(c *gin.Context)`:
* Mengambil parameter ID dari URL `c.Param("id")`, memanggil fungsi `repositories.DeleteCard(...)`, dan jika sukses mengembalikan response JSON `200 OK`:
  ```json
  {
    "id": "CARD-001",
    "message": "Data Berhasil Dihapus"
  }
  ```

---

### 3.6. `routers/card_bridge_routers.go` (Pendaftaran Alamat URL & CORS)

**Tugas File:** Mengatur daftar alamat URL (endpoint), menentukan HTTP Method (`GET`, `POST`, `DELETE`), dan mengatur akses izin domain eksternal (*CORS*).

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
	router.DELETE("/card/delete/:id", controllers.DeleteCard)

	return router
}
```

#### Penjelasan Sintaks, Variabel, & Fungsi:
1. **`func StartServer() *gin.Engine`**: Fungsi yang membuat dan mengembalikan objek mesin router Gin (`*gin.Engine`).
2. **`router := gin.Default()`**: Membuat router Gin standar yang otomatis dilengkapi 2 middleware bawaan:
   * **Logger**: Mencatat setiap aktivitas request yang masuk di terminal.
   * **Recovery**: Mencegah server mati (*crash*) saat terjadi kegagalan/panic yang tak terduga.
3. **`router.Use(cors.New(cors.Config{ ... }))`**: Memasang middleware CORS (Cross-Origin Resource Sharing).
   * `AllowOrigins: []string{"*"}`: Mengizinkan semua domain/aplikasi web mana pun mengakses API ini.
   * `AllowMethods`: Mengizinkan jenis metode HTTP `GET`, `POST`, `DELETE`, `PUT`, `OPTIONS`.
   * `AllowHeaders`: Mengizinkan header HTTP seperti `Content-Type` dan `Authorization`.
   * `MaxAge: 12 * time.Hour`: Menyimpan cache izin CORS di browser selama 12 jam.
4. **Pendaftaran Endpoint**:
   * **`router.GET("/cards", controllers.GetCards)`**: Jika ada request `GET` ke `/cards`, jalankan fungsi `controllers.GetCards`.
   * **`router.POST("/card/input/:id", controllers.InsertCard)`**: Jika ada request `POST` ke `/card/input/:id` (`:id` adalah variabel dinamis), jalankan `controllers.InsertCard`.
   * **`router.DELETE("/card/delete/:id", controllers.DeleteCard)`**: Jika ada request `DELETE` ke `/card/delete/:id`, jalankan `controllers.DeleteCard`.
5. **`return router`**: Mengembalikan objek router yang siap dipakai oleh `main.go`.

---

### 3.7. `main.go` (Pintu Masuk Utama Aplikasi)

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
6. **`psqlInfo := fmt.Sprintf(...)`**: Merangkai teks *Connection String* (DSN) format PostgreSQL:
   `host=localhost port=5432 user=postgres password=nickytirta54 dbname=mcsbab6 sslmode=disable` (`sslmode=disable` mematikan SSL karena dijalankan di server lokal).
7. **`DB, err = sql.Open("postgres", psqlInfo)`**: Inisialisasi awal koneksi database PostgreSQL. Jika terjadi kesalahan format parameter, `log.Fatal` akan mencetak pesan error dan menghentikan program.
8. **`database.DBMigration(DB)`**: Memanggil fungsi migrasi dari package `database`. Memastikan tabel `card` otomatis dibuat di PostgreSQL jika belum ada.
9. **`defer DB.Close()`**: Menjamin koneksi database ditutup dengan aman saat program berhenti.
10. **`routers.StartServer().Run(PORT)`**: Memanggil fungsi `StartServer()` dari `routers` dan menyalakan server di port `:50000`. Perintah ini bersifat *blocking* (terus berjalan mendengarkan request masuk sampai pengguna menekan `Ctrl+C`).

---

## 📖 4. Glosarium Istilah Pemula

| Istilah | Penjelasan Sederhana untuk Pemula |
| :--- | :--- |
| **Golang / Go** | Bahasa pemrograman buatan Google yang cepat, ringan, dan sangat cocok untuk membuat REST API & backend server. |
| **Gin Gonic** | Framework (alat bantu) di Go yang mempermudah pembuatan web server dan pengaturan URL/endpoint API. |
| **PostgreSQL** | Database relasional (RDBMS) populer untuk menyimpan data secara permanen di komputer/server. |
| **Struct** | Fitur Go untuk membuat wadah tipe data sendiri yang berisi campuran variabel (seperti objek). |
| **Pointer (`*`)** | Menyimpan alamat lokasi data di memori komputer, bukan menduplikat nilai datanya. |
| **Package** | Cara Go mengelompokkan file-file kode dalam satu folder agar rapi dan bisa dipakai ulang. |
| **Migration** | Cara otomatis untuk membuat/mengubah struktur tabel database menggunakan skrip SQL teratur. |
| **CORS** | Sistem keamanan browser yang mengatur apakah domain luar boleh mengakses API kita atau tidak. |
| **JSON** | Format teks standar berbentuk kunci-nilai (`"key": "value"`) yang digunakan untuk pertukaran data di internet. |

---

## 🚀 5. Panduan Menjalankan Proyek Langkah demi Langkah

### Langkah 1: Siapkan Database PostgreSQL
Buka aplikasi **pgAdmin** atau terminal `psql`, lalu jalankan perintah SQL berikut untuk membuat database:
```sql
CREATE DATABASE mcsbab6;
```

### Langkah 2: Periksa Konfigurasi Koneksi Database
Buka file `main.go`, pastikan konstanta `user`, `password`, dan `dbName` sudah sesuai dengan pengaturan PostgreSQL di komputer Anda:
```go
const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "PASSWORD_POSTGRES_ANDA"
	dbName   = "mcsbab6"
)
```

### Langkah 3: Unduh Dependensi Library Go
Buka terminal di folder proyek ini (`testbab6`), lalu ketik perintah:
```bash
go mod tidy
```

### Langkah 4: Jalankan Aplikasi
Jalankan aplikasi dengan perintah:
```bash
go run main.go
```

### Langkah 5: Tampilan terminal jika berhasil:
```text
Migrasi berhasil dilakukan 1 &{...}
[GIN-debug] GET    /cards                    --> bab6/controllers.GetCards (4 handlers)
[GIN-debug] POST   /card/input/:id           --> bab6/controllers.InsertCard (4 handlers)
[GIN-debug] DELETE /card/delete/:id          --> bab6/controllers.DeleteCard (4 handlers)
[GIN-debug] Listening and serving HTTP on :50000
```

---

## 🧪 6. Pengujian Endpoint API (Menggunakan cURL / Postman)

Base URL Server: **`http://localhost:50000`**

### 1. Tambah Kartu Baru (`POST`)
* **URL Endpoint**: `http://localhost:50000/card/input/CARD-001`
* **Metode**: `POST`
* **Perintah cURL**:
  ```bash
  curl -X POST http://localhost:50000/card/input/CARD-001
  ```
* **Hasil Response (`200 OK`)**:
  ```json
  {
    "id": "CARD-001"
  }
  ```

### 2. Tampilkan Semua Kartu (`GET`)
* **URL Endpoint**: `http://localhost:50000/cards`
* **Metode**: `GET`
* **Perintah cURL**:
  ```bash
  curl -X GET http://localhost:50000/cards
  ```
* **Hasil Response (`200 OK`)**:
  ```json
  {
    "result": [
      {
        "id": "CARD-001"
      }
    ]
  }
  ```

### 3. Hapus Kartu Berdasarkan ID (`DELETE`)
* **URL Endpoint**: `http://localhost:50000/card/delete/CARD-001`
* **Metode**: `DELETE`
* **Perintah cURL**:
  ```bash
  curl -X DELETE http://localhost:50000/card/delete/CARD-001
  ```
* **Hasil Response (`200 OK`)**:
  ```json
  {
    "id": "CARD-001",
    "message": "Data Berhasil Dihapus"
  }
  ```
